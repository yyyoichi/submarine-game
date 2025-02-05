package actions

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func IsKeyAJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyA)
}

func IsKeyBJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyBackspace) || inpututil.IsKeyJustPressed(ebiten.KeyB)
}

func IsKeyRightJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyRight)
}

func IsKeyLeftJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyLeft)
}

func IsKeyUpJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyUp)
}

func IsKeyDownJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyDown)
}
