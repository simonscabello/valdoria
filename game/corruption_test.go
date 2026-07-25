package game

import "testing"

// Só a travessia pela base corrompe o reino. A gárgula entra por um lado e sai
// pelo outro por design — punir isso seria punir o jogador por uma regra do
// próprio inimigo.
func TestOnlyBottomEscapesCorrupt(t *testing.T) {
	var c Corruption

	c.add(kindGargoyle)
	if c.value != 0 || c.escaped != 0 {
		t.Fatalf("saída lateral não deveria corromper: valor %v, fugas %d", c.value, c.escaped)
	}

	c.add(kindCrow)
	if c.value != corruptionCrow || c.escaped != 1 {
		t.Fatalf("fuga de corvo deveria corromper %v, foi %v", corruptionCrow, c.value)
	}
}

// Ameaças pesadas custam mais caro ao escapar que um corvo.
func TestCorruptionWeightsAreOrdered(t *testing.T) {
	order := []enemyKind{kindCrow, kindHarpy, kindMage, kindWyvern}
	for i := 1; i < len(order); i++ {
		prev, cur := corruptionWeight(order[i-1]), corruptionWeight(order[i])
		if cur < prev {
			t.Errorf("%s (%v) deveria corromper ao menos tanto quanto %s (%v)",
				EnemyKindName(order[i]), cur, EnemyKindName(order[i-1]), prev)
		}
	}
	if corruptionWeight(kindGargoyle) != 0 {
		t.Error("a gárgula não corrompe: ela sai pela lateral por design")
	}
}

func TestCorruptionTiersAdvanceAndCap(t *testing.T) {
	var c Corruption
	if c.tier != tierSteady {
		t.Fatal("uma partida começa com o reino firme")
	}
	for i := 0; i < 500; i++ {
		c.add(kindWyvern)
	}
	if c.value != maxCorruption {
		t.Errorf("a corrupção deveria saturar em %v, foi %v", maxCorruption, c.value)
	}
	if c.tier != tierFall || !c.fallen() {
		t.Errorf("corrupção máxima deveria ser a Queda de Valdoria, foi %q", c.tierName())
	}
}

// A troca central do jogo: o mundo fica pior **e** paga melhor. Se qualquer um
// dos dois lados não subir, a corrupção deixa de ser uma aposta.
func TestCorruptionRaisesRiskAndReward(t *testing.T) {
	for tier := tierShadow; tier < corruptionTierCount; tier++ {
		prev, cur := corruptionTable[tier-1], corruptionTable[tier]
		if cur.scoreMul <= prev.scoreMul {
			t.Errorf("faixa %q deveria pagar mais que %q (%v vs %v)",
				cur.name, prev.name, cur.scoreMul, prev.scoreMul)
		}
		if cur.enemyHealthMul < prev.enemyHealthMul {
			t.Errorf("faixa %q deveria ser ao menos tão perigosa quanto %q", cur.name, prev.name)
		}
	}
}

// Subir de faixa é um evento: precisa ser anunciado ao jogador.
func TestTierChangeIsAnnouncedOnce(t *testing.T) {
	var c Corruption
	changed := false
	for i := 0; i < 200 && !changed; i++ {
		changed = c.add(kindCrow)
	}
	if !changed {
		t.Fatal("acumular corrupção deveria acabar mudando de faixa")
	}
	if c.announceTimer <= 0 || c.announce == "" {
		t.Error("a mudança de faixa deveria anunciar a nova condição do reino")
	}
	// Um acréscimo dentro da mesma faixa não anuncia de novo.
	if c.add(kindCrow) {
		t.Error("acréscimo dentro da mesma faixa não deveria reanunciar")
	}
}

// A corrupção só sobe no instante do spawn: inimigos já em cena não podem ficar
// mais fortes debaixo do jogador.
func TestSpawnAppliesCorruptionState(t *testing.T) {
	g := simGame()
	g.corruption.value = 90 // Colapso
	g.corruption.tier = tierFor(90)

	e := newCrow(100)
	base := e.health
	g.spawn(e)

	if !e.mutated {
		t.Error("em corrupção alta os inimigos deveriam nascer corrompidos")
	}
	if e.health <= base {
		t.Errorf("o inimigo corrompido deveria ser mais resistente: %d -> %d", base, e.health)
	}
	if len(g.enemies) != 1 {
		t.Fatalf("o inimigo deveria entrar em jogo, tem %d", len(g.enemies))
	}
}

// O corvo corrompido deixa de ser inofensivo — é a consequência mais sentida.
func TestShadowCrowFires(t *testing.T) {
	e := newCrow(100)
	mutateEnemy(e)

	var shots []*EnemyBullet
	ctx := &enemyContext{playerX: 100, playerY: 300, bullets: &shots}
	for i := 0; i < 400 && len(shots) == 0; i++ {
		e.update(ctx)
	}
	if len(shots) == 0 {
		t.Fatal("o corvo corrompido deveria disparar ao descer")
	}
	// Uma única vez: ele é uma ameaça pontual, não uma torre.
	for i := 0; i < 400; i++ {
		e.update(ctx)
	}
	if len(shots) != 1 {
		t.Errorf("o corvo corrompido deveria atirar uma única vez, atirou %d", len(shots))
	}
}

// Deixar o reino cair troca o confronto final: é o gancho de rejogo.
func TestFallenKingdomSummonsAscendedBoss(t *testing.T) {
	g := simGame()
	g.startBoss()
	normal := g.boss

	g2 := simGame()
	g2.corruption.value = maxCorruption
	g2.corruption.tier = tierFall
	g2.startBoss()
	ascended := g2.boss

	if !ascended.ascended {
		t.Fatal("com o reino caído deveria surgir Vharak Ascendido")
	}
	if ascended.maxHealth <= normal.maxHealth {
		t.Errorf("o chefe ascendido deveria ser mais resistente: %d vs %d",
			ascended.maxHealth, normal.maxHealth)
	}
	if ascended.name() == normal.name() {
		t.Error("o confronto verdadeiro deveria se apresentar com outro nome")
	}
	for phase := bossPhase1; phase <= bossPhase3; phase++ {
		ascended.phase, normal.phase = phase, phase
		if len(ascended.patterns()) <= len(normal.patterns()) {
			t.Errorf("fase %d: o ascendido deveria ter mais padrões", phase)
		}
	}
}

// A corrupção multiplica os pontos — é o outro lado da aposta.
func TestCorruptionMultipliesScore(t *testing.T) {
	base := simGame()
	basePoints := base.registerKill(100)

	rich := simGame()
	rich.corruption.value = maxCorruption
	rich.corruption.tier = tierFall
	richPoints := rich.registerKill(100)

	if richPoints <= basePoints {
		t.Errorf("corrupção alta deveria pagar mais: %d vs %d", richPoints, basePoints)
	}
}

// Reiniciar tem que limpar o medidor: uma partida nova começa com o reino em pé.
func TestResetClearsCorruption(t *testing.T) {
	g := New()
	g.corruption.value = 80
	g.corruption.tier = tierFor(80)
	g.corruption.escaped = 12

	g.reset()

	if g.corruption.value != 0 || g.corruption.escaped != 0 || g.corruption.tier != tierSteady {
		t.Errorf("uma sessão nova deveria zerar a corrupção: %v, %d fugas, faixa %q",
			g.corruption.value, g.corruption.escaped, g.corruption.tierName())
	}
}
