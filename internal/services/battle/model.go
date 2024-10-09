package battle

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"math"
	"time"

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

	err = writeUUID(&buf, m.GameId)
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
	dist.GameId, err = readUUID(r)
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

	err = writeUUID(&buf, m.GameId)
	if err != nil {
		return nil, err
	}

	err = writeInt64(&buf, int64(m.ReverseUnixNano))
	if err != nil {
		return nil, err
	}

	err = writeUUID(&buf, m.PlayerId)
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
	err = writeUUID(&buf, m.GameId)
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

	dist.GameId, err = readUUID(r)
	if err != nil {
		return dist, err
	}

	reverse, err := readInt64(r)
	if err != nil {
		return dist, err
	}
	dist.ReverseUnixNano = reverseUnixNano(reverse)

	dist.PlayerId, err = readUUID(r)
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

var lenUUID = 36

func writeUUID(w io.Writer, s string) error {
	var b = make([]byte, lenUUID)
	if lenUUID < len(s) {
		// 末尾から書き込み
		_ = copy(b, []byte(s)[len(s)-lenUUID:])
	} else {
		_ = copy(b, []byte(s))
	}
	_, err := w.Write(b)
	return err
}

func readUUID(r io.Reader) (string, error) {
	var b = make([]byte, lenUUID)
	_, err := r.Read(b)
	if err != nil {
		return "", err
	}
	return string(bytes.Trim(b, "\x00")), nil
}

func writeInt64(w io.Writer, i int64) error {
	return binary.Write(w, binary.BigEndian, uint64(i))
}

func readInt64(r io.Reader) (int64, error) {
	var b = make([]byte, 8)
	_, err := r.Read(b)
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(b)), nil
}
