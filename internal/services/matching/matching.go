package matching

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yyyoichi/hookdb"
)

type MatchingService struct {
	waitPlayer string
	mu         sync.Mutex

	db                 *hookdb.HookDB
	waitTickerDuration time.Duration
}

func New() *MatchingService {
	var matching MatchingService
	matching.init()
	return &matching
}

func (s *MatchingService) Leave(playerId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.waitPlayer == playerId {
		s.waitPlayer = ""
		return nil
	}
	return s.db.RemoveHook([]byte(fmt.Sprintf("WP%s", playerId)))
}

func (s *MatchingService) Match(ctx context.Context, cancel func(error)) (string, <-chan string) {
	s.init()
	s.mu.Lock()
	defer s.mu.Unlock()

	ch := make(chan string)
	playerId := uuid.NewString()

	k := fmt.Sprintf("WP%s", playerId)

	if s.waitPlayer == "" {
		s.waitPlayer = playerId
		innerCh := make(chan string)

		s.db.AppendHook([]byte(k), func(k, v []byte) (removeHook bool) {
			select {
			case <-ctx.Done():
			default:
				innerCh <- string(v)
			}
			return true
		})
		go func() {
			defer close(ch)
			defer close(innerCh)
			select {
			case <-ctx.Done():
				s.mu.Lock()
				defer s.mu.Unlock()
				if playerId == s.waitPlayer {
					s.waitPlayer = ""
				}
				_ = s.db.RemoveHook([]byte(fmt.Sprintf("WP%s", playerId)))
				return
			case gameId, ok := <-innerCh:
				if !ok {
					return
				}
				ch <- gameId
			}
		}()
	} else {
		gameId := uuid.NewString()
		err := s.db.Put([]byte(fmt.Sprintf("WP%s", s.waitPlayer)), []byte(gameId))
		if err != nil {
			cancel(err)
			close(ch)
		} else {
			s.waitPlayer = ""
		}
		go func() {
			defer close(ch)
			select {
			case <-ctx.Done():
			case ch <- gameId:
			}
		}()
	}
	return playerId, ch

}

func (s *MatchingService) init() {
	if s.waitTickerDuration == 0 {
		s.waitTickerDuration = time.Duration(time.Millisecond * 200)
	}
	if s.db == nil {
		s.db = hookdb.New()
	}
}
