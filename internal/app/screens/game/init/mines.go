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

type initMinesConfig struct {
	w, h               int // ボード数
	islandSectors      []int
	enableCursor       func(int, core.Direction) (int, bool)
	pushSelectedSector func(int) int // 機雷を敷設設定をして総数を返す
	popSelectedSector  func() *int   // 機雷設定の最後の一つを覗いてそれを返す。
	selectedSectors    *[]int        // 機雷設定
	mineCount          int           // 機雷設定の総数
}

type minesScreen struct {
	// views
	overlay *components.Overlay
	cw      *components.CommandWindow
	ocean   *components.Ocean

	// local state
	cursol    *int
	setCursol state.SetStateFunc[int]

	config initMinesConfig
}

func newMinesScreen(c initMinesConfig) *minesScreen {
	overlay := components.SimpleLeadOverlay(components.SimpleLeadOverlayConfig{
		LeadText: "制御機雷の設置",
	})
	cw := components.DefaultCommandWindow(components.CommandWindowConfig{
		Commands: []components.Text{{
			Text:  "機雷設置",
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
	red := components.DefaultSectorImage(components.SectorImageConfig{
		Width:  w,
		Height: h,
		Color:  components.Red,
	})
	red.Clear()

	for sector := range c.w * c.h {
		if slices.Contains(c.islandSectors, sector) {
			continue
		}
		ocean.SectorImages[sector] = red
	}

	var initCursol int
	if l := len(*c.selectedSectors); l > 0 {
		initCursol = (*c.selectedSectors)[l-1]
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

	return &minesScreen{
		overlay: overlay,
		cw:      cw,
		ocean:   ocean,

		cursol:    cursol,
		setCursol: setCursol,

		config: c,
	}
}

func (s *minesScreen) Update() error {
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
			if len(*s.config.selectedSectors) == s.config.mineCount {
				// すべて選択済みなら
				// 最後の一つ取り除く
				s.config.popSelectedSector()
			}
			s.setCursol(i)
		}
	}

	// 選択
	if actions.IsKeyAJustPressed() {
		s.config.pushSelectedSector(*s.cursol)
		nextCursor := func() *int {
			for _, sector := range []core.Direction{core.North, core.East, core.South, core.West} {
				if i, ok := s.config.enableCursor(*s.cursol, sector); ok {
					return &i
				}
			}
			for {
				i := rand.IntN(s.config.w * s.config.h)
				if !slices.Contains(s.config.islandSectors, i) {
					return &i
				}
			}
		}()

		s.setCursol(*nextCursor)
	}
	// 取り消し
	if actions.IsKeyBJustPressed() {
		s.config.popSelectedSector()
	}
	return nil
}

func (s *minesScreen) Draw(screen *ebiten.Image) {
	bg := ebiten.NewImage(config.ScreenWidth, config.ScreenHeight)
	bg.Fill(color.RGBA{28, 25, 37, 255})
	screen.DrawImage(bg, nil)

	s.ocean.ClearOverSectorImage()
	if s.cursol != nil {
		_, w, _ := s.ocean.SectorImage()
		s.ocean.OverSectorImage(*s.cursol, &images.CenterImage{Size: w, P: float64(w) * 0.2, Img: images.Mine})
	}
	for _, sector := range *s.config.selectedSectors {
		_, w, _ := s.ocean.SectorImage()
		s.ocean.OverSectorImage(sector, &images.CenterImage{Size: w, P: float64(w) * 0.2, Img: images.Mine})
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(s.ocean.DrawImageGeoM())
	screen.DrawImage(s.ocean.Image(), op)

	screen.DrawImage(s.cw.Image(), s.cw.DrawImageOptions())

	screen.DrawImage(s.overlay.Image(), s.overlay.DrawImageOptions())
}

func (s *minesScreen) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
