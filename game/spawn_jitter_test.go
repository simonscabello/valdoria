package game

import "testing"

func TestSpawnJitterEasyIsDeterministic(t *testing.T) {
	t.Cleanup(func() { setDifficulty(diffNormal) })
	setDifficulty(diffEasy)
	SetSeed(42)
	resetRNG()

	base := 100.0
	for i := 0; i < 20; i++ {
		if got := jitterSpawnX(base); got != base {
			t.Fatalf("no Fácil X deveria ser fixo, veio %v", got)
		}
		if got := jitterSpawnInterval(35); got != 35 {
			t.Fatalf("no Fácil intervalo deveria ser fixo, veio %d", got)
		}
	}
}

func TestSpawnJitterHardVariesButReproducible(t *testing.T) {
	t.Cleanup(func() { setDifficulty(diffNormal) })
	setDifficulty(diffHard)

	SetSeed(99)
	resetRNG()
	a := make([]float64, 8)
	for i := range a {
		a[i] = jitterSpawnX(100)
	}

	SetSeed(99)
	resetRNG()
	for i := range a {
		if got := jitterSpawnX(100); got != a[i] {
			t.Fatalf("com mesma seed a sequência deveria repetir (i=%d: %v vs %v)", i, got, a[i])
		}
	}

	// Em Difícil, nem todos os sorteios ficam no X scriptado.
	varied := false
	for _, x := range a {
		if x != 100 {
			varied = true
			break
		}
	}
	if !varied {
		t.Fatal("esperava variação de X no Difícil")
	}

	for _, x := range a {
		if x < 100-28 || x > 100+28 {
			t.Fatalf("X %v fora do jitter ±28", x)
		}
	}
}

func TestSpawnIntervalJitterBounds(t *testing.T) {
	t.Cleanup(func() { setDifficulty(diffNormal) })
	setDifficulty(diffHard)
	SetSeed(7)
	resetRNG()

	for i := 0; i < 40; i++ {
		got := jitterSpawnInterval(20)
		floor := (20 * 2) / 3
		if got < floor || got > 20+16 {
			t.Fatalf("intervalo %d fora dos limites [%d, %d]", got, floor, 20+16)
		}
	}
	if jitterSpawnInterval(0) != 0 {
		t.Fatal("intervalo 0 deve permanecer 0")
	}
}

func TestFairSpawnXKeepsMargin(t *testing.T) {
	if got := fairSpawnX(0, crowSize); got != float64(spawnFairMargin) {
		t.Fatalf("borda esquerda: got %v want %v", got, spawnFairMargin)
	}
	right := ScreenWidth - crowSize
	wantMax := ScreenWidth - crowSize - spawnFairMargin
	if got := fairSpawnX(float64(right), crowSize); got != float64(wantMax) {
		t.Fatalf("borda direita: got %v want %v", got, wantMax)
	}
	if got := fairSpawnX(115, crowSize); got != 115 {
		t.Fatalf("centro deveria permanecer: got %v", got)
	}
}

func TestStage1OpeningSpawnsAreReachable(t *testing.T) {
	// Os primeiros spawns da fase 1 não devem nascer colados nas bordas.
	def := stage1()
	minX := float64(spawnFairMargin)
	maxX := float64(ScreenWidth - crowSize - spawnFairMargin)
	for i, w := range def.waves {
		if w.startTick > 500 {
			break
		}
		if w.kind != kindCrow {
			continue
		}
		x := fairSpawnX(w.spawnX, crowSize)
		if x < minX || x > maxX {
			t.Errorf("onda %d (tick %d): spawnX %.0f fora da faixa justa [%.0f, %.0f]",
				i, w.startTick, x, minX, maxX)
		}
	}
}

func TestWaveStartTickUnchangedByJitter(t *testing.T) {
	// O startTick da onda continua exato; só o ritmo interno é jitterado.
	t.Cleanup(func() { setDifficulty(diffNormal) })
	setDifficulty(diffHard)
	SetSeed(1)
	resetRNG()

	l := &Level{
		sections: phase1Sections(),
		events: []*waveEvent{
			{startTick: 10, kind: kindCrow, count: 1, formation: formationSingle, spawnX: 20},
		},
	}
	if spawned := stepLevel(l, 9); len(spawned) != 0 {
		t.Fatalf("antes do tick 10: %d", len(spawned))
	}
	if spawned := stepLevel(l, 1); len(spawned) != 1 {
		t.Fatalf("no tick 10 deveria nascer 1, veio %d", len(spawned))
	}
}
