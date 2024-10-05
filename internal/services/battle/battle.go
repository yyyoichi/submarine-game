package battle

import (
	"errors"
	"fmt"
	"time"

	"github.com/yyyoichi/submarine-game/internal/core"
	"github.com/yyyoichi/submarine-game/internal/store"
)

type BattleService struct {
	Store *store.Store
}

// 新しいゲームをセットする
func (s *BattleService) setGame(game core.Game) error {
	game.Timestamp = time.Now()

	var models store.Models[gameModel]
	models.Append(gameModel{Game: game})
	return s.Store.Set(&models)
}

// ゲームを取得する
func (s *BattleService) getGame(gameId string) (*gameModel, error) {
	var models store.Models[gameModel]
	models.Append(newGameModel(gameId))
	err := s.Store.Get(&models)
	if err != nil {
		if errors.Is(err, store.ErrKeyNotFound) {
			err = fmt.Errorf("%w: gameId[%s]: %w", ErrGameNotFound, gameId, err)
		}
		return nil, err
	}
	values := models.GetValues()
	if len(values) == 0 {
		return nil, fmt.Errorf("%w: gameId[%s]", ErrGameNotFound, gameId)
	}
	return &values[0], nil
}

func (s *BattleService) deleteGame(gameId string) error {
	var models store.Models[gameModel]
	models.Append(newGameModel(gameId))
	return s.Store.Delete(&models)
}

// 行動を記録する
func (s *BattleService) appendAction(action core.Action) error {
	model := actionModel{Action: action}
	model.Action.Timestamp = model.ReverseUnixNano.setTimestamp()

	var models store.Models[actionModel]
	models.Append(model)
	return s.Store.Set(&models, store.WithTTL(time.Duration(time.Minute*30)))
}

// 各プレイヤーの最後の行動を取得する
func (s *BattleService) getPrevActions(gameId string) (map[string]*actionModel, error) {
	var got = make(map[string]struct{}, 2)

	var models store.Models[actionModel]
	models.Append(newActionModel(gameId))
	models.IsQueryTarget = func(am actionModel) (is bool, end bool) {
		_, found := got[am.PlayerId]
		if !found {
			got[am.PlayerId] = struct{}{}
			is = true
		}
		end = len(got) == 2
		return
	}
	err := s.Store.Query(&models, store.WithReverse(false), store.WithPrefetchValues(false))
	if err != nil {
		if errors.Is(err, store.ErrKeyNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var resp = make(map[string]*actionModel, 2)
	for _, v := range models.GetValues() {
		v.Timestamp = v.ReverseUnixNano.restore()
		resp[v.PlayerId] = &v
	}
	return resp, nil
}

// ゲームの最後の行動を取得する
func (s *BattleService) getLatestAction(gameId string) (*actionModel, error) {
	var models store.Models[actionModel]
	models.Append(newActionModel(gameId))
	models.IsQueryTarget = func(am actionModel) (is bool, end bool) {
		return true, true
	}
	err := s.Store.Query(&models, store.WithReverse(false), store.WithPrefetchValues(false))
	if err != nil {
		if errors.Is(err, store.ErrKeyNotFound) {
			return nil, nil
		}
		return nil, err
	}
	values := models.GetValues()
	if len(values) == 0 {
		return nil, nil
	}
	v := values[0]
	v.Timestamp = v.ReverseUnixNano.restore()
	return &v, nil
}

// playerIdの最後の行動を取得する
func (s *BattleService) getPrevAction(gameId, playerId string) (*actionModel, error) {
	var models store.Models[actionModel]
	models.Append(newActionModel(gameId))
	models.IsQueryTarget = func(am actionModel) (is bool, end bool) {
		if am.PlayerId == playerId {
			return true, true
		}
		return false, false
	}
	err := s.Store.Query(&models, store.WithReverse(false), store.WithPrefetchValues(true))
	if err != nil {
		if errors.Is(err, store.ErrKeyNotFound) {
			return nil, nil
		}
		return nil, err
	}
	values := models.GetValues()
	if len(values) == 0 {
		return nil, nil
	}
	v := values[0]
	v.Timestamp = v.ReverseUnixNano.restore()
	return &v, nil
}
