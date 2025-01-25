package components

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yyyoichi/submarine-game/internal/app/animation"
	"github.com/yyyoichi/submarine-game/internal/app/fonts"
)

type (
	TrunTextConfig struct {
		IsMe bool
	}
)

var (
	TurnText = func(config TrunTextConfig) *Text {
		txt := &Text{
			Font:  fonts.TrainOneRegularSource,
			Text:  "行動開始",
			Size:  48,
			Color: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
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
		}
		if !config.IsMe {
			txt.Text = "敵艦の行動"
		}
		return txt
	}
)

type (
	Text struct {
		Font       *text.GoTextFaceSource
		Text       string
		Size       float64
		Color      color.NRGBA
		Animeation TextAmimation
	}
	TextAmimation struct {
		TranslatesX *animation.Animation
		TranslatesY *animation.Animation
		Fade        *animation.Animation
	}
)

func (t *Text) Clear() {
	if t.Animeation.TranslatesX != nil {
		t.Animeation.TranslatesX.Clear()
	}
	if t.Animeation.TranslatesY != nil {
		t.Animeation.TranslatesY.Clear()
	}
	if t.Animeation.Fade != nil {
		t.Animeation.Fade.Clear()
	}
}

func (t *Text) Zero() {
	if t.Animeation.TranslatesX != nil {
		t.Animeation.TranslatesX.Zero()
	}
	if t.Animeation.TranslatesY != nil {
		t.Animeation.TranslatesY.Zero()
	}
	if t.Animeation.Fade != nil {
		t.Animeation.Fade.Zero()
	}
}

func (t *Text) GoTextFace() *text.GoTextFace {
	return &text.GoTextFace{
		Source: t.Font,
		Size:   t.Size,
	}
}
func (t *Text) DrawOptions() *text.DrawOptions {
	op := &text.DrawOptions{}
	op.GeoM.Translate(t.translates())
	op.ColorScale.ScaleWithColor(t.color())
	return op
}

func (t *Text) translates() (tx float64, ty float64) {
	if t.Animeation.TranslatesX != nil {
		tx = float64(t.Animeation.TranslatesX.Value())
	}
	if t.Animeation.TranslatesY != nil {
		ty = float64(t.Animeation.TranslatesY.Value())
	}
	return
}

func (t *Text) color() color.NRGBA {
	return color.NRGBA{R: t.Color.R, G: t.Color.G, B: t.Color.B, A: uint8(float32(t.Color.A) * t.Animeation.Fade.Value())}
}
