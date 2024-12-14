package state

type Options[T any] struct {
	InitialValue    *T            // default new(T)
	SetValidateFunc func(T) error // defult return nil
}

func (o *Options[T]) getInitialValue() *T {
	if o.InitialValue != nil {
		return o.InitialValue
	}
	return new(T)
}

func (o *Options[T]) valid(v T) error {
	if o.SetValidateFunc != nil {
		return o.SetValidateFunc(v)
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
