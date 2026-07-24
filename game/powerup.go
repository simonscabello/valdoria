package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type powerupType int

const (
	powerLight powerupType = iota
	powerFire
	powerIce
	powerHeal
	powerShield
)

type Powerup struct {
	kind     powerupType
	x, y     float64
	dead     bool
	animTick int
}

func newPowerup(kind powerupType, centerX, centerY float64) *Powerup {
	return &Powerup{kind: kind, x: centerX - powerupSize/2, y: centerY - powerupSize/2}
}

func (p *Powerup) update() {
	p.y += powerupSpeed
	p.animTick++
}

func (p *Powerup) offScreen() bool {
	return p.y > ScreenHeight
}

func (p *Powerup) centerX() float64 { return p.x + powerupSize/2 }
func (p *Powerup) centerY() float64 { return p.y + powerupSize/2 }

func (p *Powerup) color() color.RGBA {
	switch p.kind {
	case powerLight:
		return lightColor
	case powerFire:
		return flameColor
	case powerIce:
		return iceColor
	case powerHeal:
		return color.RGBA{0x4a, 0xe0, 0x6a, 0xff}
	default:
		return color.RGBA{0xe0, 0xe0, 0xf0, 0xff}
	}
}

func powerupSpriteName(k powerupType) string {
	switch k {
	case powerFire:
		return "power_fire"
	case powerIce:
		return "power_ice"
	case powerHeal:
		return "power_heal"
	case powerShield:
		return "power_shield"
	default:
		return "power_light"
	}
}

func (p *Powerup) draw(screen *ebiten.Image) {
	name := powerupSpriteName(p.kind)
	bob := math.Sin(float64(p.animTick)*0.2) * 1.2
	if sw, sh, ok := spriteBounds(name); ok {
		scale := math.Max(powerupSize/float64(sw), powerupSize/float64(sh)) * 1.25
		drawSprite(screen, name, p.centerX(), p.centerY()+bob, scale, false, 0, false)
		return
	}

	// Fallback geométrico.
	vector.DrawFilledRect(screen, float32(p.x), float32(p.y+bob), powerupSize, powerupSize, p.color(), false)
	vector.DrawFilledRect(screen, float32(p.x+3), float32(p.y+3+bob), powerupSize-6, powerupSize-6, color.RGBA{0x10, 0x10, 0x18, 0xff}, false)
}
