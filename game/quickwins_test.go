package game

import (
	"image/color"
	"math"
	"testing"
)

func TestNormalizeDiagonalKeepsSpeed(t *testing.T) {
	// Reto: sem alteração.
	if dx, dy := normalizeDiagonal(1, 0); dx != 1 || dy != 0 {
		t.Errorf("movimento reto não deveria mudar, obtive (%v, %v)", dx, dy)
	}
	// Diagonal: módulo deve ser 1 (mesma velocidade da reta), não √2.
	dx, dy := normalizeDiagonal(1, 1)
	if mag := math.Hypot(dx, dy); math.Abs(mag-1) > 1e-6 {
		t.Errorf("diagonal deveria ter módulo 1, foi %v", mag)
	}
	// Parado permanece parado.
	if dx, dy := normalizeDiagonal(0, 0); dx != 0 || dy != 0 {
		t.Errorf("parado deveria continuar parado, obtive (%v, %v)", dx, dy)
	}
}

func TestBossDeathClearsArenaAndProtectsPlayer(t *testing.T) {
	g := New()
	g.startNewGame()
	g.startBoss()
	g.boss.phase = bossPhase2
	g.boss.invulnerable = false
	g.boss.health = 3
	// Cena poluída no instante da morte: inimigos invocados e projéteis na tela.
	g.enemies = append(g.enemies, newCrow(50), newHarpy(80))
	g.enemyBullets = append(g.enemyBullets, newEnemyBullet(60, 60, 0, 2))
	g.player.invincible = 0

	g.boss.takeDamage(10) // zera a vida -> startDying -> justDied
	g.updateBoss()

	if len(g.enemies) != 0 {
		t.Errorf("a arena deveria ser limpa ao morrer o chefe, restaram %d inimigos", len(g.enemies))
	}
	if len(g.enemyBullets) != 0 {
		t.Errorf("os projéteis inimigos deveriam ser limpos, restaram %d", len(g.enemyBullets))
	}
	if g.player.invincible < bossDeathDuration {
		t.Errorf("o jogador deveria ficar invulnerável durante a morte do chefe, foi %d", g.player.invincible)
	}
}

func TestScorePopupSpawnsOnKillAndResets(t *testing.T) {
	g := &Game{player: newPlayer(), formations: map[int]*formationTracker{}}
	e := newCrow(100)
	e.x, e.y = 100, 100
	g.enemies = append(g.enemies, e)
	g.bullets = append(g.bullets, newBullet(100, 100, 0, -1, e.health, 0, lightColor))

	g.bulletEnemyCollisions()

	if len(g.popups) != 1 {
		t.Fatalf("abater um inimigo deveria gerar um popup de pontuação, veio %d", len(g.popups))
	}
	if g.popups[0].text != itoa(crowScore) {
		t.Errorf("popup deveria mostrar %s, mostrou %s", itoa(crowScore), g.popups[0].text)
	}

	g.startNewGame()
	if len(g.popups) != 0 {
		t.Errorf("reiniciar deveria limpar os popups, restaram %d", len(g.popups))
	}
}

func TestShieldBreakGivesNoDamageAndKeepsCombo(t *testing.T) {
	g := New()
	g.startNewGame()
	g.player.applyPowerup(powerShield)
	g.combo = 8
	g.comboTimer = 100
	before := g.player.health

	g.hitPlayer(1)

	if g.player.health != before {
		t.Errorf("o escudo não deveria deixar cair a vida, foi %d -> %d", before, g.player.health)
	}
	if g.player.hasShield() {
		t.Error("o escudo deveria ser consumido")
	}
	if g.combo == 0 {
		t.Error("absorver com escudo não deveria zerar o combo")
	}
	if g.sectionDamaged {
		t.Error("absorver com escudo não deveria marcar o trecho como danificado")
	}
}

func TestHitPlayerRealDamageResetsCombo(t *testing.T) {
	g := New()
	g.startNewGame()
	g.combo = 8
	g.comboTimer = 100
	before := g.player.health

	g.hitPlayer(1)

	if g.player.health != before-1 {
		t.Errorf("sem escudo o dano deveria reduzir a vida, foi %d -> %d", before, g.player.health)
	}
	if g.combo != 0 || !g.sectionDamaged {
		t.Error("dano real deveria zerar o combo e marcar o trecho")
	}
}

func TestEnemyBulletColorDistinctFromPlayerWeapons(t *testing.T) {
	// Guarda contra regressão de leitura: o projétil inimigo precisa ser bem
	// diferente de todas as cores de arma do jogador.
	weapons := []color.RGBA{lightColor, flameColor, iceColor}
	for _, w := range weapons {
		if colorDistance(enemyBulletColor, w) < 120 {
			t.Errorf("cor do projétil inimigo %v está próxima demais da arma %v", enemyBulletColor, w)
		}
	}
}

func colorDistance(a, b color.RGBA) float64 {
	dr := float64(a.R) - float64(b.R)
	dg := float64(a.G) - float64(b.G)
	db := float64(a.B) - float64(b.B)
	return math.Sqrt(dr*dr + dg*dg + db*db)
}

func TestParticleBurstRespectsCap(t *testing.T) {
	g := &Game{}
	for i := 0; i < maxParticles+50; i++ {
		g.particles = append(g.particles, &particle{life: 10, maxLife: 10})
	}
	before := len(g.particles)
	g.spawnBurst(10, 10, explosionParticles, 2, 20, 2, false, color.RGBA{})
	if len(g.particles) != before {
		t.Errorf("no teto de partículas nenhuma nova deveria ser criada, %d -> %d", before, len(g.particles))
	}
}

func TestEnemyTelegraphsBeforeFiring(t *testing.T) {
	e := newHarpy(50)
	var bullets []*EnemyBullet
	ctx := &enemyContext{bullets: &bullets}

	sawTelegraphBeforeShot := false
	for i := 0; i < harpyFireInterval+1; i++ {
		e.update(ctx)
		if len(bullets) > 0 {
			break
		}
		if e.telegraph > 0 {
			sawTelegraphBeforeShot = true
		}
	}

	if len(bullets) == 0 {
		t.Fatal("a harpia deveria ter disparado dentro do intervalo")
	}
	if !sawTelegraphBeforeShot {
		t.Error("o inimigo deveria avisar (telegraph) antes de disparar")
	}
}

func TestFormationBonusSpawnsPopup(t *testing.T) {
	g := &Game{formations: map[int]*formationTracker{}}
	g.formations[1] = &formationTracker{total: 2}
	for i := 0; i < 2; i++ {
		e := newCrow(0)
		e.formationID = 1
		e.dead = true
		g.enemies = append(g.enemies, e)
	}

	g.removeDead()

	if len(g.popups) == 0 {
		t.Error("destruir uma formação completa deveria gerar um popup de bônus")
	}
}

func TestSectionBonusSpawnsPopup(t *testing.T) {
	g := New()
	g.startNewGame()
	g.lastSection = 0
	g.sectionDamaged = false
	g.level.section = 1

	g.checkSectionBonus()

	if len(g.popups) == 0 {
		t.Error("concluir um trecho sem dano deveria gerar um popup de bônus")
	}
}

func TestHitStopFreezesGameplay(t *testing.T) {
	g := New()
	g.startNewGame()
	g.hitStop = 2
	tickBefore := g.level.tick
	elapsedBefore := g.elapsedTicks

	g.updatePlay()

	if g.level.tick != tickBefore {
		t.Error("durante o hit-stop a linha do tempo não deveria avançar")
	}
	if g.elapsedTicks != elapsedBefore {
		t.Error("durante o hit-stop o tempo da fase não deveria avançar")
	}
	if g.hitStop != 1 {
		t.Errorf("o hit-stop deveria decrementar, foi %d", g.hitStop)
	}
}
