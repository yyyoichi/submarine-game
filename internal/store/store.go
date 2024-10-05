package store

import (
	"errors"
	"fmt"
	"iter"

	"github.com/dgraph-io/badger/v4"
)

type (
	Store struct {
		*badger.DB
	}
	// 実装者は意識する必要はない。
	models interface {
		keys() iter.Seq2[[]byte, error]
		setValue([]byte) error
		prefixKeys() iter.Seq2[[]byte, error]
		setValues(iter.Seq2[*badger.Item, error]) error

		values() ([][]byte, [][]byte, error)
	}
)

func New() (*Store, error) {
	opt := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}
	return &Store{DB: db}, nil
}

func (s *Store) Get(models models) error {
	err := s.View(func(txn *badger.Txn) error {
		for key, err := range models.keys() {
			if err != nil {
				return err
			}
			item, err := txn.Get(key)
			if err != nil {
				if errors.Is(err, badger.ErrKeyNotFound) {
					return fmt.Errorf("%w: key[%s]: %w", ErrKeyNotFound, key, err)
				}
				return err
			}
			err = item.Value(models.setValue)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (s *Store) Query(models models, options ...QueryOption) error {
	var opt QueryOptionos
	for _, o := range options {
		_ = o(&opt)
	}
	var itScan iter.Seq2[*badger.Item, error] = func(yield func(*badger.Item, error) bool) {
		err := s.View(func(txn *badger.Txn) error {
			for prefix, err := range models.prefixKeys() {
				if err != nil {
					return err
				}
				opts := badger.DefaultIteratorOptions
				opts.PrefetchValues = opt.prefetchValues
				opts.Reverse = opt.reverse
				opts.Prefix = prefix
				it := txn.NewIterator(opts)
				defer it.Close()
				for it.Rewind(); it.Valid(); it.Next() {
					item := it.Item()
					if ok := yield(item, nil); !ok {
						return nil
					}
				}
			}
			return nil
		})
		if err != nil {
			_ = yield(nil, err)
		}
	}
	return models.setValues(itScan)
}

func (s *Store) Set(models models, options ...SetOption) error {
	var opt SetOptions
	for _, o := range options {
		_ = o(&opt)
	}
	ks, vs, err := models.values()
	if err != nil {
		return err
	}
	err = s.Update(func(txn *badger.Txn) error {
		for i := range ks {
			k, v := ks[i], vs[i]
			entry := badger.NewEntry(k, v)
			if opt.ttl != nil {
				entry.WithTTL(*opt.ttl)
			}
			err := txn.SetEntry(entry)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (s *Store) Delete(models models) error {
	err := s.Update(func(txn *badger.Txn) error {
		for k, err := range models.keys() {
			if err != nil {
				return err
			}
			err := txn.Delete(k)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
