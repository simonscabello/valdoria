package game

import "testing"

// TestFullSessionReachesVictoryWithoutPanic dirige uma sessão inteira apenas
// pela lógica (sem renderização): ondas, movimento, colisões e o chefe até a
// morte. Serve como teste de fumaça contra panics e estados inválidos.
func TestFullSessionReachesVictoryWithoutPanic(t *testing.T) {
	t.Cleanup(resetDevState)
	dev.invincible = true // o jogador não deve morrer para o teste percorrer todo o fluxo

	g := New()
	g.startNewGame()

	// Trecho de jogo normal: exercita spawns, movimento, disparos inimigos,
	// colisões, partículas e popups por várias centenas de frames.
	for i := 0; i < 800 && g.state == statePlaying; i++ {
		if err := g.Update(); err != nil {
			t.Fatalf("Update retornou erro no frame %d: %v", i, err)
		}
	}

	// Força o confronto com o chefe e o conduz até a morte aplicando dano
	// sempre que ele estiver vulnerável.
	g.startBoss()
	reachedVictory := false
	for i := 0; i < 6000; i++ {
		if err := g.Update(); err != nil {
			t.Fatalf("Update (chefe) retornou erro no frame %d: %v", i, err)
		}
		if g.boss != nil && !g.boss.invulnerable {
			g.boss.takeDamage(5)
		}
		if g.state == stateVictory {
			reachedVictory = true
			break
		}
		if g.state == stateGameOver {
			t.Fatal("com invencibilidade o jogador não deveria chegar a Game Over")
		}
	}

	if !reachedVictory {
		t.Fatal("a sessão deveria alcançar a vitória após derrotar o chefe")
	}
	if len(g.enemies) != 0 || len(g.enemyBullets) != 0 {
		t.Errorf("após a vitória não deveriam restar inimigos (%d) nem projéteis (%d)",
			len(g.enemies), len(g.enemyBullets))
	}
	if g.score <= 0 {
		t.Errorf("a pontuação final deveria ser positiva, foi %d", g.score)
	}
}

// TestBossReachedNaturallyAfterWaves confirma que, esgotadas as ondas e sem
// inimigos em tela, o jogo efetivamente entra no combate com o chefe.
func TestBossReachedNaturallyAfterWaves(t *testing.T) {
	t.Cleanup(resetDevState)
	dev.invincible = true

	g := New()
	g.startNewGame()

	reachedBoss := false
	for i := 0; i < 40000; i++ {
		if err := g.Update(); err != nil {
			t.Fatalf("Update retornou erro no frame %d: %v", i, err)
		}
		// Limpa quaisquer inimigos remanescentes para não travar em uma onda
		// esperando o jogador (que aqui está imóvel) abater os retardatários.
		for _, e := range g.enemies {
			e.takeDamage(999)
		}
		if g.state == stateBoss {
			reachedBoss = true
			break
		}
	}

	if !reachedBoss {
		t.Fatal("depois de todas as ondas o chefe deveria aparecer")
	}
}
