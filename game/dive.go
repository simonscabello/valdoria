package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// O Mergulho do Grifo.
//
// Todo shoot'em up move uma nave. Aqui o jogador monta um animal, e o mergulho
// é o que torna isso perceptível: o grifo se lança para frente, atravessa a
// parede de projéteis e volta. É a única mecânica do jogo que avança **contra**
// a rolagem da tela.
//
// O desenho segue a regra do dodge roll de Enter the Gungeon: invulnerabilidade
// legível, custo visível e uma janela de vulnerabilidade logo depois. Sem esses
// três, um botão de invencibilidade apaga o desafio inteiro.

// diveState descreve em que ponto do gesto o grifo está.
type diveState int

const (
	diveIdle diveState = iota
	diveActive
	diveRecovering
)

// tryDive inicia um mergulho se houver fôlego e o grifo não estiver no meio de
// outro. Devolve true quando o mergulho começou.
func (p *Player) tryDive() bool {
	if p.diveState != diveIdle || p.stamina < diveStaminaCost {
		return false
	}
	p.diveState = diveActive
	p.diveTimer = diveDuration
	p.stamina -= diveStaminaCost
	p.diveHit = p.diveHit[:0]

	// Direção: para onde o jogador aponta; sem direção, para frente (cima).
	dx, dy := moveVector()
	if dx == 0 && dy == 0 {
		dy = -1
	}
	n := math.Hypot(dx, dy)
	p.diveVX, p.diveVY = dx/n, dy/n
	return true
}

// updateDive avança o gesto e o fôlego. Devolve true enquanto o mergulho move o
// jogador — nesse caso o movimento normal é ignorado.
func (p *Player) updateDive() bool {
	switch p.diveState {
	case diveActive:
		p.diveTimer--
		// Desacelera no fim: o mergulho termina planando, não freando no ar.
		t := float64(p.diveTimer) / diveDuration
		speed := diveSpeed * (0.35 + 0.65*t)
		p.moveBy(p.diveVX, p.diveVY, speed)
		if p.diveTimer <= 0 {
			p.diveState = diveRecovering
			p.diveTimer = diveRecovery
		}
		return true
	case diveRecovering:
		p.diveTimer--
		p.stamina += staminaRegenSlow
		if p.diveTimer <= 0 {
			p.diveState = diveIdle
		}
	default:
		p.stamina += staminaRegen
	}
	if p.stamina > staminaMax {
		p.stamina = staminaMax
	}
	return false
}

// diving indica que o grifo está atravessando: invulnerável e causando dano.
func (p *Player) diving() bool { return p.diveState == diveActive }

// canDive diz se há fôlego para mais um mergulho (usado pelo HUD).
func (p *Player) canDive() bool {
	return p.diveState == diveIdle && p.stamina >= diveStaminaCost
}

// staminaRatio devolve o fôlego de 0 a 1.
func (p *Player) staminaRatio() float64 { return p.stamina / staminaMax }

// alreadyDove evita que o mesmo mergulho fira o mesmo inimigo várias vezes.
func (p *Player) alreadyDove(e *Enemy) bool {
	for _, hit := range p.diveHit {
		if hit == e {
			return true
		}
	}
	return false
}

// handleDive lê a ação e dispara o gesto, com o feedback que o vende.
func (g *Game) handleDive() {
	if g.state != statePlaying && g.state != stateBoss {
		return
	}
	if !actionJustPressed(actDive) || !g.player.tryDive() {
		return
	}
	g.audio.playSFX(sfxDive)
	g.addShake(shakeDiveMagnitude)
	// Penas soltas no ponto de partida: o rastro do arranque.
	g.spawnBurst(g.player.centerX(), g.player.centerY(), 8, 2.4, 22, 2, true,
		color.RGBA{0xf0, 0xbc, 0x4c, 0xff})
}

// diveCollisions aplica o dano de contato do mergulho. O grifo atravessa: cada
// inimigo é ferido uma única vez por mergulho.
func (g *Game) diveCollisions() {
	if !g.player.diving() {
		return
	}
	px, py := g.player.x, g.player.y
	for _, e := range g.enemies {
		if e.dead || g.player.alreadyDove(e) {
			continue
		}
		if !collides(px, py, playerSize, playerSize, e.x, e.y, e.w, e.h) {
			continue
		}
		g.player.diveHit = append(g.player.diveHit, e)
		e.takeDamage(diveDamage)
		g.hitStop = hitStopBig
		g.spawnBurst(e.centerX(), e.centerY(), 6, 2.6, 18, 2, false, lightColor)
		if e.dead {
			points := g.registerKill(e.score)
			g.spawnDrop(e)
			g.spawnExplosion(e.centerX(), e.centerY(), e.color)
			g.spawnScorePopup(e.centerX(), e.y, points)
		}
	}
	if g.boss != nil && collides(px, py, playerSize, playerSize, g.boss.x, g.boss.y, g.boss.w, g.boss.h) {
		if !g.player.diveHitBoss {
			g.player.diveHitBoss = true
			g.boss.takeDamage(diveDamage)
			g.hitStop = hitStopBig
		}
	}
}

// diveCameraOffset empurra a câmera na direção do mergulho, para o gesto
// parecer avanço e não teletransporte.
func (g *Game) diveCameraOffset() (float64, float64) {
	p := g.player
	if p.diveState != diveActive {
		return 0, 0
	}
	t := float64(p.diveTimer) / diveDuration
	push := diveCameraPush * math.Sin(t*math.Pi)
	return -p.diveVX * push, -p.diveVY * push
}

// drawDiveTrail desenha o rastro dourado do mergulho.
func (p *Player) drawDiveTrail(dst *ebiten.Image) {
	if p.diveState != diveActive {
		return
	}
	cx, cy := p.centerX(), p.centerY()
	for i := 1; i <= 4; i++ {
		d := float64(i) * 5
		a := uint8(150 - i*30)
		vector.DrawFilledRect(dst,
			float32(cx-p.diveVX*d-3), float32(cy-p.diveVY*d-3),
			6, 6, color.RGBA{0xff, 0xe8, 0x92, a}, false)
	}
}

// drawStaminaBar mostra o fôlego do grifo no rodapé, ao lado das vidas.
// Sem essa leitura o mergulho vira tentativa e erro.
func (g *Game) drawStaminaBar(screen *ebiten.Image) {
	const w, h = 44, 3
	x, y := float32(6), float32(ScreenHeight-34)

	vector.DrawFilledRect(screen, x, y, w, h, color.RGBA{0x18, 0x14, 0x0c, 0xc0}, false)
	fill := color.RGBA{0xd0, 0xb0, 0x5a, 0xff}
	if !g.player.canDive() {
		fill = color.RGBA{0x6a, 0x5a, 0x3a, 0xff} // sem fôlego para mergulhar
	}
	vector.DrawFilledRect(screen, x, y, w*float32(g.player.staminaRatio()), h, fill, false)

	// Marca do custo de um mergulho: o jogador vê quando pode usar de novo.
	tick := float32(w) * diveStaminaCost / staminaMax
	vector.DrawFilledRect(screen, x+tick, y-1, 1, h+2, withAlpha(uiChipEdge, 200), false)
}
