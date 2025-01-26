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
}
