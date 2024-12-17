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

func BackgoundImage(maxX, maxY float64) iter.Seq2[*ebiten.Image, *ebiten.DrawImageOptions] {
	img, _, err := ebitenutil.NewImageFromFileSystem(assets, "assets/maptile_umi_01.png")
	b := img.Bounds()
	x, y := float64(b.Dx()), float64(b.Dy())
	sx, sy := cellsize/x, cellsize/y
	return func(yield func(*ebiten.Image, *ebiten.DrawImageOptions) bool) {
		if err != nil {
			return
		}
		for ty := 0.0; ty < maxY; ty += y * sy {
			for tx := 0.0; tx < maxX; tx += x * sx {
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
