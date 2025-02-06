package gameinit

import (
	"image/color"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yyyoichi/submarine-game/internal/app/components"
	"github.com/yyyoichi/submarine-game/internal/app/config"
	"github.com/yyyoichi/submarine-game/internal/app/fonts"
)

type initSectorConfig struct {
	w, h              int // ボード数
	islandSectors     []int
	setSelectedSector func(int)
	selectedSector    *int
}

type sectorScreen struct {
	overlay *components.Overlay
	cw      *components.CommandWindow
	ocean   *components.Ocean
}

func newSectorScreen(c initSectorConfig) *sectorScreen {
	overlay := components.SimpleLeadOverlay(components.SimpleLeadOverlayConfig{
		LeadText: "作戦開始海域の選定",
	})
	cw := components.DefaultCommandWindow(components.CommandWindowConfig{
		Commands: []components.Text{{
			Text:  "作戦開始",
			Color: components.Black,
			Font:  fonts.DotGothic16RegularSource,
		}},
	})
	ocean := components.DefaultOcean(components.OceanConfig{
		IslandSectors: c.islandSectors,
	})
	cw.Clear()
	overlay.Clear()
	_, w, h := ocean.SectorImage()
	green := components.DefaultSectorImage(components.SectorImageConfig{
		Width:  w,
		Height: h,
		Color:  components.Green,
	})
	green.Clear()

	for sector := range c.w * c.h {
		if slices.Contains(c.islandSectors, sector) {
			continue
		}
		ocean.SectorImages[sector] = green
	}

	return &sectorScreen{
		overlay: overlay,
		cw:      cw,
		ocean:   ocean,
	}
}

func (s *sectorScreen) Update() error {
	return nil
}

func (s *sectorScreen) Draw(screen *ebiten.Image) {

	bg := ebiten.NewImage(config.ScreenWidth, config.ScreenHeight)
	bg.Fill(color.RGBA{28, 25, 37, 255})
	screen.DrawImage(bg, nil)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(s.ocean.DrawImageGeoM())
	screen.DrawImage(s.ocean.Image(), op)

	screen.DrawImage(s.cw.Image(), s.cw.DrawImageOptions())

	screen.DrawImage(s.overlay.Image(), s.overlay.DrawImageOptions())
}

func (s *sectorScreen) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
