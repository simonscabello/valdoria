package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type EnemyBullet struct {
	x, y   float64
	vx, vy float64
	dead   bool
}

func newEnemyBullet(x, y, vx, vy float64) *EnemyBullet {
	return &EnemyBullet{x: x - enemyBulletSize/2, y: y, vx: vx, vy: vy}
}

func (b *EnemyBullet) update() {
	b.x += b.vx
	b.y += b.vy
}

func (b *EnemyBullet) offScreen() bool {
	return b.y > ScreenHeight || b.y+enemyBulletSize < 0 ||
		b.x > ScreenWidth || b.x+enemyBulletSize < 0
}

func (b *EnemyBullet) draw(screen *ebiten.Image) {
	c := color.RGBA{0xff, 0x6a, 0x3a, 0xff}
	vector.DrawFilledRect(screen, float32(b.x), float32(b.y), enemyBulletSize, enemyBulletSize, c, false)
}
