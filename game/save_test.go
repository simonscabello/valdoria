package game

import "testing"

func TestSaveRoundTrip(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	saveEnabled = true
	t.Cleanup(func() { saveEnabled = false })

	want := saveData{
		HighScore: 1234, SurvivalBest: 567, Difficulty: diffHard,
		Shake: shakeReduced, Muted: true,
		MusicVolume: 0.4, SFXVolume: 0.9, VolumesSet: true,
	}
	want.write()

	got := loadSave()
	if got.HighScore != 1234 || got.SurvivalBest != 567 || got.Difficulty != diffHard || got.Shake != shakeReduced || !got.Muted {
		t.Errorf("round-trip do save falhou: %+v", got)
	}
	if got.MusicVolume != 0.4 || got.SFXVolume != 0.9 || !got.VolumesSet {
		t.Errorf("volumes do save falharam: %+v", got)
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
	if got.MusicVolume != defaultMusicVolume || got.SFXVolume != defaultSFXVolume {
		t.Errorf("volumes padrão esperados %.1f/%.1f, foi %.1f/%.1f",
			defaultMusicVolume, defaultSFXVolume, got.MusicVolume, got.SFXVolume)
	}
}

func TestOldSaveGetsDefaultVolumes(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	saveEnabled = true
	t.Cleanup(func() { saveEnabled = false })

	// Simula save antigo sem volumes_set.
	legacy := saveData{HighScore: 10, Difficulty: diffNormal, Shake: shakeFull}
	legacy.write()

	got := loadSave()
	if got.MusicVolume != defaultMusicVolume || got.SFXVolume != defaultSFXVolume {
		t.Errorf("save antigo deveria receber volumes padrão, foi música=%v sfx=%v", got.MusicVolume, got.SFXVolume)
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
