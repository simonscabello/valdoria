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

func TestLightStunsAndChains(t *testing.T) {
	g := &Game{player: newPlayer(), formations: map[int]*formationTracker{}}
	a := newCrow(100)
	a.x, a.y = 100, 80
	a.health = 20
	b := newCrow(100)
	b.x, b.y = 120, 80
	b.health = 20
	g.enemies = append(g.enemies, a, b)

	shot := newBullet(100, 80, 0, -1, 4, 0, lightColor)
	shot.element = weaponLight
	g.bullets = append(g.bullets, shot)
	g.bulletEnemyCollisions()

	if a.stun <= 0 {
		t.Fatal("luz deveria atordoar o alvo")
	}
	if b.health >= 20 {
		t.Fatal("luz deveria encadear dano no vizinho")
	}
	if b.stun <= 0 {
		t.Fatal("cadeia deveria atordoar o vizinho")
	}
}

func TestBossHasMoreHealth(t *testing.T) {
	if bossMaxHealth < 350 {
		t.Fatalf("chefe deveria ser mais resistente, HP=%d", bossMaxHealth)
	}
	b := newBoss()
	if b.health != bossMaxHealth {
		t.Fatalf("HP inicial %d != %d", b.health, bossMaxHealth)
	}
}
