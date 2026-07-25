package game

import "testing"

// TestGameDemandsDodging é a guarda contra o defeito relatado em playtest: "o
// jogo está muito fácil, mesmo no Normal".
//
// A causa medida era concreta — uma salva da Lança de Luz Nv3 (12 de dano)
// matava **qualquer** inimigo do jogo, então eles morriam antes de disparar: a
// campanha inteira gerava 57 projéteis inimigos, e um jogador que nunca
// desviava a terminava perdendo 2 pontos de vida.
//
// O critério é simples e não depende de opinião: um jogador com poder de fogo
// máximo e **zero** habilidade de desvio não pode terminar a campanha.
func TestGameDemandsDodging(t *testing.T) {
	prev := currentDifficulty
	defer setDifficulty(prev)
	setDifficulty(diffNormal)
	SetSeed(42)

	for _, m := range measurePressure() {
		if !m.Died {
			t.Errorf("modelo %q terminou a campanha sem desviar (%d golpes em %.0fs): jogo fácil demais",
				m.Model, m.Hits, m.SurvivedSeconds)
		}
	}
}

// TestEnemiesThreatenWithoutDragging guarda as **duas** falhas opostas, que
// aparecem iguais numa tabela de vida e são problemas contrários:
//
//	inofensivo -> morre antes de disparar (o defeito original)
//	esponja    -> fica viva absorvendo tiros (o defeito da primeira correção)
//
// Não dá para checar isso por vida: o tempo-para-matar depende também da
// largura do inimigo e do quanto ele se esquiva. A salva da Luz cobre 18px,
// então um wyvern de 24px morre mais rápido que uma gárgula de 18px com metade
// da vida dele.
func TestEnemiesThreatenWithoutDragging(t *testing.T) {
	prev := currentDifficulty
	defer setDifficulty(prev)
	setDifficulty(diffNormal)
	SetSeed(42)

	for _, m := range measureThreat() {
		if m.TTK > maxEnemyTTK {
			t.Errorf("%s vira esponja: %.2fs para morrer (teto %.1fs)",
				m.Kind, m.TTK, maxEnemyTTK)
		}
		// O corvo não atira de propósito: ele é a massa descartável do jogo.
		if m.Kind == EnemyKindName(kindCrow) {
			continue
		}
		if m.Shots == 0 {
			t.Errorf("%s morre sem disparar (%.2fs): é alvo, não ameaça", m.Kind, m.TTK)
		}
	}
}

// O projétil inimigo precisa ser mais rápido que o jogador — senão dá para
// simplesmente sair andando na frente dele, e nenhum padrão de tiro importa.
func TestEnemyBulletsOutrunThePlayer(t *testing.T) {
	prev := currentDifficulty
	defer setDifficulty(prev)
	setDifficulty(diffNormal)

	if enemyBulletSpeed <= playerSpeed {
		t.Errorf("projétil inimigo (%v) deveria ser mais rápido que o jogador (%v)",
			enemyBulletSpeed, playerSpeed)
	}
}

// As três dificuldades precisam jogar de forma diferente. Antes, elas mexiam só
// na vida dos inimigos e nos drops: a ameaça real — projéteis na tela — era
// idêntica nas três, e os presets eram quase indistinguíveis.
func TestDifficultiesChangeRealPressure(t *testing.T) {
	prev := currentDifficulty
	defer setDifficulty(prev)

	easy := difficultyTable[diffEasy]
	normal := difficultyTable[diffNormal]
	hard := difficultyTable[diffHard]

	axes := []struct {
		name             string
		e, n, h          float64
		harderMeansLower bool
	}{
		{"vida dos inimigos", easy.enemyHealthMul, normal.enemyHealthMul, hard.enemyHealthMul, false},
		{"velocidade dos projeteis", easy.bulletSpeedMul, normal.bulletSpeedMul, hard.bulletSpeedMul, false},
		{"cadencia de tiro", easy.fireRateMul, normal.fireRateMul, hard.fireRateMul, false},
		{"invencibilidade", easy.invincibilityMul, normal.invincibilityMul, hard.invincibilityMul, true},
	}
	for _, a := range axes {
		ok := a.e < a.n && a.n < a.h
		if a.harderMeansLower {
			ok = a.e > a.n && a.n > a.h
		}
		if !ok {
			t.Errorf("%s não separa as dificuldades: facil %v, normal %v, dificil %v",
				a.name, a.e, a.n, a.h)
		}
	}
}

// A dificuldade tem que se traduzir em pressão medida, não só em constantes.
func TestHarderDifficultyKillsFaster(t *testing.T) {
	prev := currentDifficulty
	defer setDifficulty(prev)

	survival := map[difficultyLevel]float64{}
	for _, d := range []difficultyLevel{diffEasy, diffNormal, diffHard} {
		setDifficulty(d)
		SetSeed(42)
		survival[d] = simulatePressure(modelPassive).SurvivedSeconds
	}
	if !(survival[diffEasy] > survival[diffNormal] && survival[diffNormal] > survival[diffHard]) {
		t.Errorf("a sobrevivência deveria cair com a dificuldade: facil %.0fs, normal %.0fs, dificil %.0fs",
			survival[diffEasy], survival[diffNormal], survival[diffHard])
	}
}
