package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Projéteis inimigos usam magenta/roxo com núcleo claro: uma cor "ameaçadora"
// deliberadamente distinta de todas as armas do jogador (luz amarela, chamas
// laranjas, gelo ciano), para leitura imediata de aliado × inimigo.
var (
	enemyBulletColor   = color.RGBA{0xc8, 0x28, 0xdc, 0xff}
	enemyBulletCore    = color.RGBA{0xff, 0xda, 0xff, 0xff}
	enemyBulletOutline = color.RGBA{0x20, 0x08, 0x28, 0xff}
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
	// Halo escuro (contorno) + corpo magenta + núcleo claro: destaca o projétil
	// sobre qualquer fundo e o diferencia dos disparos do jogador.
	x, y := float32(b.x), float32(b.y)
	vector.DrawFilledRect(screen, x-1, y-1, enemyBulletSize+2, enemyBulletSize+2, enemyBulletOutline, false)
	vector.DrawFilledRect(screen, x, y, enemyBulletSize, enemyBulletSize, enemyBulletColor, false)
	vector.DrawFilledRect(screen, x+1, y+1, enemyBulletSize-2, enemyBulletSize-2, enemyBulletCore, false)
}
