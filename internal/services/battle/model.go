package battle

import (
	"bytes"
	"encoding/json"
	"math"
	"time"

	"github.com/yyyoichi/submarine-game/internal/core"
	"github.com/yyyoichi/submarine-game/internal/store"
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

	err = store.WriteUUID(&buf, m.GameId)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m gameModel) PrefixKey() ([]byte, error) {
	return m.Key()
}

func (gameModel) ParseKey(k []byte) (gameModel, error) {
	var dist gameModel

	r := bytes.NewReader(k)
	_, err := r.ReadByte()
	if err != nil {
		return dist, err
	}
	dist.GameId, err = store.ReadUUID(r)
	return dist, err
}

func (gameModel) Parse(v []byte) (dist gameModel, err error) {
	err = json.Unmarshal(v, &dist)
	return
}

// implements sotre.model
type actionModel struct {
	core.Action
	ReverseUnixNano reverseUnixNano
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

	err = store.WriteUUID(&buf, m.GameId)
	if err != nil {
		return nil, err
	}

	err = store.WriteInt64(&buf, int64(m.ReverseUnixNano))
	if err != nil {
		return nil, err
	}

	err = store.WriteUUID(&buf, m.PlayerId)
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
	err = store.WriteUUID(&buf, m.GameId)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m actionModel) ParseKey(k []byte) (actionModel, error) {
	var dist actionModel

	r := bytes.NewReader(k)
	_, err := r.ReadByte()
	if err != nil {
		return dist, err
	}

	dist.GameId, err = store.ReadUUID(r)
	if err != nil {
		return dist, err
	}

	reverse, err := store.ReadInt64(r)
	if err != nil {
		return dist, err
	}
	dist.ReverseUnixNano = reverseUnixNano(reverse)

	dist.PlayerId, err = store.ReadUUID(r)
	if err != nil {
		return dist, err
	}

	return dist, nil
}

func (m actionModel) Parse(v []byte) (dist actionModel, err error) {
	err = json.Unmarshal(v, &dist)
	if err != nil {
		return
	}
	dist.Timestamp = dist.ReverseUnixNano.restore()
	return
}

// NOTE badgerのReverseイテレーションができないので応急処置
// MaxInt64から現在時刻を引いた行動時刻
type reverseUnixNano int64

func (u *reverseUnixNano) setTimestamp() time.Time {
	now := time.Now()
	*u = reverseUnixNano(math.MaxInt64 - now.UnixNano())
	return now
}
func (u *reverseUnixNano) restore() time.Time {
	return time.Unix(0, int64(math.MaxInt64-*u))
}
