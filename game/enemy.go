package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type enemyKind int

const (
	kindCrow enemyKind = iota
	kindHarpy
	kindGargoyle
	kindWyvern
	kindBallista
	kindMage
)

// enemyContext expõe ao comportamento apenas o que ele precisa do mundo:
// a posição do jogador (para mira) e a lista de projéteis inimigos.
type enemyContext struct {
	playerX, playerY float64
	bullets          *[]*EnemyBullet
}

type enemyBehavior interface {
	update(e *Enemy, ctx *enemyContext)
}

type Enemy struct {
	kind        enemyKind
	x, y        float64
	w, h        float64
	health      int
	score       int
	damage      int
	dead        bool
	escaped     bool
	hitFlash    int
	telegraph   int
	hasDrop     bool
	drop        powerupType
	formationID int
	color       color.RGBA
	accent      color.RGBA
	behavior    enemyBehavior
	animTick    int

	// Status elementais.
	slow      int     // frames restantes de gelo
	slowSkip  float64 // acumulador para rodar behavior a iceSlowFactor
	burn      int
	burnTick  int
	stun      int
}

func (e *Enemy) update(ctx *enemyContext) {
	if e.hitFlash > 0 {
		e.hitFlash--
	}
	if e.telegraph > 0 {
		e.telegraph--
	}
	e.animTick++
	e.tickStatusEffects()
	if e.dead {
		return
	}
	if e.stun > 0 {
		return
	}
	if e.slow > 0 {
		e.slowSkip += iceSlowFactor
		if e.slowSkip < 1 {
			return
		}
		e.slowSkip--
	}
	e.behavior.update(e, ctx)
}

func (e *Enemy) takeDamage(n int) {
	if e.burn > 0 && burnBonusPct > 0 {
		n += n * burnBonusPct / 100
	}
	e.health -= n
	e.hitFlash = hitFlashDuration
	if e.health <= 0 {
		e.dead = true
	}
}

// aimTelegraph agenda o aviso visual de disparo: quando faltarem
// enemyTelegraphFrames para atirar, o inimigo pisca um contorno de alerta,
// dando ao jogador tempo para reagir (evita tiros "surpresa").
func aimTelegraph(e *Enemy, fireCD int) {
	if fireCD == enemyTelegraphFrames {
		e.telegraph = enemyTelegraphFrames
	}
}

func (e *Enemy) centerX() float64 { return e.x + e.w/2 }
func (e *Enemy) centerY() float64 { return e.y + e.h/2 }

// crossedBottom indica que o inimigo passou pela base da tela (fuga vertical).
// Saídas laterais (gárgula) não contam — só a fuga "por baixo" pune o jogador.
func (e *Enemy) crossedBottom() bool {
	return e.y > ScreenHeight
}

// offScreen remove qualquer inimigo que saia bem além da área visível por
// qualquer lado, cobrindo tanto quem desce quanto a gárgula que sai pela lateral.
func (e *Enemy) offScreen() bool {
	const margin = 48
	return e.x < -e.w-margin || e.x > ScreenWidth+margin ||
		e.y < -e.h-margin || e.y > ScreenHeight+margin
}

func (e *Enemy) draw(screen *ebiten.Image) {
	e.drawStatusAura(screen)

	// Aviso de disparo: contorno de alerta piscando antes de o inimigo atirar.
	if e.telegraph > 0 && (e.telegraph/3)%2 == 0 {
		warn := color.RGBA{0xff, 0xe0, 0x40, 0xff}
		vector.StrokeRect(screen, float32(e.x-2), float32(e.y-2), float32(e.w+4), float32(e.h+4), 1, warn, false)
	}

	// Sprite (pixel art com contorno próprio); voadores batem asas em 2 frames.
	bob := 0.0
	phase := float64(e.animTick) * enemyWingSpeed
	if e.kind != kindGargoyle && e.kind != kindBallista {
		bob = math.Sin(phase) * 1.0
	}
	base := enemySpriteName(e.kind)
	name := base
	if e.kind != kindGargoyle && e.kind != kindBallista && e.kind != kindMage {
		name = wingFrameName(base, phase)
	}
	// Escala sempre pelo frame base, para a batida de asa não pulsar o tamanho.
	if sw, sh, ok := spriteBounds(base); ok {
		scale := math.Max(e.w/float64(sw), e.h/float64(sh)) * 1.1
		drawSprite(screen, name, e.centerX(), e.centerY()+bob, scale, false, 0, e.hitFlash > 0)
		return
	}

	// Fallback geométrico (quando não há sprite disponível).
	vector.DrawFilledRect(screen, float32(e.x-1), float32(e.y-1), float32(e.w+2), float32(e.h+2), enemyOutline, false)

	body := e.color
	if e.hitFlash > 0 {
		body = color.RGBA{0xff, 0xff, 0xff, 0xff}
	}
	x, y := float32(e.x), float32(e.y)
	w, h := float32(e.w), float32(e.h)
	// Batida de asas: as asas sobem e descem em ciclo próprio de cada inimigo.
	flap := float32(math.Sin(float64(e.animTick)*enemyWingSpeed) * 2)

	switch e.kind {
	case kindCrow:
		// Corvo: corpo pequeno e escuro com asas pontudas curtas.
		vector.DrawFilledRect(screen, x, y, w, h, body, false)
		vector.DrawFilledRect(screen, x-3, y+2-flap, 3, 2, e.accent, false)
		vector.DrawFilledRect(screen, x+w, y+2-flap, 3, 2, e.accent, false)
	case kindHarpy:
		// Harpia: asas grandes e claras, corpo mais estreito.
		vector.DrawFilledRect(screen, x-4, y-flap, 4, h*0.6, e.accent, false)
		vector.DrawFilledRect(screen, x+w, y-flap, 4, h*0.6, e.accent, false)
		vector.DrawFilledRect(screen, x+2, y, w-4, h, body, false)
		vector.DrawFilledRect(screen, x+w/2-1, y-3, 2, 3, e.accent, false)
	case kindGargoyle:
		// Gárgula: corpo robusto de pedra com chifres e asas curtas rígidas.
		vector.DrawFilledRect(screen, x, y, w, h, body, false)
		vector.DrawFilledRect(screen, x-3, y+4, 3, h-8, e.accent, false)
		vector.DrawFilledRect(screen, x+w, y+4, 3, h-8, e.accent, false)
		vector.DrawFilledRect(screen, x+1, y-3, 3, 3, e.accent, false)
		vector.DrawFilledRect(screen, x+w-4, y-3, 3, 3, e.accent, false)
	case kindWyvern:
		// Wyvern: grande, asas membranosas amplas, focinho e cauda visíveis.
		vector.DrawFilledRect(screen, x-6, y+2-flap, 6, h-6, e.accent, false)
		vector.DrawFilledRect(screen, x+w, y+2-flap, 6, h-6, e.accent, false)
		vector.DrawFilledRect(screen, x, y, w, h, body, false)
		vector.DrawFilledRect(screen, x+w/2-2, y-4, 4, 4, body, false)
		vector.DrawFilledRect(screen, x+w/2-2, y+h, 4, 5, e.accent, false)
	case kindBallista:
		// Balista: base larga e baixa de madeira/pedra com um braço de disparo.
		vector.DrawFilledRect(screen, x, y+h*0.4, w, h*0.6, body, false)
		vector.DrawFilledRect(screen, x+w*0.15, y+h*0.7, w*0.7, h*0.3, e.accent, false)
		vector.DrawFilledRect(screen, x+w*0.4, y, w*0.2, h*0.6, e.accent, false)
		vector.DrawFilledRect(screen, x+w*0.2, y+h*0.2, w*0.6, 2, body, false)
	case kindMage:
		// Feiticeiro: manto escuro, capuz e um núcleo mágico brilhante pulsante.
		vector.DrawFilledRect(screen, x, y+2, w, h-2, body, false)
		vector.DrawFilledRect(screen, x+w*0.25, y, w*0.5, 3, body, false)
		glow := float32((math.Sin(float64(e.animTick)*0.2) + 1) * 1.5)
		vector.DrawFilledRect(screen, x+w/2-2-glow/2, y+h*0.4, 4+glow, 4+glow, e.accent, false)
	}
}

// enemySpriteName mapeia o tipo de inimigo ao nome do seu sprite.
func enemySpriteName(k enemyKind) string {
	switch k {
	case kindHarpy:
		return "harpy"
	case kindGargoyle:
		return "gargoyle"
	case kindWyvern:
		return "wyvern"
	case kindBallista:
		return "ballista"
	case kindMage:
		return "mage"
	default:
		return "crow"
	}
}

// aimVelocity devolve a velocidade de um projétil apontado da origem ao alvo,
// com módulo igual a speed. Alvo coincidente com a origem cai para baixo.
func aimVelocity(fromX, fromY, toX, toY, speed float64) (vx, vy float64) {
	dx := toX - fromX
	dy := toY - fromY
	dist := math.Hypot(dx, dy)
	if dist == 0 {
		return 0, speed
	}
	return dx / dist * speed, dy / dist * speed
}

func newCrow(x float64) *Enemy {
	return &Enemy{
		kind:     kindCrow,
		x:        x,
		y:        -crowSize,
		w:        crowSize,
		h:        crowSize,
		health:   scaledEnemyHealth(crowHealth),
		score:    crowScore,
		damage:   crowDamage,
		color:    color.RGBA{0x6b, 0x4a, 0x8a, 0xff},
		accent:   color.RGBA{0x2a, 0x1a, 0x3a, 0xff},
		behavior: &crowBehavior{speed: crowSpeed},
	}
}

func newHarpy(x float64) *Enemy {
	return &Enemy{
		kind:     kindHarpy,
		x:        x,
		y:        -harpySize,
		w:        harpySize,
		h:        harpySize,
		health:   scaledEnemyHealth(harpyHealth),
		score:    harpyScore,
		damage:   harpyDamage,
		color:    color.RGBA{0x3a, 0xb0, 0x9e, 0xff},
		accent:   color.RGBA{0xe0, 0xe0, 0xf0, 0xff},
		behavior: &harpyBehavior{baseX: x, fireCD: harpyFireInterval},
	}
}

func newGargoyle(fromLeft bool, targetX, y float64) *Enemy {
	x := float64(-gargoyleSize)
	if !fromLeft {
		x = ScreenWidth
	}
	return &Enemy{
		kind:     kindGargoyle,
		x:        x,
		y:        y,
		w:        gargoyleSize,
		h:        gargoyleSize,
		health:   scaledEnemyHealth(gargoyleHealth),
		score:    gargoyleScore,
		damage:   gargoyleDamage,
		color:    color.RGBA{0x7a, 0x7d, 0x85, 0xff},
		accent:   color.RGBA{0x3a, 0x3d, 0x45, 0xff},
		behavior: &gargoyleBehavior{fromLeft: fromLeft, targetX: targetX, attackTimer: gargoyleAttackDuration, fireCD: gargoyleFireInterval},
	}
}

func newWyvern(x float64) *Enemy {
	return &Enemy{
		kind:     kindWyvern,
		x:        x,
		y:        -wyvernSize,
		w:        wyvernSize,
		h:        wyvernSize,
		health:   scaledEnemyHealth(wyvernHealth),
		score:    wyvernScore,
		damage:   wyvernDamage,
		color:    color.RGBA{0x4a, 0x8a, 0x3a, 0xff},
		accent:   color.RGBA{0x2a, 0x5a, 0x22, 0xff},
		behavior: &wyvernBehavior{fireCD: wyvernFireInterval},
	}
}

func newBallista(x float64) *Enemy {
	return &Enemy{
		kind:     kindBallista,
		x:        x,
		y:        -ballistaSize,
		w:        ballistaSize,
		h:        ballistaSize * 0.7,
		health:   scaledEnemyHealth(ballistaHealth),
		score:    ballistaScore,
		damage:   ballistaDamage,
		color:    color.RGBA{0x6a, 0x52, 0x36, 0xff},
		accent:   color.RGBA{0x3a, 0x2c, 0x1c, 0xff},
		behavior: &ballistaBehavior{fireCD: ballistaFireInterval},
	}
}

func newMage(x float64) *Enemy {
	return &Enemy{
		kind:     kindMage,
		x:        x,
		y:        -mageSize,
		w:        mageSize,
		h:        mageSize,
		health:   scaledEnemyHealth(mageHealth),
		score:    mageScore,
		damage:   mageDamage,
		color:    color.RGBA{0x7a, 0x3a, 0xb0, 0xff},
		accent:   color.RGBA{0xd0, 0xa0, 0xff, 0xff},
		behavior: &mageBehavior{baseX: x, fireCD: mageFireInterval},
	}
}

// Corvo corrompido: desce em linha reta e não dispara.
type crowBehavior struct {
	speed float64
}

func (b *crowBehavior) update(e *Enemy, _ *enemyContext) {
	e.y += b.speed
}

// Harpia: desce em zigue-zague e dispara para baixo em intervalos regulares.
type harpyBehavior struct {
	baseX  float64
	phase  float64
	fireCD int
}

func (b *harpyBehavior) update(e *Enemy, ctx *enemyContext) {
	b.phase += harpyFreq
	e.y += harpySpeed
	e.x = b.baseX + math.Sin(b.phase)*harpyAmplitude

	b.fireCD--
	aimTelegraph(e, b.fireCD)
	if b.fireCD <= 0 {
		b.fireCD = harpyFireInterval
		*ctx.bullets = append(*ctx.bullets, newEnemyBullet(e.centerX(), e.y+e.h, 0, enemyBulletSpeed))
	}
}

type gargoylePhase int

const (
	gargoyleEntering gargoylePhase = iota
	gargoyleAttacking
	gargoyleLeaving
)

// Gárgula: entra pela lateral, para em uma posição, ataca por alguns segundos
// e sai pela lateral oposta.
type gargoyleBehavior struct {
	fromLeft    bool
	targetX     float64
	phase       gargoylePhase
	attackTimer int
	fireCD      int
}

func (b *gargoyleBehavior) update(e *Enemy, ctx *enemyContext) {
	switch b.phase {
	case gargoyleEntering:
		if b.fromLeft {
			e.x += gargoyleSpeed
			if e.x >= b.targetX {
				e.x = b.targetX
				b.phase = gargoyleAttacking
			}
			return
		}
		e.x -= gargoyleSpeed
		if e.x <= b.targetX {
			e.x = b.targetX
			b.phase = gargoyleAttacking
		}
	case gargoyleAttacking:
		b.attackTimer--
		b.fireCD--
		aimTelegraph(e, b.fireCD)
		if b.fireCD <= 0 {
			b.fireCD = gargoyleFireInterval
			*ctx.bullets = append(*ctx.bullets, newEnemyBullet(e.centerX(), e.y+e.h, 0, enemyBulletSpeed))
		}
		if b.attackTimer <= 0 {
			b.phase = gargoyleLeaving
		}
	case gargoyleLeaving:
		if b.fromLeft {
			e.x += gargoyleSpeed
			return
		}
		e.x -= gargoyleSpeed
	}
}

// Wyvern: desce devagar tentando se alinhar ao jogador e dispara projéteis
// mirados na posição do jogador no instante do disparo.
type wyvernBehavior struct {
	fireCD int
}

func (b *wyvernBehavior) update(e *Enemy, ctx *enemyContext) {
	e.y += wyvernDescend

	diff := ctx.playerX - e.centerX()
	if diff > wyvernAlign {
		e.x += wyvernAlign
	} else if diff < -wyvernAlign {
		e.x -= wyvernAlign
	}

	b.fireCD--
	aimTelegraph(e, b.fireCD)
	if b.fireCD <= 0 {
		b.fireCD = wyvernFireInterval
		vx, vy := aimVelocity(e.centerX(), e.y+e.h, ctx.playerX, ctx.playerY, enemyBulletSpeed)
		*ctx.bullets = append(*ctx.bullets, newEnemyBullet(e.centerX(), e.y+e.h, vx, vy))
	}
}

// Balista corrompida: desce lenta como o cenário e solta rajadas de virotes
// mirados, sempre telegrafadas. Ameaça "terrestre" resistente.
type ballistaBehavior struct {
	fireCD    int
	burstLeft int
	burstCD   int
}

func (b *ballistaBehavior) update(e *Enemy, ctx *enemyContext) {
	e.y += ballistaDescend

	if b.burstLeft > 0 {
		b.burstCD--
		if b.burstCD <= 0 {
			vx, vy := aimVelocity(e.centerX(), e.y+e.h, ctx.playerX, ctx.playerY, enemyBulletSpeed)
			*ctx.bullets = append(*ctx.bullets, newEnemyBullet(e.centerX(), e.y+e.h, vx, vy))
			b.burstLeft--
			b.burstCD = ballistaBurstGap
		}
		return
	}

	b.fireCD--
	aimTelegraph(e, b.fireCD)
	if b.fireCD <= 0 {
		b.fireCD = ballistaFireInterval
		b.burstLeft = ballistaBurst
		b.burstCD = 0
	}
}

// Feiticeiro corrompido: mergulha até o alto da tela, desliza de lado e libera
// anéis completos de projéteis. Frágil, mas perigoso se ignorado.
type mageBehavior struct {
	baseX  float64
	phase  float64
	fireCD int
}

func (b *mageBehavior) update(e *Enemy, ctx *enemyContext) {
	if e.y < mageStopY {
		e.y += mageDescend
	} else {
		b.phase += 0.03
		e.x = b.baseX + math.Sin(b.phase)*mageDrift*20
	}

	b.fireCD--
	aimTelegraph(e, b.fireCD)
	if b.fireCD <= 0 {
		b.fireCD = mageFireInterval
		for i := 0; i < mageRingCount; i++ {
			a := float64(i) / float64(mageRingCount) * 2 * math.Pi
			vx := math.Sin(a) * mageBulletSpeed
			vy := math.Cos(a) * mageBulletSpeed
			*ctx.bullets = append(*ctx.bullets, newEnemyBullet(e.centerX(), e.centerY(), vx, vy))
		}
	}
}
