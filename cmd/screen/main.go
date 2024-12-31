package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yyyoichi/submarine-game/internal/app"
)

func main() {
	ebiten.SetWindowTitle("Hello, World!!")
	ebiten.SetTPS(60)
	if err := ebiten.RunGame(app.New()); err != nil {
		log.Fatal(err)
	}
}
