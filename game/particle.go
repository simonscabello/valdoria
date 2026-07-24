package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type particle struct {
	x, y    float64
	vx, vy  float64
	life    int
	maxLife int
	size    float64
	gravity bool
	col     color.RGBA
}

func (p *particle) update() {
	p.x += p.vx
	p.y += p.vy
	p.vx *= particleDrag
	p.vy *= particleDrag
	if p.gravity {
		p.vy += particleGravity
	}
	p.life--
}

func (p *particle) dead() bool { return p.life <= 0 }

func (p *particle) draw(screen *ebiten.Image) {
	alpha := uint8(255 * p.life / p.maxLife)
	size := p.size * float64(p.life) / float64(p.maxLife)
	if size < 1 {
		size = 1
	}
	vector.DrawFilledRect(screen, float32(p.x), float32(p.y), float32(size), float32(size), withAlpha(p.col, alpha), false)
}

func (g *Game) updateParticles() {
	for _, p := range g.particles {
		p.update()
	}
	g.particles = filterAlive(g.particles, func(p *particle) bool { return !p.dead() })
}

func (g *Game) drawParticles(screen *ebiten.Image) {
	for _, p := range g.particles {
		p.draw(screen)
	}
}

// spawnBurst lança partículas radiais a partir de um ponto (explosões, dano).
func (g *Game) spawnBurst(x, y float64, count int, speed float64, life int, size float64, gravity bool, col color.RGBA) {
	if len(g.particles) >= maxParticles {
		return
	}
	for i := 0; i < count; i++ {
		angle := fxRand.Float64() * 2 * math.Pi
		mag := speed * (0.4 + fxRand.Float64()*0.6)
		g.particles = append(g.particles, &particle{
			x: x, y: y,
			vx:      math.Cos(angle) * mag,
			vy:      math.Sin(angle) * mag,
			life:    life,
			maxLife: life,
			size:    size,
			gravity: gravity,
			col:     col,
		})
	}
}

func (g *Game) spawnExplosion(x, y float64, col color.RGBA) {
	g.spawnBurst(x, y, explosionParticles, explosionSpeed, explosionLife, 3, true, col)
	g.audio.playSFX(sfxEnemyDown)
}

// spawnCollectRing distribui partículas em círculo ao coletar um power-up.
func (g *Game) spawnCollectRing(x, y float64, col color.RGBA) {
	if len(g.particles) >= maxParticles {
		return
	}
	for i := 0; i < collectParticles; i++ {
		angle := float64(i) / float64(collectParticles) * 2 * math.Pi
		g.particles = append(g.particles, &particle{
			x: x, y: y,
			vx:      math.Cos(angle) * collectSpeed,
			vy:      math.Sin(angle) * collectSpeed,
			life:    collectLife,
			maxLife: collectLife,
			size:    2,
			col:     col,
		})
	}
}

// spawnTrail deixa um rastro tênue atrás de um projétil especial.
func (g *Game) spawnTrail(x, y float64, col color.RGBA) {
	if len(g.particles) >= maxParticles {
		return
	}
	g.particles = append(g.particles, &particle{
		x: x, y: y,
		life:    trailLife,
		maxLife: trailLife,
		size:    2,
		col:     col,
	})
}
