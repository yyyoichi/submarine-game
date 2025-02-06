package images

import (
	"embed"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed src/*
var src embed.FS

var (
	MaptileUmi01        *ebiten.Image
	MaptileRangaBlack02 *ebiten.Image
	MyLocation          *ebiten.Image
)

func new(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFileSystem(src, path)
	if err != nil {
		panic(err)
	}
	return img
}

func init() {
	MaptileUmi01 = new("src/maptile_umi_01.png")
	MaptileRangaBlack02 = new("src/maptile_renga_black_02.png")
	// https://fonts.google.com/icons?selected=Material+Symbols+Outlined:my_location:FILL@0;wght@400;GRAD@0;opsz@24&icon.query=location&icon.size=24&icon.color=%23e8eaed
	MyLocation = new("src/my_location_24dp_E8EAED_FILL0_wght400_GRAD0_opsz24.png")
}

type CenterImage struct {
	Size int     // 一片の長さ。正方形前提
	P    float64 // 余白
	*ebiten.Image
}

func (ci *CenterImage) DrawImageOptions() *ebiten.DrawImageOptions {
	size := float64(ci.Size)
	b := ci.Image.Bounds()
	x, y := float64(b.Dx()), float64(b.Dy())
	sx, sy := size/x, size/y

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(ci.P/2, ci.P/2)
	op.GeoM.Scale(sx, sy)
	return op
}
