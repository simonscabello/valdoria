package game

import "testing"

func TestSurvivalSpawnsEnemiesAndNeverBosses(t *testing.T) {
	t.Cleanup(resetDevState)
	dev.invincible = true // o jogador não deve morrer durante a amostragem

	g := New()
	g.mode = modeSurvival
	g.startNewGame()

	for i := 0; i < 400; i++ {
		g.updatePlay()
	}

	if len(g.enemies) == 0 {
		t.Error("o modo sobrevivência deveria gerar inimigos ao longo do tempo")
	}
	if g.state != statePlaying {
		t.Errorf("sobrevivência não deveria entrar em chefe/vitória, estado %d", g.state)
	}
}

func TestSurvivalTracksOwnBestScore(t *testing.T) {
	g := New()
	g.mode = modeSurvival
	g.startNewGame()

	g.highScore = 5000 // recorde da campanha não deve ser afetado
	g.score = 800
	g.updateBestScore()

	if g.survivalBest != 800 {
		t.Errorf("recorde de sobrevivência deveria ser 800, foi %d", g.survivalBest)
	}
	if g.highScore != 5000 {
		t.Errorf("recorde da campanha não deveria mudar na sobrevivência, foi %d", g.highScore)
	}
}

func TestSurvivalModeResetsWithSession(t *testing.T) {
	g := New()
	g.mode = modeSurvival
	g.startNewGame()
	if g.mode != modeSurvival {
		t.Fatal("startNewGame deveria preservar o modo sobrevivência")
	}
	if g.boss != nil {
		t.Error("sobrevivência não deveria iniciar com chefe")
	}
}
