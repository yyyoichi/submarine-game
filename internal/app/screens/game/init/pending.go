package gameinit

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yyyoichi/submarine-game/internal/app/components"
	"github.com/yyyoichi/submarine-game/internal/app/config"
)

type pendingScreen struct {
	demo *components.Overlay
}

func newPendingScreen() *pendingScreen {
	s := components.SimpleLeadOverlay(components.SimpleLeadOverlayConfig{
		LeadText: "出撃...",
	})
	s.Clear()
	return &pendingScreen{
		demo: s,
	}
}

func (s *pendingScreen) Update() error {
	return nil
}

func (s *pendingScreen) Draw(screen *ebiten.Image) {
	screen.DrawImage(s.demo.Image(), s.demo.DrawImageOptions())
}

func (s *pendingScreen) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
