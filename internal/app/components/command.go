package components

import (
	"image/color"
	"time"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yyyoichi/submarine-game/internal/app/animation"
	"github.com/yyyoichi/submarine-game/internal/app/config"
	"github.com/yyyoichi/submarine-game/internal/app/fonts"
)

var (
	MineCommandText = func() Text {
		return Text{
			Text:  "機雷",
			Color: Red,
			Font:  fonts.DotGothic16RegularSource,
		}
	}
	TorpedoCommandText = func() Text {
		return Text{
			Text:  "魚雷",
			Color: Orange,
			Font:  fonts.DotGothic16RegularSource,
		}
	}
	MoveCommandText = func() Text {
		return Text{
			Text:  "潜航",
			Color: Green,
			Font:  fonts.DotGothic16RegularSource,
		}
	}
)

type (
	CommandWindowConfig struct {
		Commands []Text
	}
)

var (
	DefaultCommandWindow = func(c CommandWindowConfig) *CommandWindow {
		cw := &CommandWindow{
			Commands:               c.Commands,
			Width:                  config.ScreenWidth / 6,
			Height:                 config.ScreenHeight / 10,
			SelfTranslateX:         float64(config.ScreenWidth) * 15 / 20,
			SelfTranslateY:         float64(config.ScreenHeight) * (1 - 0.6) / 2,
			commandWindowViewPoint: newCommandWindowViewPoint(len(c.Commands)),
		}
		for i, command := range cw.Commands {
			switch l := utf8.RuneCountInString(command.Text); l {
			case 1, 2:
				cw.Commands[i].Size = 55
			case 3:
				cw.Commands[i].Size = 45
			case 4:
				cw.Commands[i].Size = 35
			default:
				cw.Commands[i].Size = 24
			}
		}
		return cw
	}
)

type CommandWindow struct {
	Commands                       []Text
	Width                          int
	Height                         int
	SelfTranslateX, SelfTranslateY float64 // 自身を描画するときの位置
	commandWindowViewPoint
}

func NewCommandWindow(texts ...string) *CommandWindow {
	commands := make([]Text, len(texts))
	for i, t := range texts {
		commands[i] = Text{Text: t, Font: fonts.DotGothic16RegularSource, Size: 24 * 3, Color: color.NRGBA{R: 0, G: 0, B: 0, A: 255}}
	}
	commands = append(commands, Text{
		Text:  "機雷",
		Size:  24 * 3,
		Color: color.NRGBA{226, 123, 10, 255},
		Font:  fonts.DotGothic16RegularSource,
	})
	return &CommandWindow{
		Commands:               commands,
		Width:                  config.ScreenWidth / 6,
		Height:                 config.ScreenHeight / 10,
		SelfTranslateX:         float64(config.ScreenWidth) * 15 / 20,
		SelfTranslateY:         float64(config.ScreenHeight) * (1 - 0.6) / 2,
		commandWindowViewPoint: newCommandWindowViewPoint(len(commands)),
	}
}

func (cw *CommandWindow) Clear() {
	for i := range cw.Commands {
		cw.Commands[i].Clear()
	}
}
func (cw *CommandWindow) Zero() {
	for i := range cw.Commands {
		cw.Commands[i].Zero()
	}
}
func (cw *CommandWindow) Image() *ebiten.Image {
	if cw.max == 0 {
		cw.max = len(cw.Commands)
	}

	parent := ebiten.NewImage(cw.Width, cw.Height*len(cw.Commands))
	for i, cmd := range cw.Commands {
		cimg := ebiten.NewImage(cw.Width, cw.Height)
		cimg.Fill(color.White)
		op := cmd.DrawOptions()
		op.PrimaryAlign = text.AlignCenter
		op.SecondaryAlign = text.AlignCenter
		// 気持ち上に描画
		op.GeoM.Translate(float64(cw.Width)/2, float64(cw.Height)*0.95/2)
		text.Draw(cimg, cmd.Text, cmd.GoTextFace(), op)
		// parentに描画
		iop := &ebiten.DrawImageOptions{}
		iop.GeoM.Translate(0, float64(cw.Height*i))
		parent.DrawImage(cimg, iop)
	}

	// parentをimgに描画
	img := ebiten.NewImage(cw.Width, cw.Height)
	img.Fill(color.Black)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(cw.translates(cw.Height))
	img.DrawImage(parent, op)
	// 枠に薄い黒を描画する

	// imgに最大maxShadow(<255)の濃さ影を描画する。中心に向けて徐々に薄くなる。
	// point ... 1:上 2:右 4:下 8:左
	drawShadow := func(img *ebiten.Image, rate float64, maxShadow float64, point int8) {
		// 影が落ちるサイズ
		bounds := img.Bounds()
		h, w := float64(bounds.Dy()), float64(bounds.Dx())

		top, bottom, rigth, left := point&1 == 1, point&4 == 4, point&2 == 2, point&8 == 8
		if top || bottom {
			scale := h * rate // 影が落ちるサイズ
			for i := range int(scale) {
				// 1ずつ影を落とす
				fi := float64(i)
				im := ebiten.NewImage(int(w), 1)
				im.Fill(color.NRGBA{R: 0, G: 0, B: 0, A: uint8((scale - fi) / scale * maxShadow)})

				if top {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(0, fi)
					img.DrawImage(im, op)
				}
				if bottom {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(0, h-fi)
					img.DrawImage(im, op)
				}
			}
		}
		if rigth || left {
			scale := w * rate // 影が落ちるサイズ
			for i := range int(scale) {
				// 1ずつ影を落とす
				fi := float64(i)
				im := ebiten.NewImage(1, int(h))
				im.Fill(color.NRGBA{R: 0, G: 0, B: 0, A: uint8((scale - fi) / scale * maxShadow)})

				if rigth {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(w-fi, 0)
					img.DrawImage(im, op)
				}
				if left {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(fi, 0)
					img.DrawImage(im, op)
				}
			}
		}
	}
	drawShadow(img, 0.35, 200, 5)  // 上下
	drawShadow(img, 0.02, 120, 10) // 左右
	return img
}

func (cw *CommandWindow) DrawImageOptions() *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(cw.SelfTranslateX, cw.SelfTranslateY)
	return op
}

type commandWindowViewPoint struct {
	max             int
	point           int
	scrollAnimation animation.Animation
}

func newCommandWindowViewPoint(l int) commandWindowViewPoint {
	return commandWindowViewPoint{
		max:   l,
		point: 0,
	}
}

func (p *commandWindowViewPoint) Up() {
	if p.max == p.point+1 {
		return
	}
	p.point++
	p.scrollAnimation = animation.Animation{
		ByPercentage: [][2]float32{{0, 1}, {0.9, 0}, {0.95, -0.1}, {1, 0}},
		Duration:     time.Duration(time.Millisecond * 500),
	}
	p.scrollAnimation.Clear()
}
func (p *commandWindowViewPoint) Down() {
	if p.point == 0 {
		return
	}
	p.point--
	p.scrollAnimation = animation.Animation{
		ByPercentage: [][2]float32{{0, -1}, {0.9, 0}, {0.95, 0.1}, {1, 0}},
		Duration:     time.Duration(time.Millisecond * 500),
	}
	p.scrollAnimation.Clear()
}

func (p *commandWindowViewPoint) translates(height int) (tx float64, ty float64) {
	ty = -float64(height*p.point) + float64(p.scrollAnimation.Value()*float32(height))
	return
}
