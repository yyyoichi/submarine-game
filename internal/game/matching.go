package game

import (
	"sync"

	"github.com/google/uuid"
)

type Matching struct {
	waitUser string
	mu       sync.Mutex
}

func NewMatching() *Matching {
	m := Matching{
		waitUser: "",
		mu:       sync.Mutex{},
	}
	return &m
}

// 参加してユーザIdを取得する。マッチした場合は対戦相手も返す。
func (m *Matching) Join() (string, *[2]string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u1 := uuid.NewString()
	if m.waitUser == "" {
		m.waitUser = u1
		return u1, nil
	}
	u2 := m.waitUser
	m.waitUser = ""
	return u1, &[2]string{u1, u2}
}

func (m *Matching) Leave(user string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.waitUser == user {
		m.waitUser = ""
	}
}
