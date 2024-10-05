package store

import "time"

type QueryOptionos struct {
	reverse        bool
	prefetchValues bool
}
type QueryOption func(o *QueryOptionos) error

func WithReverse(v bool) QueryOption {
	return func(o *QueryOptionos) error {
		o.reverse = v
		return nil
	}
}

func WithPrefetchValues(v bool) QueryOption {
	return func(o *QueryOptionos) error {
		o.prefetchValues = v
		return nil
	}
}

type SetOptions struct {
	ttl *time.Duration
}
type SetOption func(o *SetOptions) error

func WithTTL(v time.Duration) SetOption {
	return func(o *SetOptions) error {
		o.ttl = &v
		return nil
	}
}
