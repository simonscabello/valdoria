package game

import (
	"os"
	"strconv"
)

// devConfig reúne as opções de desenvolvimento. Tudo fica desligado por padrão,
// então o build padrão não expõe nenhuma dessas ferramentas.
type devConfig struct {
	enabled      bool
	showHUD      bool // FPS, contagens e tempo da fase
	showHitboxes bool
	invincible   bool

	startSection  int
	startBoss     bool
	timeScale     int
	fastTimeScale int
}

var dev = devConfig{
	timeScale:     1,
	fastTimeScale: 6,
}

// Options configura o jogo de forma programática (útil para testes e ferramentas).
type Options struct {
	Dev          bool
	StartSection int
	StartBoss    bool
	Seed         int64
	HasSeed      bool
}

// Configure aplica as opções recebidas ao estado global do pacote.
func Configure(o Options) {
	dev.enabled = o.Dev
	dev.showHUD = o.Dev
	dev.startSection = o.StartSection
	dev.startBoss = o.StartBoss
	if o.HasSeed {
		SetSeed(o.Seed)
	}
}

// InitFromEnv lê variáveis de ambiente e configura o jogo. É o gancho chamado
// pelo main para ativar o modo de desenvolvimento sem afetar o build padrão:
//
//	VALDORIA_DEV=1        ativa o modo de desenvolvimento
//	VALDORIA_SECTION=2    inicia a fase no trecho 2 (0 a 3)
//	VALDORIA_BOSS=1       inicia direto no chefe
//	VALDORIA_SEED=42      fixa a semente da jogabilidade
func InitFromEnv() Options {
	o := Options{}
	if v := os.Getenv("VALDORIA_DEV"); v == "1" || v == "true" {
		o.Dev = true
	}
	if v := os.Getenv("VALDORIA_SECTION"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			o.StartSection = n
		}
	}
	if v := os.Getenv("VALDORIA_BOSS"); v == "1" || v == "true" {
		o.StartBoss = true
	}
	if v := os.Getenv("VALDORIA_SEED"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			o.Seed = n
			o.HasSeed = true
		}
	}
	Configure(o)
	return o
}
