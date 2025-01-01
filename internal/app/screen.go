package app

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/qmuntal/stateless"
	"github.com/yyyoichi/submarine-game/internal/app/animation"
	"github.com/yyyoichi/submarine-game/internal/app/components"
	"github.com/yyyoichi/submarine-game/internal/app/fonts"
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

var (
	fadeAni = animation.Animation{
		TPS:          60,
		ByPercentage: [][2]float32{{0, 0}, {0.1, 1}, {0.8, 1}, {1, 0}},
		Duration:     time.Duration(time.Millisecond * 1200),
	}
	swipAni = animation.Animation{
		TPS:          60,
		ByPercentage: [][2]float32{{0, -10}, {1, 0}},
		Duration:     time.Duration(time.Millisecond * 300),
	}
)

var titleFontSource *text.GoTextFaceSource

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.TrainOneRegular))
	if err != nil {
		log.Fatal(err)
	}
	titleFontSource = s
}

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
	if g.count%(60*4) == 0 {
		fadeAni.Clear()
		swipAni.Clear()
	}
	g.count++
	if g.count/60%10 == 0 {
		g.state.Fire(Play, GamePlay)
	}
	if g.count/60%10 == 5 {
		g.state.Fire(Init, GameStart)
	}
	fadeAni.Update()
	swipAni.Update()
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

	m.ClearOceanMap(2, g.count/60%(6*6))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(10, 10)
	screen.DrawImage(m.Src, op)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("count: %d, state: %s", g.count, g.state.String()))

	sf := components.FullScreen{
		Width:  screenWidth,
		Height: screenHeight,
		Fade:   fadeAni.Value,
	}
	sf.Clear()
	img := sf.Templete()
	{
		op := &text.DrawOptions{}
		op.GeoM.Translate(20, float64(swipAni.Value()))
		op.ColorScale.ScaleWithColor(color.NRGBA{255, 255, 255, uint8(255 * fadeAni.Value())})
		text.Draw(img, "相手のターン", &text.GoTextFace{
			Source: titleFontSource,
			Size:   48,
		}, op)
	}
	sf.Draw(img)
	screen.DrawImage(sf.Src, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func Draw(screen *ebiten.Image) {
	for img, options := range components.BackgoundImage(screenWidth, screenHeight) {
		screen.DrawImage(img, options)
	}
}
