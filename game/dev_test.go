package game

import "testing"

// resetDevState devolve o estado global de dev/seed ao padrão entre testes.
func resetDevState() {
	dev = devConfig{timeScale: 1, fastTimeScale: 6}
	seedFixed = false
}

func TestConfigureEnablesDevAndSeed(t *testing.T) {
	t.Cleanup(resetDevState)
	Configure(Options{Dev: true, StartSection: 2, StartBoss: true, Seed: 42, HasSeed: true})

	if !dev.enabled || !dev.showHUD {
		t.Error("Configure deveria ativar o modo dev e o HUD de diagnóstico")
	}
	if dev.startSection != 2 || !dev.startBoss {
		t.Error("Configure deveria aplicar trecho inicial e início no chefe")
	}
	if CurrentSeed() != 42 || !seedFixed {
		t.Errorf("Configure deveria fixar a semente 42, foi %d (fixa=%v)", CurrentSeed(), seedFixed)
	}
}

func TestDevInvincibilityBlocksHits(t *testing.T) {
	t.Cleanup(resetDevState)
	p := newPlayer()

	dev.invincible = true
	if p.canBeHit() {
		t.Error("com invencibilidade dev o jogador não pode ser atingido")
	}
	dev.invincible = false
	if !p.canBeHit() {
		t.Error("sem invencibilidade dev o jogador pode ser atingido")
	}
}

func TestClearEntitiesRemovesAll(t *testing.T) {
	g := &Game{}
	g.enemies = append(g.enemies, newCrow(0))
	g.bullets = append(g.bullets, newBullet(0, 0, 0, -1, 1, 0, lightColor))
	g.enemyBullets = append(g.enemyBullets, newEnemyBullet(0, 0, 0, 1))

	g.clearEntities()

	if len(g.enemies) != 0 || len(g.bullets) != 0 || len(g.enemyBullets) != 0 {
		t.Error("clearEntities deveria remover inimigos e projéteis")
	}
}

func TestSeedReproducibility(t *testing.T) {
	t.Cleanup(resetDevState)
	draw := func() []float64 {
		out := make([]float64, 8)
		for i := range out {
			out[i] = rng.Float64()
		}
		return out
	}

	SetSeed(1234)
	first := draw()
	SetSeed(1234)
	second := draw()

	for i := range first {
		if first[i] != second[i] {
			t.Fatalf("mesma semente deveria repetir a sequência (índice %d)", i)
		}
	}
}

func TestResetRNGReproducesFixedSeed(t *testing.T) {
	t.Cleanup(resetDevState)
	draw := func() []float64 {
		out := make([]float64, 5)
		for i := range out {
			out[i] = rng.Float64()
		}
		return out
	}

	SetSeed(99)
	first := draw()
	resetRNG()
	second := draw()

	for i := range first {
		if first[i] != second[i] {
			t.Fatalf("resetRNG com semente fixa deveria repetir a sequência (índice %d)", i)
		}
	}
}

func TestRandomRuneDeterministicWithSeed(t *testing.T) {
	t.Cleanup(resetDevState)
	SetSeed(7)
	a := []powerupType{randomRune(), randomRune(), randomRune()}
	SetSeed(7)
	b := []powerupType{randomRune(), randomRune(), randomRune()}
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("randomRune deveria ser determinístico por semente (índice %d)", i)
		}
	}
}
