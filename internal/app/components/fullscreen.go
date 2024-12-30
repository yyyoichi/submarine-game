package components

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type FullScreen struct {
	Width, Height float64 //
	Fade          interface {
		Opacity() float32
	}
	Src *ebiten.Image
}

func (s *FullScreen) Clear() {
	s.Src = ebiten.NewImage(int(s.Width), int(s.Height))
}

func (s *FullScreen) Draw(img *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.SetA(s.Fade.Opacity())
	s.Src.DrawImage(s.Src, op)
}
