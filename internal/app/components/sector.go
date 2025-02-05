package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yyyoichi/submarine-game/internal/app/animation"
)

type SectorImageConfig struct {
	Width  int
	Height int
	Color  color.Color
}

var (
	DefaultSectorImage = func(c SectorImageConfig) *SectorImage {
		si := SectorImage{
			Width:  c.Width,
			Height: c.Height,
			Color:  c.Color,
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
