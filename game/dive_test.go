package game

import "testing"

func TestDiveCostsStaminaAndGrantsInvulnerability(t *testing.T) {
	p := newPlayer()
	if !p.canDive() {
		t.Fatal("o grifo deveria começar com fôlego para mergulhar")
	}

	full := p.stamina
	if !p.tryDive() {
		t.Fatal("o mergulho deveria iniciar com fôlego cheio")
	}
	if p.stamina >= full {
		t.Errorf("o mergulho deveria custar fôlego: %v -> %v", full, p.stamina)
	}
	if !p.diving() || p.canBeHit() {
		t.Error("durante o mergulho o grifo atravessa: invulnerável")
	}
}

// O mergulho é um gesto, não um estado: ele termina sozinho e devolve o
// controle. Sem isso viraria invencibilidade sob demanda.
func TestDiveEndsAndRecovers(t *testing.T) {
	p := newPlayer()
	p.tryDive()

	for i := 0; i < diveDuration; i++ {
		p.updateDive()
	}
	if p.diving() {
		t.Error("o mergulho deveria acabar depois de diveDuration frames")
	}
	if p.canBeHit() == false && p.invincible == 0 {
		t.Error("terminado o mergulho, o grifo volta a ser atingível")
	}

	// Durante a recuperação não dá para emendar outro mergulho.
	if p.canDive() {
		t.Error("não deveria dar para mergulhar de novo durante a recuperação")
	}
	for i := 0; i < diveRecovery+1; i++ {
		p.updateDive()
	}
	if p.diveState != diveIdle {
		t.Error("a recuperação deveria terminar e liberar o gesto")
	}
}

// Sem fôlego não há mergulho — é o custo que impede o abuso.
func TestDiveBlockedWithoutStamina(t *testing.T) {
	p := newPlayer()
	p.stamina = diveStaminaCost - 1
	if p.tryDive() {
		t.Error("sem fôlego suficiente o mergulho não deveria iniciar")
	}
}

func TestStaminaRegeneratesToMax(t *testing.T) {
	p := newPlayer()
	p.tryDive()
	for i := 0; i < diveDuration+diveRecovery; i++ {
		p.updateDive()
	}
	drained := p.stamina
	for i := 0; i < 600; i++ {
		p.updateDive()
	}
	if p.stamina <= drained {
		t.Errorf("o fôlego deveria se recuperar: %v -> %v", drained, p.stamina)
	}
	if p.stamina > staminaMax {
		t.Errorf("o fôlego não deveria passar do máximo, foi %v", p.stamina)
	}
}

// Atravessar um inimigo fere uma única vez por mergulho: o grifo passa por
// dentro, não fica triturando.
func TestDiveDamagesEachEnemyOnce(t *testing.T) {
	g := simGame()
	e := newWyvern(100)
	e.x, e.y = g.player.x, g.player.y
	e.health = 999
	g.enemies = append(g.enemies, e)

	g.player.tryDive()
	before := e.health
	g.diveCollisions()
	afterFirst := e.health
	g.diveCollisions()

	if afterFirst >= before {
		t.Fatalf("o mergulho deveria ferir o inimigo atravessado: %d -> %d", before, afterFirst)
	}
	if e.health != afterFirst {
		t.Errorf("o mesmo inimigo não pode ser ferido duas vezes no mesmo mergulho: %d", e.health)
	}
}

// A direção segue o que o jogador aponta; sem direção, o grifo avança.
func TestDiveDefaultsForward(t *testing.T) {
	p := newPlayer()
	p.tryDive()
	if p.diveVY >= 0 {
		t.Errorf("sem direção apontada o mergulho deveria ir para frente (cima), vy = %v", p.diveVY)
	}
}
