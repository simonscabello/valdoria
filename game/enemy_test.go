package game

import (
	"math"
	"testing"
)

func TestEnemyTakeDamageReducesHealth(t *testing.T) {
	e := newWyvern(0)
	start := e.health
	e.takeDamage(3)
	if e.health != start-3 {
		t.Errorf("health = %d, quero %d", e.health, start-3)
	}
	if e.dead {
		t.Error("wyvern não deveria morrer com dano parcial")
	}
	if e.hitFlash == 0 {
		t.Error("dano deveria acionar o hitFlash")
	}
}

func TestEnemyDiesWhenHealthDepleted(t *testing.T) {
	e := newCrow(0)
	e.takeDamage(crowHealth)
	if !e.dead {
		t.Error("corvo deveria estar morto após perder toda a vida")
	}
}

func TestBulletKillGrantsScore(t *testing.T) {
	g := &Game{player: newPlayer()}
	crow := newCrow(100)
	crow.x, crow.y = 100, 100
	g.enemies = append(g.enemies, crow)
	// Dano suficiente para abater: o corvo deixou de morrer com 1 de dano
	// quando a vida dos inimigos foi recalibrada.
	g.bullets = append(g.bullets, newBullet(crow.x, crow.y, 0, -1, crow.health, 0, lightColor))

	g.handleCollisions()

	if !crow.dead {
		t.Fatal("corvo deveria morrer ao ser atingido")
	}
	if g.score != crowScore {
		t.Errorf("score = %d, quero %d", g.score, crowScore)
	}
}

func TestCrowMovesStraightDown(t *testing.T) {
	e := newCrow(50)
	ctx := &enemyContext{}
	startX, startY := e.x, e.y

	e.update(ctx)

	if e.x != startX {
		t.Errorf("corvo não deveria mudar de x: %v -> %v", startX, e.x)
	}
	if e.y != startY+crowSpeed {
		t.Errorf("y = %v, quero %v", e.y, startY+crowSpeed)
	}
}

func TestHarpyFiresDownwardAfterInterval(t *testing.T) {
	e := newHarpy(50)
	var bullets []*EnemyBullet
	ctx := &enemyContext{bullets: &bullets}

	for i := 0; i < harpyFireInterval; i++ {
		e.update(ctx)
	}

	if len(bullets) == 0 {
		t.Fatal("harpia deveria disparar após o intervalo")
	}
	b := bullets[0]
	if b.vx != 0 || b.vy <= 0 {
		t.Errorf("projétil da harpia deveria descer reto, vx=%v vy=%v", b.vx, b.vy)
	}
}

func TestAimVelocityPointsToTarget(t *testing.T) {
	vx, vy := aimVelocity(0, 0, 3, 4, 10)
	if math.Abs(vx-6) > 1e-9 || math.Abs(vy-8) > 1e-9 {
		t.Errorf("aimVelocity = (%v, %v), quero (6, 8)", vx, vy)
	}
	if speed := math.Hypot(vx, vy); math.Abs(speed-10) > 1e-9 {
		t.Errorf("módulo = %v, quero 10", speed)
	}
}

func TestAimVelocityZeroDistanceFallsDown(t *testing.T) {
	vx, vy := aimVelocity(5, 5, 5, 5, 2)
	if vx != 0 || vy != 2 {
		t.Errorf("alvo coincidente deveria descer reto, obtive (%v, %v)", vx, vy)
	}
}

func TestMageFiresFullRing(t *testing.T) {
	e := newMage(100)
	e.y = mageStopY // já posicionado para atirar
	var bullets []*EnemyBullet
	ctx := &enemyContext{playerX: 100, playerY: 300, bullets: &bullets}

	for i := 0; i < mageFireInterval; i++ {
		e.update(ctx)
	}

	if len(bullets) != mageRingCount {
		t.Fatalf("o feiticeiro deveria soltar um anel de %d projéteis, veio %d", mageRingCount, len(bullets))
	}
}

func TestBallistaFiresAimedBurst(t *testing.T) {
	e := newBallista(100)
	var bullets []*EnemyBullet
	ctx := &enemyContext{playerX: 100, playerY: 300, bullets: &bullets}

	// Intervalo + a rajada completa espaçada.
	for i := 0; i < ballistaFireInterval+ballistaBurst*ballistaBurstGap+2; i++ {
		e.update(ctx)
	}

	if len(bullets) < ballistaBurst {
		t.Fatalf("a balista deveria disparar uma rajada de %d virotes, veio %d", ballistaBurst, len(bullets))
	}
	// Todos os virotes descem em direção ao jogador.
	for i, b := range bullets {
		if b.vy <= 0 {
			t.Errorf("virote %d deveria descer, vy=%v", i, b.vy)
		}
	}
}

func TestScaledEnemyHealthByDifficulty(t *testing.T) {
	t.Cleanup(func() { setDifficulty(diffNormal) })

	setDifficulty(diffNormal)
	base := scaledEnemyHealth(wyvernHealth)
	if base != wyvernHealth {
		t.Errorf("no Normal a vida deveria ser a base %d, foi %d", wyvernHealth, base)
	}

	setDifficulty(diffHard)
	hard := scaledEnemyHealth(wyvernHealth)
	if hard <= wyvernHealth {
		t.Errorf("no Difícil a vida deveria ser maior que %d, foi %d", wyvernHealth, hard)
	}

	setDifficulty(diffEasy)
	easy := scaledEnemyHealth(wyvernHealth)
	if easy >= wyvernHealth {
		t.Errorf("no Fácil a vida deveria ser menor que %d, foi %d", wyvernHealth, easy)
	}
}
