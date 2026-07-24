package game

import "testing"

func stepLevel(l *Level, ticks int) []*Enemy {
	var all []*Enemy
	for i := 0; i < ticks; i++ {
		all = append(all, l.update()...)
	}
	return all
}

func TestEventActivatesAtCorrectTick(t *testing.T) {
	l := &Level{
		sections: phase1Sections(),
		events: []*waveEvent{
			{startTick: 10, kind: kindCrow, count: 1, formation: formationSingle, spawnX: 20},
		},
	}

	if spawned := stepLevel(l, 9); len(spawned) != 0 {
		t.Fatalf("nenhum inimigo deveria surgir antes do tick 10, veio %d", len(spawned))
	}
	if spawned := stepLevel(l, 1); len(spawned) != 1 {
		t.Fatalf("um inimigo deveria surgir no tick 10, veio %d", len(spawned))
	}
}

func TestWaveCompletes(t *testing.T) {
	l := &Level{
		sections: phase1Sections(),
		events: []*waveEvent{
			{startTick: 1, kind: kindCrow, count: 3, interval: 5, formation: formationSingle},
		},
	}

	spawned := stepLevel(l, 40)
	if len(spawned) != 3 {
		t.Errorf("deveriam surgir 3 inimigos, vieram %d", len(spawned))
	}
	if !l.events[0].done {
		t.Error("o evento deveria estar concluído")
	}
	if !l.allEventsDone() {
		t.Error("todos os eventos deveriam estar concluídos")
	}
}

func TestSectionChanges(t *testing.T) {
	l := &Level{sections: phase1Sections()}
	l.tick = l.sections[1].startTick - 1

	l.update()

	if l.section != 1 {
		t.Errorf("section = %d, quero 1", l.section)
	}
	if l.announceTimer <= 0 {
		t.Error("a troca de trecho deveria exibir um aviso")
	}
	if l.announce != l.sections[1].warning {
		t.Errorf("aviso = %q, quero %q", l.announce, l.sections[1].warning)
	}
}

func TestPauseHaltsTimeline(t *testing.T) {
	g := New()

	g.state = statePlaying
	before := g.level.tick
	g.Update()
	if g.level.tick != before+g.timeScale {
		t.Fatalf("jogando: tick deveria avançar %d, foi %d", g.timeScale, g.level.tick-before)
	}

	g.state = statePaused
	paused := g.level.tick
	g.Update()
	if g.level.tick != paused {
		t.Errorf("pausado: tick não deveria avançar, avançou %d", g.level.tick-paused)
	}
}

func TestReadyForBoss(t *testing.T) {
	l := &Level{
		sections: phase1Sections(),
		events: []*waveEvent{
			{startTick: 1, kind: kindCrow, count: 1, formation: formationSingle},
		},
	}

	if l.readyForBoss(0) {
		t.Error("não deveria preparar o chefe com evento pendente")
	}

	stepLevel(l, 5)
	if !l.readyForBoss(0) {
		t.Error("deveria preparar o chefe: eventos concluídos e sem inimigos")
	}
	if l.readyForBoss(2) {
		t.Error("não deveria preparar o chefe com inimigos em tela")
	}
}
