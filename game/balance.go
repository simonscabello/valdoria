package game

import "math"

// Ferramenta de medição de balanceamento.
//
// Roda a lógica real do jogo (armas, projéteis, colisões, status elementais,
// chefe e linha do tempo das fases) sem janela, sem renderização e sem áudio,
// para que decisões de balanceamento sejam **medidas** e não estimadas.
//
// Motivação: o stun-lock da Lança de Luz atravessou 120 testes intactos porque
// nenhum teste pergunta se o jogo é justo — só se ele funciona. Esta ferramenta
// é a resposta a isso.
//
// Uso: go run ./cmd/balance

// staticBehavior mantém o alvo parado. Usado só na medição, para isolar o dano
// da arma do movimento do inimigo.
type staticBehavior struct{}

func (staticBehavior) update(*Enemy, *enemyContext) {}

// simTarget cria um alvo imóvel e praticamente indestrutível, para medir dano
// aplicado ao longo do tempo sem que ele morra no meio da amostra.
func simTarget(x, y, w, h float64) *Enemy {
	const huge = 1 << 30
	return &Enemy{
		kind: kindCrow, x: x, y: y, w: w, h: h,
		health: huge, behavior: staticBehavior{},
	}
}

// simGame monta um Game mínimo — sem áudio, sem cenário, sem HUD — capaz de
// rodar as colisões reais entre projéteis e inimigos.
func simGame(targets ...*Enemy) *Game {
	return &Game{
		player:     newPlayer(),
		formations: map[int]*formationTracker{},
		enemies:    targets,
	}
}

// WeaponMetrics reúne as medições de uma arma em um nível.
type WeaponMetrics struct {
	Weapon string
	Level  int

	// Os três cenários correspondem às três identidades de arma. Uma arma
	// assimétrica não deve ser igual às outras em tudo — ela deve vencer no
	// seu cenário e perder nos demais.

	// FocusDPS: alvo único e largo, longe (situação de chefe). Domínio da Luz.
	FocusDPS float64
	// SwarmDPS: cinco alvos pequenos lado a lado (formação). Domínio das Chamas.
	SwarmDPS float64
	// ColumnDPS: quatro alvos empilhados na vertical (fila descendo). Domínio
	// do Gelo, onde a perfuração paga.
	ColumnDPS float64

	// Cooldown em frames e projéteis por salva, para leitura rápida.
	Cooldown int
	Shots    int
}

// Peak devolve o melhor desempenho da arma entre os três cenários — é a medida
// justa para comparar armas com identidades diferentes.
func (m WeaponMetrics) Peak() float64 {
	return math.Max(m.FocusDPS, math.Max(m.SwarmDPS, m.ColumnDPS))
}

// BestScenario nomeia o cenário em que a arma rende mais.
func (m WeaponMetrics) BestScenario() string {
	switch m.Peak() {
	case m.FocusDPS:
		return "foco"
	case m.SwarmDPS:
		return "formacao"
	default:
		return "coluna"
	}
}

// simulateWeapon mede o dano de uma arma contra um conjunto de alvos parados
// durante `frames`, devolvendo o dano total aplicado.
func simulateWeapon(w weaponType, level int, targets []*Enemy, frames int) int {
	g := simGame(targets...)
	before := 0
	for _, t := range targets {
		before += t.health
	}

	ctx := &enemyContext{playerX: g.player.centerX(), playerY: g.player.centerY()}
	cooldown := 0
	for i := 0; i < frames; i++ {
		if cooldown <= 0 {
			g.bullets = append(g.bullets, fireWeapon(w, level, g.player.centerX(), g.player.y)...)
			cooldown = weaponCooldown(w, level)
		}
		cooldown--

		for _, b := range g.bullets {
			b.update()
			if b.offScreen() {
				b.dead = true
			}
		}
		g.bulletEnemyCollisions()
		// O update do inimigo aplica a queimadura ao longo do tempo.
		for _, e := range targets {
			e.update(ctx)
		}
		g.bullets = filterAlive(g.bullets, func(b *Bullet) bool { return !b.dead })
	}

	after := 0
	for _, t := range targets {
		after += t.health
	}
	return before - after
}

// measureWeapons devolve o perfil completo das três armas em todos os níveis.
func measureWeapons() []WeaponMetrics {
	const frames = 600 // 10 segundos
	out := make([]WeaponMetrics, 0, 9)

	for _, w := range []weaponType{weaponLight, weaponFlame, weaponIce} {
		for level := 1; level <= maxWeaponLevel; level++ {
			// Foco: um alvo do tamanho do chefe, no alto da tela.
			focus := []*Enemy{simTarget(ScreenWidth/2-bossW/2, bossY, bossW, bossH)}
			focusDmg := simulateWeapon(w, level, focus, frames)

			// Formação: cinco alvos pequenos lado a lado, na altura em que um
			// bando é normalmente interceptado.
			swarm := make([]*Enemy, 0, 5)
			for i := 0; i < 5; i++ {
				x := ScreenWidth/2 - crowSize/2 + float64(i-2)*formationGapX
				swarm = append(swarm, simTarget(x, 140, crowSize, crowSize))
			}
			swarmDmg := simulateWeapon(w, level, swarm, frames)

			// Coluna: quatro alvos empilhados na mesma faixa vertical — o caso
			// em que a perfuração converte um tiro em vários acertos.
			column := make([]*Enemy, 0, 4)
			for i := 0; i < 4; i++ {
				y := 60 + float64(i)*formationGapY
				column = append(column, simTarget(ScreenWidth/2-crowSize/2, y, crowSize, crowSize))
			}
			columnDmg := simulateWeapon(w, level, column, frames)

			out = append(out, WeaponMetrics{
				Weapon:    weaponName(w),
				Level:     level,
				FocusDPS:  float64(focusDmg) * 60 / frames,
				SwarmDPS:  float64(swarmDmg) * 60 / frames,
				ColumnDPS: float64(columnDmg) * 60 / frames,
				Cooldown:  weaponCooldown(w, level),
				Shots:     len(fireWeapon(w, level, 120, 200)),
			})
		}
	}
	return out
}

// BossMetrics mede quanto tempo o confronto com o chefe realmente dura e se ele
// chega a executar seus padrões.
type BossMetrics struct {
	Weapon string
	// FullUptime: segundos com o gatilho segurado o tempo todo (pior caso de
	// duração — o jogador ideal que nunca precisa desviar).
	FullUptime float64
	// Realistic: segundos com 50% de aproveitamento (o jogador real, que passa
	// metade do tempo desviando).
	Realistic float64
	// Patterns: quantas vezes o chefe iniciou um padrão de ataque antes de
	// morrer. Zero significa que o combate foi apagado.
	Patterns int
	// ReachedPhase3: se o combate chegou à fase final (cristais expostos).
	ReachedPhase3 bool
}

// simulateBoss roda o confronto completo com o chefe usando uma arma no nível
// máximo. `uptime` é a fração de frames em que o jogador dispara.
func simulateBoss(w weaponType, uptime float64, ascended bool) BossMetrics {
	g := simGame()
	g.boss = newBoss(ascended)
	g.player.weapon = w
	g.player.weaponLevel = maxWeaponLevel

	var enemies []*Enemy
	ctx := &bossContext{
		playerX: g.player.centerX(), playerY: g.player.centerY(),
		bullets: &g.enemyBullets, enemies: &enemies,
	}

	cooldown := 0
	patterns := 0
	lastAction := actionCooldown
	combatFrames := 0
	reachedP3 := false

	const limit = 60 * 60 * 10 // teto de 10 minutos, para nunca travar
	for i := 0; i < limit && g.boss.phase != bossDead; i++ {
		g.boss.update(ctx)
		if g.boss.phase == bossDying {
			continue
		}
		fighting := !g.boss.invulnerable
		if fighting {
			combatFrames++
			if g.boss.action == actionFiring && lastAction != actionFiring {
				patterns++
			}
			lastAction = g.boss.action
			if g.boss.phase == bossPhase3 {
				reachedP3 = true
			}
		}

		// Modelo de aproveitamento: dispara nos primeiros `uptime` de cada
		// segundo e desvia no resto.
		firing := fighting && float64(i%60) < uptime*60
		if cooldown > 0 {
			cooldown--
		}
		if firing && cooldown <= 0 {
			g.bullets = append(g.bullets, fireWeapon(w, maxWeaponLevel, g.player.centerX(), g.player.y)...)
			cooldown = weaponCooldown(w, maxWeaponLevel)
		}
		for _, b := range g.bullets {
			b.update()
			if b.offScreen() {
				b.dead = true
			}
		}
		g.bulletBossCollisions()
		g.bullets = filterAlive(g.bullets, func(b *Bullet) bool { return !b.dead })
	}

	return BossMetrics{
		Weapon:        weaponName(w),
		Realistic:     float64(combatFrames) / 60,
		Patterns:      patterns,
		ReachedPhase3: reachedP3,
	}
}

// measureBoss mede o confronto com cada arma, em aproveitamento total e real.
func measureBoss(ascended bool) []BossMetrics {
	out := make([]BossMetrics, 0, 3)
	for _, w := range []weaponType{weaponLight, weaponFlame, weaponIce} {
		full := simulateBoss(w, 1.0, ascended)
		real := simulateBoss(w, 0.5, ascended)
		name := weaponName(w)
		if ascended {
			name += " *"
		}
		out = append(out, BossMetrics{
			Weapon:        name,
			FullUptime:    full.Realistic,
			Realistic:     real.Realistic,
			Patterns:      real.Patterns,
			ReachedPhase3: real.ReachedPhase3,
		})
	}
	return out
}

// --- Ameaça: o inimigo é perigoso ou só demorado? ---

// ThreatMetrics separa as duas formas de um inimigo ser "difícil".
//
// Um inimigo **ameaçador** dispara uma ou duas vezes antes de morrer: ele te
// obriga a reagir. Uma **esponja** fica viva absorvendo tiros, o que não gera
// tensão nenhuma — só faz o jogador segurar o gatilho por mais tempo.
//
// Sem esta medição, "aumentar a dificuldade" degenera sempre em inflar vida, que
// é a forma preguiçosa e menos divertida de dificultar um jogo.
type ThreatMetrics struct {
	Kind   string
	Health int
	// TTK: segundos até morrer sob fogo focado da Luz no nível máximo. É um
	// piso — no jogo real o jogador divide o fogo e demora mais.
	TTK float64
	// Shots: quantos projéteis o inimigo consegue disparar antes de morrer.
	Shots int
	// Lifetime: segundos que ele levaria para atravessar a tela sem ser tocado.
	Lifetime float64
	// Pressure: fração do tempo de vida em que ele ainda estava vivo sob fogo.
	Pressure float64
}

// measureThreat mede cada inimigo isolado contra um jogador alinhado com fogo
// máximo — o melhor caso para o jogador, e portanto o piso do tempo de vida.
func measureThreat() []ThreatMetrics {
	out := make([]ThreatMetrics, 0, len(EnemyKinds()))
	for _, k := range EnemyKinds() {
		out = append(out, measureThreatFor(k))
	}
	return out
}

func measureThreatFor(k enemyKind) ThreatMetrics {
	m := ThreatMetrics{Kind: EnemyKindName(k)}

	// Tempo de travessia sem interferência, para saber quanto da vida do
	// inimigo o jogador realmente consome.
	free := spawnEnemy(k, ScreenWidth/2, 0, true)
	m.Lifetime = float64(framesUntilGone(free)) / 60

	e := spawnEnemy(k, ScreenWidth/2, 0, true)
	m.Health = e.health

	g := simGame(e)
	g.player.weaponLevel = maxWeaponLevel
	var shots []*EnemyBullet
	ctx := &enemyContext{playerX: g.player.centerX(), playerY: g.player.centerY(), bullets: &shots}

	cooldown := 0
	frames := 0
	const limit = 60 * 60
	for ; frames < limit && !e.dead; frames++ {
		// O jogador acompanha o alvo na horizontal: fogo focado, o melhor caso.
		g.player.x = e.centerX() - playerSize/2
		g.player.clampToScreen()

		if cooldown > 0 {
			cooldown--
		}
		if cooldown <= 0 {
			g.bullets = append(g.bullets,
				fireWeapon(weaponLight, maxWeaponLevel, g.player.centerX(), g.player.y)...)
			cooldown = weaponCooldown(weaponLight, maxWeaponLevel)
		}
		for _, b := range g.bullets {
			b.update()
			if b.offScreen() {
				b.dead = true
			}
		}
		g.bulletEnemyCollisions()
		g.bullets = filterAlive(g.bullets, func(b *Bullet) bool { return !b.dead })
		e.update(ctx)
		if e.offScreen() {
			break
		}
	}
	m.TTK = float64(frames) / 60
	m.Shots = len(shots)
	if m.Lifetime > 0 {
		m.Pressure = m.TTK / m.Lifetime
	}
	return m
}

// framesUntilGone conta quantos frames um inimigo intocado leva para sair de
// cena. Feiticeiros param no alto e nunca saem: devolvem o teto.
func framesUntilGone(e *Enemy) int {
	var shots []*EnemyBullet
	ctx := &enemyContext{playerX: ScreenWidth / 2, playerY: ScreenHeight - 40, bullets: &shots}
	const limit = 60 * 30
	for i := 0; i < limit; i++ {
		e.update(ctx)
		if e.offScreen() {
			return i
		}
	}
	return limit
}

// --- Pressão: quão difícil o jogo realmente é ---

// PlayerModel descreve como o jogador simulado se comporta. Não é uma IA boa —
// é o piso: se até o modelo mais burro sobrevive, o jogo é fácil demais.
type PlayerModel int

const (
	// modelPassive nunca se move. Mede a pressão bruta que chega sozinha até o
	// jogador: um jogo minimamente exigente não deveria deixá-lo terminar.
	modelPassive PlayerModel = iota
	// modelAverage desliza na horizontal atrás do inimigo mais próximo, sem
	// nunca desviar de projéteis — o jogador casual dos primeiros minutos.
	modelAverage
)

func (m PlayerModel) String() string {
	if m == modelAverage {
		return "mediano"
	}
	return "passivo"
}

// PressureMetrics é o retrato da dificuldade sentida por um modelo de jogador.
type PressureMetrics struct {
	Model      string
	Hits       int // golpes que realmente tiraram vida
	Lives      int // vidas perdidas
	Died       bool
	Escapes    int
	Corruption float64
	// PeakBullets: maior número de projéteis inimigos simultâneos.
	PeakBullets int
	// BulletsFired: total de projéteis inimigos disparados na campanha.
	BulletsFired int
	// SurvivedSeconds: até onde o modelo chegou.
	SurvivedSeconds float64
}

// simulatePressure roda a campanha inteira com um modelo de jogador, usando a
// lógica real de movimento de inimigos, disparos, colisões e invencibilidade.
//
// O jogador simulado atira sem parar com a Luz no nível máximo — ou seja, é
// generoso com ele em poder de fogo e implacável em habilidade de desvio.
// Se mesmo assim ele não morre, a dificuldade não vem do design, vem do acaso.
func simulatePressure(model PlayerModel) PressureMetrics {
	m := PressureMetrics{Model: model.String()}

	g := simGame()
	g.player.weaponLevel = maxWeaponLevel
	g.lives = diffParams().startLives
	g.mode = modeCampaign

	cooldown := 0
	tick := 0
	for _, def := range campaignStages() {
		l := newLevelFromStage(def)
		for i := 0; i <= l.duration; i++ {
			tick++
			for _, e := range l.update() {
				g.spawn(e)
			}

			if model == modelAverage {
				trackNearestEnemy(g)
			}
			if g.player.invincible > 0 {
				g.player.invincible--
			}

			if cooldown > 0 {
				cooldown--
			}
			if cooldown <= 0 {
				g.bullets = append(g.bullets,
					fireWeapon(g.player.weapon, g.player.weaponLevel, g.player.centerX(), g.player.y)...)
				cooldown = weaponCooldown(g.player.weapon, g.player.weaponLevel)
			}

			ctx := &enemyContext{
				playerX: g.player.centerX(), playerY: g.player.centerY(),
				bullets: &g.enemyBullets,
			}
			before := len(g.enemyBullets)
			for _, e := range g.enemies {
				e.update(ctx)
				if e.offScreen() {
					e.escaped, e.escapedBottom, e.dead = true, e.crossedBottom(), true
				}
			}
			m.BulletsFired += len(g.enemyBullets) - before

			for _, b := range g.bullets {
				b.update()
				if b.offScreen() {
					b.dead = true
				}
			}
			for _, b := range g.enemyBullets {
				b.update()
				if b.offScreen() {
					b.dead = true
				}
			}

			g.bulletEnemyCollisions()
			if g.player.canBeHit() {
				if hit := pressureHit(g); hit {
					m.Hits++
					if g.player.health <= 0 {
						m.Lives++
						g.lives--
						if g.lives <= 0 {
							m.Died = true
							m.SurvivedSeconds = float64(tick) / 60
							m.Escapes = g.corruption.escaped
							m.Corruption = g.corruption.value
							m.PeakBullets = maxInt(m.PeakBullets, len(g.enemyBullets))
							return m
						}
						g.player.respawn()
						g.enemyBullets = g.enemyBullets[:0]
					}
				}
			}

			for _, e := range g.enemies {
				if e.dead && e.escapedBottom {
					g.corruption.add(e.kind)
				}
			}
			g.bullets = filterAlive(g.bullets, func(b *Bullet) bool { return !b.dead })
			g.enemies = filterAlive(g.enemies, func(e *Enemy) bool { return !e.dead })
			g.enemyBullets = filterAlive(g.enemyBullets, func(b *EnemyBullet) bool { return !b.dead })
			m.PeakBullets = maxInt(m.PeakBullets, len(g.enemyBullets))
		}
	}

	m.SurvivedSeconds = float64(tick) / 60
	m.Escapes = g.corruption.escaped
	m.Corruption = g.corruption.value
	return m
}

// pressureHit aplica dano de projétil e de colisão, devolvendo true se a vida
// realmente caiu. Reusa as mesmas regras do jogo (um golpe por frame).
func pressureHit(g *Game) bool {
	hx, hy, hw, hh := g.player.hitbox()
	for _, b := range g.enemyBullets {
		if b.dead {
			continue
		}
		if collides(hx, hy, hw, hh, b.x, b.y, enemyBulletSize, enemyBulletSize) {
			b.dead = true
			return g.player.hit(enemyBulletDamage)
		}
	}
	for _, e := range g.enemies {
		if e.dead {
			continue
		}
		if collides(hx, hy, hw, hh, e.x, e.y, e.w, e.h) {
			e.dead = true
			return g.player.hit(e.damage)
		}
	}
	return false
}

// trackNearestEnemy desliza o jogador na horizontal atrás do alvo mais próximo,
// como faz quem ainda não aprendeu a desviar: mira, mas não foge.
func trackNearestEnemy(g *Game) {
	var target *Enemy
	best := 1e9
	for _, e := range g.enemies {
		if e.dead {
			continue
		}
		if d := math.Abs(e.centerX() - g.player.centerX()); d < best {
			best, target = d, e
		}
	}
	if target == nil {
		return
	}
	dx := target.centerX() - g.player.centerX()
	if math.Abs(dx) < playerSpeed {
		return
	}
	if dx > 0 {
		g.player.x += playerSpeed
	} else {
		g.player.x -= playerSpeed
	}
	g.player.clampToScreen()
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// measurePressure mede a campanha contra os dois modelos de jogador.
func measurePressure() []PressureMetrics {
	return []PressureMetrics{
		simulatePressure(modelPassive),
		simulatePressure(modelAverage),
	}
}

// CorruptionMetrics descreve para onde o Medidor de Corrupção leva uma partida,
// dado um perfil de jogador (a fração de inimigos que ele deixa passar).
type CorruptionMetrics struct {
	// LeakRate: fração dos inimigos que atravessam a base.
	LeakRate float64
	// Final: corrupção ao fim da campanha.
	Final float64
	Tier  string
	// Escaped: quantos inimigos passaram.
	Escaped int
	// Mutated: quantos inimigos nasceram corrompidos por causa disso.
	Mutated int
	// ScoreMul: multiplicador de pontos vigente no fim da campanha.
	ScoreMul float64
	// AscendedBoss: se o jogador desbloqueou o confronto verdadeiro.
	AscendedBoss bool
	// TierAt: em que segundo cada faixa foi atingida (-1 = nunca).
	TierAt []float64
}

// simulateCorruption roda a campanha deixando escapar `leak` dos inimigos que
// podem corromper, e reporta onde o medidor termina.
//
// É a medição que calibra os pesos: um jogador competente precisa terminar
// baixo, e um jogador que ignora o jogo precisa colapsar.
func simulateCorruption(leak float64) CorruptionMetrics {
	m := CorruptionMetrics{LeakRate: leak, TierAt: make([]float64, corruptionTierCount)}
	for i := range m.TierAt {
		m.TierAt[i] = -1
	}
	m.TierAt[tierSteady] = 0

	var c Corruption
	tick := 0
	for _, def := range campaignStages() {
		l := newLevelFromStage(def)
		for i := 0; i <= l.duration; i++ {
			tick++
			for _, e := range l.update() {
				if c.mutates() {
					m.Mutated++
				}
				if corruptionWeight(e.kind) <= 0 || rng.Float64() >= leak {
					continue
				}
				before := c.tier
				c.add(e.kind)
				m.Escaped++
				if c.tier != before && m.TierAt[c.tier] < 0 {
					m.TierAt[c.tier] = float64(tick) / 60
				}
			}
		}
	}
	m.Final = c.value
	m.Tier = c.tierName()
	m.ScoreMul = c.scoreMul()
	m.AscendedBoss = c.fallen()
	return m
}

// measureCorruption traça a curva do medidor para perfis de jogador diferentes.
func measureCorruption() []CorruptionMetrics {
	// 0% = jogador perfeito · 30% = jogador que ignora metade das ameaças.
	rates := []float64{0.02, 0.05, 0.10, 0.15, 0.25, 0.40}
	out := make([]CorruptionMetrics, 0, len(rates))
	for _, r := range rates {
		out = append(out, simulateCorruption(r))
	}
	return out
}

// StageMetrics descreve o que uma fase realmente gera, e não o que ela parece
// gerar ao ler as ondas.
type StageMetrics struct {
	Name        string
	Waves       int
	Enemies     int
	TotalHealth int
	ByKind      map[enemyKind]int
	LastSpawn   int
	Seconds     float64
	// PeakOnScreen: maior número de inimigos vivos simultaneamente, medido
	// deixando todos sobreviverem (aproximação do pior caso de densidade).
	PeakOnScreen int
}

// measureStage roda a linha do tempo completa de uma fase, contando o que ela
// gera por tipo e acompanhando a ocupação da tela.
func measureStage(def *stageDef) StageMetrics {
	l := newLevelFromStage(def)
	m := StageMetrics{
		Name:   def.name,
		Waves:  len(def.waves),
		ByKind: map[enemyKind]int{},
	}

	var alive []*Enemy
	var shots []*EnemyBullet // descartado: aqui só interessa o fluxo de inimigos
	ctx := &enemyContext{playerX: ScreenWidth / 2, playerY: ScreenHeight - 40, bullets: &shots}
	for tick := 0; tick <= l.duration; tick++ {
		spawned := l.update()
		if len(spawned) > 0 {
			m.LastSpawn = tick
		}
		for _, e := range spawned {
			m.Enemies++
			m.TotalHealth += e.health
			m.ByKind[e.kind]++
		}
		alive = append(alive, spawned...)

		// Deixa os inimigos viverem e saírem naturalmente, para medir ocupação.
		for _, e := range alive {
			e.update(ctx)
			if e.offScreen() {
				e.dead = true
			}
		}
		alive = filterAlive(alive, func(e *Enemy) bool { return !e.dead })
		if len(alive) > m.PeakOnScreen {
			m.PeakOnScreen = len(alive)
		}
	}
	m.Seconds = float64(m.LastSpawn) / 60
	return m
}

// ProgressionMetrics mede a curva de poder dentro de uma partida: quanto tempo
// o jogador leva para chegar ao teto e quantos power-ups caem depois disso.
type ProgressionMetrics struct {
	TotalEnemies int
	// GuaranteedDrops: runas garantidas por onda (a espinha da progressão).
	GuaranteedDrops int
	// RunesToMax: runas do elemento equipado necessárias para o nível máximo.
	RunesToMax int
	// KillsToMax / SecondsToMax: onde, na campanha, o teto de poder é atingido.
	KillsToMax   int
	SecondsToMax float64
	// TotalDrops: runas efetivamente coletadas numa run completa.
	TotalDrops int
	// WastedRunes: runas elementais coletadas já no teto (sem efeito).
	WastedRunes     int
	CampaignSeconds float64
}

// measureProgression simula a campanha inteira coletando os drops reais — os
// garantidos por onda e os aleatórios — sobre um Player de verdade, e registra
// em que ponto o teto de poder é atingido.
//
// A conta algébrica não serve aqui: trocar de elemento zera a carga parcial, e
// as runas garantidas alternam de elemento de propósito. Só simulando dá para
// saber quando o jogador realmente chega ao topo.
func measureProgression(stages []StageMetrics) ProgressionMetrics {
	m := ProgressionMetrics{RunesToMax: (maxWeaponLevel - 1) * runesPerLevel}
	for _, s := range stages {
		m.CampaignSeconds += s.Seconds
	}

	g := simGame()
	p := g.player
	tick, maxTick := 0, -1
	chance := scaledDropChance(dropChance)

	for _, def := range campaignStages() {
		for _, w := range def.waves {
			if w.hasDrop {
				m.GuaranteedDrops++
			}
		}
		l := newLevelFromStage(def)
		for i := 0; i <= l.duration; i++ {
			tick++
			// Aproximação: cada inimigo gerado é abatido (é a curva de poder de
			// um jogador que não deixa nada escapar — o melhor caso).
			for _, e := range l.update() {
				m.TotalEnemies++

				kind, dropped := e.drop, e.hasDrop
				if !dropped && rng.Float64() < chance {
					kind, dropped = g.randomRune(), true
				}
				if !dropped {
					continue
				}
				m.TotalDrops++
				// Inerte de verdade é a runa do elemento **já equipado** no teto
				// de poder: ela não sobe nível nem troca de magia. Uma runa de
				// outro elemento no teto ainda vale — é como o jogador escolhe
				// com qual magia vai encarar o chefe.
				if p.weaponLevel >= maxWeaponLevel && isRuneOfEquipped(p, kind) {
					m.WastedRunes++
				}
				p.applyPowerup(kind)
				if maxTick < 0 && p.weaponLevel >= maxWeaponLevel {
					maxTick, m.KillsToMax = tick, m.TotalEnemies
				}
			}
		}
	}
	if maxTick >= 0 {
		m.SecondsToMax = float64(maxTick) / 60
	} else {
		m.SecondsToMax = m.CampaignSeconds // nunca chegou ao teto
	}
	return m
}

// isElementRune diz se a runa afeta a arma equipada.
func isElementRune(k powerupType) bool {
	return k == powerLight || k == powerFire || k == powerIce
}

// isRuneOfEquipped diz se a runa é do mesmo elemento que o jogador já usa.
func isRuneOfEquipped(p *Player, k powerupType) bool {
	if !isElementRune(k) {
		return false
	}
	switch k {
	case powerFire:
		return p.weapon == weaponFlame
	case powerIce:
		return p.weapon == weaponIce
	default:
		return p.weapon == weaponLight
	}
}

// BalanceReport é o relatório completo de uma medição.
type BalanceReport struct {
	Difficulty   string
	Weapons      []WeaponMetrics
	Boss         []BossMetrics
	BossAscended []BossMetrics
	Stages       []StageMetrics
	Progression  ProgressionMetrics
	Corruption   []CorruptionMetrics
	Pressure     []PressureMetrics
	Threat       []ThreatMetrics
	Totals       StageMetrics // agregado da campanha
}

// MeasureBalance roda a bateria completa de medições na dificuldade indicada.
// É determinística: fixa a semente antes de começar.
func MeasureBalance(diff difficultyLevel, seed int64) BalanceReport {
	prevDiff := currentDifficulty
	setDifficulty(diff)
	SetSeed(seed)
	defer setDifficulty(prevDiff)

	r := BalanceReport{
		Difficulty:   difficultyName(diff),
		Weapons:      measureWeapons(),
		Boss:         measureBoss(false),
		BossAscended: measureBoss(true),
	}

	total := StageMetrics{Name: "CAMPANHA", ByKind: map[enemyKind]int{}}
	for _, def := range campaignStages() {
		s := measureStage(def)
		r.Stages = append(r.Stages, s)
		total.Waves += s.Waves
		total.Enemies += s.Enemies
		total.TotalHealth += s.TotalHealth
		total.Seconds += s.Seconds
		for k, n := range s.ByKind {
			total.ByKind[k] += n
		}
		if s.PeakOnScreen > total.PeakOnScreen {
			total.PeakOnScreen = s.PeakOnScreen
		}
	}
	r.Totals = total
	r.Progression = measureProgression(r.Stages)
	r.Corruption = measureCorruption()
	r.Pressure = measurePressure()
	r.Threat = measureThreat()
	return r
}

// EnemyKindName devolve o nome legível de um tipo de inimigo (para relatórios).
func EnemyKindName(k enemyKind) string {
	switch k {
	case kindHarpy:
		return "Harpia"
	case kindGargoyle:
		return "Gargula"
	case kindWyvern:
		return "Wyvern"
	case kindBallista:
		return "Balista"
	case kindMage:
		return "Feiticeiro"
	default:
		return "Corvo"
	}
}

// EnemyKinds devolve os tipos na ordem canônica, para relatórios estáveis.
func EnemyKinds() []enemyKind {
	return []enemyKind{kindCrow, kindHarpy, kindGargoyle, kindWyvern, kindBallista, kindMage}
}

// Share devolve a fração que um tipo representa no total de inimigos gerados.
func (m StageMetrics) Share(k enemyKind) float64 {
	if m.Enemies == 0 {
		return 0
	}
	return float64(m.ByKind[k]) / float64(m.Enemies)
}

// ParseDifficulty converte o nome de uma dificuldade no preset correspondente.
func ParseDifficulty(s string) (difficultyLevel, bool) {
	switch s {
	case "facil", "fácil", "easy":
		return diffEasy, true
	case "normal":
		return diffNormal, true
	case "dificil", "difícil", "hard":
		return diffHard, true
	}
	return diffNormal, false
}

// targetShare é o alvo de composição do bestiário definido na direção do jogo:
// nenhum inimigo domina, e os quatro tipos "pesados" deixam de ser decorativos.
// Serve de fonte única para o relatório e para o teste de regressão.
var targetShare = map[enemyKind]float64{
	kindCrow:     0.35,
	kindHarpy:    0.22,
	kindGargoyle: 0.12,
	kindWyvern:   0.10,
	kindMage:     0.10,
	kindBallista: 0.08,
}

// TargetShareLabel devolve a participação-alvo de um tipo, para relatórios.
func TargetShareLabel(k enemyKind) string {
	return itoa(int(targetShare[k]*100+0.5)) + "%"
}

// Check é um critério objetivo de balanceamento, verificado sobre a medição.
type Check struct {
	Name   string
	Detail string
	Pass   bool
}

// Bandas de aceitação do balanceamento. Ficam aqui, e não espalhadas por
// testes, para que "o que é equilibrado" seja uma decisão explícita e única.
const (
	maxPeakSpread   = 1.6  // pico da arma mais forte / pico da mais fraca
	minBossSeconds  = 60.0 // confronto curto demais = clímax apagado
	maxBossSeconds  = 180.0
	minBossPatterns = 8
	maxCrowShare    = 0.40 // acima disso a campanha vira um inimigo só
	minHeavyShare   = 0.06 // abaixo disso o inimigo é decorativo
	// A curva de poder precisa cobrir a run: o teto não pode chegar nos primeiros
	// segundos (chegava aos 3% — 7s de 210s), e as runas não podem virar lixo
	// depois dele.
	//
	// 20% é o teto realista *desta versão*: com três níveis de arma e ~30 runas
	// por partida, nenhuma calibragem de custo faz a progressão durar mais. O
	// alvo sobe para 60% quando as Runas Fundidas (v0.6) acrescentarem níveis.
	minPowerCurvePct = 0.20
	maxWastedRunes   = 0.20 // fração das runas coletadas que caem inertes

	// Um jogador que deixa passar até 5% dos inimigos não pode ser empurrado
	// para além da primeira faixa: a corrupção precisa ser consequência de
	// negligência real, nunca de imprecisão normal.
	maxGoodCorruption = 26.0

	// Teto de tempo-para-matar sob fogo focado. Acima disso o inimigo deixa de
	// ser uma ameaça e vira uma tarefa: o jogador segura o gatilho e espera.
	maxEnemyTTK = 2.0
)

// Checks avalia o relatório contra os critérios de design da direção do jogo.
func Checks(r BalanceReport) []Check {
	out := []Check{}

	spread := PeakSpread(r.Weapons, maxWeaponLevel)
	out = append(out, Check{
		Name:   "Toda arma brilha no seu cenario",
		Detail: fmtRatio("desequilibrio entre os picos", spread, maxPeakSpread),
		Pass:   spread > 0 && spread <= maxPeakSpread,
	})

	out = append(out, checkIdentities(r), checkNoDominance(r))

	for _, b := range r.Boss {
		pass := b.Realistic >= minBossSeconds && b.Realistic <= maxBossSeconds &&
			b.Patterns >= minBossPatterns && b.ReachedPhase3
		out = append(out, Check{
			Name: "Chefe sustenta o confronto com " + b.Weapon,
			Detail: itoaFixed(b.Realistic) + "s realistas (alvo " +
				itoaFixed(minBossSeconds) + "-" + itoaFixed(maxBossSeconds) + "s), " +
				itoa(b.Patterns) + " padroes (min " + itoa(minBossPatterns) + ")",
			Pass: pass,
		})
	}

	crow := r.Totals.Share(kindCrow)
	out = append(out, Check{
		Name:   "A campanha nao e um inimigo so",
		Detail: "corvo em " + pct(crow) + " (teto " + pct(maxCrowShare) + ")",
		Pass:   crow <= maxCrowShare,
	})

	for _, k := range []enemyKind{kindGargoyle, kindWyvern, kindMage, kindBallista} {
		s := r.Totals.Share(k)
		out = append(out, Check{
			Name:   EnemyKindName(k) + " tem presenca real",
			Detail: pct(s) + " da campanha (minimo " + pct(minHeavyShare) + ", alvo " + TargetShareLabel(k) + ")",
			Pass:   s >= minHeavyShare,
		})
	}

	p := r.Progression
	frac := 0.0
	if p.CampaignSeconds > 0 {
		frac = p.SecondsToMax / p.CampaignSeconds
	}
	out = append(out, Check{
		Name:   "A curva de poder cobre a run",
		Detail: "teto de poder em " + pct(frac) + " da run (minimo " + pct(minPowerCurvePct) + ")",
		Pass:   frac >= minPowerCurvePct,
	})

	out = append(out, checkCorruptionCurve(r)...)
	out = append(out, checkPressure(r)...)
	out = append(out, checkThreat(r)...)

	wasted := 0.0
	if p.TotalDrops > 0 {
		wasted = float64(p.WastedRunes) / float64(p.TotalDrops)
	}
	out = append(out, Check{
		Name: "As runas continuam valendo algo",
		Detail: itoa(p.WastedRunes) + " de " + itoa(p.TotalDrops) + " runas caem inertes (" +
			pct(wasted) + ", teto " + pct(maxWastedRunes) + ")",
		Pass: wasted <= maxWastedRunes,
	})

	return out
}

// checkPressure valida que o jogo exige do jogador.
//
// Existe porque a facilidade é o defeito mais difícil de perceber lendo código:
// a campanha inteira gerava 57 projéteis inimigos e um jogador que nunca
// desviava a terminava sem perder uma vida — os inimigos morriam antes de
// conseguir atirar. Um alvo que não revida não é um inimigo.
func checkPressure(r BalanceReport) []Check {
	out := make([]Check, 0, len(r.Pressure))
	for _, p := range r.Pressure {
		out = append(out, Check{
			Name: "Quem nao desvia nao termina a campanha (" + p.Model + ")",
			Detail: itoa(p.Hits) + " golpes, " + itoa(p.Lives) + " vidas perdidas, " +
				"parou aos " + itoaFixed(p.SurvivedSeconds) + "s de " +
				itoaFixed(r.Progression.CampaignSeconds) + "s",
			Pass: p.Died,
		})
	}
	return out
}

// checkThreat separa ameaça de esponja.
//
// Existe por causa de um segundo relato de playtest: depois de corrigir a
// facilidade, "ficou difícil, mas eles demoram". As duas coisas têm a mesma
// aparência num relatório de vida dos inimigos e são problemas opostos — só
// medindo tiros disparados e tempo-para-matar dá para saber qual é qual.
func checkThreat(r BalanceReport) []Check {
	out := make([]Check, 0, 2)

	slowest := ThreatMetrics{}
	silent := []string{}
	for _, t := range r.Threat {
		if t.TTK > slowest.TTK {
			slowest = t
		}
		// O corvo é a massa descartável do jogo: ele não atira de propósito.
		if t.Shots == 0 && t.Kind != EnemyKindName(kindCrow) {
			silent = append(silent, t.Kind)
		}
	}

	out = append(out, Check{
		Name:   "Nenhum inimigo vira esponja",
		Detail: "o mais demorado e " + slowest.Kind + " com " + itoaFixed(slowest.TTK) + "s (teto " + itoaFixed(maxEnemyTTK) + "s)",
		Pass:   slowest.TTK <= maxEnemyTTK,
	})

	detail := "todos os inimigos que atiram conseguem disparar antes de morrer"
	if len(silent) > 0 {
		detail = "morrem sem disparar: " + joinNames(silent)
	}
	out = append(out, Check{
		Name:   "Todo inimigo armado chega a ameacar",
		Detail: detail,
		Pass:   len(silent) == 0,
	})
	return out
}

func joinNames(names []string) string {
	out := ""
	for i, n := range names {
		if i > 0 {
			out += ", "
		}
		out += n
	}
	return out
}

// checkCorruptionCurve valida a aposta central do jogo: um jogador competente
// termina o reino em pé, um jogador desatento vê o mundo desabar, e existe uma
// faixa intermediária larga onde deixar passar é uma escolha, não um acidente.
func checkCorruptionCurve(r BalanceReport) []Check {
	var good, bad *CorruptionMetrics
	for i := range r.Corruption {
		c := &r.Corruption[i]
		if c.LeakRate <= 0.05 && (good == nil || c.LeakRate > good.LeakRate) {
			good = c
		}
		if c.LeakRate >= 0.40 && (bad == nil || c.LeakRate < bad.LeakRate) {
			bad = c
		}
	}
	out := []Check{}
	if good != nil {
		out = append(out, Check{
			Name:   "Jogador competente mantem o reino em pe",
			Detail: "com " + pct(good.LeakRate) + " de fugas termina em " + itoaFixed(good.Final) + "% (" + good.Tier + ")",
			Pass:   good.Final <= maxGoodCorruption,
		})
	}
	if bad != nil {
		out = append(out, Check{
			Name:   "Ignorar as ameacas derruba o reino",
			Detail: "com " + pct(bad.LeakRate) + " de fugas termina em " + itoaFixed(bad.Final) + "% (" + bad.Tier + ")",
			Pass:   bad.Final >= maxCorruption,
		})
	}
	reachable := 0
	for _, c := range r.Corruption {
		if c.AscendedBoss {
			reachable++
		}
	}
	out = append(out, Check{
		Name:   "O final alternativo e alcancavel, mas nao acidental",
		Detail: itoa(reachable) + " de " + itoa(len(r.Corruption)) + " perfis chegam a Vharak Ascendido",
		Pass:   reachable >= 1 && reachable <= len(r.Corruption)/2,
	})
	return out
}

// checkIdentities verifica que cada cenário tem um dono diferente — é o que
// separa "três armas" de "uma arma e duas alternativas piores".
func checkIdentities(r BalanceReport) Check {
	owners := scenarioOwners(r)
	distinct := map[string]bool{}
	for _, name := range owners {
		distinct[name] = true
	}
	return Check{
		Name: "Cada cenario tem uma arma dona",
		Detail: "foco: " + owners["foco"] + " | formacao: " + owners["formacao"] +
			" | coluna: " + owners["coluna"],
		Pass: len(distinct) == 3,
	}
}

// checkNoDominance rejeita a situação em que uma arma é a melhor escolha em
// qualquer situação — o defeito que a Lança de Luz tinha.
func checkNoDominance(r BalanceReport) Check {
	for _, w := range AtLevel(r.Weapons, maxWeaponLevel) {
		wins := 0
		for _, scenario := range scenarioOwners(r) {
			if scenario == w.Weapon {
				wins++
			}
		}
		if wins >= 3 {
			return Check{
				Name:   "Nenhuma arma e a resposta certa para tudo",
				Detail: w.Weapon + " vence nos tres cenarios",
				Pass:   false,
			}
		}
	}
	return Check{
		Name:   "Nenhuma arma e a resposta certa para tudo",
		Detail: "nenhuma arma vence em todos os cenarios",
		Pass:   true,
	}
}

// scenarioOwners devolve, para cada cenário, a arma de maior dano no nível máximo.
func scenarioOwners(r BalanceReport) map[string]string {
	owners := map[string]string{"foco": "-", "formacao": "-", "coluna": "-"}
	top := map[string]float64{}
	for _, w := range AtLevel(r.Weapons, maxWeaponLevel) {
		for name, v := range map[string]float64{
			"foco": w.FocusDPS, "formacao": w.SwarmDPS, "coluna": w.ColumnDPS,
		} {
			if v > top[name] {
				top[name], owners[name] = v, w.Weapon
			}
		}
	}
	return owners
}

func fmtRatio(label string, got, max float64) string {
	return label + " " + itoaFixed(got) + "x (teto " + itoaFixed(max) + "x)"
}

// itoaFixed formata com uma casa decimal, sem depender de fmt.
func itoaFixed(v float64) string {
	n := int(v*10 + 0.5)
	return itoa(n/10) + "." + itoa(n%10)
}

func pct(v float64) string { return itoa(int(v*100+0.5)) + "%" }

// PeakSpread mede o desequilíbrio entre a arma que mais rende no seu melhor
// cenário e a que menos rende no dela. É a comparação justa entre armas
// assimétricas: 1.0 significa que todas brilham igualmente onde deveriam.
func PeakSpread(ms []WeaponMetrics, level int) float64 {
	best, worst := 0.0, math.MaxFloat64
	for _, m := range ms {
		if m.Level != level {
			continue
		}
		p := m.Peak()
		if p > best {
			best = p
		}
		if p < worst {
			worst = p
		}
	}
	if worst <= 0 || worst == math.MaxFloat64 {
		return 0
	}
	return best / worst
}

// AtLevel filtra as medições de um nível.
func AtLevel(ms []WeaponMetrics, level int) []WeaponMetrics {
	out := make([]WeaponMetrics, 0, 3)
	for _, m := range ms {
		if m.Level == level {
			out = append(out, m)
		}
	}
	return out
}
