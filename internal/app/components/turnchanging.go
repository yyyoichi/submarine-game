package components

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yyyoichi/submarine-game/internal/app/animation"
	"github.com/yyyoichi/submarine-game/internal/app/config"
	"github.com/yyyoichi/submarine-game/internal/app/fonts"
)

var fade = &animation.Animation{
	ByPercentage: [][2]float32{{0, 0}, {0.1, 1}, {0.8, 1}, {1, 0}},
	Duration:     time.Duration(time.Millisecond * 1200),
}
var textAnime = &animation.Animation{
	ByPercentage: [][2]float32{{0, -10}, {1, 0}},
	Duration:     time.Duration(time.Millisecond * 300),
}

type TrunChanging struct {
	isMe bool
}

func (tc *TrunChanging) Clear(isMe bool) {
	fade.Clear()
	textAnime.Clear()
	tc.isMe = isMe
}

func (tc *TrunChanging) Image() *ebiten.Image {
	sf := FullScreen{
		Width:  float64(config.ScreenWidth),
		Height: float64(config.ScreenHeight),
		Fade:   fade.Value,
	}
	sf.Clear()
	img := sf.Templete()
	{
		var txtString = "あなたのターン"
		if !tc.isMe {
			txtString = "相手のターン"
		}
		fontSize := 48
		op := &text.DrawOptions{}
		op.GeoM.Translate(20, float64(textAnime.Value()+float32(config.ScreenWidth/2-fontSize)))
		op.ColorScale.ScaleWithColor(color.NRGBA{255, 255, 255, uint8(255 * fade.Value())})
		text.Draw(img, txtString, &text.GoTextFace{
			Source: fonts.TrainOneRegularSource,
			Size:   float64(fontSize),
		}, op)
		sf.Draw(img)
	}
	return sf.Src
}

func (tc *TrunChanging) Value() string {
	return fmt.Sprintf("TurnChanging{isMe: %v, fade: %s, textAnim}", tc.isMe, fade.String(), textAnime.String())
}
