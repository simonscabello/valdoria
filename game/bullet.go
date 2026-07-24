package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	bulletWidth  = 3
	bulletHeight = 8
	bulletSpeed  = 5
)

type Bullet struct {
	x, y float64
	dead bool
}

func newBullet(x, y float64) *Bullet {
	return &Bullet{x: x, y: y}
}

func (b *Bullet) update() {
	b.y -= bulletSpeed
}

func (b *Bullet) offScreen() bool {
	return b.y+bulletHeight < 0
}

func (b *Bullet) draw(screen *ebiten.Image) {
	c := color.RGBA{0xff, 0xf0, 0x8c, 0xff}
	vector.DrawFilledRect(screen, float32(b.x), float32(b.y), bulletWidth, bulletHeight, c, false)
}
