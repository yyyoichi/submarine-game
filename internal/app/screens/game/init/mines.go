package gameinit

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yyyoichi/submarine-game/internal/app/components"
	"github.com/yyyoichi/submarine-game/internal/app/config"
)

type minesScreen struct {
	demo *components.Overlay
}

func newMinesScreen() *minesScreen {
	s := components.SimpleLeadOverlay(components.SimpleLeadOverlayConfig{
		LeadText: "制御機雷の設置",
	})
	s.Clear()
	return &minesScreen{
		demo: s,
	}
}

func (s *minesScreen) Update() error {
	return nil
}

func (s *minesScreen) Draw(screen *ebiten.Image) {
	screen.DrawImage(s.demo.Image(), s.demo.DrawImageOptions())
}

func (s *minesScreen) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
