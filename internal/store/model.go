package store

import (
	"encoding/json"
	"iter"
	"sync"

	"github.com/dgraph-io/badger/v4"
)

type (
	Models[T model[T]] struct {
		src  []T
		dist []T
		mu   sync.Mutex

		IsQueryTarget func(T) (is bool, end bool)
	}
	// 実装はこのインターフェースを満たす必要がある。
	model[T any] interface {
		// 完全一致キー
		Key() ([]byte, error)
		// Query時に利用する前方一致キー
		PrefixKey() ([]byte, error)
		// キーパーサ
		ParseKey([]byte) (T, error)
	}
)

func (m *Models[T]) Append(src T) {
	m.mu.Lock()
	defer m.mu.Unlock()

	l := len(m.src)
	var appended = make([]T, l+1)
	_ = copy(appended, m.src)
	appended[l] = src
	m.src = appended
}

func (m *Models[T]) GetValues() []T {
	return m.dist
}

func (m *Models[T]) keys() iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		for _, m := range m.src {
			v, err := m.Key()
			if ok := yield(v, err); !ok {
				return
			}
		}
	}
}

func (m *Models[T]) setValue(v []byte) error {
	var dist T
	err := json.Unmarshal(v, &dist)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	l := len(m.dist)
	var appended = make([]T, l+1)
	_ = copy(appended, m.dist)
	appended[l] = dist
	m.dist = appended
	return nil
}

func (m *Models[T]) prefixKeys() iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		for _, m := range m.src {
			v, err := m.PrefixKey()
			if ok := yield(v, err); !ok {
				return
			}
		}
	}
}

func (m *Models[T]) setValues(items iter.Seq2[*badger.Item, error]) error {
	for item, err := range items {
		if err != nil {
			return err
		}
		if m.IsQueryTarget == nil {
			err = item.Value(m.setValue)
			if err != nil {
				return err
			}
			continue
		}
		// filter target value
		var model T
		model, err = model.ParseKey(item.Key())
		if err != nil {
			return err
		}
		ok, end := m.IsQueryTarget(model)
		if ok {
			err = item.Value(m.setValue)
			if err != nil {
				return err
			}
		}
		if end {
			break
		}
	}
	return nil
}

func (m *Models[T]) values() ([][]byte, [][]byte, error) {
	var ks = make([][]byte, len(m.src))
	var vs = make([][]byte, len(m.src))
	for i, model := range m.src {
		v, err := json.Marshal(model)
		if err != nil {
			return nil, nil, err
		}
		k, err := model.Key()
		if err != nil {
			return nil, nil, err
		}
		vs[i] = v
		ks[i] = k
	}
	return ks, vs, nil
}
