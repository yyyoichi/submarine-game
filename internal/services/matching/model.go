package matching

import (
	"bytes"
	"encoding/json"

	"github.com/yyyoichi/submarine-game/internal/store"
)

// implements sotre.model
type matchModel struct {
	// 先んじてwaitしていたプレイヤ
	PlayerId string
	GameId   string
	// 合流したプレイや
	EnemyId string
}

func (m matchModel) Key() ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.WriteString("M")
	if err != nil {
		return nil, err
	}

	err = store.WriteUUID(&buf, m.PlayerId)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m matchModel) PrefixKey() ([]byte, error) {
	return m.Key()
}

func (matchModel) ParseKey(k []byte) (matchModel, error) {
	var dist matchModel

	r := bytes.NewReader(k)
	_, err := r.ReadByte()
	if err != nil {
		return dist, err
	}
	dist.PlayerId, err = store.ReadUUID(r)
	return dist, err
}

func (matchModel) Parse(v []byte) (dist matchModel, err error) {
	err = json.Unmarshal(v, &dist)
	return
}
