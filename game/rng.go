package game

import (
	"math/rand"
	"time"
)

// rng é a fonte de aleatoriedade que afeta a jogabilidade (drops, posições de
// spawn). Manter apenas o essencial aqui torna as sessões reproduzíveis por seed.
// fxRand cobre efeitos puramente visuais (partículas, vibração, cenário) e não
// interfere na reprodução de bugs de jogabilidade.
var (
	gameplaySeed = time.Now().UnixNano()
	seedFixed    = false
	rng          = rand.New(rand.NewSource(gameplaySeed))
	fxRand       = rand.New(rand.NewSource(time.Now().UnixNano() + 1))
)

// SetSeed fixa a semente da jogabilidade para reproduzir uma sessão. A partir
// daqui, cada nova sessão reinicia a sequência a partir desta semente.
func SetSeed(seed int64) {
	gameplaySeed = seed
	seedFixed = true
	rng = rand.New(rand.NewSource(seed))
}

// CurrentSeed devolve a semente de jogabilidade em uso.
func CurrentSeed() int64 { return gameplaySeed }

// resetRNG reinicia a sequência de jogabilidade quando a semente é fixa,
// garantindo que reiniciar a sessão reproduza exatamente os mesmos sorteios.
func resetRNG() {
	if seedFixed {
		rng = rand.New(rand.NewSource(gameplaySeed))
	}
}
