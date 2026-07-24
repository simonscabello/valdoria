package game

import "testing"

func TestCampaignHasBossOnlyOnLastStage(t *testing.T) {
	stages := campaignStages()
	if len(stages) < 2 {
		t.Fatalf("a campanha deveria ter várias fases, tem %d", len(stages))
	}
	for i, s := range stages {
		last := i == len(stages)-1
		if s.hasBoss != last {
			t.Errorf("fase %d: hasBoss=%v, esperado %v (só a última tem chefe)", i, s.hasBoss, last)
		}
	}
}

func TestStageDurationCoversLastWave(t *testing.T) {
	for _, s := range campaignStages() {
		last := 0
		for _, w := range s.waves {
			if w.startTick > last {
				last = w.startTick
			}
		}
		if s.duration() <= last {
			t.Errorf("fase %q: duração %d deveria passar da última onda %d", s.name, s.duration(), last)
		}
	}
}

// stageSpawnCount roda a linha do tempo completa de uma fase e conta quantos
// inimigos ela gera.
func stageSpawnCount(def *stageDef) int {
	l := newLevelFromStage(def)
	total := 0
	for i := 0; i <= l.duration; i++ {
		total += len(l.update())
	}
	return total
}

func TestStagesAreDense(t *testing.T) {
	// Guarda contra regressão de densidade: cada fase deve manter um fluxo
	// contínuo de inimigos (o problema relatado era "poucos monstros").
	for _, def := range campaignStages() {
		got := stageSpawnCount(def)
		if got < 100 {
			t.Errorf("fase %q gera poucos inimigos (%d); densidade insuficiente", def.name, got)
		}
	}
}

func TestDensityScaleIncreasesCount(t *testing.T) {
	if scaleWaveCount(5) <= 5 {
		t.Errorf("count 5 deveria crescer com densidade, foi %d", scaleWaveCount(5))
	}
	if scaleWaveCount(1) != 1 {
		t.Errorf("count 1 (formação única) não deveria mudar, foi %d", scaleWaveCount(1))
	}
	if scaleWaveInterval(40) >= 40 {
		t.Errorf("intervalo 40 deveria encolher, foi %d", scaleWaveInterval(40))
	}
}

func TestSectionsHaveBiomeMusic(t *testing.T) {
	for _, def := range campaignStages() {
		for _, sec := range def.sections {
			if sec.music == musicNone || sec.music == musicMenu || sec.music == musicBoss {
				t.Errorf("seção %q da fase %q sem trilha de bioma (music=%d)", sec.name, def.name, sec.music)
			}
		}
	}
}

func TestNewLevelFromStageBuildsEvents(t *testing.T) {
	def := stage1()
	l := newLevelFromStage(def)
	if len(l.events) != len(def.waves) {
		t.Errorf("eventos = %d, esperado %d", len(l.events), len(def.waves))
	}
	if l.duration != def.duration() {
		t.Errorf("duração = %d, esperado %d", l.duration, def.duration())
	}
}

func TestAdvanceStageMovesToNextStage(t *testing.T) {
	g := New()
	g.startNewGame()
	if g.stageIndex != 0 {
		t.Fatalf("a campanha deveria começar na fase 0, foi %d", g.stageIndex)
	}
	g.score = 1000
	g.player.weaponLevel = 3

	g.advanceStage()

	if g.stageIndex != 1 {
		t.Errorf("advanceStage deveria ir para a fase 1, foi %d", g.stageIndex)
	}
	if g.state != statePlaying {
		t.Error("avançar de fase deveria continuar jogando")
	}
	if g.score != 1000 || g.player.weaponLevel != 3 {
		t.Error("avançar de fase deveria preservar pontuação e poder")
	}
	if g.level.announceTimer <= 0 {
		t.Error("avançar de fase deveria anunciar a nova região")
	}
}

func TestLastStageStartsBoss(t *testing.T) {
	g := New()
	g.startNewGame()
	g.stageIndex = len(g.stages) - 1
	g.level = newLevelFromStage(g.stages[g.stageIndex])
	if !g.level.hasBoss {
		t.Fatal("a última fase deveria ter chefe")
	}

	// Sem eventos pendentes e sem inimigos, deve preparar o chefe.
	for _, ev := range g.level.events {
		ev.done = true
	}
	g.enemies = g.enemies[:0]
	g.updatePlay()

	if g.state != stateBoss {
		t.Errorf("concluir a última fase deveria iniciar o chefe, estado %d", g.state)
	}
}
