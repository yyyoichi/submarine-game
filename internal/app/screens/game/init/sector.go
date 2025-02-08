package gameinit

import (
	"image/color"
	"math/rand/v2"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yyyoichi/submarine-game/internal/app/actions"
	"github.com/yyyoichi/submarine-game/internal/app/components"
	"github.com/yyyoichi/submarine-game/internal/app/config"
	"github.com/yyyoichi/submarine-game/internal/app/fonts"
	"github.com/yyyoichi/submarine-game/internal/app/images"
	"github.com/yyyoichi/submarine-game/internal/app/state"
	"github.com/yyyoichi/submarine-game/internal/core"
)

type initSectorConfig struct {
	w, h              int // ボード数
	islandSectors     []int
	enableCursor      func(int, core.Direction) (int, bool)
	setSelectedSector func(*int)
	selectedSector    *int
}

type sectorScreen struct {
	// views
	overlay *components.Overlay
	cw      *components.CommandWindow
	ocean   *components.Ocean

	// local state
	cursol    *int
	setCursol state.SetStateFunc[int]

	config initSectorConfig
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
	if c.selectedSector != nil {
		initCursol = *c.selectedSector
	} else {
		for {
			i := rand.IntN(c.w * c.h)
			if !slices.Contains(c.islandSectors, i) {
				initCursol = i
				break
			}
		}
	}

	cursol, setCursol := state.UseState(state.WithInitialValue(initCursol))

	return &sectorScreen{
		overlay: overlay,
		cw:      cw,
		ocean:   ocean,

		cursol:    cursol,
		setCursol: setCursol,

		config: c,
	}
}

func (s *sectorScreen) Update() error {
	var moveTo *core.Direction
	if actions.IsKeyUpJustPressed() {
		moveTo = &core.North
	}
	if actions.IsKeyDownJustPressed() {
		moveTo = &core.South
	}
	if actions.IsKeyRightJustPressed() {
		moveTo = &core.East
	}
	if actions.IsKeyLeftJustPressed() {
		moveTo = &core.West
	}
	if moveTo != nil {
		// 移動可能なら選択中を解除して、カーソルを動かす
		if i, ok := s.config.enableCursor(*s.cursol, *moveTo); ok {
			s.config.setSelectedSector(nil)
			s.setCursol(i)
		}
	}

	// 選択
	if actions.IsKeyAJustPressed() {
		s.config.setSelectedSector(s.cursol)
	}
	return nil
}

func (s *sectorScreen) Draw(screen *ebiten.Image) {

	bg := ebiten.NewImage(config.ScreenWidth, config.ScreenHeight)
	bg.Fill(color.RGBA{28, 25, 37, 255})
	screen.DrawImage(bg, nil)

	if s.cursol != nil {
		s.ocean.ClearOverSectorImage()
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
