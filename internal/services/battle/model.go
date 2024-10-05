package battle

import (
	"bufio"
	"bytes"
	"encoding/binary"

	"github.com/yyyoichi/submarine-game/internal/core"
)

// implements sotre.model
type gameModel struct {
	core.Game
}

func newGameModel(gameId string) gameModel {
	return gameModel{Game: core.Game{GameId: gameId}}
}

func (m gameModel) Key() ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.WriteString("G")
	if err != nil {
		return nil, err
	}
	_, err = buf.WriteString("\n")
	if err != nil {
		return nil, err
	}
	_, err = buf.WriteString(m.GameId)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m gameModel) PrefixKey() ([]byte, error) {
	return m.Key()
}

func (m gameModel) ParseKey(k []byte) (dist gameModel, err error) {
	sc := bufio.NewScanner(bytes.NewReader(k))

	for i := 0; sc.Scan(); i++ {
		switch i {
		case 0:
			_ = sc.Bytes()
		case 1:
			dist.GameId = sc.Text()
		}
	}
	return
}

// implements sotre.model
type actionModel struct {
	core.Action
}

func newActionModel(gameId string) actionModel {
	return actionModel{Action: core.Action{GameId: gameId}}
}

func (m actionModel) Key() ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.WriteString("A")
	if err != nil {
		return nil, err
	}
	_, err = buf.WriteString("\n")
	if err != nil {
		return nil, err
	}
	_, err = buf.WriteString(m.GameId)
	if err != nil {
		return nil, err
	}
	_, err = buf.WriteString("\n")
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buf, binary.BigEndian, uint64(m.ReverseUnixNano))
	if err != nil {
		return nil, err
	}
	_, err = buf.WriteString("\n")
	if err != nil {
		return nil, err
	}
	_, err = buf.WriteString(m.PlayerId)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m actionModel) PrefixKey() ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.WriteString("A")
	if err != nil {
		return nil, err
	}
	_, err = buf.WriteString("\n")
	if err != nil {
		return nil, err
	}
	_, err = buf.WriteString(m.GameId)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m actionModel) ParseKey(k []byte) (dist actionModel, err error) {
	sc := bufio.NewScanner(bytes.NewReader(k))

	for i := 0; sc.Scan(); i++ {
		switch i {
		case 0:
			_ = sc.Bytes()
		case 1:
			dist.GameId = sc.Text()
		case 2:
			reverse := binary.BigEndian.Uint64(sc.Bytes())
			dist.ReverseUnixNano = int64(reverse)
		case 3:
			dist.PlayerId = sc.Text()
		}
	}
	return
}
