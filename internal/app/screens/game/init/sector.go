package gameinit

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yyyoichi/submarine-game/internal/app/components"
	"github.com/yyyoichi/submarine-game/internal/app/config"
)

type sectorScreen struct {
	demo *components.Overlay
	cw   *components.CommandWindow

	count int
}

func NewSec() *sectorScreen {
	s := &sectorScreen{
		demo: components.SimpleLeadOverlay(components.SimpleLeadOverlayConfig{
			LeadText: "作戦開始海域の選定",
		}),
	}
	s.demo.Clear()
	return s
}

func newSectorScreen() *sectorScreen {
	o := components.SimpleLeadOverlay(components.SimpleLeadOverlayConfig{
		LeadText: "作戦開始海域の選定",
	})
	cw := components.NewCommandWindow("機雷", "魚雷", "潜航")
	cw.Clear()
	o.Clear()
	return &sectorScreen{
		demo: o,
		cw:   cw,
	}
}

func (s *sectorScreen) Update() error {
	s.count++
	if s.count%100 == 0 {
		s.cw.Up()
	}
	return nil
}

func (s *sectorScreen) Draw(screen *ebiten.Image) {
	screen.DrawImage(s.demo.Image(), s.demo.DrawImageOptions())

	screen.DrawImage(s.cw.Image(), nil)
}

func (s *sectorScreen) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
