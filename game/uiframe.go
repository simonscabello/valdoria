package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Molduras: painéis de menu (ferro/pergaminho) vs chips de HUD (placa leve).
var (
	uiFrameOuter = color.RGBA{0x12, 0x0e, 0x0a, 0xff}
	uiFrameIron  = color.RGBA{0x7a, 0x62, 0x3e, 0xff}
	uiFrameGold  = color.RGBA{0xd0, 0xb0, 0x5a, 0xff}
	uiFrameFill  = color.RGBA{0x28, 0x20, 0x16, 0xd8}

	uiChipFill   = color.RGBA{0x08, 0x0a, 0x14, 0xb0} // fundo escuro translúcido
	uiChipEdge   = color.RGBA{0xe0, 0xc4, 0x6a, 0xa0} // filete dourado suave
	uiChipShadow = color.RGBA{0x00, 0x00, 0x00, 0x50}
)

// drawUIPanel desenha um painel com moldura de ferro e interior de pergaminho
// (menus / pausa — mais presença visual).
func drawUIPanel(dst *ebiten.Image, x, y, w, h float32) {
	vector.DrawFilledRect(dst, x+2, y+2, w, h, uiChipShadow, false)
	vector.DrawFilledRect(dst, x, y, w, h, uiFrameOuter, false)
	vector.DrawFilledRect(dst, x+1, y+1, w-2, h-2, uiFrameIron, false)
	vector.DrawFilledRect(dst, x+2, y+2, w-4, h-4, uiFrameFill, false)
	vector.DrawFilledRect(dst, x+2, y+2, w-4, 1, uiFrameGold, false)
	vector.DrawFilledRect(dst, x+2, y+h-3, w-4, 1, withAlpha(uiFrameGold, 100), false)
}

// drawUIChip desenha uma placa leve para o HUD: sombra + fill translúcido +
// um único filete dourado (sem moldura tripla que pesava a tela).
func drawUIChip(dst *ebiten.Image, x, y, w, h float32) {
	vector.DrawFilledRect(dst, x+1, y+1, w, h, uiChipShadow, false)
	vector.DrawFilledRect(dst, x, y, w, h, uiChipFill, false)
	vector.StrokeRect(dst, x, y, w, h, 1, uiChipEdge, false)
}

// drawIronBar desenha uma barra (vida do chefe / progresso) com trilho discreto.
func drawIronBar(dst *ebiten.Image, x, y, w, h float32, fill color.RGBA, ratio float64) {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	vector.DrawFilledRect(dst, x, y, w, h, color.RGBA{0x10, 0x0c, 0x18, 0xc0}, false)
	vector.StrokeRect(dst, x, y, w, h, 1, withAlpha(uiChipEdge, 140), false)
	innerW := (w - 2) * float32(ratio)
	if innerW > 0 {
		vector.DrawFilledRect(dst, x+1, y+1, innerW, h-2, fill, false)
	}
}
