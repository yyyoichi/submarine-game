package components

import (
	_ "image/png"
	"iter"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yyyoichi/submarine-game/internal/app/images"
)

var (
	cellsize = 50.0
)

func FillWithImage(width, height float64, img *ebiten.Image) iter.Seq2[*ebiten.Image, *ebiten.DrawImageOptions] {
	b := img.Bounds()
	x, y := float64(b.Dx()), float64(b.Dy())
	sx, sy := cellsize/x, cellsize/y
	return func(yield func(*ebiten.Image, *ebiten.DrawImageOptions) bool) {
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
	return FillWithImage(width, height, images.MaptileUmi01)
}
