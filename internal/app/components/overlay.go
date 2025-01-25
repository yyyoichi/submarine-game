package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yyyoichi/submarine-game/internal/app/animation"
	"github.com/yyyoichi/submarine-game/internal/app/config"
)

type (
	SimpleTrunOverlayConfig struct {
		TrunTextConfig
	}
)

var (
	SimpleTrunOverlay = func(config SimpleTrunOverlayConfig) *Overlay {
		o := &Overlay{
			RightButtomText: TurnText(config.TrunTextConfig),
			Animation: OverlayAnimation{
				Fade: animation.DefaultFade(),
			},
		}
		return o
	}
)

type (
	Overlay struct {
		RightButtomText *Text
		Animation       OverlayAnimation
	}
	OverlayAnimation struct {
		Fade *animation.Animation
	}
)

func (o *Overlay) Clear() {
	if o.Animation.Fade != nil {
		o.Animation.Fade.Clear()
	}
	if o.RightButtomText != nil {
		o.RightButtomText.Clear()
	}
}

func (o *Overlay) Zero() {
	if o.Animation.Fade != nil {
		o.Animation.Fade.Zero()
	}
	if o.RightButtomText != nil {
		o.RightButtomText.Clear()
	}
}

func (o *Overlay) Image() *ebiten.Image {
	img := ebiten.NewImage(config.ScreenWidth, config.ScreenHeight)
	img.Fill(color.Black)

	o.drawRBText(img)

	return img
}

func (o *Overlay) drawRBText(img *ebiten.Image) {
	if o.RightButtomText == nil {
		return
	}
	baseGeoM := ebiten.GeoM{}
	// 左に一文字分、下に1/2文字分の余白
	baseGeoM.Translate(o.RightButtomText.Size, float64(config.ScreenHeight)-o.RightButtomText.Size*1.5)
	op := o.RightButtomText.DrawOptions()
	op.GeoM.Concat(baseGeoM)
	text.Draw(img, o.RightButtomText.Text, o.RightButtomText.GoTextFace(), op)
}

func (o *Overlay) DrawImageOptions() *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	if o.Animation.Fade != nil {
		op.ColorScale.SetA(o.Animation.Fade.Value())
	}
	return op
}
