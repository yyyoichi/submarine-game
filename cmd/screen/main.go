package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"net/http"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yyyoichi/submarine-game/internal/app/state"
)

type Game struct {
}

func (g *Game) Update() error {
	resp, err := http.Get("/api/echo")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	_, set := state.UseGlobalState("word", state.WithInitialValue("..."))
	set(buf.String())
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	v, _ := state.UseGlobalState("word", state.WithInitialValue("..."))
	screen.Fill(color.RGBA{0xff, 0, 0, 0xff})
	vector.StrokeRect(screen, 0, 0, 100, 100, 1, color.RGBA{0, 0xff, 0, 0xff}, true)
	ebitenutil.DebugPrint(screen, *v)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowTitle("Hello, World!!")
	ebiten.SetTPS(10)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
