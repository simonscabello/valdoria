package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	griffinBody   = color.RGBA{0xd0, 0xa8, 0x5a, 0xff}
	griffinWing   = color.RGBA{0x9a, 0x78, 0x3a, 0xff}
	griffinBeak   = color.RGBA{0xff, 0xd7, 0x4d, 0xff}
	knightArmor   = color.RGBA{0x4d, 0xd0, 0xff, 0xff}
	shieldRing    = color.RGBA{0x6a, 0xd0, 0xff, 0xff}
	precisionMark = color.RGBA{0xff, 0x3a, 0x3a, 0xff}
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
	wingPhase   float64
	tilt        float64
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

	dx, dy = normalizeDiagonal(dx, dy)
	p.x += dx * speed
	p.y += dy * speed
	p.clampToScreen()

	p.wingPhase += playerWingSpeed
	p.tilt += (dx*playerTiltMax - p.tilt) * playerTiltRate

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

// normalizeDiagonal mantém a mesma velocidade nas oito direções: a diagonal é
// reduzida por 1/√2 para não ficar ~41% mais rápida que o movimento reto.
func normalizeDiagonal(dx, dy float64) (float64, float64) {
	if dx != 0 && dy != 0 {
		const invSqrt2 = 0.70710678
		dx *= invSqrt2
		dy *= invSqrt2
	}
	return dx, dy
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

// gainWeapon aplica uma runa de arma. Trocar de arma nunca rebaixa o poder: o
// nível atual é preservado ao trocar (o nível funciona como uma potência
// compartilhada), e coletar a runa da arma já equipada sobe um nível até o
// máximo. Assim, todo power-up é desejável — nunca uma armadilha a ser evitada.
func (p *Player) gainWeapon(w weaponType) {
	if p.weapon != w {
		p.weapon = w
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
	if dev.invincible {
		return false
	}
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
		p.drawGriffin(screen)
	}

	if p.muzzleFlash > 0 {
		flash := color.RGBA{0xff, 0xff, 0xff, 0xff}
		vector.DrawFilledRect(screen, float32(p.x+bulletInset-1), float32(p.y-2), bulletWidth+2, 3, flash, false)
		vector.DrawFilledRect(screen, float32(p.x+playerSize-bulletInset-bulletWidth-1), float32(p.y-2), bulletWidth+2, 3, flash, false)
	}

	if p.precision {
		hx, hy, hw, hh := p.hitbox()
		vector.DrawFilledRect(screen, float32(hx), float32(hy), float32(hw), float32(hh), precisionMark, false)
	}

	if p.shieldTimer > 0 {
		vector.StrokeRect(screen, float32(p.x-2), float32(p.y-2), playerSize+4, playerSize+4, 1, shieldRing, false)
	}
}

// drawGriffin desenha a silhueta do grifo com o cavaleiro montado, asas em
// batida e leve inclinação ao mover para os lados.
func (p *Player) drawGriffin(screen *ebiten.Image) {
	// Sprite do grifo/cavaleiro: batida de asa (2 frames), inclinação e bob.
	rot := p.tilt * 0.035
	bob := math.Sin(p.wingPhase) * 0.6
	name := wingFrameName("player", p.wingPhase)
	if drawSprite(screen, name, p.centerX(), p.centerY()+bob, 1, false, rot, false) {
		return
	}

	// Fallback geométrico.
	cx := float32(p.centerX())
	top := float32(p.y)
	flap := float32(math.Sin(p.wingPhase) * 3)
	tilt := float32(p.tilt)

	const wingW, wingH = 7, 4
	vector.DrawFilledRect(screen, cx-2-wingW, top+4-flap-tilt, wingW, wingH, griffinWing, false)
	vector.DrawFilledRect(screen, cx+2, top+4-flap+tilt, wingW, wingH, griffinWing, false)

	// Contorno escuro atrás do corpo, para o grifo destacar-se do cenário.
	vector.DrawFilledRect(screen, cx-4, top+1, 8, playerSize-3, enemyOutline, false)
	vector.DrawFilledRect(screen, cx-3, top+2, 6, playerSize-5, griffinBody, false)
	vector.DrawFilledRect(screen, cx-1, top+playerSize-3, 2, 4, griffinWing, false)
	vector.DrawFilledRect(screen, cx-2, top, 4, 3, griffinBody, false)
	vector.DrawFilledRect(screen, cx-1, top-3, 2, 3, griffinBeak, false)

	rider := cx - 2 + tilt*0.5
	vector.DrawFilledRect(screen, rider, top+5, 4, 5, knightArmor, false)
	vector.DrawFilledRect(screen, rider+1.5, top+1, 1, 4, knightArmor, false)
}
