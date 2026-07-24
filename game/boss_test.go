package game

import "testing"

func combatBoss() *Boss {
	b := newBoss()
	b.phase = bossPhase1
	b.invulnerable = false
	return b
}

func TestBossPhaseChangesByHealth(t *testing.T) {
	b := combatBoss()

	b.health = int(0.8 * bossMaxHealth)
	b.refreshPhase()
	if b.phase != bossPhase1 {
		t.Errorf("acima de 65%% deveria ser fase 1, foi %d", b.phase)
	}

	b.health = int(0.5 * bossMaxHealth)
	b.refreshPhase()
	if b.phase != bossPhase2 {
		t.Errorf("entre 65%% e 30%% deveria ser fase 2, foi %d", b.phase)
	}

	b.health = int(0.2 * bossMaxHealth)
	b.refreshPhase()
	if b.phase != bossPhase3 {
		t.Errorf("abaixo de 30%% deveria ser fase 3, foi %d", b.phase)
	}
	if !b.crystalsActive {
		t.Error("a fase 3 deveria ativar os cristais")
	}
}

func TestBossPatternSelectionCycles(t *testing.T) {
	b := combatBoss()
	pats := bossPatterns[bossPhase1]

	if b.currentPattern().kind != pats[0].kind {
		t.Error("primeiro padrão incorreto")
	}
	b.nextPattern()
	if b.currentPattern().kind != pats[1].kind {
		t.Error("segundo padrão incorreto")
	}
	b.nextPattern()
	if b.currentPattern().kind != pats[0].kind {
		t.Error("seleção deveria voltar ao primeiro padrão")
	}
}

func TestBossInvulnerableDuringEntry(t *testing.T) {
	b := newBoss()
	if !b.invulnerable {
		t.Fatal("o chefe deveria entrar invulnerável")
	}
	before := b.health
	b.takeDamage(50)
	if b.health != before {
		t.Errorf("dano não deveria contar na entrada, vida %d -> %d", before, b.health)
	}
}

func TestBossStopsAttackingAfterDeath(t *testing.T) {
	b := combatBoss()
	b.action = actionFiring
	b.actionTimer = 30
	b.health = 3

	b.takeDamage(5)
	if b.phase != bossDying {
		t.Fatalf("vida zerada deveria iniciar a morte, fase %d", b.phase)
	}
	if !b.invulnerable {
		t.Error("o chefe deveria ficar invulnerável ao morrer")
	}

	var bullets []*EnemyBullet
	var enemies []*Enemy
	ctx := &bossContext{bullets: &bullets, enemies: &enemies}
	for i := 0; i < bossDeathDuration+5; i++ {
		b.update(ctx)
	}
	if len(bullets) != 0 {
		t.Errorf("o chefe morto não deveria disparar, veio %d", len(bullets))
	}
	if b.phase != bossDead {
		t.Errorf("após a sequência de morte deveria estar morto, fase %d", b.phase)
	}
}

func TestBossDefeatTransitionsToVictory(t *testing.T) {
	g := New()
	g.startBoss()
	g.boss.phase = bossDead

	before := g.score
	g.checkBossDefeat()

	if g.state != stateVictory {
		t.Errorf("derrotar o chefe deveria levar à vitória, estado %d", g.state)
	}
	expected := bossScore + g.lives*lifeBonus + g.bombCharges*bombBonusPoints
	if g.score != before+expected {
		t.Errorf("vitória deveria conceder %d pontos, deu %d", expected, g.score-before)
	}
	if g.victoryLifeBonus != g.lives*lifeBonus {
		t.Errorf("bônus de vidas incorreto: %d", g.victoryLifeBonus)
	}
}

func TestBossFiresWhileAttacking(t *testing.T) {
	b := combatBoss()
	b.action = actionFiring
	b.actionTimer = 60
	b.patternIndex = 0 // patternAimedFire na fase 1

	var bullets []*EnemyBullet
	var enemies []*Enemy
	ctx := &bossContext{playerX: 100, playerY: 300, bullets: &bullets, enemies: &enemies}
	for i := 0; i < 40; i++ {
		b.update(ctx)
	}
	if len(bullets) == 0 {
		t.Error("o chefe atacando deveria gerar projéteis")
	}
}
