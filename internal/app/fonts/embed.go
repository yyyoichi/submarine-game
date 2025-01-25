package fonts

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	//go:embed TrainOne-Regular.ttf
	trainOneRegular []byte
	//go:embed DotGothic16-Regular.ttf
	dotGothic16Regular       []byte
	TrainOneRegularSource    *text.GoTextFaceSource
	DotGothic16RegularSource *text.GoTextFaceSource
)

func new(ttf []byte) *text.GoTextFaceSource {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(ttf))
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func init() {
	TrainOneRegularSource = new(trainOneRegular)
	DotGothic16RegularSource = new(dotGothic16Regular)
}
