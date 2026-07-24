package game

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const starCount = 60

type star struct {
	x, y  float64
	speed float64
	size  float64
}

func (g *Game) initStars() {
	g.stars = make([]*star, starCount)
	for i := range g.stars {
		g.stars[i] = &star{
			x:     rand.Float64() * ScreenWidth,
			y:     rand.Float64() * ScreenHeight,
			speed: 0.5 + rand.Float64()*1.5,
			size:  1 + rand.Float64(),
		}
	}
}

func (g *Game) updateStars() {
	for _, s := range g.stars {
		s.y += s.speed
		if s.y > ScreenHeight {
			s.y = 0
			s.x = rand.Float64() * ScreenWidth
		}
	}
}

func (g *Game) drawStars(screen *ebiten.Image) {
	c := color.RGBA{0x9a, 0x9a, 0xd0, 0xff}
	for _, s := range g.stars {
		vector.DrawFilledRect(screen, float32(s.x), float32(s.y), float32(s.size), float32(s.size), c, false)
	}
}
