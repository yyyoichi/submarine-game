package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yyyoichi/submarine-game/internal/app/components"
)

const (
	screenWidth  = 430
	screenHeight = 930
)

type Game struct {
}

func (g *Game) Update() error {
	// resp, err := http.Get("/api/echo")
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()
	// if resp.StatusCode != http.StatusOK {
	// 	return fmt.Errorf("unexpected status: %s", resp.Status)
	// }
	// var buf bytes.Buffer
	// buf.ReadFrom(resp.Body)
	// _, set := state.UseGlobalState("word", state.WithInitialValue("..."))
	// set(buf.String())
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for img, options := range components.BackgoundImage(screenWidth, screenHeight) {
		screen.DrawImage(img, options)
	}
	m := components.OceanMap{
		Width:  screenWidth - 20,
		Height: screenWidth - 20,
		W:      6,
		H:      6,
		P:      8,
	}
	m.ClearOceanMap(2, 10)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(10, 200)
	screen.DrawImage(m.Src, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func Draw(screen *ebiten.Image) {
	for img, options := range components.BackgoundImage(screenWidth, screenHeight) {
		screen.DrawImage(img, options)
	}
}
