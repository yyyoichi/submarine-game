package battle

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/yyyoichi/submarine-game/internal/core"
	"github.com/yyyoichi/submarine-game/internal/store"
)

type BattleService struct {
	Store           *store.Store
	timeoutDuration time.Duration
}

func (s *BattleService) DeploySubmarineAndMines(ctx context.Context, input *DeploySubmarineAndMinesInput) error {
	s.init()

	// 存在するゲームか
	game, err := s.getGame(input.GameId)
	if err != nil {
		return fmt.Errorf("%w: cannot get game[%s]: %w", ErrGameNotFound, input.GameId, err)
	}

	// プレイヤーIdは存在しているか
	submarine, found := game.Submarines[input.PlayerId]
	if !found {
		// ゲームが存在しないことにする
		return fmt.Errorf("%w: player[%s] is not found", ErrGameNotFound, input.PlayerId)
	}
	// タイムアウトしていないか
	if game.Since() > s.timeoutDuration {
		return fmt.Errorf("cannot deploy submarine: %w", ErrTimeout)
	}

	// 機雷の数は正しいか
	if submarine.MineCount < uint8(len(input.Mines)) {
		return fmt.Errorf("%w: %v", ErrInvalidMineCount, input.Mines)
	}
	// 初回行動か
	prev, err := s.getPrevAction(input.GameId, input.PlayerId)
	if err != nil {
		return fmt.Errorf("cannot get player[%s] prev action of game[%s]", input.PlayerId, input.GameId)
	}
	if prev != nil {
		return fmt.Errorf("%w: already deploye", ErrInvalidActionType)
	}

	action := core.Action{
		GameId:   input.GameId,
		PlayerId: input.PlayerId,
		At:       core.Sector(input.At),
		Mines:    make([]core.Sector, len(input.Mines)),
	}
	// 位置は正しいか
	if slices.Contains(game.Islands, action.At) {
		return fmt.Errorf("%w: deployment sector[%v] is island", ErrInvalidSector, action.At)
	}
	for i, m := range input.Mines {
		action.Mines[i] = core.Sector(m)
		if slices.Contains(game.Islands, action.Mines[i]) {
			return fmt.Errorf("%w: mine sector[%v] is island", ErrInvalidSector, action.Mines[i])
		}
	}

	err = s.appendAction(action)
	if err != nil {
		return fmt.Errorf("cannot append action: %w", err)
	}
	return nil
}

func (s *BattleService) Move(ctx context.Context, input *ActionInput) error {
	s.init()

	// 存在するゲームか
	game, err := s.getGame(input.GameId)
	if err != nil {
		return fmt.Errorf("%w: cannot get game[%s]: %w", ErrGameNotFound, input.GameId, err)
	}
	at := core.Sector(input.At)

	me, _, err := s.GetValidPrevActions(ctx, &GetValidPrevActionsInput{
		Game:            game.Game,
		PlayerId:        input.PlayerId,
		ExpSectorStatus: core.CanMove,
		At:              at,
	})
	if err != nil {
		return fmt.Errorf("cannot get valid actions: %w", err)
	}
	action := core.Action{
		GameId:   input.GameId,
		PlayerId: input.PlayerId,
		At:       at,
		T:        core.MoveAction,
		To:       at,
		Mines:    me.Mines[:],
	}
	err = s.appendAction(action)
	if err != nil {
		return fmt.Errorf("cannot append action: %w", err)
	}
	return nil
}

func (s *BattleService) FireTorpedo(ctx context.Context, input *ActionInput) error {
	s.init()

	// 存在するゲームか
	game, err := s.getGame(input.GameId)
	if err != nil {
		return fmt.Errorf("%w: cannot get game[%s]: %w", ErrGameNotFound, input.GameId, err)
	}
	at := core.Sector(input.At)

	me, enemy, err := s.GetValidPrevActions(ctx, &GetValidPrevActionsInput{
		Game:            game.Game,
		PlayerId:        input.PlayerId,
		ExpSectorStatus: core.CanFireTorpedo,
		At:              at,
	})
	if err != nil {
		return fmt.Errorf("cannot get valid actions: %w", err)
	}
	action := core.Action{
		GameId:       input.GameId,
		PlayerId:     input.PlayerId,
		At:           me.At,
		T:            core.TorpedoFireAction,
		ActionResult: game.ActionResult(enemy.At, at),
		To:           at,
		Mines:        me.Mines[:],
	}
	err = s.appendAction(action)
	if err != nil {
		return fmt.Errorf("cannot append action: %w", err)
	}
	return nil
}

func (s *BattleService) TriggerMine(ctx context.Context, input *ActionInput) error {
	s.init()

	// 存在するゲームか
	game, err := s.getGame(input.GameId)
	if err != nil {
		return fmt.Errorf("%w: cannot get game[%s]: %w", ErrGameNotFound, input.GameId, err)
	}
	at := core.Sector(input.At)

	me, enemy, err := s.GetValidPrevActions(ctx, &GetValidPrevActionsInput{
		Game:            game.Game,
		PlayerId:        input.PlayerId,
		ExpSectorStatus: core.CanTriggerMine,
		At:              at,
	})
	if err != nil {
		return fmt.Errorf("cannot get valid actions: %w", err)
	}
	action := core.Action{
		GameId:       input.GameId,
		PlayerId:     input.PlayerId,
		At:           me.At,
		T:            core.MineTriggerAction,
		ActionResult: game.ActionResult(enemy.At, at),
		To:           at,
		Mines: slices.DeleteFunc(me.Mines, func(s core.Sector) bool {
			return s == at
		}),
	}
	err = s.appendAction(action)
	if err != nil {
		return fmt.Errorf("cannot append action: %w", err)
	}
	return nil
}

func (s *BattleService) GetLogs(ctx context.Context, input *GetLogsInput) (*GetLogsOutput, error) {
	return nil, nil
}

// 初回行動以降の味方/敵の前回行動を取得する。
// input.Atへのinput.EnableStatusの許可を期待する。
func (s *BattleService) GetValidPrevActions(_ context.Context, input *GetValidPrevActionsInput) (*core.Action, *core.Action, error) {
	s.init()

	// プレイヤーIdは存在しているか
	_, found := input.Game.Submarines[input.PlayerId]
	if !found {
		// ゲームが存在しないことにする
		return nil, nil, fmt.Errorf("%w: player[%s] is not found", ErrGameNotFound, input.PlayerId)
	}

	// 初回行動済みか
	prevs, err := s.getPrevActions(input.Game.GameId)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot get prevs actions: %w", err)
	}
	if len(prevs) < 2 {
		return nil, nil, fmt.Errorf("%w: cannot move, prev action is no initialize", ErrInvalidActionType)
	}
	prev, found := prevs[input.PlayerId]
	if !found {
		return nil, nil, fmt.Errorf("%w: cannot move, prev action is none", ErrInvalidActionType)
	}
	enemy, found := prevs[input.Game.Enemy(input.PlayerId)]
	if !found {
		return nil, nil, fmt.Errorf("%w: cannot move, prev action is none", ErrInvalidActionType)
	}

	// 自分のが最近行動している
	if enemy.Timestamp.UnixNano() < prev.Timestamp.UnixNano() {
		return nil, nil, fmt.Errorf("%w: exp player < enemy", ErrNotInTurn)
	}
	// タイムアウト
	if s.timeoutDuration < enemy.Since() {
		return nil, nil, fmt.Errorf("%w: too much time has passed since the last enemy action", ErrTimeout)
	}

	// TODO ゲーム終了判定

	// セクターの利用可能ステータスの確認
	enables := input.Game.SectorStatus(input.PlayerId, prev.Action, input.At)
	status, found := enables[input.At]
	if !found || !slices.Contains(status, input.ExpSectorStatus) {
		return nil, nil, fmt.Errorf("%w: ", ErrInvalidActionType)
	}
	return &prev.Action, &enemy.Action, nil
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

func (s *BattleService) init() {
	if s.timeoutDuration == 0 {
		s.timeoutDuration = time.Duration(time.Second*30 + time.Millisecond*500)
	}
}
