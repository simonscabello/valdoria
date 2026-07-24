package game

import (
	"math"
	"testing"
)

func TestBulletDoesNotDamageSameEnemyTwice(t *testing.T) {
	g := &Game{player: newPlayer()}
	e := newWyvern(0)
	e.x, e.y, e.health = 100, 100, 20
	g.enemies = append(g.enemies, e)
	b := newBullet(100, 100, 0, -1, 1, 5, iceColor) // perfura, mas não repete no mesmo alvo
	g.bullets = append(g.bullets, b)

	g.handleCollisions()
	g.handleCollisions()

	if e.health != 19 {
		t.Errorf("o mesmo projétil deveria acertar o inimigo uma vez, vida %d", e.health)
	}
}

func TestPlayerTakesOnlyOneHitPerFrame(t *testing.T) {
	g := &Game{player: newPlayer()}
	g.player.x, g.player.y = 100, 100
	for i := 0; i < 3; i++ {
		e := newCrow(0)
		e.x, e.y = 100, 100
		g.enemies = append(g.enemies, e)
	}
	before := g.player.health

	g.handleCollisions()

	if g.player.health != before-1 {
		t.Errorf("o jogador deveria perder apenas 1 de vida por frame, perdeu %d", before-g.player.health)
	}
}

func TestRemoveDeadCompactsAndCountsKills(t *testing.T) {
	g := &Game{formations: map[int]*formationTracker{}}
	killed := newCrow(0)
	killed.dead = true
	escaped := newCrow(0)
	escaped.dead, escaped.escaped = true, true
	alive := newCrow(0)
	g.enemies = append(g.enemies, killed, escaped, alive)

	g.removeDead()

	if len(g.enemies) != 1 || g.enemies[0] != alive {
		t.Fatalf("removeDead deveria manter só o inimigo vivo, restaram %d", len(g.enemies))
	}
	if g.enemiesDefeated != 1 {
		t.Errorf("apenas o inimigo abatido (não o que escapou) deveria contar, contou %d", g.enemiesDefeated)
	}
}

func TestOffscreenBulletIsRemoved(t *testing.T) {
	g := &Game{}
	up := newBullet(10, 0, 0, -20, 1, 0, lightColor)
	g.bullets = append(g.bullets, up)

	g.updateBullets()
	g.removeDead()

	if len(g.bullets) != 0 {
		t.Error("projétil que sai da tela deveria ser removido")
	}
}

func TestOffscreenEnemyEscapesWithoutScore(t *testing.T) {
	g := &Game{player: newPlayer(), formations: map[int]*formationTracker{}}
	e := newCrow(50)
	e.y = ScreenHeight + 100
	g.enemies = append(g.enemies, e)

	g.updateEnemies()
	if !e.escaped || !e.dead {
		t.Fatal("inimigo fora da tela deveria escapar e ser marcado para remoção")
	}
	g.removeDead()
	if g.enemiesDefeated != 0 {
		t.Error("inimigo que escapou não deveria contar como abatido")
	}
	if len(g.enemies) != 0 {
		t.Error("inimigo que escapou deveria ser removido")
	}
}

func TestEnemyMovesDeterministicallyPerTick(t *testing.T) {
	e := newCrow(50)
	ctx := &enemyContext{}
	startY := e.y
	for i := 0; i < 10; i++ {
		e.update(ctx)
	}
	if math.Abs(e.y-(startY+10*crowSpeed)) > 1e-9 {
		t.Errorf("movimento por tick deveria ser previsível, y=%v", e.y)
	}
}

func TestStartNewGameClearsAllState(t *testing.T) {
	g := New()
	g.startNewGame()

	g.score = 500
	g.enemiesDefeated = 9
	g.elapsedTicks = 1200
	g.maxMult = 4
	g.shakeMag = 5
	g.damageFlash = 8
	g.particles = append(g.particles, &particle{life: 5, maxLife: 5})
	g.enemies = append(g.enemies, newCrow(0))
	g.enemyBullets = append(g.enemyBullets, newEnemyBullet(0, 0, 0, 1))
	g.startBoss()

	g.startNewGame()

	if g.score != 0 || g.enemiesDefeated != 0 || g.elapsedTicks != 0 || g.maxMult != 1 {
		t.Error("contadores da sessão deveriam reiniciar")
	}
	if g.shakeMag != 0 || g.damageFlash != 0 || len(g.particles) != 0 {
		t.Error("efeitos deveriam reiniciar")
	}
	if len(g.enemies) != 0 || len(g.enemyBullets) != 0 || g.boss != nil {
		t.Error("nenhuma entidade anterior deveria permanecer")
	}
	if g.state != statePlaying || g.lives != startingLives {
		t.Error("a nova sessão deveria começar jogando com as vidas iniciais")
	}
}
