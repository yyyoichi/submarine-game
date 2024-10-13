package matching

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yyyoichi/submarine-game/internal/store"
)

type MatchingService struct {
	waitPlayer string
	mu         sync.Mutex

	store              *store.Store
	waitTickerDuration time.Duration
}

func (s *MatchingService) Leave(playerId string) error {
	s.init()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.waitPlayer == playerId {
		s.waitPlayer = ""
	}
	return nil
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
		output.WaitMatching = func(ctx context.Context) <-chan error {
			return nil
		}
		return &output, nil
	}
	output.WaitMatching = func(ctx context.Context) <-chan error {
		ch := make(chan error)
		go func() {
			defer close(ch)
			for {
				tick := time.NewTicker(s.waitTickerDuration)
				defer tick.Stop()

				m, err := s.found(output.PlayerId)
				if err != nil {
					ch <- err
					return
				}
				if m != nil {
					output.EnemyId = m.EnemyId
					output.GameId = m.GameId
					output.Matched = true
					return
				}

				select {
				case <-ctx.Done():
					ch <- context.Cause(ctx)
					return
				case <-tick.C:

				}
			}
		}()

		return ch
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

func (s *MatchingService) match(playerId string) (string, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.waitPlayer == "" {
		s.waitPlayer = playerId
		return "", "", nil
	}

	gameId := uuid.NewString()
	// waitPlayerが検索できるようにゲーム情報を渡しておく。
	waitPlayer := s.waitPlayer
	var models store.Models[matchModel]
	models.Append(matchModel{
		GameId:   gameId,
		PlayerId: waitPlayer,
		EnemyId:  playerId,
	})
	err := s.store.Set(&models)
	if err != nil {
		return "", "", err
	}
	s.waitPlayer = ""
	return gameId, waitPlayer, err
}

func (s *MatchingService) init() {
	if s.waitTickerDuration == 0 {
		s.waitTickerDuration = time.Duration(time.Millisecond * 200)
	}
}
