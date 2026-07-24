package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const starCount = 60

type star struct {
	x, y  float64
	speed float64
	size  float64
}

// bgItem é um elemento genérico de uma camada de parallax (nuvem, colina ou
// estrutura). Cada camada rola verticalmente em velocidade própria.
type bgItem struct {
	x, y float64
	size float64
}

func (g *Game) initStars() {
	g.stars = make([]*star, starCount)
	for i := range g.stars {
		g.stars[i] = &star{
			x:     fxRand.Float64() * ScreenWidth,
			y:     fxRand.Float64() * ScreenHeight,
			speed: 0.5 + fxRand.Float64()*1.5,
			size:  1 + fxRand.Float64(),
		}
	}
	g.initParallax()
}

func (g *Game) initParallax() {
	g.clouds = spawnLayer(cloudCount, 10, 16)
	g.hills = spawnLayer(hillCount, 40, 30)
	g.structures = spawnLayer(structureCount, 16, 10)
	g.dust = make([]*star, 0, starCount)
	for i := 0; i < starCount; i++ {
		g.dust = append(g.dust, &star{
			x:     fxRand.Float64() * ScreenWidth,
			y:     fxRand.Float64() * ScreenHeight,
			speed: dustLayerSpeed * (0.6 + fxRand.Float64()*0.8),
			size:  1 + fxRand.Float64(),
		})
	}
}

func spawnLayer(count int, minSize, extra float64) []*bgItem {
	items := make([]*bgItem, count)
	for i := range items {
		items[i] = &bgItem{
			x:    fxRand.Float64() * ScreenWidth,
			y:    fxRand.Float64() * ScreenHeight,
			size: minSize + fxRand.Float64()*extra,
		}
	}
	return items
}

func (g *Game) updateStars() {
	for _, s := range g.stars {
		s.y += s.speed
		if s.y > ScreenHeight {
			s.y = 0
			s.x = fxRand.Float64() * ScreenWidth
		}
	}
}

func (g *Game) updateParallax() {
	advanceLayer(g.clouds, cloudLayerSpeed)
	advanceLayer(g.hills, hillLayerSpeed)
	advanceLayer(g.structures, structureLayerSpeed)
	for _, s := range g.dust {
		s.y += s.speed
		if s.y > ScreenHeight {
			s.y = -s.size
			s.x = fxRand.Float64() * ScreenWidth
		}
	}
}

func advanceLayer(items []*bgItem, speed float64) {
	for _, it := range items {
		it.y += speed
		if it.y > ScreenHeight+it.size {
			it.y = -it.size
			it.x = fxRand.Float64() * ScreenWidth
		}
	}
}

func (g *Game) drawStars(screen *ebiten.Image) {
	c := color.RGBA{0x9a, 0x9a, 0xd0, 0xff}
	for _, s := range g.stars {
		vector.DrawFilledRect(screen, float32(s.x), float32(s.y), float32(s.size), float32(s.size), c, false)
	}
}

// drawParallax desenha as quatro camadas do cenário respeitando o tema do trecho:
// nuvens altas, colinas distantes, estruturas/árvores e partículas de primeiro plano.
func (g *Game) drawParallax(screen *ebiten.Image, theme bgTheme) {
	for _, c := range g.clouds {
		drawCloud(screen, c, theme.cloud)
	}
	for _, h := range g.hills {
		drawHill(screen, h, theme.hill)
	}
	for _, s := range g.structures {
		drawStructure(screen, s, theme)
	}
	for _, d := range g.dust {
		vector.DrawFilledRect(screen, float32(d.x), float32(d.y), float32(d.size), float32(d.size), withAlpha(theme.dust, 0xb0), false)
	}
}

func drawCloud(screen *ebiten.Image, c *bgItem, col color.RGBA) {
	soft := withAlpha(col, 0x90)
	x, y, s := float32(c.x), float32(c.y), float32(c.size)
	vector.DrawFilledRect(screen, x, y, s*1.6, s*0.6, soft, false)
	vector.DrawFilledRect(screen, x+s*0.3, y-s*0.3, s, s*0.6, soft, false)
}

func drawHill(screen *ebiten.Image, h *bgItem, col color.RGBA) {
	x, y, s := float32(h.x), float32(h.y), float32(h.size)
	vector.DrawFilledRect(screen, x, y, s*2, s, col, false)
	vector.DrawFilledRect(screen, x+s*0.4, y-s*0.4, s*1.2, s*0.6, col, false)
}

// drawStructure desenha a silhueta correta para o tema: árvores nos campos,
// casas em chamas na vila, muralhas e torres do castelo.
func drawStructure(screen *ebiten.Image, it *bgItem, theme bgTheme) {
	x, y, s := float32(it.x), float32(it.y), float32(it.size)
	switch theme.style {
	case styleVillageFire:
		vector.DrawFilledRect(screen, x, y, s, s, theme.structure, false)
		vector.DrawFilledRect(screen, x-1, y-s*0.4, s+2, s*0.4, theme.structure, false)
		vector.DrawFilledRect(screen, x+s*0.3, y-s*0.7, s*0.4, s*0.3, theme.accent, false)
	case styleWalls:
		vector.DrawFilledRect(screen, x, y, s*1.6, s, theme.structure, false)
		for i := 0; i < 3; i++ {
			vector.DrawFilledRect(screen, x+float32(i)*s*0.6, y-s*0.3, s*0.35, s*0.3, theme.structure, false)
		}
	case styleCastle:
		vector.DrawFilledRect(screen, x, y-s, s*0.8, s*2, theme.structure, false)
		vector.DrawFilledRect(screen, x, y-s*1.3, s*0.8, s*0.3, theme.structure, false)
		vector.DrawFilledRect(screen, x+s*0.25, y-s*0.4, s*0.3, s*0.4, theme.accent, false)
	default: // styleFields: árvores
		vector.DrawFilledRect(screen, x+s*0.35, y, s*0.3, s*0.8, theme.structure, false)
		vector.DrawFilledRect(screen, x, y-s*0.6, s, s*0.7, theme.accent, false)
		vector.DrawFilledRect(screen, x+s*0.2, y-s, s*0.6, s*0.5, theme.accent, false)
	}
}
