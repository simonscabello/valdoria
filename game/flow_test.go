package game

import "testing"

func TestNewStartsInMenu(t *testing.T) {
	g := New()
	if g.state != stateMenu {
		t.Errorf("o jogo deveria iniciar no menu, estado %d", g.state)
	}
}

func TestStartNewGameCreatesFreshSession(t *testing.T) {
	g := New()
	g.startNewGame()
	g.score = 500
	g.enemies = append(g.enemies, newCrow(50))
	g.enemyBullets = append(g.enemyBullets, newEnemyBullet(10, 10, 0, 1))
	g.powerups = append(g.powerups, newPowerup(powerHeal, 10, 10))
	g.bullets = append(g.bullets, newBullet(0, 0, 0, -1, 1, 0, lightColor))
	g.startBoss()

	g.startNewGame()

	if g.state != statePlaying {
		t.Errorf("nova sessão deveria estar jogando, estado %d", g.state)
	}
	if g.score != 0 {
		t.Errorf("pontuação deveria zerar, foi %d", g.score)
	}
	if len(g.enemies) != 0 || len(g.enemyBullets) != 0 || len(g.powerups) != 0 || len(g.bullets) != 0 {
		t.Error("nenhuma entidade anterior deveria permanecer")
	}
	if g.boss != nil {
		t.Error("o chefe anterior não deveria permanecer")
	}
	if g.lives != startingLives {
		t.Errorf("vidas deveriam reiniciar em %d, foi %d", startingLives, g.lives)
	}
}

func TestPauseIgnoredOutsideGameplay(t *testing.T) {
	g := New()
	g.handlePauseToggle()
	if g.state != stateMenu {
		t.Error("pausa não deveria funcionar no menu")
	}
}

func TestGameOverDoesNotAdvancePhase(t *testing.T) {
	g := New()
	g.startNewGame()
	g.state = stateGameOver
	before := g.level.tick

	g.Update()

	if g.level.tick != before {
		t.Error("Game Over não deveria avançar a linha do tempo")
	}
}

func TestMenuNavigationWraps(t *testing.T) {
	g := New()
	g.moveMenuSelection(menuItemCount)
	if g.menuIndex != 0 {
		t.Errorf("índice inicial deveria permanecer 0 sem entrada, foi %d", g.menuIndex)
	}
}

func TestCycleShakeWraps(t *testing.T) {
	g := New()
	if g.shakeSetting != shakeFull {
		t.Fatalf("vibração deveria iniciar cheia, foi %d", g.shakeSetting)
	}
	g.cycleShake()
	if g.shakeSetting != shakeReduced {
		t.Errorf("após um ciclo deveria ser reduzida, foi %d", g.shakeSetting)
	}
	g.cycleShake()
	g.cycleShake()
	if g.shakeSetting != shakeFull {
		t.Errorf("três ciclos deveriam voltar à cheia, foi %d", g.shakeSetting)
	}
}

func TestShakeOffZeroesOffset(t *testing.T) {
	g := New()
	g.shakeSetting = shakeOff
	g.addShake(shakeBombMagnitude)
	if g.shakeMag != 0 {
		t.Error("com vibração desligada não deveria acumular magnitude")
	}
	ox, oy := g.shakeOffset()
	if ox != 0 || oy != 0 {
		t.Errorf("offset deveria ser zero, foi (%v, %v)", ox, oy)
	}
}
