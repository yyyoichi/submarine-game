package gameinit

import (
	"image/color"
	"math/rand/v2"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yyyoichi/submarine-game/internal/app/components"
	"github.com/yyyoichi/submarine-game/internal/app/config"
	"github.com/yyyoichi/submarine-game/internal/app/fonts"
	"github.com/yyyoichi/submarine-game/internal/app/images"
	"github.com/yyyoichi/submarine-game/internal/app/state"
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

	cursol    *int
	setCursol state.SetStateFunc[int]
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

	var initCursol int
	for {
		i := rand.IntN(c.w * c.h)
		if !slices.Contains(c.islandSectors, i) {
			initCursol = i
			break
		}
	}

	cursol, setCursol := state.UseState(state.WithInitialValue(initCursol))

	return &sectorScreen{
		overlay: overlay,
		cw:      cw,
		ocean:   ocean,

		cursol:    cursol,
		setCursol: setCursol,
	}
}

func (s *sectorScreen) Update() error {
	return nil
}

func (s *sectorScreen) Draw(screen *ebiten.Image) {

	bg := ebiten.NewImage(config.ScreenWidth, config.ScreenHeight)
	bg.Fill(color.RGBA{28, 25, 37, 255})
	screen.DrawImage(bg, nil)

	if s.cursol != nil {
		_, w, _ := s.ocean.SectorImage()
		s.ocean.OverSectorImage(*s.cursol, &images.CenterImage{Size: w, P: float64(w) * 0.2, Img: images.MyLocation})
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(s.ocean.DrawImageGeoM())
	screen.DrawImage(s.ocean.Image(), op)

	screen.DrawImage(s.cw.Image(), s.cw.DrawImageOptions())

	screen.DrawImage(s.overlay.Image(), s.overlay.DrawImageOptions())
}

func (s *sectorScreen) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
