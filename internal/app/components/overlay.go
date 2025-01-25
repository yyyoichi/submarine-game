package components

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yyyoichi/submarine-game/internal/app/animation"
	"github.com/yyyoichi/submarine-game/internal/app/config"
	"github.com/yyyoichi/submarine-game/internal/app/fonts"
)

type (
	SimpleTrunOverlayConfig struct {
		TrunTextConfig
	}
	SimpleLeadOverlayConfig struct {
		LeadText string
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
	SimpleLeadOverlay = func(config SimpleLeadOverlayConfig) *Overlay {
		o := &Overlay{
			CenterText: &Text{
				Text:  config.LeadText,
				Font:  fonts.DotGothic16RegularSource,
				Size:  48,
				Color: color.NRGBA{R: 100, G: 100, B: 0, A: 255},
				Animeation: TextAmimation{
					TranslatesX: &animation.Animation{
						ByPercentage: [][2]float32{{0, -10}, {1, 0}},
						Duration:     time.Duration(time.Millisecond * 300),
					},
					TranslatesY: &animation.Animation{
						ByPercentage: [][2]float32{{0, 2}, {1, 0}},
						Duration:     time.Duration(time.Millisecond * 300),
					},
					Fade: animation.DefaultFade(),
				},
			},
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
		CenterText      *Text
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
	if o.CenterText != nil {
		o.CenterText.Clear()
	}
}

func (o *Overlay) Zero() {
	if o.Animation.Fade != nil {
		o.Animation.Fade.Zero()
	}
	if o.RightButtomText != nil {
		o.RightButtomText.Zero()
	}
	if o.CenterText != nil {
		o.CenterText.Zero()
	}
}

func (o *Overlay) Image() *ebiten.Image {
	img := ebiten.NewImage(config.ScreenWidth, config.ScreenHeight)
	img.Fill(color.RGBA{0xff, 0, 0, 0xff})

	o.drawRBText(img)
	o.drawCText(img)
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

func (o *Overlay) drawCText(img *ebiten.Image) {
	if o.CenterText == nil {
		return
	}
	// 左右均等、上下中央に配置
	op := o.CenterText.DrawOptions()
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignCenter
	// 気持ち上に描画
	op.GeoM.Translate(float64(config.ScreenWidth)/2, float64(config.ScreenHeight)*0.95/2)
	text.Draw(img, o.CenterText.Text, o.CenterText.GoTextFace(), op)
}

func (o *Overlay) DrawImageOptions() *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	if o.Animation.Fade != nil {
		op.ColorScale.ScaleAlpha(o.Animation.Fade.Value())
	}
	return op
}
