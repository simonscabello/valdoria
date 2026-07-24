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
type difficultyParams struct {
	name           string
	startLives     int
	startBombs     int
	enemyHealthMul float64
	dropChanceMul  float64
	// spawnXJitter: desvio máximo em ± pixels na posição X do spawn.
	spawnXJitter int
	// spawnTickJitter: desvio máximo em ± frames no intervalo entre aparições
	// da mesma onda (não atrasa o startTick da onda).
	spawnTickJitter int
}

var difficultyTable = [difficultyCount]difficultyParams{
	diffEasy: {
		name: "Facil", startLives: 4, startBombs: 3,
		enemyHealthMul: 0.8, dropChanceMul: 1.4,
		spawnXJitter: 0, spawnTickJitter: 0,
	},
	diffNormal: {
		name: "Normal", startLives: startingLives, startBombs: bombStartCharges,
		enemyHealthMul: 1.0, dropChanceMul: 1.0,
		spawnXJitter: 14, spawnTickJitter: 8,
	},
	diffHard: {
		name: "Dificil", startLives: 2, startBombs: 1,
		enemyHealthMul: 1.35, dropChanceMul: 0.7,
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
