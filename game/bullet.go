package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	bulletWidth  = 3
	bulletHeight = 8
)

type Bullet struct {
	x, y       float64
	vx, vy     float64
	w, h       float64
	damage     int
	pierce     int // inimigos adicionais que ainda pode atravessar
	dead       bool
	trail      bool // deixa rastro de partículas (projéteis especiais)
	hitEnemies []*Enemy
	color      color.RGBA
}

func newBullet(x, y, vx, vy float64, damage, pierce int, c color.RGBA) *Bullet {
	return &Bullet{
		x: x, y: y, vx: vx, vy: vy,
		w: bulletWidth, h: bulletHeight,
		damage: damage, pierce: pierce, color: c,
	}
}

func (b *Bullet) update() {
	b.x += b.vx
	b.y += b.vy
}

func (b *Bullet) offScreen() bool {
	return b.y+b.h < 0 || b.y > ScreenHeight || b.x+b.w < 0 || b.x > ScreenWidth
}

func (b *Bullet) alreadyHit(e *Enemy) bool {
	for _, hit := range b.hitEnemies {
		if hit == e {
			return true
		}
	}
	return false
}

func (b *Bullet) draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(b.w), float32(b.h), b.color, false)
}
