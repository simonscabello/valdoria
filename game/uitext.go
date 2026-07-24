package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	etext "github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

// uiFace é a fonte própria da interface: uma bitmap embutida (sem asset externo)
// que substitui a fonte de depuração, permitindo texto colorido e medição exata
// de largura para alinhamentos corretos.
var uiFace = etext.NewGoXFace(basicfont.Face7x13)

// Cores da interface, com identidade de "pergaminho medieval".
var (
	uiInk       = color.RGBA{0xf0, 0xe8, 0xcf, 0xff} // texto principal (pergaminho)
	uiInkDim    = color.RGBA{0xb8, 0xac, 0x90, 0xff} // texto secundário
	uiHighlight = color.RGBA{0xff, 0xd7, 0x4d, 0xff} // seleção/destaque (dourado)
)

// textWidth mede a largura em pixels de um texto na fonte da interface.
func textWidth(s string) float64 {
	w, _ := etext.Measure(s, uiFace, 0)
	return w
}

// drawText desenha um texto colorido no canto superior esquerdo (x, y).
func drawText(dst *ebiten.Image, s string, x, y float64, col color.Color) {
	op := &etext.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(col)
	etext.Draw(dst, s, uiFace, op)
}

// drawTextCentered desenha um texto horizontalmente centralizado na tela.
func drawTextCentered(dst *ebiten.Image, s string, y float64, col color.Color) {
	drawText(dst, s, (ScreenWidth-textWidth(s))/2, y, col)
}
