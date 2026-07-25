package game

import "math"

// difficultyLevel seleciona o preset de dificuldade da campanha. É um estado
// global de pacote (como dev/seed), lido nos pontos de spawn e no reset.
type difficultyLevel int

const (
	diffEasy difficultyLevel = iota
	diffNormal
	diffHard
	difficultyCount
)

// difficultyParams reúne os multiplicadores e recursos iniciais de cada preset.
//
// A dificuldade mexe nos **quatro eixos que definem a pressão** de um shmup:
// quanta ameaça existe (vida dos inimigos), com que rapidez ela chega
// (velocidade dos projéteis), com que frequência ela chega (cadência de tiro) e
// quanto o jogador aguenta (vidas, invencibilidade, drops).
//
// Antes, os presets mexiam só na vida dos inimigos, nos drops e no jitter de
// spawn: os três modos jogavam quase igual, porque a ameaça real — projéteis na
// tela — era idêntica nos três.
type difficultyParams struct {
	name           string
	startLives     int
	startBombs     int
	enemyHealthMul float64
	dropChanceMul  float64

	// bulletSpeedMul e fireRateMul controlam a pressão de fogo. Cadência acima
	// de 1 significa intervalos **menores** entre disparos.
	bulletSpeedMul float64
	fireRateMul    float64
	// invincibilityMul estica ou encurta o perdão após levar um golpe.
	invincibilityMul float64

	// spawnXJitter: desvio máximo em ± pixels na posição X do spawn.
	spawnXJitter int
	// spawnTickJitter: desvio máximo em ± frames no intervalo entre aparições
	// da mesma onda (não atrasa o startTick da onda).
	spawnTickJitter int
}

var difficultyTable = [difficultyCount]difficultyParams{
	diffEasy: {
		name: "Facil", startLives: 4, startBombs: 3,
		enemyHealthMul: 0.7, dropChanceMul: 1.4,
		bulletSpeedMul: 0.85, fireRateMul: 0.75, invincibilityMul: 1.35,
		spawnXJitter: 0, spawnTickJitter: 0,
	},
	diffNormal: {
		name: "Normal", startLives: startingLives, startBombs: bombStartCharges,
		enemyHealthMul: 1.0, dropChanceMul: 1.0,
		bulletSpeedMul: 1.0, fireRateMul: 1.0, invincibilityMul: 1.0,
		spawnXJitter: 14, spawnTickJitter: 8,
	},
	diffHard: {
		name: "Dificil", startLives: 2, startBombs: 1,
		enemyHealthMul: 1.3, dropChanceMul: 0.7,
		bulletSpeedMul: 1.15, fireRateMul: 1.35, invincibilityMul: 0.8,
		spawnXJitter: 28, spawnTickJitter: 16,
	},
}

var currentDifficulty = diffNormal

func setDifficulty(d difficultyLevel) {
	if d < 0 || d >= difficultyCount {
		d = diffNormal
	}
	currentDifficulty = d
}

func diffParams() difficultyParams { return difficultyTable[currentDifficulty] }

func difficultyName(d difficultyLevel) string {
	if d < 0 || d >= difficultyCount {
		return "Normal"
	}
	return difficultyTable[d].name
}

// scaledEnemyHealth aplica o multiplicador de vida do preset ao valor-base do
// inimigo, nunca deixando cair abaixo de 1.
func scaledEnemyHealth(base int) int {
	h := int(math.Round(float64(base) * diffParams().enemyHealthMul))
	if h < 1 {
		h = 1
	}
	return h
}

// scaledDropChance ajusta a chance de drop conforme o preset.
func scaledDropChance(base float64) float64 {
	return base * diffParams().dropChanceMul
}

// scaledBulletSpeed ajusta a velocidade dos projéteis inimigos. É o eixo que
// mais muda a sensação de dificuldade: no Difícil os tiros são bem mais rápidos
// que o jogador, no Fácil dá para sair andando na frente deles.
func scaledBulletSpeed(base float64) float64 {
	return base * diffParams().bulletSpeedMul
}

// scaledFireInterval encurta (Difícil) ou alonga (Fácil) o intervalo entre
// disparos inimigos. Nunca desce abaixo de 1 frame.
func scaledFireInterval(base int) int {
	mul := diffParams().fireRateMul
	if mul <= 0 {
		return base
	}
	n := int(math.Round(float64(base) / mul))
	if n < 1 {
		return 1
	}
	return n
}

// scaledInvincibility ajusta o perdão após levar um golpe.
func scaledInvincibility(base int) int {
	n := int(math.Round(float64(base) * diffParams().invincibilityMul))
	if n < 1 {
		return 1
	}
	return n
}

// jitterSpawnX desloca a posição scriptada em ±spawnXJitter (0 no Fácil).
func jitterSpawnX(base float64) float64 {
	j := diffParams().spawnXJitter
	if j <= 0 {
		return base
	}
	return base + float64(rng.Intn(2*j+1)-j)
}

// jitterSpawnInterval devolve o intervalo da onda com um atraso/adianto
// aleatório. Nunca acelera abaixo de ~2/3 do intervalo base (evita rajadas
// impossíveis no início).
func jitterSpawnInterval(base int) int {
	j := diffParams().spawnTickJitter
	if j <= 0 || base <= 0 {
		return base
	}
	next := base + rng.Intn(2*j+1) - j
	floor := (base * 2) / 3
	if floor < 1 {
		floor = 1
	}
	if next < floor {
		return floor
	}
	return next
}
