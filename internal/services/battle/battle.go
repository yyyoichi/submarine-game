package battle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/yyyoichi/hookdb"
	"github.com/yyyoichi/submarine-game/internal/core"
	"github.com/yyyoichi/submarine-game/internal/store"
)

type BattleService struct {
	Store              *store.Store
	db                 hookdb.HookDB
	timeoutDuration    time.Duration
	waitTickerDuration time.Duration
}

func New(s *store.Store) *BattleService {
	var battle BattleService
	battle.Store = s
	battle.init()
	return &battle
}

func (s *BattleService) NewGame(gameId string, playerIds [2]string) error {
	game := core.Game{
		GameId:    gameId,
		PlayerIds: playerIds,
		OceanMap:  core.DefaultOceanMap,
		Islands:   make([]core.Sector, int(core.DefaultOceanMap.IslandCount)),
		Submarines: map[string]core.Submarine{
			playerIds[0]: core.DefaultSubmarine,
			playerIds[1]: core.DefaultSubmarine,
		},
	}
	l := int(core.DefaultOceanMap.H * core.DefaultOceanMap.W)
	game.Islands[0] = core.Sector(rand.IntN(l))
	game.Islands[1] = core.Sector(rand.IntN(l))
	err := s.setGame(game)
	if err != nil {
		return fmt.Errorf("cannot create new game: %w", err)
	}
	return nil
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
		Game:            *game,
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
		From:     me.At,
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
		Game:            *game,
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
		From:         me.At,
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
		Game:            *game,
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
		From:         me.At,
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
	s.init()

	// 存在するゲームか
	game, err := s.getGame(input.GameId)
	if err != nil {
		return nil, fmt.Errorf("%w: cannot get game[%s]: %w", ErrGameNotFound, input.GameId, err)
	}

	actions, err := s.getAllAction(input.GameId)
	if err != nil {
		return nil, fmt.Errorf("cannot get all actions of game[%s]", input.GameId)
	}

	var resp = GetLogsOutput{
		Game:                *game,
		NumTurn:             (len(actions) + 1) / 2,
		TimeoutDurationMSec: s.timeoutDuration.Milliseconds(),
		// exp SET
		SectorActionsMap: make(map[core.Sector][]core.ActionType),
		Actions:          make([]LogAction, len(actions)),
	}
	// プレイヤの前回行動
	for _, action := range actions {
		if action.PlayerId == input.PlayerId {
			resp.Prev = action.Action
			break
		}
	}

	resp.GameOver = s.gameIsOver(*game, actions)

	if len(actions) == 0 || (len(actions) == 1 && actions[0].PlayerId != input.PlayerId) {
		// 行動がないか、あっても一つで相手の行動のみの場合
		resp.RequireAction = true
		resp.RequireDeployAction = true
	} else {
		latest := actions[0]
		resp.RequireAction = latest.PlayerId != input.PlayerId
		resp.RequireDeployAction = false
		resp.Timeout = latest.Timestamp.Add(s.timeoutDuration)
	}
	if len(actions) < 2 {
		resp.Timeout = game.Timestamp.Add(s.timeoutDuration)
	}

	// 終了時0値
	if resp.GameOver != nil {
		resp.RequireAction = false
		resp.RequireDeployAction = false
		resp.Timeout = time.Time{}
	}

	// 行動ログ
	resp.Actions = make([]LogAction, len(actions))
	for i, action := range actions {
		resp.Actions[i] = LogAction{
			PlayerId:     action.PlayerId,
			T:            action.T,
			ActionResult: action.ActionResult,
			Turn:         (len(actions) - i - 1) / 2,
			At:           action.At,
			From:         action.From,
			To:           action.To,
			Direction:    core.UnknownDirection,
		}
		if resp.GameOver == nil && action.PlayerId != input.PlayerId {
			// 決着ついておらず、相手の行動の場合、行動位置をマスクする。
			resp.Actions[i].At = -1
			if action.T == core.MoveAction {
				resp.Actions[i].To = -1
				resp.Actions[i].From = -1
			}
		}
		if action.T == core.MoveAction {
			resp.Actions[i].Direction = game.Direction(action.From, action.To)
		}
	}
	// 行動可能海域計算
	if resp.GameOver == nil {
		sectors := game.SectorStatus(input.PlayerId, resp.Prev)
		for sector, ss := range sectors {
			acts := make([]core.ActionType, 0, len(ss))
			for _, s := range ss {
				var actionType core.ActionType
				switch s {
				case core.CanMove:
					actionType = core.MoveAction
				case core.CanFireTorpedo:
					actionType = core.TorpedoFireAction
				case core.CanTriggerMine:
					actionType = core.MineTriggerAction
				}
				acts = append(acts, actionType)
			}
			resp.SectorActionsMap[sector] = acts
		}
	}

	// 最新の行動が自分である場合
	return &resp, nil
}

// 自分のターンまで待機する
func (s *BattleService) WaitTurn(ctx context.Context, input *WaitTurnInput) (<-chan struct{}, error) {
	s.init()

	// 存在するゲームか
	game, err := s.getGame(input.GameId)
	if err != nil {
		return nil, fmt.Errorf("%w: cannot get game[%s]: %w", ErrGameNotFound, input.GameId, err)
	}
	// プレイヤーIdは存在しているか
	_, found := game.Submarines[input.PlayerId]
	if !found {
		// ゲームが存在しないことにする
		return nil, fmt.Errorf("%w: player[%s] is not found", ErrGameNotFound, input.PlayerId)
	}

	ch := make(chan struct{})
	go func() {
		defer close(ch)
		tick := time.NewTicker(s.waitTickerDuration)
		defer tick.Stop()

		for {
			latest, err := s.getLatestAction(input.GameId)
			if err != nil {
				continue
			}
			if gemeOver := s.gameOver(&latest.Action); gemeOver != nil {
				return
			}
			if latest.PlayerId != input.PlayerId {
				return
			}
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
			}
		}
	}()
	return ch, nil
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

	// セクターの利用可能ステータスの確認
	enables := input.Game.SectorStatus(input.PlayerId, prev.Action, input.At)
	status, found := enables[input.At]
	if !found || !slices.Contains(status, input.ExpSectorStatus) {
		return nil, nil, fmt.Errorf("%w: ", ErrInvalidActionType)
	}
	return &prev.Action, &enemy.Action, nil
}

func (s *BattleService) gameIsOver(game core.Game, actions []actionModel) *GameOver {
	var resp GameOver
	switch l := len(actions); l {
	case 0:
		if time.Now().After(game.Timestamp.Add(s.timeoutDuration)) {
			resp.Reason = core.Timeout
			return &resp
		}
		return nil
	case 1:
		latest := actions[0].Action
		if time.Now().After(game.Timestamp.Add(s.timeoutDuration)) {
			resp.Winner = latest.PlayerId
			resp.Reason = core.Timeout
			return &resp
		}
		return nil
	}
	return s.gameOver(&actions[0].Action)
}

func (s *BattleService) gameOver(latest *core.Action) *GameOver {
	if latest == nil {
		return nil
	}
	var resp GameOver
	if latest.ActionResult == core.Hit {
		resp.Winner = latest.PlayerId
		if latest.T == core.TorpedoFireAction {
			resp.Reason = core.TorpedoHit
		}
		if latest.T == core.MineTriggerAction {
			resp.Reason = core.MineHit
		}
		return &resp
	}
	if time.Now().After(latest.Timestamp.Add(s.timeoutDuration)) {
		resp.Winner = latest.PlayerId
		resp.Reason = core.Timeout
	}

	if resp.Winner == "" {
		return nil
	}
	return &resp
}

// 新しいゲームをセットする
func (s *BattleService) setGame(game core.Game) error {
	game.Timestamp = time.Now()
	key, err := getGameModelKey(game.GameId)
	if err != nil {
		return fmt.Errorf("cannot get key from game model: %w", err)
	}
	v, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("cannot marshal game model: %w", err)
	}

	err = s.db.Put(key, v)
	if err != nil {
		return fmt.Errorf("cannot put game model: %w", err)
	}

	return nil
}

// ゲームを取得する
func (s *BattleService) getGame(gameId string) (*core.Game, error) {
	key, err := getGameModelKey(gameId)
	if err != nil {
		return nil, fmt.Errorf("cannot get key from game model: %w", err)
	}
	v, err := s.db.Get(key)
	if err != nil {
		if errors.Is(err, hookdb.ErrKeyNotFound) {
			err = fmt.Errorf("%w: gameId[%s]: %w", ErrGameNotFound, gameId, err)
		}
		return nil, fmt.Errorf("cannot get game model: %w", err)
	}
	var output core.Game
	err = json.Unmarshal(v, &output)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal game model: %w", err)
	}
	return &output, nil
}

func (s *BattleService) deleteGame(gameId string) error {
	key, err := getGameModelKey(gameId)
	if err != nil {
		return fmt.Errorf("cannot get key from game model: %w", err)
	}
	err = s.db.Delete(key)
	if err != nil {
		if errors.Is(err, hookdb.ErrKeyNotFound) {
			return fmt.Errorf("%w: gameId[%s]: %w", ErrGameNotFound, gameId, err)
		}
		return fmt.Errorf("cannot delete game model: %w", err)
	}
	return nil
}

// 行動を記録する
func (s *BattleService) appendAction(action core.Action) error {
	action.Timestamp = time.Now()
	key, err := getActionModelKey(action.GameId, action.PlayerId, action.Timestamp)
	if err != nil {
		return fmt.Errorf("cannot get key from action model: %w", err)
	}
	v, err := json.Marshal(action)
	if err != nil {
		return fmt.Errorf("cannot marshal action model: %w", err)
	}
	err = s.db.Put(key, v)
	if err != nil {
		return fmt.Errorf("cannot put action model: %w", err)
	}
	return nil
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
	return &values[0], nil
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
	return &values[0], nil
}

// 新しい順に行動を全件取得する。
func (s *BattleService) getAllAction(gameId string) ([]actionModel, error) {
	var models store.Models[actionModel]
	models.Append(newActionModel(gameId))
	models.IsQueryTarget = func(am actionModel) (is bool, end bool) {
		return true, false
	}
	err := s.Store.Query(&models)
	if err != nil {
		if errors.Is(err, store.ErrKeyNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return models.GetValues(), nil
}

func (s *BattleService) init() {
	if s.timeoutDuration == 0 {
		s.timeoutDuration = time.Duration(time.Second*30 + time.Millisecond*500)
	}
	if s.waitTickerDuration == 0 {
		s.waitTickerDuration = time.Duration(time.Millisecond * 200)
	}
}
