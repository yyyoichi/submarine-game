package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yyyoichi/submarine-game/internal/app/animation"
)

type SectorImageConfig struct {
	Width  int
	Height int
}

var (
	OrangeSectorImage = func(c SectorImageConfig) *SectorImage {
		si := SectorImage{
			Width:  c.Width,
			Height: c.Height,
			// Orange
			Color: color.RGBA{226, 123, 10, 1},
			Animation: SectorAnimation{
				Flashing: animation.DefaultLoop(),
			},
		}
		return &si
	}
	RedSectorImage = func(c SectorImageConfig) *SectorImage {
		si := SectorImage{
			Width:  c.Width,
			Height: c.Height,
			// 濃い赤
			Color: color.RGBA{190, 40, 22, 1},
			Animation: SectorAnimation{
				Flashing: animation.DefaultLoop(),
			},
		}
		return &si
	}
	GreenSectorImage = func(c SectorImageConfig) *SectorImage {
		si := SectorImage{
			Width:  c.Width,
			Height: c.Height,
			// 濃い緑
			Color: color.RGBA{22, 150, 63, 1},
			Animation: SectorAnimation{
				Flashing: animation.DefaultLoop(),
			},
		}
		return &si
	}
)

type (
	SectorImage struct {
		Width     int
		Height    int
		Color     color.Color
		Animation SectorAnimation
	}
	SectorAnimation struct {
		Flashing *animation.LoopAnimation
	}
)

func (si *SectorImage) Image() *ebiten.Image {
	img := ebiten.NewImage(si.Width, si.Height)
	img.Fill(si.Color)
	return img
}

func (si *SectorImage) DrawImageOptions() *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	if si.Animation.Flashing != nil {
		op.ColorScale.ScaleAlpha(si.Animation.Flashing.Value())
	}
	return op
}

func (si *SectorImage) Clear() {
	if si.Animation.Flashing != nil {
		si.Animation.Flashing.Clear()
	}
}

func (si *SectorImage) Zero() {
	if si.Animation.Flashing != nil {
		si.Animation.Flashing.Zero()
	}
}
