package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Resolução interna pensada para pixel art. A janela é escalada no main.
const (
	ScreenWidth  = 240
	ScreenHeight = 320
)

type state int

const (
	stateMenu state = iota
	stateControls
	statePlaying
	statePaused
	stateBoss
	stateVictory
	stateGameOver
)

type Game struct {
	state state

	player       *Player
	bullets      []*Bullet
	enemies      []*Enemy
	enemyBullets []*EnemyBullet
	powerups     []*Powerup
	stars        []*star

	clouds     []*bgItem
	hills      []*bgItem
	structures []*bgItem
	dust       []*star
	particles  []*particle
	popups     []*scorePopup

	shakeMag     float64
	shakeSetting shakeLevel
	damageFlash  int
	hitStop      int
	scene        *ebiten.Image

	level       *Level
	timeScale   int
	resumeState state

	lives           int
	bombCharges     int
	bombEffectTimer int

	score      int
	highScore  int
	combo      int
	comboTimer int

	sectionDamaged bool
	lastSection    int
	formations     map[int]*formationTracker

	boss *Boss

	menuIndex int
	fade      int

	enemiesDefeated  int
	maxMult          int
	elapsedTicks     int
	victoryLifeBonus int
	victoryBombBonus int
	victoryTime      int

	audio *AudioManager
}

const menuItemCount = 4

var gameOverOptions = []string{"Tentar novamente", "Voltar ao menu"}
var victoryOptions = []string{"Jogar novamente", "Voltar ao menu"}

type formationTracker struct {
	total   int
	killed  int
	escaped int
}

func New() *Game {
	g := &Game{}
	g.audio = getAudio()
	g.reset()
	g.initStars()
	g.enterMenu()
	return g
}

// reset cria uma sessão completamente nova, sem qualquer estado anterior.
func (g *Game) reset() {
	resetRNG()
	g.state = statePlaying
	g.player = newPlayer()
	g.bullets = g.bullets[:0]
	g.enemies = g.enemies[:0]
	g.enemyBullets = g.enemyBullets[:0]
	g.powerups = g.powerups[:0]
	g.level = newLevel()
	g.timeScale = dev.timeScale

	g.lives = startingLives
	g.bombCharges = bombStartCharges
	g.bombEffectTimer = 0
	g.score = 0
	g.combo = 0
	g.comboTimer = 0
	g.sectionDamaged = false
	g.lastSection = g.level.section
	g.formations = map[int]*formationTracker{}
	g.boss = nil

	g.enemiesDefeated = 0
	g.maxMult = 1
	g.elapsedTicks = 0
	g.victoryLifeBonus = 0
	g.victoryBombBonus = 0
	g.victoryTime = 0

	g.particles = g.particles[:0]
	g.popups = g.popups[:0]
	g.shakeMag = 0
	g.damageFlash = 0
	g.hitStop = 0

	if dev.startBoss {
		g.startBoss()
	}
}

func (g *Game) enterMenu() {
	g.state = stateMenu
	g.menuIndex = 0
	g.fade = fadeDuration
}

func (g *Game) startNewGame() {
	g.reset()
	g.fade = fadeDuration
}

func (g *Game) startBoss() {
	g.enemies = g.enemies[:0]
	g.enemyBullets = g.enemyBullets[:0]
	g.boss = newBoss()
	g.state = stateBoss
}

func (g *Game) Update() error {
	if g.fade > 0 {
		g.fade--
	}
	g.updateStars()

	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.audio.toggleMute()
	}
	g.updateMusicForState()
	g.audio.update()
	g.handleDevToggles()

	switch g.state {
	case stateMenu:
		return g.updateMenu()
	case stateControls:
		g.updateControls()
		return nil
	case statePaused:
		g.handlePauseToggle()
		return nil
	case stateGameOver:
		g.updateEndScreen(gameOverOptions)
		return nil
	case stateVictory:
		g.updateEndScreen(victoryOptions)
		return nil
	}

	g.updatePlay()
	return nil
}

// updateMusicForState escolhe a trilha conforme a tela atual. A pausa mantém a
// trilha do estado que estava em andamento.
func (g *Game) updateMusicForState() {
	st := g.state
	if st == statePaused {
		st = g.resumeState
	}
	switch st {
	case statePlaying:
		g.audio.playMusic(musicPhase)
	case stateBoss:
		g.audio.playMusic(musicBoss)
	case stateGameOver, stateVictory:
		g.audio.playMusic(musicNone)
	default:
		g.audio.playMusic(musicMenu)
	}
}

// updateMenu navega o menu inicial. Sair encerra o jogo sem erro.
func (g *Game) updateMenu() error {
	g.moveMenuSelection(menuItemCount)
	if !inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return nil
	}
	g.audio.playSFX(sfxMenu)
	switch g.menuIndex {
	case 0:
		g.startNewGame()
	case 1:
		g.state = stateControls
		g.fade = fadeDuration
	case 2:
		g.cycleShake()
	case 3:
		return ebiten.Termination
	}
	return nil
}

func (g *Game) updateControls() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.enterMenu()
	}
}

func (g *Game) updateEndScreen(options []string) {
	g.moveMenuSelection(len(options))
	if !inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return
	}
	g.audio.playSFX(sfxMenu)
	if g.menuIndex == 0 {
		g.startNewGame()
		return
	}
	g.enterMenu()
}

func (g *Game) moveMenuSelection(count int) {
	moved := false
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.menuIndex = (g.menuIndex - 1 + count) % count
		moved = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.menuIndex = (g.menuIndex + 1) % count
		moved = true
	}
	if moved {
		g.audio.playSFX(sfxMenu)
	}
}

func (g *Game) updatePlay() {
	g.handlePauseToggle()
	if g.state == statePaused {
		return
	}

	// Impact freeze: por poucos frames a lógica congela para dar peso aos golpes
	// fortes, mas os efeitos visuais seguem animando para não parecer travamento.
	if g.hitStop > 0 {
		g.hitStop--
		g.updateParticles()
		g.updatePopups()
		g.updateShake()
		if g.damageFlash > 0 {
			g.damageFlash--
		}
		return
	}

	g.elapsedTicks++
	g.player.update()
	if shots := g.player.tryShoot(); len(shots) > 0 {
		g.bullets = append(g.bullets, shots...)
		g.audio.playSFX(sfxShoot)
	}
	g.handleDebugSpawns()
	g.handleBomb()

	if g.bombEffectTimer > 0 {
		g.bombEffectTimer--
	}
	if g.comboTimer > 0 {
		g.comboTimer--
		if g.comboTimer == 0 {
			g.combo = 0
		}
	}

	if g.state == statePlaying {
		g.handleTimeScaleToggle()
		for i := 0; i < g.timeScale; i++ {
			g.spawnFromLevel()
		}
	}
	if g.state == stateBoss {
		g.updateBoss()
	}

	g.updateBullets()
	g.updateEnemies()
	g.updateEnemyBullets()
	g.updatePowerups()
	g.handleCollisions()
	g.collectPowerups()
	g.removeDead()

	g.updateParallax()
	g.updateParticles()
	g.updatePopups()
	g.updateShake()
	if g.damageFlash > 0 {
		g.damageFlash--
	}

	if g.player.health <= 0 {
		g.loseLife()
	}
	if g.score > g.highScore {
		g.highScore = g.score
	}
	if g.state == statePlaying {
		g.checkSectionBonus()
		if g.level.readyForBoss(len(g.enemies)) {
			if !g.sectionDamaged {
				g.score += waveNoDamageBonus
			}
			g.startBoss()
		}
	}
	if g.state == stateBoss {
		g.checkBossDefeat()
	}
}

func (g *Game) updateBoss() {
	ctx := &bossContext{
		playerX: g.player.centerX(),
		playerY: g.player.centerY(),
		bullets: &g.enemyBullets,
		enemies: &g.enemies,
	}
	g.boss.update(ctx)

	if g.boss.justEntered {
		g.enemies = g.enemies[:0]
		g.enemyBullets = g.enemyBullets[:0]
		g.boss.justEntered = false
	}
	if g.boss.phaseChanged {
		g.enemyBullets = g.enemyBullets[:0]
		g.addShake(shakeBossMagnitude)
		g.hitStop = hitStopBig
		g.boss.phaseChanged = false
	}
	// Ao iniciar a morte, a arena é limpa e o jogador fica invulnerável durante
	// toda a sequência: nenhum projétil ou inimigo residual pode causar uma
	// derrota injusta no exato instante da vitória.
	if g.boss.justDied {
		g.enemies = g.enemies[:0]
		g.enemyBullets = g.enemyBullets[:0]
		if g.player.invincible < bossDeathDuration {
			g.player.invincible = bossDeathDuration
		}
		g.hitStop = hitStopBig
		g.boss.justDied = false
	}
	if g.boss.phase == bossDying {
		g.bossDeathEffects()
	}
}

// bossDeathEffects espalha explosões pelo corpo do chefe durante a morte.
func (g *Game) bossDeathEffects() {
	if g.boss.tick%6 != 0 {
		return
	}
	x := g.boss.x + fxRand.Float64()*g.boss.w
	y := g.boss.y + fxRand.Float64()*g.boss.h
	g.spawnBurst(x, y, bossDeathParticles/2, 2.6, 26, 3, true, color.RGBA{0xff, 0xc0, 0x50, 0xff})
	g.addShake(shakeBossMagnitude)
}

func (g *Game) checkBossDefeat() {
	if g.boss.phase != bossDead {
		return
	}
	g.victoryTime = g.elapsedTicks / 60
	g.victoryLifeBonus = g.lives * lifeBonus
	g.victoryBombBonus = g.bombCharges * bombBonusPoints
	g.score += bossScore + g.victoryLifeBonus + g.victoryBombBonus
	if g.score > g.highScore {
		g.highScore = g.score
	}
	g.state = stateVictory
	g.menuIndex = 0
	g.fade = fadeDuration
	g.audio.playSFX(sfxVictory)
}

func (g *Game) spawnFromLevel() {
	for _, e := range g.level.update() {
		if e.formationID > 0 {
			g.registerFormationMember(e.formationID)
		}
		g.enemies = append(g.enemies, e)
	}
}

func (g *Game) registerFormationMember(id int) {
	t := g.formations[id]
	if t == nil {
		t = &formationTracker{}
		g.formations[id] = t
	}
	t.total++
}

func (g *Game) loseLife() {
	g.lives--
	if g.lives <= 0 {
		g.state = stateGameOver
		g.menuIndex = 0
		g.fade = fadeDuration
		g.audio.playSFX(sfxGameOver)
		return
	}
	g.respawn()
}

func (g *Game) respawn() {
	g.enemyBullets = g.enemyBullets[:0]
	g.player.respawn()
	g.combo = 0
	g.comboTimer = 0
}

func (g *Game) checkSectionBonus() {
	if g.level.section == g.lastSection {
		return
	}
	if !g.sectionDamaged {
		g.score += waveNoDamageBonus
	}
	g.sectionDamaged = false
	g.lastSection = g.level.section
}

func (g *Game) handleBomb() {
	if g.state != statePlaying && g.state != stateBoss {
		return
	}
	if g.bombCharges <= 0 {
		return
	}
	if !inpututil.IsKeyJustPressed(ebiten.KeyX) && !inpututil.IsKeyJustPressed(ebiten.KeyControl) {
		return
	}
	g.useBomb()
}

func (g *Game) useBomb() {
	g.bombCharges--
	g.enemyBullets = g.enemyBullets[:0]
	for _, e := range g.enemies {
		if e.dead {
			continue
		}
		e.takeDamage(bombDamage)
		if e.dead {
			g.score += e.score
			g.spawnExplosion(e.centerX(), e.centerY(), e.color)
		}
	}
	if g.boss != nil {
		g.boss.takeDamage(bombDamage / 20) // dano relevante, mas não instantâneo
	}
	g.player.invincible = bombInvincibility
	g.bombEffectTimer = bombEffectDuration
	g.addShake(shakeBombMagnitude)
	g.audio.playSFX(sfxBomb)
}

func (g *Game) handlePauseToggle() {
	if !inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return
	}
	if g.state == statePaused {
		g.state = g.resumeState
		return
	}
	if g.state == statePlaying || g.state == stateBoss {
		g.resumeState = g.state
		g.state = statePaused
	}
}

func (g *Game) handleTimeScaleToggle() {
	if !dev.enabled || !inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		return
	}
	if g.timeScale == dev.fastTimeScale {
		g.timeScale = dev.timeScale
		return
	}
	g.timeScale = dev.fastTimeScale
}

func (g *Game) updateBullets() {
	for _, b := range g.bullets {
		b.update()
		if b.trail {
			g.spawnTrail(b.x, b.y, b.color)
		}
		if b.offScreen() {
			b.dead = true
		}
	}
}

func (g *Game) spawnCrowFormation() {
	count := 3 + rng.Intn(2)
	const gap = crowSize + 6
	startX := rng.Float64() * (ScreenWidth - float64(count)*gap)
	for i := 0; i < count; i++ {
		g.enemies = append(g.enemies, newCrow(startX+float64(i)*gap))
	}
}

func (g *Game) spawnGargoyle() {
	fromLeft := rng.Intn(2) == 0
	targetX := 40 + rng.Float64()*(ScreenWidth-80-gargoyleSize)
	y := 30 + rng.Float64()*60
	g.enemies = append(g.enemies, newGargoyle(fromLeft, targetX, y))
}

func randX(size float64) float64 {
	return rng.Float64() * (ScreenWidth - size)
}

// handleDebugSpawns concentra as teclas do modo de desenvolvimento e só age
// quando ele está ativo, mantendo o build padrão livre dessas ferramentas.
func (g *Game) handleDebugSpawns() {
	if !dev.enabled {
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDigit1) {
		g.spawnCrowFormation()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDigit2) {
		g.enemies = append(g.enemies, newHarpy(randX(harpySize)))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDigit3) {
		g.spawnGargoyle()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDigit4) {
		g.enemies = append(g.enemies, newWyvern(randX(wyvernSize)))
	}
	g.debugSpawnPowerup(ebiten.KeyZ, powerLight)
	g.debugSpawnPowerup(ebiten.KeyF, powerFire)
	g.debugSpawnPowerup(ebiten.KeyC, powerIce)
	g.debugSpawnPowerup(ebiten.KeyV, powerHeal)
	g.debugSpawnPowerup(ebiten.KeyB, powerShield)

	if inpututil.IsKeyJustPressed(ebiten.KeyK) && g.boss != nil {
		g.boss.takeDamage(devBossDamage)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.clearEntities()
	}
}

// clearEntities remove inimigos e projéteis sem contá-los como abatidos.
func (g *Game) clearEntities() {
	g.enemies = g.enemies[:0]
	g.enemyBullets = g.enemyBullets[:0]
	g.bullets = g.bullets[:0]
}

// handleDevToggles alterna as sobreposições de diagnóstico do modo dev.
func (g *Game) handleDevToggles() {
	if !dev.enabled {
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		dev.showHUD = !dev.showHUD
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		dev.showHitboxes = !dev.showHitboxes
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		dev.invincible = !dev.invincible
	}
}

func (g *Game) debugSpawnPowerup(key ebiten.Key, kind powerupType) {
	if inpututil.IsKeyJustPressed(key) {
		g.powerups = append(g.powerups, newPowerup(kind, g.player.centerX(), 24))
	}
}

func (g *Game) updateEnemies() {
	ctx := &enemyContext{
		playerX: g.player.centerX(),
		playerY: g.player.centerY(),
		bullets: &g.enemyBullets,
	}
	for _, e := range g.enemies {
		e.update(ctx)
		if e.offScreen() {
			e.escaped = true
			e.dead = true
		}
	}
}

func (g *Game) updateEnemyBullets() {
	for _, b := range g.enemyBullets {
		b.update()
		if b.offScreen() {
			b.dead = true
		}
	}
}

func (g *Game) handleCollisions() {
	g.bulletEnemyCollisions()
	if g.boss != nil {
		g.bulletBossCollisions()
	}

	if g.player.canBeHit() {
		g.playerEnemyCollisions()
	}
	if g.player.canBeHit() {
		g.playerBulletCollisions()
	}
	if g.boss != nil && g.player.canBeHit() {
		g.playerBossCollision()
	}
}

func (g *Game) bulletBossCollisions() {
	boss := g.boss
	for _, b := range g.bullets {
		if b.dead {
			continue
		}
		if !collides(b.x, b.y, b.w, b.h, boss.x, boss.y, boss.w, boss.h) {
			continue
		}
		damage := b.damage
		if boss.hitCrystal(b.x, b.y, b.w, b.h) {
			damage *= crystalBonus
		}
		boss.takeDamage(damage)
		b.dead = true
	}
}

func (g *Game) playerBossCollision() {
	boss := g.boss
	if boss.phase == bossDying || boss.phase == bossDead {
		return
	}
	hx, hy, hw, hh := g.player.hitbox()
	if collides(hx, hy, hw, hh, boss.x, boss.y, boss.w, boss.h) {
		g.hitPlayer(bossContactDamage)
	}
}

func (g *Game) bulletEnemyCollisions() {
	for _, b := range g.bullets {
		if b.dead {
			continue
		}
		for _, e := range g.enemies {
			if e.dead || b.alreadyHit(e) {
				continue
			}
			if !collides(b.x, b.y, b.w, b.h, e.x, e.y, e.w, e.h) {
				continue
			}
			e.takeDamage(b.damage)
			b.hitEnemies = append(b.hitEnemies, e)
			if e.dead {
				points := g.registerKill(e.score)
				g.spawnDrop(e)
				g.spawnExplosion(e.centerX(), e.centerY(), e.color)
				g.spawnScorePopup(e.centerX(), e.y, points)
				if e.kind == kindWyvern {
					g.hitStop = hitStopBig // peso extra ao abater um inimigo grande
				}
			}
			if b.pierce > 0 {
				b.pierce--
				continue
			}
			b.dead = true
			break
		}
	}
}

func (g *Game) playerEnemyCollisions() {
	hx, hy, hw, hh := g.player.hitbox()
	for _, e := range g.enemies {
		if e.dead {
			continue
		}
		if collides(hx, hy, hw, hh, e.x, e.y, e.w, e.h) {
			e.dead = true
			g.spawnExplosion(e.centerX(), e.centerY(), e.color)
			g.hitPlayer(e.damage)
			return
		}
	}
}

func (g *Game) playerBulletCollisions() {
	hx, hy, hw, hh := g.player.hitbox()
	for _, b := range g.enemyBullets {
		if b.dead {
			continue
		}
		if collides(hx, hy, hw, hh, b.x, b.y, enemyBulletSize, enemyBulletSize) {
			b.dead = true
			g.hitPlayer(enemyBulletDamage)
			return
		}
	}
}

// registerKill contabiliza uma eliminação por disparo, aplica o multiplicador de
// combo e devolve os pontos efetivamente concedidos (para exibição).
func (g *Game) registerKill(base int) int {
	g.combo++
	g.comboTimer = comboWindow
	mult := g.multiplier()
	points := base * mult
	g.score += points
	if mult > g.maxMult {
		g.maxMult = mult
	}
	return points
}

func (g *Game) multiplier() int {
	m := 1 + g.combo/comboStep
	if m > maxMultiplier {
		return maxMultiplier
	}
	return m
}

// hitPlayer centraliza a aplicação de dano ao jogador nas três origens
// (inimigo, projétil, chefe). Distingue o dano real da absorção pelo escudo,
// garantindo o feedback correto em cada caso.
func (g *Game) hitPlayer(damage int) {
	hadShield := g.player.hasShield()
	if g.player.hit(damage) {
		g.onPlayerDamaged()
		return
	}
	if hadShield {
		g.onShieldBroken()
	}
}

func (g *Game) onPlayerDamaged() {
	g.combo = 0
	g.comboTimer = 0
	g.sectionDamaged = true
	g.triggerDamageFlash()
	g.addShake(shakeHitMagnitude)
	g.audio.playSFX(sfxPlayerHit)
	if g.player != nil {
		g.spawnBurst(g.player.centerX(), g.player.centerY(), 6, 2.0, 18, 2, false, color.RGBA{0xff, 0x60, 0x60, 0xff})
	}
}

// onShieldBroken dá feedback claro de que o escudo absorveu um golpe (som,
// anel de partículas e leve vibração), sem contar como dano recebido.
func (g *Game) onShieldBroken() {
	g.addShake(shakeHitMagnitude)
	g.audio.playSFX(sfxShield)
	if g.player != nil {
		g.spawnCollectRing(g.player.centerX(), g.player.centerY(), shieldRing)
	}
}

func (g *Game) spawnDrop(e *Enemy) {
	if e.hasDrop {
		g.powerups = append(g.powerups, newPowerup(e.drop, e.centerX(), e.centerY()))
		return
	}
	if rng.Float64() < dropChance {
		g.powerups = append(g.powerups, newPowerup(randomRune(), e.centerX(), e.centerY()))
	}
}

func randomRune() powerupType {
	switch rng.Intn(5) {
	case 0:
		return powerLight
	case 1:
		return powerFire
	case 2:
		return powerIce
	case 3:
		return powerHeal
	default:
		return powerShield
	}
}

func (g *Game) updatePowerups() {
	for _, p := range g.powerups {
		p.update()
		if p.offScreen() {
			p.dead = true
		}
	}
}

func (g *Game) collectPowerups() {
	for _, p := range g.powerups {
		if p.dead {
			continue
		}
		if collides(g.player.x, g.player.y, playerSize, playerSize, p.x, p.y, powerupSize, powerupSize) {
			g.player.applyPowerup(p.kind)
			p.dead = true
			g.spawnCollectRing(p.x+powerupSize/2, p.y+powerupSize/2, p.color())
			g.audio.playSFX(sfxPickup)
		}
	}
}

func (g *Game) removeDead() {
	for _, e := range g.enemies {
		if !e.dead {
			continue
		}
		if !e.escaped {
			g.enemiesDefeated++
		}
		if e.formationID > 0 {
			g.accountFormation(e)
		}
	}

	g.bullets = filterAlive(g.bullets, func(b *Bullet) bool { return !b.dead })
	g.enemies = filterAlive(g.enemies, func(e *Enemy) bool { return !e.dead })
	g.enemyBullets = filterAlive(g.enemyBullets, func(b *EnemyBullet) bool { return !b.dead })
	g.powerups = filterAlive(g.powerups, func(p *Powerup) bool { return !p.dead })
}

// accountFormation concede bônus quando toda uma formação é destruída sem
// que nenhum de seus membros escape da tela.
func (g *Game) accountFormation(e *Enemy) {
	t := g.formations[e.formationID]
	if t == nil {
		return
	}
	if e.escaped {
		t.escaped++
	} else {
		t.killed++
	}
	if t.killed+t.escaped < t.total {
		return
	}
	if t.escaped == 0 {
		g.score += formationBonus
	}
	delete(g.formations, e.formationID)
}

// filterAlive mantém apenas os itens aprovados por keep, reaproveitando o array
// e liberando as posições restantes para que o GC recolha itens removidos.
func filterAlive[T any](items []*T, keep func(*T) bool) []*T {
	alive := items[:0]
	for _, it := range items {
		if keep(it) {
			alive = append(alive, it)
		}
	}
	for i := len(alive); i < len(items); i++ {
		items[i] = nil
	}
	return alive
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case stateMenu:
		g.drawMenu(screen)
	case stateControls:
		g.drawControls(screen)
	case stateGameOver:
		g.drawGameOver(screen)
	case stateVictory:
		g.drawVictory(screen)
	default:
		g.drawWorld(screen)
	}
	g.drawFade(screen)
}

// drawWorld renderiza a cena em uma imagem própria, aplica a vibração ao
// compô-la na tela e desenha o HUD por cima (sem tremer), preservando a leitura.
func (g *Game) drawWorld(screen *ebiten.Image) {
	if g.scene == nil {
		g.scene = ebiten.NewImage(ScreenWidth, ScreenHeight)
	}
	g.scene.Clear()
	g.drawScene(g.scene)

	screen.Fill(g.level.background())
	ox, oy := g.shakeOffset()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(ox, oy)
	screen.DrawImage(g.scene, op)

	if g.damageFlash > 0 {
		alpha := uint8(damageFlashAlpha * g.damageFlash / damageFlashDuration)
		vector.DrawFilledRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0xff, 0x20, 0x20, alpha}, false)
	}

	g.drawOverlayHUD(screen)
}

// drawScene desenha o cenário e as entidades. Projéteis são desenhados por
// último para permanecerem visíveis acima das partículas e efeitos.
func (g *Game) drawScene(dst *ebiten.Image) {
	theme := g.level.theme()
	dst.Fill(theme.sky)
	g.drawParallax(dst, theme)
	g.drawParticles(dst)

	for _, p := range g.powerups {
		p.draw(dst)
	}
	for _, e := range g.enemies {
		e.draw(dst)
	}
	if g.boss != nil {
		g.boss.draw(dst)
	}
	g.player.draw(dst)

	for _, b := range g.enemyBullets {
		b.draw(dst)
	}
	for _, b := range g.bullets {
		b.draw(dst)
	}

	if g.bombEffectTimer > 0 {
		g.drawBombEffect(dst)
	}

	g.drawPopups(dst)

	if dev.showHitboxes {
		g.drawHitboxes(dst)
	}
}

// drawHitboxes desenha os contornos de colisão para depuração visual.
func (g *Game) drawHitboxes(dst *ebiten.Image) {
	green := color.RGBA{0x40, 0xff, 0x40, 0xff}
	yellow := color.RGBA{0xff, 0xff, 0x40, 0xff}

	hx, hy, hw, hh := g.player.hitbox()
	vector.StrokeRect(dst, float32(hx), float32(hy), float32(hw), float32(hh), 1, green, false)
	for _, e := range g.enemies {
		vector.StrokeRect(dst, float32(e.x), float32(e.y), float32(e.w), float32(e.h), 1, green, false)
	}
	for _, b := range g.bullets {
		vector.StrokeRect(dst, float32(b.x), float32(b.y), float32(b.w), float32(b.h), 1, green, false)
	}
	for _, b := range g.enemyBullets {
		vector.StrokeRect(dst, float32(b.x), float32(b.y), enemyBulletSize, enemyBulletSize, 1, yellow, false)
	}
	if g.boss != nil {
		vector.StrokeRect(dst, float32(g.boss.x), float32(g.boss.y), float32(g.boss.w), float32(g.boss.h), 1, green, false)
	}
}

func (g *Game) drawOverlayHUD(screen *ebiten.Image) {
	g.drawHUD(screen)
	if g.boss != nil {
		g.drawBossHUD(screen)
	}
	if dev.enabled && dev.showHUD {
		g.drawDevHUD(screen)
	}

	if g.level.announceTimer > 0 && g.level.announce != "" {
		ebitenutil.DebugPrintAt(screen, g.level.announce, ScreenWidth/2-len(g.level.announce)*3, 70)
	}

	if g.state == statePaused {
		ebitenutil.DebugPrintAt(screen, "PAUSADO", ScreenWidth/2-24, ScreenHeight/2-4)
	}
}

// drawDevHUD mostra FPS, contagens, tempo da fase e a semente em uso.
func (g *Game) drawDevHUD(screen *ebiten.Image) {
	lines := []string{
		"FPS " + itoa(int(ebiten.ActualFPS()+0.5)),
		"INIM " + itoa(len(g.enemies)),
		"PROJ " + itoa(len(g.bullets)+len(g.enemyBullets)),
		"PART " + itoa(len(g.particles)),
		"TEMPO " + itoa(g.elapsedTicks/60) + "s",
		"TS x" + itoa(g.timeScale),
		"SEED " + itoa(int(CurrentSeed())),
	}
	if dev.invincible {
		lines = append(lines, "INVENCIVEL")
	}
	for i, l := range lines {
		ebitenutil.DebugPrintAt(screen, l, 4, 44+i*12)
	}
}

func (g *Game) drawStarsBackground(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x0a, 0x0a, 0x1e, 0xff})
	g.drawStars(screen)
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	g.drawStarsBackground(screen)
	drawCentered(screen, "ASAS DE VALDORIA", 60)
	drawCentered(screen, "Um shoot 'em up medieval fantastico", 78)
	options := []string{
		"Iniciar jogo",
		"Controles",
		"Vibracao: " + shakeLevelName(g.shakeSetting),
		"Sair",
	}
	g.drawOptions(screen, options, 150)
	drawCentered(screen, "Setas/WASD move  Enter confirma", ScreenHeight-16)
	ebitenutil.DebugPrintAt(screen, "v"+Version, 4, ScreenHeight-14)
}

func (g *Game) drawControls(screen *ebiten.Image) {
	g.drawStarsBackground(screen)
	drawCentered(screen, "CONTROLES", 30)
	lines := []string{
		"Setas ou WASD: mover",
		"Shift: precisao",
		"Espaco: disparar",
		"X ou Ctrl: invocacao ancestral",
		"Escape: pausar",
		"Enter: confirmar",
	}
	for i, l := range lines {
		ebitenutil.DebugPrintAt(screen, l, 24, 60+i*18)
	}
	drawCentered(screen, "Enter/Escape para voltar", ScreenHeight-20)
}

// reachedLabel descreve até onde o jogador avançou ao perder.
func (g *Game) reachedLabel() string {
	if g.boss != nil {
		return "Chefe: Vharak"
	}
	return g.level.sectionName()
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	g.drawStarsBackground(screen)
	drawCentered(screen, "GAME OVER", 40)
	stats := []string{
		"Pontuacao: " + itoa(g.score),
		"Trecho: " + g.reachedLabel(),
		"Inimigos derrotados: " + itoa(g.enemiesDefeated),
		"Maior multiplicador: x" + itoa(g.maxMult),
	}
	for i, s := range stats {
		ebitenutil.DebugPrintAt(screen, s, 30, 80+i*16)
	}
	g.drawOptions(screen, gameOverOptions, 170)
}

func (g *Game) drawVictory(screen *ebiten.Image) {
	g.drawStarsBackground(screen)
	drawCentered(screen, "FASE CONCLUIDA!", 34)
	stats := []string{
		"Pontuacao final: " + itoa(g.score),
		"Bonus por vidas: " + itoa(g.victoryLifeBonus),
		"Bonus por bombas: " + itoa(g.victoryBombBonus),
		"Tempo: " + itoa(g.victoryTime) + "s",
		"Inimigos derrotados: " + itoa(g.enemiesDefeated),
	}
	for i, s := range stats {
		ebitenutil.DebugPrintAt(screen, s, 30, 70+i*16)
	}
	g.drawOptions(screen, victoryOptions, 180)
}

func (g *Game) drawOptions(screen *ebiten.Image, options []string, y int) {
	for i, opt := range options {
		text := "  " + opt
		if i == g.menuIndex {
			text = "> " + opt
		}
		ebitenutil.DebugPrintAt(screen, text, ScreenWidth/2-len(text)*3, y+i*16)
	}
}

func (g *Game) drawFade(screen *ebiten.Image) {
	if g.fade <= 0 {
		return
	}
	alpha := uint8(255 * g.fade / fadeDuration)
	vector.DrawFilledRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, alpha}, false)
}

func drawCentered(screen *ebiten.Image, text string, y int) {
	ebitenutil.DebugPrintAt(screen, text, ScreenWidth/2-len(text)*3, y)
}

func (g *Game) drawBossHUD(screen *ebiten.Image) {
	b := g.boss
	if b.phase == bossDead {
		return
	}
	if b.phase == bossEntering {
		// Aviso de entrada pulsante para leitura clara do início do combate.
		if (b.tick/16)%2 == 0 {
			drawCentered(screen, "! ALERTA !", ScreenHeight/2-14)
		}
		ebitenutil.DebugPrintAt(screen, "VHARAK, O DRAGAO CORROMPIDO", ScreenWidth/2-81, ScreenHeight/2)
		return
	}

	const barY = 12
	const margin = 10
	const barH = 6
	width := float32(ScreenWidth - 2*margin)

	// Moldura + trilho escuro.
	vector.DrawFilledRect(screen, margin-1, barY-1, width+2, barH+2, color.RGBA{0x10, 0x08, 0x08, 0xff}, false)
	vector.DrawFilledRect(screen, margin, barY, width, barH, color.RGBA{0x30, 0x10, 0x10, 0xff}, false)

	// Preenchimento colorido pela fase atual (vermelho → laranja → vermelho vivo).
	fill := color.RGBA{0xc0, 0x30, 0x30, 0xff}
	switch b.phase {
	case bossPhase2:
		fill = color.RGBA{0xd0, 0x70, 0x20, 0xff}
	case bossPhase3:
		fill = color.RGBA{0xff, 0x40, 0x40, 0xff}
	}
	vector.DrawFilledRect(screen, margin, barY, width*float32(b.healthRatio()), barH, fill, false)

	// Marcas nos limiares de troca de fase (65% e 30%), para o jogador antecipar.
	tick := color.RGBA{0x20, 0x10, 0x10, 0xff}
	for _, r := range []float64{0.65, 0.30} {
		x := float32(margin) + width*float32(r)
		vector.DrawFilledRect(screen, x, barY, 1, barH, tick, false)
	}
	ebitenutil.DebugPrintAt(screen, "VHARAK", margin, 0)

	if b.action == actionWarning {
		ebitenutil.DebugPrintAt(screen, "!", int(b.centerX())-2, int(b.y+b.h)+4)
	}
}

// drawBombEffect desenha o dragão ancestral atravessando a tela de baixo p/ cima.
func (g *Game) drawBombEffect(screen *ebiten.Image) {
	progress := 1 - float64(g.bombEffectTimer)/float64(bombEffectDuration)
	y := float32(ScreenHeight - progress*(ScreenHeight+60))
	body := color.RGBA{0xff, 0xd0, 0x60, 0xc0}
	vector.DrawFilledRect(screen, 0, y, ScreenWidth, 24, body, false)
	wing := color.RGBA{0xff, 0xa0, 0x30, 0xa0}
	vector.DrawFilledRect(screen, 0, y-10, ScreenWidth, 6, wing, false)
	vector.DrawFilledRect(screen, 0, y+28, ScreenWidth, 6, wing, false)
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "SCORE "+itoa(g.score), 4, 4)
	ebitenutil.DebugPrintAt(screen, "HIGH "+itoa(g.highScore), 4, 16)
	if g.multiplier() > 1 {
		ebitenutil.DebugPrintAt(screen, "x"+itoa(g.multiplier()), 4, 28)
	}

	g.drawHealthBar(screen)
	ebitenutil.DebugPrintAt(screen, "VIDAS "+itoa(g.lives), 4, ScreenHeight-28)
	ebitenutil.DebugPrintAt(screen, "BOMBA "+itoa(g.bombCharges), 4, ScreenHeight-16)

	weapon := weaponName(g.player.weapon) + " Nv" + itoa(g.player.weaponLevel)
	ebitenutil.DebugPrintAt(screen, weapon, ScreenWidth-len(weapon)*6-4, ScreenHeight-16)
	if g.player.hasShield() {
		ebitenutil.DebugPrintAt(screen, "ESCUDO", ScreenWidth-46, ScreenHeight-28)
	}

	// Barra de progresso discreta na borda superior.
	barWidth := float32(g.level.progress() * ScreenWidth)
	vector.DrawFilledRect(screen, 0, 0, ScreenWidth, 2, color.RGBA{0x22, 0x22, 0x33, 0xff}, false)
	vector.DrawFilledRect(screen, 0, 0, barWidth, 2, color.RGBA{0xd0, 0xc0, 0x50, 0xff}, false)

	name := g.level.sectionName()
	ebitenutil.DebugPrintAt(screen, name, ScreenWidth-len(name)*6-4, 4)
}

// drawHealthBar mostra os pontos de vida como pequenos blocos.
func (g *Game) drawHealthBar(screen *ebiten.Image) {
	const pip = 8
	for i := 0; i < maxHealth; i++ {
		c := color.RGBA{0x40, 0x20, 0x20, 0xff}
		if i < g.player.health {
			c = color.RGBA{0xe0, 0x50, 0x50, 0xff}
		}
		vector.DrawFilledRect(screen, float32(60+i*(pip+2)), 4, pip, pip, c, false)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// collides testa sobreposição entre dois retângulos AABB.
func collides(ax, ay, aw, ah, bx, by, bw, bh float64) bool {
	return ax < bx+bw && ax+aw > bx && ay < by+bh && ay+ah > by
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
