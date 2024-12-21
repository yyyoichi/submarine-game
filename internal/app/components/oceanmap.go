package components

import (
	"image/color"
	"iter"

	"github.com/hajimehoshi/ebiten/v2"
)

type OceanMap struct {
	Width, Height float64       // 幅、高さ
	W, H          int           // 横、縦の分割数
	P             float64       // 各セクターの余白
	Src           *ebiten.Image // (Width, Height)の画像

	sectW, sectH float64 // セクタの幅、高さ
	fillW, fillH int     // セクタのFillの幅、高さ
}

func (m *OceanMap) ClearOceanMap() {
	//全辺のパディング1/2を除いたのがセクタのサイズ
	m.sectW, m.sectH = (m.Width-m.P)/float64(m.W), (m.Height-m.P)/float64(m.H)
	m.fillW, m.fillH = int(m.sectW-m.P), int(m.sectH-m.P)

	// scrを黒いブロック画像で埋める
	m.Src = ebiten.NewImage(int(m.Width), int(m.Height))
	for img, op := range getFillImageParams(m.Width, m.Height, path_maptile_ranga_black_02) {
		m.Src.DrawImage(img, op)
	}
	// 各セクターごとにBlendして穴を開ける。
	white := ebiten.NewImage(m.fillW, m.fillH)
	white.Fill(color.White) // 重なった部分を透過するだけなので何色でもいい。
	for sector := range m.generateSector() {
		op := &ebiten.DrawImageOptions{}
		m.applayTxTy(op, sector)
		op.Blend = ebiten.BlendDestinationOut
		m.Src.DrawImage(white, op)
	}
}

func (m *OceanMap) IsLand(sector int) {
	img := ebiten.NewImage(m.fillW, m.fillH)
	img.Fill(color.White)
	op := &ebiten.DrawImageOptions{}
	m.applayTxTy(op, sector)
	m.Src.DrawImage(img, op)
}

func (m *OceanMap) applayTxTy(op *ebiten.DrawImageOptions, sector int) {
	i, j := sector%m.W, sector/m.W
	x, y := m.sectW*float64(i), m.sectH*float64(j)
	// 全辺分のパディング1/2+セクター分のパティング1/2 -> 隣接したティングが合わさってm.Pになる
	op.GeoM.Translate(x+m.P, y+m.P)
}

func (m *OceanMap) generateSector() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; i < m.W*m.H; i++ {
			if !yield(i) {
				return
			}
		}
	}
}
