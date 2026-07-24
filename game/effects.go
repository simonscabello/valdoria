package game

// addShake reforça a vibração atual, respeitando a preferência do jogador.
func (g *Game) addShake(mag float64) {
	if g.shakeSetting == shakeOff {
		return
	}
	if mag > g.shakeMag {
		g.shakeMag = mag
	}
}

func (g *Game) updateShake() {
	g.shakeMag *= shakeDecay
	if g.shakeMag < shakeMinMagnitude {
		g.shakeMag = 0
	}
}

// shakeOffset devolve o deslocamento aleatório aplicado à cena neste frame.
func (g *Game) shakeOffset() (float64, float64) {
	m := g.shakeMag * g.shakeSetting.scale()
	if m <= 0 {
		return 0, 0
	}
	return clampShake((fxRand.Float64()*2 - 1) * m), clampShake((fxRand.Float64()*2 - 1) * m)
}

func clampShake(v float64) float64 {
	if v > shakeMaxOffset {
		return shakeMaxOffset
	}
	if v < -shakeMaxOffset {
		return -shakeMaxOffset
	}
	return v
}

func (g *Game) triggerDamageFlash() {
	g.damageFlash = damageFlashDuration
}

// cycleShake alterna entre vibração cheia, reduzida e desligada.
func (g *Game) cycleShake() {
	g.shakeSetting = (g.shakeSetting + 1) % 3
	if g.shakeSetting == shakeOff {
		g.shakeMag = 0
	}
}
