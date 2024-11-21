package matching

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yyyoichi/hookdb"
	"github.com/yyyoichi/submarine-game/internal/store"
)

type MatchingService struct {
	waitPlayer string
	mu         sync.Mutex

	store              *store.Store
	db                 hookdb.HookDB
	waitTickerDuration time.Duration
}

func New(s *store.Store) *MatchingService {
	var matching MatchingService
	matching.store = s
	matching.init()
	return &matching
}

// [2]string{gameId, enemeyId}
func (s *MatchingService) Wait(ctx context.Context, playerId string) <-chan [2]string {
	s.init()

	ch := make(chan [2]string)
	go func() {
		defer func() {
			close(ch)
			s.mu.Lock()
			defer s.mu.Unlock()
			if s.waitPlayer == playerId {
				s.waitPlayer = ""
			}
		}()
		for {
			tick := time.NewTicker(s.waitTickerDuration)
			defer tick.Stop()

			m, err := s.found(playerId)
			if err != nil {
				continue
			}
			if m != nil {
				ch <- [2]string{m.GameId, m.EnemyId}
				return
			}
			if !s.waitIsMe(playerId) {
				return
			}

			select {
			case <-ctx.Done():
				return
			case <-tick.C:
			}
		}
	}()

	return ch
}

func (s *MatchingService) Join() (*JoinOutput, error) {
	s.init()

	var output JoinOutput
	output.PlayerId = uuid.NewString()
	gameId, enemyId, err := s.match(output.PlayerId)
	if err != nil {
		return nil, err
	}
	if gameId != "" {
		// matched
		output.Matched = true
		output.GameId = gameId
		output.EnemyId = enemyId
	}
	return &output, nil
}

func (s *MatchingService) found(playerId string) (*matchModel, error) {
	var models store.Models[matchModel]
	models.Append(matchModel{
		PlayerId: playerId,
	})
	err := s.store.Get(&models)
	if err != nil && !errors.Is(err, store.ErrKeyNotFound) {
		return nil, err
	}
	values := models.GetValues()
	if len(values) == 0 {
		return nil, nil
	}
	err = s.store.Delete(&models)
	return &values[0], err
}

func (s *MatchingService) waitIsMe(playerId string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.waitPlayer == playerId
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

	// go func() {
	// 	if s.waitPlayer == "" {
	// 		s.waitPlayer = playerId
	// 		s.db.AppendHook([]byte(fmt.Sprintf("WP%s", playerId)), func(k, v []byte) (removeHook bool) {
	// 			defer close(ch)
	// 			select {
	// 			case <-ctx.Done():
	// 				return true
	// 			default:
	// 				ch <- string(v)
	// 			}
	// 			return true
	// 		})

	// 		return
	// 	}
	// 	defer close(ch)
	// 	// ユーザが見つかった場合
	// 	gameId := uuid.NewString()
	// 	// waitPlayerが検索できるようにゲーム情報を渡しておく。
	// 	err := s.db.Put([]byte(fmt.Sprintf("WP%s", s.waitPlayer)), []byte(gameId))
	// 	if err != nil {
	// 		cancel(err)
	// 		return
	// 	}
	// 	s.waitPlayer = ""
	// 	ch <- gameId
	// }()
	// return playerId, ch
}

func (s *MatchingService) match(playerId string) (string, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.waitPlayer == "" {
		s.waitPlayer = playerId
		s.db.AppendHook([]byte(fmt.Sprintf("WP%s", playerId)), func(k, v []byte) (removeHook bool) {

			return true
		})
		return "", "", nil
	}

	gameId := uuid.NewString()
	// waitPlayerが検索できるようにゲーム情報を渡しておく。
	waitPlayer := s.waitPlayer
	tx := s.db.TransactionWithLock()
	err := tx.Put([]byte(fmt.Sprintf("WP%s", waitPlayer)), []byte(gameId))
	if err != nil {
		tx.Rollback()
		return "", "", err
	}
	if err = tx.Commit(); err != nil {
		return "", "", err
	}
	s.waitPlayer = ""
	return gameId, waitPlayer, nil
}

func (s *MatchingService) init() {
	if s.waitTickerDuration == 0 {
		s.waitTickerDuration = time.Duration(time.Millisecond * 200)
	}
}
