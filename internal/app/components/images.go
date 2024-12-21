package components

import (
	"embed"
	_ "image/png"
	"iter"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed assets/*
var assets embed.FS
var (
	cellsize = 25.0
)

const (
	path_maptile_umi_01         = "assets/maptile_umi_01.png"
	path_maptile_ranga_black_02 = "assets/maptile_renga_black_02.png"
)

func getFillImageParams(width, height float64, path string) iter.Seq2[*ebiten.Image, *ebiten.DrawImageOptions] {
	img, _, err := ebitenutil.NewImageFromFileSystem(assets, path)
	b := img.Bounds()
	x, y := float64(b.Dx()), float64(b.Dy())
	sx, sy := cellsize/x, cellsize/y
	return func(yield func(*ebiten.Image, *ebiten.DrawImageOptions) bool) {
		if err != nil {
			return
		}
		for ty := 0.0; ty < height; ty += y * sy {
			for tx := 0.0; tx < width; tx += x * sx {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(sx, sy)
				op.GeoM.Translate(tx, ty)
				if !yield(img, op) {
					return
				}
			}
		}
	}
}

// 全画面背景
func BackgoundImage(width, height float64) iter.Seq2[*ebiten.Image, *ebiten.DrawImageOptions] {
	return getFillImageParams(width, height, path_maptile_umi_01)
}
