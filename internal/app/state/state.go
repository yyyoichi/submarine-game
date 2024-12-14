package state

import "sync"

var statesStore = store{states: make(map[string]*storeData)}

type (
	store struct {
		mu     sync.RWMutex
		states map[string]*storeData
	}
	storeData struct {
		mu    sync.Mutex
		valid func(v any) error
		value any // must pointer
	}
)

type initConfig interface {
	getInitialValue() any // must return pointer
	valid(v any) error    // recieve non-pointer
}

func (s *store) putStateKey(key string, config initConfig) {
	s.mu.RLock()
	if _, found := s.states[key]; found {
		s.mu.RUnlock()
		return
	}
	s.mu.RUnlock()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[key] = &storeData{
		value: config.getInitialValue(),
		valid: config.valid,
	}
}
