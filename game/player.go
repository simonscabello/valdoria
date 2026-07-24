package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
	x, y        float64
	health      int
	cooldown    int
	invincible  int
	muzzleFlash int
	precision   bool
}

func newPlayer() *Player {
	return &Player{
		x:      ScreenWidth/2 - playerSize/2,
		y:      ScreenHeight - playerSize*3,
		health: initialHealth,
	}
}

func (p *Player) update() {
	p.precision = ebiten.IsKeyPressed(ebiten.KeyShift)

	speed := playerSpeed
	if p.precision {
		speed = playerPrecisionSpeed
	}

	var dx, dy float64
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		dx--
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		dx++
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		dy--
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		dy++
	}

	// Normaliza a diagonal para manter a mesma velocidade em oito direções.
	if dx != 0 && dy != 0 {
		const invSqrt2 = 0.7071
		dx *= invSqrt2
		dy *= invSqrt2
	}

	p.x += dx * speed
	p.y += dy * speed
	p.clampToScreen()

	if p.cooldown > 0 {
		p.cooldown--
	}
	if p.invincible > 0 {
		p.invincible--
	}
	if p.muzzleFlash > 0 {
		p.muzzleFlash--
	}
}

func (p *Player) clampToScreen() {
	if p.x < 0 {
		p.x = 0
	}
	if p.x > ScreenWidth-playerSize {
		p.x = ScreenWidth - playerSize
	}
	if p.y < 0 {
		p.y = 0
	}
	if p.y > ScreenHeight-playerSize {
		p.y = ScreenHeight - playerSize
	}
}

func (p *Player) tryShoot() []*Bullet {
	if !ebiten.IsKeyPressed(ebiten.KeySpace) || p.cooldown > 0 {
		return nil
	}
	p.cooldown = shootCooldown
	p.muzzleFlash = muzzleFlashDuration
	return []*Bullet{
		newBullet(p.x+bulletInset, p.y),
		newBullet(p.x+playerSize-bulletInset-bulletWidth, p.y),
	}
}

// hitbox devolve a área real de colisão do jogador, menor que o corpo desenhado.
func (p *Player) hitbox() (x, y, w, h float64) {
	const offset = (playerSize - playerHitboxSize) / 2
	return p.x + offset, p.y + offset, playerHitboxSize, playerHitboxSize
}

func (p *Player) canBeHit() bool {
	return p.invincible == 0
}

func (p *Player) hit(damage int) {
	p.health -= damage
	p.invincible = invincibilityDuration
}

func (p *Player) centerX() float64 { return p.x + playerSize/2 }
func (p *Player) centerY() float64 { return p.y + playerSize/2 }

func (p *Player) draw(screen *ebiten.Image) {
	// Durante a invencibilidade o corpo pisca alternando ciclos.
	blinkHidden := p.invincible > 0 && (p.invincible/invincibilityBlink)%2 == 0
	if !blinkHidden {
		knight := color.RGBA{0x4d, 0xd0, 0xff, 0xff}
		vector.DrawFilledRect(screen, float32(p.x), float32(p.y), playerSize, playerSize, knight, false)

		beak := color.RGBA{0xff, 0xd7, 0x4d, 0xff}
		vector.DrawFilledRect(screen, float32(p.x+playerSize/2-2), float32(p.y-4), 4, 4, beak, false)
	}

	if p.muzzleFlash > 0 {
		flash := color.RGBA{0xff, 0xff, 0xff, 0xff}
		vector.DrawFilledRect(screen, float32(p.x+bulletInset-1), float32(p.y-2), bulletWidth+2, 3, flash, false)
		vector.DrawFilledRect(screen, float32(p.x+playerSize-bulletInset-bulletWidth-1), float32(p.y-2), bulletWidth+2, 3, flash, false)
	}

	if p.precision {
		hx, hy, hw, hh := p.hitbox()
		mark := color.RGBA{0xff, 0x3a, 0x3a, 0xff}
		vector.DrawFilledRect(screen, float32(hx), float32(hy), float32(hw), float32(hh), mark, false)
	}
}
