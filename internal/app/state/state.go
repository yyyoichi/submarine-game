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
		value any
	}
)

func (s *store) putStateKey(key string, initValue any) {
	s.mu.RLock()
	if _, found := s.states[key]; found {
		s.mu.RUnlock()
		return
	}
	s.mu.RUnlock()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[key] = &storeData{
		value: initValue,
	}
}
