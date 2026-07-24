package game

import "testing"

func TestSaveRoundTrip(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	saveEnabled = true
	t.Cleanup(func() { saveEnabled = false })

	want := saveData{HighScore: 1234, SurvivalBest: 567, Difficulty: diffHard, Shake: shakeReduced, Muted: true}
	want.write()

	got := loadSave()
	if got.HighScore != 1234 || got.SurvivalBest != 567 || got.Difficulty != diffHard || got.Shake != shakeReduced || !got.Muted {
		t.Errorf("round-trip do save falhou: %+v", got)
	}
}

func TestLoadSaveDefaultsWhenMissing(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	saveEnabled = true
	t.Cleanup(func() { saveEnabled = false })

	got := loadSave()
	if got.Difficulty != diffNormal || got.HighScore != 0 || got.Shake != shakeFull {
		t.Errorf("save ausente deveria trazer padrões, foi %+v", got)
	}
}

func TestCycleDifficultyWraps(t *testing.T) {
	t.Cleanup(func() { setDifficulty(diffNormal) })
	g := New()

	setDifficulty(diffEasy)
	g.cycleDifficulty()
	if currentDifficulty != diffNormal {
		t.Errorf("após um ciclo deveria ser Normal, foi %s", difficultyName(currentDifficulty))
	}
	g.cycleDifficulty()
	g.cycleDifficulty()
	if currentDifficulty != diffEasy {
		t.Errorf("três ciclos deveriam voltar a Fácil, foi %s", difficultyName(currentDifficulty))
	}
}

func TestDifficultyChangesStartingResources(t *testing.T) {
	t.Cleanup(func() { setDifficulty(diffNormal) })
	g := New()

	setDifficulty(diffEasy)
	g.startNewGame()
	easyLives := g.lives

	setDifficulty(diffHard)
	g.startNewGame()
	if g.lives >= easyLives {
		t.Errorf("Difícil deveria começar com menos vidas que Fácil (%d vs %d)", g.lives, easyLives)
	}
}
