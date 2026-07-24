package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// drawBulletGlow desenha um halo leve (1px) atrás do projétil — só o bastante
// para leitura, sem o brilho forte que ofuscava a cena.
func drawBulletGlow(dst *ebiten.Image, x, y, w, h float32, c color.RGBA) {
	vector.DrawFilledRect(dst, x-1, y-1, w+2, h+2, withAlpha(c, 55), false)
}

// brighten empurra a cor em direção ao branco para pontas/núcleos emissivos.
func brighten(c color.RGBA, amount int) color.RGBA {
	return color.RGBA{
		R: clampByte(int(c.R) + amount),
		G: clampByte(int(c.G) + amount),
		B: clampByte(int(c.B) + amount),
		A: c.A,
	}
}

func clampByte(v int) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}
