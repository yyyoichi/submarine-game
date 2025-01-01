package fonts

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	//go:embed TrainOne-Regular.ttf
	trainOneRegular       []byte
	TrainOneRegularSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(trainOneRegular))
	if err != nil {
		log.Fatal(err)
	}
	TrainOneRegularSource = s
}
