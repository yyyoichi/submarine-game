package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yyyoichi/submarine-game/internal/app/config"
)

var (
	DefaultOcean = func(c OceanConfig) *Ocean {
		r := 0.6
		o := &Ocean{
			Width:         float64(config.ScreenHeight) * r,
			Height:        float64(config.ScreenHeight) * r,
			W:             6,
			H:             6,
			P:             5,
			IslandSectors: c.IslandSectors,
			BorderColor:   color.Black,
		}
		o.init()
		o.IslandImage = ebiten.NewImage(int(o.fillW), int(o.fillH))
		// 限りなく黒に近い緑
		o.IslandImage.Fill(color.Black)
		// o.IslandImage.Fill(color.NRGBA64{R: 0x00, G: 0x20, B: 0x00, A: 0xff})

		o.SelfTranslateY = float64(config.ScreenHeight) * (1 - r) / 2
		o.SelfTranslateX = o.SelfTranslateY + float64((config.ScreenWidth-config.ScreenHeight)/2)
		return o
	}
)

type OceanConfig struct {
	IslandSectors []int
}

type Ocean struct {
	Width, Height float64 // 幅、高さ
	W, H          int     // 横、縦の分割数
	P             float64 // 各セクターの余白

	IslandSectors []int
	IslandImage   *ebiten.Image

	BorderColor color.Color   // 枠線の色
	BorderImage *ebiten.Image // 枠線の画像

	SelfTranslateX, SelfTranslateY float64 // 自身を描画するときの位置

	sectW, sectH float64 // セクタの幅、高さ
	fillW, fillH float64 // セクタのFillの幅、高さ
}

func (o *Ocean) Image() *ebiten.Image {
	img := ebiten.NewImage(int(o.Width), int(o.Height))
	o.drawBorder(img)
	o.drawIsland(img)
	return img
}

func (o *Ocean) DrawImageGeoM() ebiten.GeoM {
	gm := ebiten.GeoM{}
	gm.Translate(o.SelfTranslateX, o.SelfTranslateY)
	return gm
}

func (o *Ocean) drawBorder(src *ebiten.Image) {
	img := ebiten.NewImage(int(o.Width), int(o.Height))
	img.Fill(o.BorderColor)
	if o.BorderImage != nil {
		for i, op := range FillWithImage(o.Width, o.Height, o.BorderImage) {
			img.DrawImage(i, op)
		}
	}

	// 各セクターごとにBlendして穴を開ける。
	// 重なった部分を透過するだけなので何色でもいい。
	white := ebiten.NewImage(int(o.fillW), int(o.fillH))
	white.Fill(color.White)
	for sector := 0; sector < o.W*o.H; sector++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(o.translates(sector))
		op.Blend = ebiten.BlendDestinationOut
		img.DrawImage(white, op)
	}
	src.DrawImage(img, nil)
}

func (o *Ocean) drawIsland(src *ebiten.Image) {
	if o.IslandImage == nil {
		return
	}
	for _, sector := range o.IslandSectors {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(o.translates(sector))
		src.DrawImage(o.IslandImage, op)
	}
}

func (o *Ocean) translates(sector int) (tx float64, ty float64) {
	i, j := sector%o.W, sector/o.W // 位置
	x, y := o.sectW*float64(i), o.sectH*float64(j)
	// 全辺分のパディング1/2+セクター分のパティング1/2 -> 隣接したティングが合わさってm.Pになる
	tx, ty = x+o.P, y+o.P
	return
}

func (o *Ocean) init() {
	//全辺のパディング1/2を除いたのがセクタのサイズ
	o.sectW, o.sectH = (o.Width-o.P)/float64(o.W), (o.Height-o.P)/float64(o.H)
	o.fillW, o.fillH = o.sectW-o.P, o.sectH-o.P
}
