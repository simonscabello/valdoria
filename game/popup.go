package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// scorePopup é um número que sobe brevemente sobre o inimigo abatido,
// reforçando visualmente o valor (já com multiplicador) de cada eliminação.
type scorePopup struct {
	x, y float64
	text string
	life int
}

func (g *Game) spawnScorePopup(x, y float64, points int) {
	g.spawnTextPopup(x, y, itoa(points))
}

// spawnBonusPopup destaca recompensas ocultas (formação, trecho sem dano) com
// um "+" para diferenciá-las das pontuações normais de eliminação.
func (g *Game) spawnBonusPopup(x, y float64, points int) {
	g.spawnTextPopup(x, y, "+"+itoa(points))
}

func (g *Game) spawnTextPopup(x, y float64, text string) {
	if len(g.popups) >= maxPopups {
		return
	}
	g.popups = append(g.popups, &scorePopup{x: x, y: y, text: text, life: popupLife})
}

func (g *Game) updatePopups() {
	for _, p := range g.popups {
		p.y -= popupRise
		p.life--
	}
	g.popups = filterAlive(g.popups, func(p *scorePopup) bool { return p.life > 0 })
}

func (g *Game) drawPopups(dst *ebiten.Image) {
	for _, p := range g.popups {
		ebitenutil.DebugPrintAt(dst, p.text, int(p.x)-len(p.text)*3, int(p.y))
	}
}
