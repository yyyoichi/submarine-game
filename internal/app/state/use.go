package state

import "github.com/google/uuid"

type SetStateFunc[T any] func(T) error

func UseGlobalState[T any](key string, opts ...OptionFunc[T]) (*T, SetStateFunc[T]) {
	var o Options[T]
	for _, opt := range opts {
		_ = opt(&o)
	}

	statesStore.putStateKey(key, &o)
	d := statesStore.states[key]
	p := d.value.(*T)
	return p, func(v T) error {
		d.mu.Lock()
		defer d.mu.Unlock()
		err := d.valid(v)
		if err != nil {
			return err
		}
		*p = v
		return nil
	}
}

func UseState[T any](opts ...OptionFunc[T]) (*T, SetStateFunc[T]) {
	key := uuid.NewString()
	return UseGlobalState(key, opts...)
}
