package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type FullScreen struct {
	Width, Height float64 //
	Fade          func() float32
	Src           *ebiten.Image
}

func (s *FullScreen) Clear() {
	s.Src = ebiten.NewImage(int(s.Width), int(s.Height))
}

func (s *FullScreen) Templete() *ebiten.Image {
	img := ebiten.NewImage(int(s.Width), int(s.Height))
	img.Fill(color.Black)
	return img
}

func (s *FullScreen) Draw(img *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.SetA(s.Fade())
	s.Src.DrawImage(img, op)
}
