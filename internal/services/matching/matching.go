package matching

import (
	"context"
	"fmt"
	"strings"
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

type MatchOutput struct {
	GameId   string
	PlayerId string
	EnemyId  string
}

func (s *MatchingService) Match(ctx context.Context) (<-chan MatchOutput, error) {
	s.init()
	s.mu.Lock()
	defer s.mu.Unlock()

	playerId := uuid.NewString()

	if s.waitPlayer == "" {
		s.waitPlayer = playerId
		output, err := s.db.Subscribe(ctx, []byte(fmt.Sprintf("WP%s", playerId)), hookdb.WithOnceSubscription())
		if err != nil {
			return nil, err
		}
		ch := make(chan MatchOutput)
		go func(me string) {
			defer close(ch)
			select {
			case <-ctx.Done():
				s.mu.Lock()
				defer s.mu.Unlock()
				if playerId == s.waitPlayer {
					s.waitPlayer = ""
				}
				_ = s.db.RemoveHook([]byte(fmt.Sprintf("WP%s", playerId)))
				return
			case v, ok := <-output:
				if !ok {
					return
				}
				vv := strings.Split(string(v), ",")
				ch <- MatchOutput{
					PlayerId: me,
					GameId:   vv[0],
					EnemyId:  vv[1],
				}
			}
		}(playerId)
		return ch, nil
	}
	gameId := uuid.NewString()
	key := []byte(fmt.Sprintf("WP%s", s.waitPlayer))
	v := []byte(fmt.Sprintf("%s,%s", gameId, playerId))
	err := s.db.Put(key, v)
	if err != nil {
		return nil, err
	}

	ch := make(chan MatchOutput)
	go func(enemy string) {
		defer close(ch)
		ch <- MatchOutput{
			GameId:   gameId,
			PlayerId: playerId,
			EnemyId:  enemy,
		}
	}(s.waitPlayer)
	s.waitPlayer = ""
	return ch, nil

}

func (s *MatchingService) init() {
	if s.waitTickerDuration == 0 {
		s.waitTickerDuration = time.Duration(time.Millisecond * 200)
	}
	if s.db == nil {
		s.db = hookdb.New()
	}
}
