package app

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/qmuntal/stateless"
	"github.com/yyyoichi/submarine-game/internal/app/components"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

type (
	State   string
	Trigger string
)

const (
	GameStart State = "GameStart"
	GamePlay  State = "GamePlay"

	Play Trigger = "Play"
	Init Trigger = "Init"
)

type Game struct {
	state *stateless.StateMachine
	count int
}

func New() *Game {
	state := stateless.NewStateMachine(GameStart)
	state.Configure(GameStart).Permit(Play, GamePlay)
	state.Configure(GamePlay).Permit(Init, GameStart)
	return &Game{state: state}
}

func (g *Game) Update() error {
	g.count++
	if g.count%10 == 0 {
		g.state.Fire(Play, GamePlay)
	}
	if g.count%10 == 5 {
		g.state.Fire(Init, GameStart)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for img, options := range components.BackgoundImage(screenWidth, screenHeight) {
		screen.DrawImage(img, options)
	}
	m := components.OceanMap{
		Width:  screenHeight - 20,
		Height: screenHeight - 20,
		W:      6,
		H:      6,
		P:      8,
	}

	m.ClearOceanMap(2, g.count%(6*6))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(10, 10)
	screen.DrawImage(m.Src, op)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("count: %d, state: %s", g.count, g.state.String()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func Draw(screen *ebiten.Image) {
	for img, options := range components.BackgoundImage(screenWidth, screenHeight) {
		screen.DrawImage(img, options)
	}
}
