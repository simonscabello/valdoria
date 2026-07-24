package main

import (
	"log"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"

	"valdoria/game"
)

// version é o número da versão do jogo. Pode ser sobrescrita em tempo de build
// com: -ldflags "-X main.version=1.2.3".
var version = "dev"

func main() {
	log.SetPrefix("[valdoria] ")
	log.SetFlags(log.Ltime)

	game.Version = version
	opts := game.InitFromEnv()

	log.Printf("Asas de Valdoria %s (%s/%s, %s)", version, runtime.GOOS, runtime.GOARCH, runtime.Version())
	if opts.Dev {
		log.Printf("modo de desenvolvimento ativo (trecho=%d chefe=%v)", opts.StartSection, opts.StartBoss)
	}

	g := game.New()
	ebiten.SetWindowSize(game.ScreenWidth*3, game.ScreenHeight*3)
	ebiten.SetWindowTitle("Asas de Valdoria v" + version)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("falha ao iniciar o jogo: %v", err)
	}
}
