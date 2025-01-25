package components

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yyyoichi/submarine-game/internal/app/animation"
	"github.com/yyyoichi/submarine-game/internal/app/fonts"
)

var (
	commandWindowWidth  = 250
	commandWindowHeight = 80
)

type CommandWindow struct {
	Commands []Text
	commandWindowViewPoint
}

func NewCommandWindow(texts ...string) *CommandWindow {
	commands := make([]Text, len(texts))
	for i, t := range texts {
		commands[i] = Text{Text: t, Font: fonts.DotGothic16RegularSource, Size: 24 * 3, Color: color.NRGBA{R: 0, G: 0, B: 0, A: 255}}
	}
	return &CommandWindow{
		Commands:               commands,
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

	parent := ebiten.NewImage(commandWindowWidth, commandWindowHeight*len(cw.Commands))
	for i, cmd := range cw.Commands {
		cimg := ebiten.NewImage(commandWindowWidth, commandWindowHeight)
		cimg.Fill(color.White)
		op := cmd.DrawOptions()
		op.PrimaryAlign = text.AlignCenter
		op.SecondaryAlign = text.AlignCenter
		// 気持ち上に描画
		op.GeoM.Translate(float64(commandWindowWidth)/2, float64(commandWindowHeight)*0.95/2)
		text.Draw(cimg, cmd.Text, cmd.GoTextFace(), op)
		// parentに描画
		iop := &ebiten.DrawImageOptions{}
		iop.GeoM.Translate(0, float64(commandWindowHeight*i))
		parent.DrawImage(cimg, iop)
	}

	// parentをimgに描画
	img := ebiten.NewImage(commandWindowWidth, commandWindowHeight)
	img.Fill(color.Black)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(cw.translates())
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

func (p *commandWindowViewPoint) translates() (tx float64, ty float64) {
	ty = -float64(commandWindowHeight*p.point) + float64(p.scrollAnimation.Value()*float32(commandWindowHeight))
	return
}
