package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	gameinit "github.com/yyyoichi/submarine-game/internal/app/screens/game/init"
)

type Game struct {
	ebiten.Game
}

func New() *Game {
	return &Game{Game: gameinit.New()}
}

// var gs = gameinit.NewSec()

// func (g *Game) Update() error {
// if g.count%120 == 1 {
// 	// gs.Update()
// 	g.overlay = components.SimpleTrunOverlay(components.SimpleTrunOverlayConfig{
// 		TrunTextConfig: components.TrunTextConfig{IsMe: true},
// 	})
// 	g.overlay.Clear()
// }
// if g.count%120 == 61 {
// 	g.overlay = components.SimpleTrunOverlay(components.SimpleTrunOverlayConfig{
// 		TrunTextConfig: components.TrunTextConfig{IsMe: false},
// 	})
// g.overlay.Clear()
// }
// 	g.count++
// 	return nil
// }

// func (g *Game) Draw(screen *ebiten.Image) {
// 	for img, options := range components.BackgoundImage(float64(config.ScreenWidth), float64(config.ScreenHeight)) {
// 		screen.DrawImage(img, options)
// 	}
// 	m := components.OceanMap{
// 		Width:  float64(config.ScreenHeight) - 20,
// 		Height: float64(config.ScreenHeight) - 20,
// 		W:      6,
// 		H:      6,
// 		P:      8,
// 	}

// 	m.ClearOceanMap(2, g.count/60%(6*6))
// 	op := &ebiten.DrawImageOptions{}
// 	op.GeoM.Translate(10, 10)
// 	screen.DrawImage(m.Src, op)
// 	ebitenutil.DebugPrint(screen, fmt.Sprintf("count: %d, state: %s, cols: %v", g.count, g.state.String(), g.overlay.DrawImageOptions().ColorScale))

// 	screen.DrawImage(g.overlay.Image(), g.overlay.DrawImageOptions())
// 	g.Game.Draw(screen)
// }

// func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return config.ScreenWidth, config.ScreenHeight
// }
