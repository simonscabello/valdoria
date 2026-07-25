package game

import "testing"

func TestIceSlowsEnemy(t *testing.T) {
	e := newCrow(100)
	start := e.y
	e.applyWeaponHit(weaponIce)
	if e.slow <= 0 {
		t.Fatal("gelo deveria aplicar lentidão")
	}
	ctx := &enemyContext{}
	// Com slow, vários updates não avançam tanto quanto sem.
	for i := 0; i < 10; i++ {
		e.update(ctx)
	}
	slowed := e.y - start

	e2 := newCrow(100)
	for i := 0; i < 10; i++ {
		e2.update(ctx)
	}
	normal := e2.y - start
	if slowed >= normal {
		t.Fatalf("inimigo lento deveria avançar menos (%v vs %v)", slowed, normal)
	}
}

func TestFireBurnsAndWeakens(t *testing.T) {
	e := newCrow(100)
	e.health = 10
	e.applyWeaponHit(weaponFlame)
	if e.burn <= 0 {
		t.Fatal("fogo deveria aplicar queimadura")
	}
	// Dano direto enquanto queima recebe bônus.
	before := e.health
	e.takeDamage(2)
	lost := before - e.health
	if lost < 3 { // 2 + 50% = 3
		t.Fatalf("queimadura deveria aumentar dano recebido, perdeu %d", lost)
	}
}

// A Lança de Luz não aplica status nem dano em cadeia.
//
// Guarda de regressão do defeito mais grave que o MVP teve: a Luz atordoava por
// 32 frames a cada acerto, com cooldown de 10 — qualquer alvo sob fogo contínuo
// ficava congelado indefinidamente, inclusive Vharak, que morria sem executar
// um único padrão. A identidade da Luz é dano focado puro.
func TestLightAppliesNoStatusOrChain(t *testing.T) {
	g := &Game{player: newPlayer(), formations: map[int]*formationTracker{}}
	alvo := newCrow(100)
	alvo.x, alvo.y = 100, 80
	alvo.health = 20
	vizinho := newCrow(100)
	vizinho.x, vizinho.y = 120, 80
	vizinho.health = 20
	g.enemies = append(g.enemies, alvo, vizinho)

	shot := newBullet(100, 80, 0, -1, 4, 0, lightColor)
	shot.element = weaponLight
	g.bullets = append(g.bullets, shot)
	g.bulletEnemyCollisions()

	if alvo.health != 16 {
		t.Fatalf("a luz deveria causar apenas dano direto: vida %d, quero 16", alvo.health)
	}
	if alvo.stun != 0 {
		t.Fatalf("a luz não pode atordoar (stun-lock), stun = %d", alvo.stun)
	}
	if alvo.slow != 0 || alvo.burn != 0 {
		t.Fatal("a luz não aplica lentidão nem queimadura")
	}
	if vizinho.health != 20 {
		t.Fatalf("a luz não encadeia dano; vizinho ficou com %d", vizinho.health)
	}
}

// O chefe também não pode ser atordoado passivamente: ele precisa executar os
// padrões que justificam o confronto.
func TestBossIsNotStunnedByLight(t *testing.T) {
	b := newBoss(false)
	b.applyWeaponHit(weaponLight)
	if b.stun != 0 {
		t.Fatalf("Vharak não pode ser atordoado pela luz, stun = %d", b.stun)
	}
}

func TestBossHasMoreHealth(t *testing.T) {
	if bossMaxHealth < 350 {
		t.Fatalf("chefe deveria ser mais resistente, HP=%d", bossMaxHealth)
	}
	b := newBoss(false)
	if b.health != bossMaxHealth {
		t.Fatalf("HP inicial %d != %d", b.health, bossMaxHealth)
	}
}
