package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Graze — raspar projéteis carrega a Invocação Ancestral.
//
// O problema que isto resolve: a bomba era o gesto mais espetacular do jogo e o
// jogador não tinha como recuperá-la. Duas cargas por partida, e a economia
// antiga ainda pagava pontos por *não* usá-las. O resultado é que a mecânica
// mais bonita do jogo quase nunca acontecia.
//
// A solução vem do Vaunt de Jamestown e do graze dos bullet-hells: chegar perto
// do perigo é recompensado. Isso inverte o instinto do jogador — em vez de
// fugir para o canto, ele passa a *querer* atravessar a chuva de tiros, que é
// exatamente a "tensão heroica" que o jogo quer produzir.
//
// A hitbox do jogador tem 4px e o raio de graze tem 13: existe uma faixa larga
// e legível entre "quase acertou" e "acertou".

// updateGraze conta os projéteis raspados neste frame e converte em carga.
func (g *Game) updateGraze() {
	if g.state != statePlaying && g.state != stateBoss {
		return
	}
	// Durante a invulnerabilidade não há mérito em chegar perto.
	if !g.player.canBeHit() {
		return
	}

	cx, cy := g.player.centerX(), g.player.centerY()
	grazed := 0
	for _, b := range g.enemyBullets {
		if b.dead || b.grazed > 0 {
			continue
		}
		if math.Hypot(b.centerX()-cx, b.centerY()-cy) > grazeRadius {
			continue
		}
		b.grazed = grazeCooldown
		grazed++
		g.spawnBurst(b.centerX(), b.centerY(), 2, 1.1, 12, 1, false, grazeColor)
	}
	if grazed == 0 {
		return
	}

	g.grazeCount += grazed
	g.audio.playSFX(sfxGraze)
	g.grazeFlash = 8
	g.addGrazeCharge(float64(grazed) * grazePerBullet)
}

// addGrazeCharge acumula carga e converte em bomba ao encher.
func (g *Game) addGrazeCharge(v float64) {
	if g.bombCharges >= bombMaxCharges {
		g.grazeCharge = 0 // com o cinto cheio, raspar deixa de acumular
		return
	}
	g.grazeCharge += v
	for g.grazeCharge >= grazeFull && g.bombCharges < bombMaxCharges {
		g.grazeCharge -= grazeFull
		g.bombCharges++
		g.onBombEarned()
	}
}

// onBombEarned celebra a conquista de uma carga. Ganhar uma bomba raspando
// projéteis é o melhor momento que o sistema pode produzir — precisa soar como
// recompensa, não como um número mudando no HUD.
func (g *Game) onBombEarned() {
	g.audio.playSFX(sfxPickup)
	g.addShake(shakeHitMagnitude)
	g.spawnCollectRing(g.player.centerX(), g.player.centerY(), grazeColor)
	g.spawnTextPopup(g.player.centerX(), g.player.y-6, "INVOCACAO!")
}

var grazeColor = color.RGBA{0x9a, 0xe8, 0xff, 0xff}

// drawGrazeAura mostra o raio de graze enquanto o jogador está raspando algo.
// Só aparece no momento em que importa: uma aura permanente viraria ruído.
func (g *Game) drawGrazeAura(dst *ebiten.Image) {
	if g.grazeFlash <= 0 {
		return
	}
	alpha := uint8(70 * g.grazeFlash / 8)
	vector.StrokeCircle(dst,
		float32(g.player.centerX()), float32(g.player.centerY()),
		float32(grazeRadius), 1, withAlpha(grazeColor, alpha), false)
}

// drawGrazeBar desenha o progresso rumo à próxima Invocação Ancestral, junto do
// contador de bombas. Sem essa leitura o jogador não descobre a mecânica.
func (g *Game) drawGrazeBar(screen *ebiten.Image) {
	const w, h = 44, 3
	x, y := float32(6), float32(ScreenHeight-30)

	vector.DrawFilledRect(screen, x, y, w, h, color.RGBA{0x0c, 0x14, 0x1c, 0xc0}, false)
	if g.bombCharges >= bombMaxCharges {
		// Cinto cheio: a barra fica sólida, sinalizando que não há o que ganhar.
		vector.DrawFilledRect(screen, x, y, w, h, withAlpha(grazeColor, 120), false)
		return
	}
	ratio := g.grazeCharge / grazeFull
	vector.DrawFilledRect(screen, x, y, w*float32(ratio), h, grazeColor, false)
}
