package state

type Options[T any] struct {
	InitialValue    *T            // default new(T). set at first time only
	SetValidateFunc func(T) error // defult return nil. set at first time only
}

// return pointer of T
func (o *Options[T]) getInitialValue() any {
	if o.InitialValue != nil {
		return o.InitialValue
	}
	return new(T)
}

// validate v
func (o *Options[T]) valid(v any) error {
	if o.SetValidateFunc != nil {
		return o.SetValidateFunc(v.(T))
	}
	return nil
}

type OptionFunc[T any] func(*Options[T]) error

func WithInitialValue[T any](v T) OptionFunc[T] {
	return func(o *Options[T]) error {
		o.InitialValue = &v
		return nil
	}
}

func WithSetValidateFunc[T any](f func(T) error) OptionFunc[T] {
	return func(o *Options[T]) error {
		o.SetValidateFunc = f
		return nil
	}
}
