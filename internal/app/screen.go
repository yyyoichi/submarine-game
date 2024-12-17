package app

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yyyoichi/submarine-game/internal/app/components"
	"github.com/yyyoichi/submarine-game/internal/app/state"
)

const (
	screenWidth  = 430
	screenHeight = 930
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
	for img, options := range components.BackgoundImage(screenWidth, screenHeight) {
		screen.DrawImage(img, options)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func Draw(screen *ebiten.Image) {
	for img, options := range components.BackgoundImage(screenWidth, screenHeight) {
		screen.DrawImage(img, options)
	}
}
