package gameinit

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yyyoichi/submarine-game/internal/app/components"
	"github.com/yyyoichi/submarine-game/internal/app/config"
)

type sectorScreen struct {
	demo *components.Overlay
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
	o.Clear()
	return &sectorScreen{
		demo: o,
	}
}

func (s *sectorScreen) Update() error {
	return nil
}

func (s *sectorScreen) Draw(screen *ebiten.Image) {
	img := ebiten.NewImage(100, 100)
	img.Fill(color.RGBA{0, 0, 0, 0xff})
	op := s.demo.DrawImageOptions()
	screen.DrawImage(s.demo.Image(), s.demo.DrawImageOptions())
	screen.DrawImage(img, nil)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("op: %v", op.ColorScale))
}

func (s *sectorScreen) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
