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
	kind     enemyKind
	x, y     float64
	w, h     float64
	health   int
	score    int
	damage   int
	dead     bool
	hitFlash int
	hasDrop  bool
	drop     powerupType
	color    color.RGBA
	accent   color.RGBA
	behavior enemyBehavior
}

func (e *Enemy) update(ctx *enemyContext) {
	if e.hitFlash > 0 {
		e.hitFlash--
	}
	e.behavior.update(e, ctx)
}

func (e *Enemy) takeDamage(n int) {
	e.health -= n
	e.hitFlash = hitFlashDuration
	if e.health <= 0 {
		e.dead = true
	}
}

func (e *Enemy) centerX() float64 { return e.x + e.w/2 }
func (e *Enemy) centerY() float64 { return e.y + e.h/2 }

// offScreen remove qualquer inimigo que saia bem além da área visível por
// qualquer lado, cobrindo tanto quem desce quanto a gárgula que sai pela lateral.
func (e *Enemy) offScreen() bool {
	const margin = 48
	return e.x < -e.w-margin || e.x > ScreenWidth+margin ||
		e.y < -e.h-margin || e.y > ScreenHeight+margin
}

func (e *Enemy) draw(screen *ebiten.Image) {
	body := e.color
	if e.hitFlash > 0 {
		body = color.RGBA{0xff, 0xff, 0xff, 0xff}
	}
	x, y := float32(e.x), float32(e.y)
	w, h := float32(e.w), float32(e.h)
	vector.DrawFilledRect(screen, x, y, w, h, body, false)

	switch e.kind {
	case kindCrow:
		vector.DrawFilledRect(screen, x-2, y+2, 2, 3, e.accent, false)
		vector.DrawFilledRect(screen, x+w, y+2, 2, 3, e.accent, false)
	case kindHarpy:
		vector.DrawFilledRect(screen, x-3, y, 3, h/2, e.accent, false)
		vector.DrawFilledRect(screen, x+w, y, 3, h/2, e.accent, false)
	case kindGargoyle:
		vector.DrawFilledRect(screen, x, y-3, 3, 3, e.accent, false)
		vector.DrawFilledRect(screen, x+w-3, y-3, 3, 3, e.accent, false)
	case kindWyvern:
		vector.DrawFilledRect(screen, x-5, y+2, 5, h-4, e.accent, false)
		vector.DrawFilledRect(screen, x+w, y+2, 5, h-4, e.accent, false)
		vector.DrawFilledRect(screen, x+w/2-2, y+h, 4, 4, e.accent, false)
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
		health:   crowHealth,
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
		health:   harpyHealth,
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
		health:   gargoyleHealth,
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
		health:   wyvernHealth,
		score:    wyvernScore,
		damage:   wyvernDamage,
		color:    color.RGBA{0x4a, 0x8a, 0x3a, 0xff},
		accent:   color.RGBA{0x2a, 0x5a, 0x22, 0xff},
		behavior: &wyvernBehavior{fireCD: wyvernFireInterval},
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
	if b.fireCD <= 0 {
		b.fireCD = wyvernFireInterval
		vx, vy := aimVelocity(e.centerX(), e.y+e.h, ctx.playerX, ctx.playerY, enemyBulletSpeed)
		*ctx.bullets = append(*ctx.bullets, newEnemyBullet(e.centerX(), e.y+e.h, vx, vy))
	}
}
