package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"valdoria/game"
)

func main() {
	g := game.New()

	ebiten.SetWindowSize(game.ScreenWidth*3, game.ScreenHeight*3)
	ebiten.SetWindowTitle("Valdoria")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
