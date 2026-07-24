package game

import "testing"

func TestLoseLifeRespawnsThenGameOver(t *testing.T) {
	g := New()
	g.lives = 2

	g.player.health = 0
	g.loseLife()
	if g.lives != 1 {
		t.Fatalf("vidas = %d, quero 1", g.lives)
	}
	if g.state != statePlaying {
		t.Error("deveria continuar jogando após reaparecer")
	}
	if g.player.health != maxHealth {
		t.Errorf("vida deveria ser restaurada, foi %d", g.player.health)
	}

	g.player.health = 0
	g.loseLife()
	if g.state != stateGameOver {
		t.Error("sem vidas restantes deveria ir para Game Over")
	}
}

func TestRespawnClearsBulletsAndReducesWeapon(t *testing.T) {
	g := New()
	g.player.weaponLevel = 3
	g.enemyBullets = append(g.enemyBullets, newEnemyBullet(50, 50, 0, 1))

	g.player.health = 0
	g.loseLife()

	if len(g.enemyBullets) != 0 {
		t.Errorf("projéteis inimigos deveriam ser limpos, restaram %d", len(g.enemyBullets))
	}
	if g.player.invincible <= 0 {
		t.Error("deveria conceder invencibilidade ao reaparecer")
	}
	if g.player.weaponLevel != 2 {
		t.Errorf("nível da arma deveria cair para 2, foi %d", g.player.weaponLevel)
	}
	if g.player.x != ScreenWidth/2-playerSize/2 {
		t.Error("jogador deveria voltar à posição segura")
	}
}

func TestRespawnKeepsWeaponAtLevelOne(t *testing.T) {
	p := newPlayer()
	p.weapon = weaponIce
	p.weaponLevel = 1
	p.respawn()
	if p.weapon != weaponIce {
		t.Error("respawn não deveria remover a arma")
	}
	if p.weaponLevel != 1 {
		t.Errorf("nível não deveria ficar abaixo de 1, foi %d", p.weaponLevel)
	}
}

func TestUseBombClearsBulletsAndDamagesEnemies(t *testing.T) {
	g := New()
	g.state = statePlaying
	before := g.bombCharges
	g.enemyBullets = append(g.enemyBullets, newEnemyBullet(10, 10, 0, 1))
	for i := 0; i < 3; i++ {
		g.enemies = append(g.enemies, newWyvern(float64(i*20)))
	}

	g.useBomb()

	if g.bombCharges != before-1 {
		t.Errorf("cargas = %d, quero %d", g.bombCharges, before-1)
	}
	if len(g.enemyBullets) != 0 {
		t.Error("a bomba deveria limpar os projéteis inimigos")
	}
	for i, e := range g.enemies {
		if !e.dead {
			t.Errorf("inimigo %d deveria ser destruído pela bomba", i)
		}
	}
	if g.player.invincible <= 0 {
		t.Error("a bomba deveria conceder invencibilidade")
	}
}

func TestMultiplierGrowsAndCaps(t *testing.T) {
	g := &Game{}
	if g.multiplier() != 1 {
		t.Errorf("combo 0 deveria dar x1, deu x%d", g.multiplier())
	}
	g.combo = comboStep
	if g.multiplier() != 2 {
		t.Errorf("combo %d deveria dar x2, deu x%d", comboStep, g.multiplier())
	}
	g.combo = comboStep * 100
	if g.multiplier() != maxMultiplier {
		t.Errorf("combo alto deveria limitar em x%d, deu x%d", maxMultiplier, g.multiplier())
	}
}

func TestRegisterKillAppliesMultiplier(t *testing.T) {
	g := &Game{}
	g.registerKill(10)
	if g.score != 10 {
		t.Errorf("primeira eliminação deveria valer 10, deu %d", g.score)
	}
	if g.comboTimer != comboWindow {
		t.Error("eliminar deveria reiniciar a janela do combo")
	}
	g.combo = 2*comboStep - 1
	g.registerKill(10) // combo vira 2*comboStep -> multiplicador 3
	if g.score != 10+30 {
		t.Errorf("score = %d, quero 40", g.score)
	}
}

func TestPlayerDamageResetsCombo(t *testing.T) {
	g := &Game{combo: 10, comboTimer: 50}
	g.onPlayerDamaged()
	if g.combo != 0 || g.comboTimer != 0 {
		t.Error("dano deveria zerar o combo")
	}
	if !g.sectionDamaged {
		t.Error("dano deveria marcar o trecho como danificado")
	}
}

func TestWaveNoDamageBonus(t *testing.T) {
	g := New()
	g.lastSection = 0
	g.sectionDamaged = false
	g.level.section = 1

	before := g.score
	g.checkSectionBonus()
	if g.score != before+waveNoDamageBonus {
		t.Errorf("concluir trecho sem dano deveria bonificar, score %d", g.score-before)
	}

	g.sectionDamaged = true
	g.level.section = 2
	before = g.score
	g.checkSectionBonus()
	if g.score != before {
		t.Error("trecho com dano não deveria bonificar")
	}
}

func TestFormationBonusOnlyWhenFullyCleared(t *testing.T) {
	g := &Game{formations: map[int]*formationTracker{}}
	g.formations[1] = &formationTracker{total: 2}
	for i := 0; i < 2; i++ {
		e := newCrow(0)
		e.formationID = 1
		e.dead = true
		g.enemies = append(g.enemies, e)
	}

	g.removeDead()
	if g.score != formationBonus {
		t.Errorf("formação destruída deveria dar bônus, score %d", g.score)
	}

	g2 := &Game{formations: map[int]*formationTracker{}}
	g2.formations[2] = &formationTracker{total: 2}
	killed := newCrow(0)
	killed.formationID = 2
	killed.dead = true
	escaped := newCrow(0)
	escaped.formationID = 2
	escaped.dead = true
	escaped.escaped = true
	g2.enemies = append(g2.enemies, killed, escaped)

	g2.removeDead()
	if g2.score != 0 {
		t.Errorf("formação com fuga não deveria bonificar, score %d", g2.score)
	}
}
