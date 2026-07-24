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
	weapon      weaponType
	weaponLevel int
	shieldTimer int
}

func newPlayer() *Player {
	return &Player{
		x:           ScreenWidth/2 - playerSize/2,
		y:           ScreenHeight - playerSize*3,
		health:      initialHealth,
		weapon:      weaponLight,
		weaponLevel: 1,
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
	if p.shieldTimer > 0 {
		p.shieldTimer--
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
	p.cooldown = weaponCooldown(p.weapon, p.weaponLevel)
	p.muzzleFlash = muzzleFlashDuration
	return fireWeapon(p.weapon, p.weaponLevel, p.centerX(), p.y)
}

func (p *Player) hasShield() bool {
	return p.shieldTimer > 0
}

// applyPowerup coleta uma runa/benefício respeitando as regras de armas.
func (p *Player) applyPowerup(kind powerupType) {
	switch kind {
	case powerLight:
		p.gainWeapon(weaponLight)
	case powerFire:
		p.gainWeapon(weaponFlame)
	case powerIce:
		p.gainWeapon(weaponIce)
	case powerHeal:
		if p.health < maxHealth {
			p.health++
		}
	case powerShield:
		p.shieldTimer = shieldDuration
	}
}

func (p *Player) gainWeapon(w weaponType) {
	if p.weapon != w {
		p.weapon = w
		p.weaponLevel = 1
		return
	}
	if p.weaponLevel < maxWeaponLevel {
		p.weaponLevel++
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

// hit aplica dano e devolve true apenas quando a vida realmente diminuiu.
func (p *Player) hit(damage int) bool {
	p.invincible = invincibilityDuration
	if p.shieldTimer > 0 {
		p.shieldTimer = 0
		return false
	}
	p.health -= damage
	return true
}

// respawn recoloca o jogador em área segura, restaura a vida, concede
// invencibilidade e reduz apenas um nível da arma atual.
func (p *Player) respawn() {
	p.x = ScreenWidth/2 - playerSize/2
	p.y = ScreenHeight - playerSize*3
	p.health = maxHealth
	p.invincible = respawnInvincibility
	p.shieldTimer = 0
	if p.weaponLevel > 1 {
		p.weaponLevel--
	}
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

	if p.shieldTimer > 0 {
		ring := color.RGBA{0x6a, 0xd0, 0xff, 0xff}
		vector.StrokeRect(screen, float32(p.x-2), float32(p.y-2), playerSize+4, playerSize+4, 1, ring, false)
	}
}
