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

// TestStagesAreDense guarda contra regressão de densidade — o problema relatado
// em playtest foi "poucos monstros / muita espera".
//
// A métrica é a **ocupação da tela**, não a contagem bruta de spawns: inimigos
// pesados (wyvern, feiticeiro, balista) permanecem muito mais tempo em cena que
// um corvo, então trocar quantidade por variedade reduz os spawns e ainda assim
// aumenta a densidade percebida. Medir spawns puniria exatamente a correção que
// tornou a campanha mais variada.
func TestStagesAreDense(t *testing.T) {
	const minPeak = 10
	for _, def := range campaignStages() {
		m := measureStage(def)
		if m.PeakOnScreen < minPeak {
			t.Errorf("fase %q: pico de %d inimigos em cena (mínimo %d); densidade insuficiente",
				def.name, m.PeakOnScreen, minPeak)
		}
		if m.Enemies < 50 {
			t.Errorf("fase %q gera só %d inimigos; fluxo insuficiente", def.name, m.Enemies)
		}
	}
}

// TestBestiaryIsBalanced guarda a composição da campanha: nenhum inimigo pode
// dominar e nenhum pode ser decorativo. Antes desta revisão, 94,6% da campanha
// eram corvos e harpias, e os outros quatro somavam 25 aparições.
func TestBestiaryIsBalanced(t *testing.T) {
	r := MeasureBalance(diffNormal, 42)
	for _, c := range Checks(r) {
		if !c.Pass {
			t.Errorf("critério de balanceamento não atendido: %s — %s", c.Name, c.Detail)
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
