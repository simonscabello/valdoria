package game

import (
	"image/color"

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
	kind powerupType
	x, y float64
	dead bool
}

func newPowerup(kind powerupType, centerX, centerY float64) *Powerup {
	return &Powerup{kind: kind, x: centerX - powerupSize/2, y: centerY - powerupSize/2}
}

func (p *Powerup) update() {
	p.y += powerupSpeed
}

func (p *Powerup) offScreen() bool {
	return p.y > ScreenHeight
}

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

func (p *Powerup) draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(p.x), float32(p.y), powerupSize, powerupSize, p.color(), false)
	// Núcleo escuro para diferenciar do restante das formas.
	vector.DrawFilledRect(screen, float32(p.x+3), float32(p.y+3), powerupSize-6, powerupSize-6, color.RGBA{0x10, 0x10, 0x18, 0xff}, false)
}
