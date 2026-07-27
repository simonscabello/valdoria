package game

import "testing"

// grazingGame monta uma cena com um projétil raspando o jogador.
func grazingGame() (*Game, *EnemyBullet) {
	g := simGame()
	g.state = statePlaying
	b := newEnemyBullet(g.player.centerX()+grazeRadius-2, g.player.centerY(), 0, 0)
	g.enemyBullets = append(g.enemyBullets, b)
	return g, b
}

func TestGrazeChargesTheBomb(t *testing.T) {
	g, _ := grazingGame()
	before := g.grazeCharge

	g.updateGraze()

	if g.grazeCharge <= before {
		t.Errorf("raspar um projétil deveria carregar a invocação: %v -> %v", before, g.grazeCharge)
	}
	if g.grazeCount != 1 {
		t.Errorf("o raspão deveria ser contabilizado, foi %d", g.grazeCount)
	}
}

// Um mesmo projétil não pode encher a bomba sozinho enquanto passa.
func TestSameBulletGrazesOnlyOnceInAWhile(t *testing.T) {
	g, b := grazingGame()
	g.updateGraze()
	first := g.grazeCharge

	g.updateGraze()
	if g.grazeCharge != first {
		t.Errorf("o mesmo projétil não deveria render carga em frames seguidos: %v", g.grazeCharge)
	}

	for i := 0; i < grazeCooldown; i++ {
		b.update()
	}
	g.updateGraze()
	if g.grazeCharge <= first {
		t.Error("passado o intervalo, o mesmo projétil pode render de novo")
	}
}

// Projétil longe não conta: o graze precisa ser sobre risco real.
func TestDistantBulletDoesNotGraze(t *testing.T) {
	g := simGame()
	g.state = statePlaying
	g.enemyBullets = append(g.enemyBullets,
		newEnemyBullet(g.player.centerX()+grazeRadius*3, g.player.centerY(), 0, 0))

	g.updateGraze()

	if g.grazeCharge != 0 || g.grazeCount != 0 {
		t.Error("um projétil distante não deveria render graze")
	}
}

// Carga cheia vira uma Invocação Ancestral: é o que faz a bomba circular em vez
// de ficar guardada até o fim da partida.
func TestGrazeGrantsBombCharge(t *testing.T) {
	g := simGame()
	g.bombCharges = 0

	g.addGrazeCharge(grazeFull)

	if g.bombCharges != 1 {
		t.Errorf("encher a carga deveria conceder uma bomba, tem %d", g.bombCharges)
	}
	if g.grazeCharge >= grazeFull {
		t.Errorf("a carga deveria zerar ao virar bomba, foi %v", g.grazeCharge)
	}
}

// O cinto tem teto: raspar com as bombas no máximo não acumula reserva infinita.
func TestGrazeRespectsBombCap(t *testing.T) {
	g := simGame()
	g.bombCharges = bombMaxCharges

	g.addGrazeCharge(grazeFull * 3)

	if g.bombCharges != bombMaxCharges {
		t.Errorf("as bombas não deveriam passar de %d, tem %d", bombMaxCharges, g.bombCharges)
	}
	if g.grazeCharge != 0 {
		t.Errorf("com o cinto cheio a carga não deveria acumular, foi %v", g.grazeCharge)
	}
}

// Durante a invulnerabilidade não há mérito em chegar perto.
func TestNoGrazeWhileInvulnerable(t *testing.T) {
	g, _ := grazingGame()
	g.player.invincible = 30

	g.updateGraze()

	if g.grazeCharge != 0 {
		t.Error("raspar durante a invencibilidade não deveria render carga")
	}
}

// O raio de graze precisa ser bem maior que a hitbox, senão a mecânica vira
// sorte: tem que existir uma faixa legível entre "quase acertou" e "acertou".
func TestGrazeRadiusIsReadable(t *testing.T) {
	if grazeRadius < playerHitboxSize*2 {
		t.Errorf("o raio de graze (%v) precisa ser bem maior que a hitbox (%v)",
			grazeRadius, float64(playerHitboxSize))
	}
}
