package gameinit

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yyyoichi/submarine-game/internal/app/components"
	"github.com/yyyoichi/submarine-game/internal/app/config"
)

type sectorScreen struct {
	demo  *components.Overlay
	cw    *components.CommandWindow
	ocean *components.Ocean

	count int
}

func newSectorScreen() *sectorScreen {
	o := components.SimpleLeadOverlay(components.SimpleLeadOverlayConfig{
		LeadText: "作戦開始海域の選定",
	})
	cw := components.DefaultCommandWindow(components.CommandWindowConfig{
		Commands: []components.Text{components.MineCommandText(), components.TorpedoCommandText(), components.MoveCommandText()},
	})
	// cw = components.NewCommandWindow("機雷", "魚雷")
	ocean := components.DefaultOcean(components.OceanConfig{
		IslandSectors: []int{9, 31},
	})
	cw.Clear()
	o.Clear()
	_, w, h := ocean.SectorImage()
	orange := components.DefaultSectorImage(components.SectorImageConfig{
		Width:  w,
		Height: h,
		Color:  components.Orange,
	})
	red := components.DefaultSectorImage(components.SectorImageConfig{
		Width:  w,
		Height: h,
		Color:  components.Red,
	})
	green := components.DefaultSectorImage(components.SectorImageConfig{
		Width:  w,
		Height: h,
		Color:  components.Green,
	})
	orange.Clear()
	red.Clear()
	green.Clear()
	ocean.SectorImages[7] = orange
	ocean.SectorImages[12] = green
	ocean.SectorImages[14] = red
	ocean.SectorImages[19] = orange
	return &sectorScreen{
		demo:  o,
		cw:    cw,
		ocean: ocean,
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

	bg := ebiten.NewImage(config.ScreenWidth, config.ScreenHeight)
	bg.Fill(color.RGBA{28, 25, 37, 1})
	screen.DrawImage(bg, nil)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(s.ocean.DrawImageGeoM())
	screen.DrawImage(s.ocean.Image(), op)

	screen.DrawImage(s.cw.Image(), s.cw.DrawImageOptions())

	screen.DrawImage(s.demo.Image(), s.demo.DrawImageOptions())
}

func (s *sectorScreen) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
