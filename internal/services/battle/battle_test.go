package battle

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yyyoichi/submarine-game/internal/core"
)

func TestBattleService(t *testing.T) {
	ctx := context.Background()

	t.Run("DeploySubmarineAndMines", func(t *testing.T) {
		t.Parallel()

		type input = DeploySubmarineAndMinesInput
		var (
			scope          = "DeploySubmarineAndMines"
			defaultIslands = []core.Sector{30, 40}
			newGame        = func(id string) core.Game {
				game := core.Game{
					GameId:    scope + id,
					PlayerIds: [2]string{"1", "2"},
					Submarines: map[string]core.Submarine{
						"1": core.DefaultSubmarine,
						"2": core.DefaultSubmarine,
					},
					OceanMap: core.DefaultOceanMap,
					Islands:  defaultIslands,
				}
				return game
			}
			newInput = func(id string, added *input) *input {
				input := &input{
					GameId:   scope + id,
					PlayerId: "1",
					At:       10,
					Mines:    []int8{1, 2},
				}
				if added == nil {
					return input
				}
				if added.PlayerId != "" {
					input.PlayerId = added.PlayerId
				}
				if added.At != 0 {
					input.At = added.At
				}
				if len(added.Mines) != 0 {
					input.Mines = added.Mines
				}
				return input
			}
		)

		test := []struct {
			preProcess func(b *BattleService)
			input      *input
			exp        error
		}{
			{func(b *BattleService) {
				b.setGame(newGame("1"))
			}, newInput("1", nil),
				nil},

			{func(b *BattleService) {
				// empty

			}, newInput("2", nil),
				ErrGameNotFound},

			{func(b *BattleService) {
				b.setGame(newGame("3"))
			}, newInput("3", &input{PlayerId: "99"}),
				ErrGameNotFound},

			{func(b *BattleService) {
				b.setGame(newGame("4"))
				b.timeoutDuration = time.Nanosecond * 1
				time.Sleep(time.Nanosecond * 2)
			}, newInput("4", nil),
				ErrTimeout},

			{func(b *BattleService) {
				b.setGame(newGame("5"))
			}, newInput("5", &input{Mines: []int8{1, 2, 3}}),
				ErrInvalidMineCount},

			{func(b *BattleService) {
				b.setGame(newGame("6"))
				b.appendAction(core.Action{
					GameId:   scope + "6",
					PlayerId: "1",
				})
			}, newInput("6", nil),
				ErrInvalidActionType},

			{func(b *BattleService) {
				b.setGame(newGame("7"))
			}, newInput("7", &input{At: int8(defaultIslands[0])}),
				ErrInvalidSector},

			{func(b *BattleService) {
				b.setGame(newGame("8"))
			}, newInput("8", &input{Mines: []int8{1, int8(defaultIslands[0])}}),
				ErrInvalidSector},
		}
		for _, tt := range test {
			var battle = New()
			tt.preProcess(battle)
			err := battle.DeploySubmarineAndMines(ctx, tt.input)
			if tt.exp == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.exp)
			}
		}
	})

	t.Run("GetValidPrevActions", func(t *testing.T) {
		t.Parallel()

		type input = GetValidPrevActionsInput
		var (
			scope   = "GetValidPrevActions"
			newGame = func(id string) core.Game {
				game := core.Game{
					GameId:    scope + id,
					PlayerIds: [2]string{"1", "2"},
					Submarines: map[string]core.Submarine{
						"1": core.DefaultSubmarine,
						"2": core.DefaultSubmarine,
					},
					OceanMap: core.DefaultOceanMap,
					Islands:  []core.Sector{30, 40},
				}
				return game
			}
			newAction = func(gameId, playerId string) core.Action {
				act := core.Action{
					GameId:    scope + gameId,
					PlayerId:  playerId,
					At:        1,
					From:      0,
					To:        1,
					T:         core.MoveAction,
					Timestamp: time.Now(),
				}
				return act
			}
			newInput = func(gameId, playerId string, at int8) *input {
				return &input{
					Game:            newGame(gameId),
					PlayerId:        playerId,
					ExpSectorStatus: core.CanMove,
					At:              core.Sector(at),
				}
			}
		)

		test := []struct {
			preProcess func(b *BattleService)
			input      *input
			exp        error
		}{
			{func(b *BattleService) {
				b.setGame(newGame("1"))
				b.appendAction(newAction("1", "1"))
				b.appendAction(newAction("1", "2"))
			}, newInput("1", "1", 0),
				nil},
			// プレイヤーが存在していない
			{func(b *BattleService) {
				b.setGame(newGame("2"))
				b.appendAction(newAction("2", "1"))
				b.appendAction(newAction("2", "2"))
			}, newInput("2", "3", 0),
				ErrGameNotFound},
			// プレイヤーが初回行動していない
			{func(b *BattleService) {
				b.setGame(newGame("3"))
				b.appendAction(newAction("3", "2"))
			}, newInput("3", "1", 0),
				ErrInvalidActionType},
			// 相手が初回行動していない
			{func(b *BattleService) {
				b.setGame(newGame("4"))
				b.appendAction(newAction("4", "1"))
			}, newInput("4", "1", 0),
				ErrInvalidActionType},
			// 順番でない
			{func(b *BattleService) {
				b.setGame(newGame("5"))
				b.appendAction(newAction("5", "1"))
				b.appendAction(newAction("5", "2"))
			}, newInput("5", "2", 0),
				ErrNotInTurn},
			// タイムアウト
			{func(b *BattleService) {
				b.setGame(newGame("6"))
				b.appendAction(newAction("6", "1"))
				b.appendAction(newAction("6", "2"))
				b.timeoutDuration = time.Nanosecond * 1
				time.Sleep(time.Nanosecond * 2)
			}, newInput("6", "1", 0),
				ErrTimeout},
			// 動ける場所でない
			{func(b *BattleService) {
				b.setGame(newGame("7"))
				b.appendAction(newAction("7", "1"))
				b.appendAction(newAction("7", "2"))
			}, newInput("7", "1", 3),
				ErrInvalidActionType},
		}
		for _, tt := range test {
			var battle = New()
			tt.preProcess(battle)
			_, _, err := battle.GetValidPrevActions(ctx, tt.input)
			if tt.exp == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.exp)
			}
		}
	})

	t.Run("Actions", func(t *testing.T) {
		t.Parallel()

		ctx = context.Background()

		var battle = New()
		// 行動はいつも13に対して行われ、常に現在位置7に対して
		// 魚雷機雷行動可能。
		// 魚雷機雷攻撃の場合、相手は常に14にいるため面舵一杯となる。
		test := []struct {
			id  string
			fn  func(context.Context, *ActionInput) error
			exp core.Action
		}{
			{"1", battle.Move, core.Action{
				GameId:   "Action1",
				PlayerId: "1",
				At:       13,
				T:        core.MoveAction,
				From:     7,
				To:       13,
				Mines:    []core.Sector{1, 13},
			}},
			{"2", battle.FireTorpedo, core.Action{
				GameId:       "Action2",
				PlayerId:     "1",
				At:           7,
				T:            core.TorpedoFireAction,
				ActionResult: core.HardToStarboard,
				From:         7,
				To:           13,
				Mines:        []core.Sector{1, 13},
			}},
			{"3", battle.TriggerMine, core.Action{
				GameId:       "Action3",
				PlayerId:     "1",
				At:           7,
				T:            core.MineTriggerAction,
				From:         7,
				ActionResult: core.HardToStarboard,
				To:           13,
				Mines:        []core.Sector{1},
			}},
		}
		for _, tt := range test {
			game := core.Game{
				GameId:    fmt.Sprintf("Action%s", tt.id),
				PlayerIds: [2]string{"1", "2"},
				Submarines: map[string]core.Submarine{
					"1": core.DefaultSubmarine,
					"2": core.DefaultSubmarine,
				},
				OceanMap: core.DefaultOceanMap,
				Islands:  []core.Sector{0, 1},
			}
			battle.setGame(game)
			battle.appendAction(core.Action{
				GameId:   game.GameId,
				PlayerId: "1",
				At:       core.Sector(7),       // !
				Mines:    []core.Sector{1, 13}, // !
			})
			battle.appendAction(core.Action{
				GameId:   game.GameId,
				PlayerId: "2",
				At:       core.Sector(14), // !
				Mines:    []core.Sector{2, 3},
			})
			input := ActionInput{
				GameId:   game.GameId,
				PlayerId: "1",
				At:       13,
			}

			err := tt.fn(ctx, &input)
			assert.NoError(t, err)

			latest, err := battle.getLatestAction(battle.db.DB, game.GameId)
			assert.NoError(t, err)
			// timestamp以外同一であることを確認する
			assert.Equal(t, tt.exp.GameId, latest.GameId)
			assert.Equal(t, tt.exp.PlayerId, latest.PlayerId)
			assert.Equal(t, tt.exp.ActionResult, latest.ActionResult)
			assert.Equal(t, tt.exp.At, latest.At)
			assert.Equal(t, tt.exp.To, latest.To)
			assert.Equal(t, tt.exp.T, latest.T)
			assert.Equal(t, tt.exp.Mines, latest.Mines)

		}
	})

	t.Run("GetLogs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		type input = GetLogsInput
		var (
			scope   = "GetLogs"
			newGame = func(id string) core.Game {
				game := core.Game{
					GameId:    scope + id,
					PlayerIds: [2]string{"1", "2"},
					Submarines: map[string]core.Submarine{
						"1": core.DefaultSubmarine,
						"2": core.DefaultSubmarine,
					},
					OceanMap: core.DefaultOceanMap,
					Islands:  []core.Sector{30, 40},
				}
				return game
			}
			newAction = func(gameId, playerId string) core.Action {
				act := core.Action{
					GameId:    scope + gameId,
					PlayerId:  playerId,
					At:        1,
					From:      0,
					To:        1,
					T:         core.MoveAction,
					Timestamp: time.Now(),
				}
				return act
			}
			newInput = func(gameId, playerId string) *input {
				return &input{
					GameId:   scope + gameId,
					PlayerId: playerId,
				}
			}
		)
		test := []struct {
			preProcess func(b *BattleService)
			input      *input
			expErr     error
			expOutput  *GetLogsOutput
		}{
			// 最後が相手
			{func(b *BattleService) {
				b.setGame(newGame("1"))
				b.appendAction(newAction("1", "1"))
				b.appendAction(newAction("1", "2"))
				b.appendAction(newAction("1", "1"))
				b.appendAction(newAction("1", "2"))
			}, newInput("1", "1"), nil, &GetLogsOutput{
				RequireAction:       true,
				RequireDeployAction: false,
				IsFirstAction:       true,
				NumTurn:             2,
				Actions: []LogAction{
					{PlayerId: "2", Turn: 1, At: -1, T: core.MoveAction, To: -1, Direction: core.East},
					{}, {}, {},
				},
			}},
			// 最後が自分
			{func(b *BattleService) {
				b.setGame(newGame("2"))
				b.appendAction(newAction("2", "1"))
				b.appendAction(newAction("2", "2"))
				b.appendAction(newAction("2", "1"))
			}, newInput("2", "1"), nil, &GetLogsOutput{
				RequireAction:       false,
				RequireDeployAction: false,
				IsFirstAction:       true,
				NumTurn:             2,
				Actions: []LogAction{
					{PlayerId: "1", Turn: 1, At: 1, T: core.MoveAction, To: 1, Direction: core.East},
					{}, {},
				},
			}},
			// 行動なし
			{func(b *BattleService) {
				b.setGame(newGame("3"))
			}, newInput("3", "1"), nil, &GetLogsOutput{
				RequireAction:       true,
				RequireDeployAction: true,
				IsFirstAction:       false,
			}},
			// 自分の初回行動あり
			{func(b *BattleService) {
				b.setGame(newGame("31"))
				b.appendAction(newAction("31", "2"))
			}, newInput("31", "2"), nil, &GetLogsOutput{
				RequireAction:       false,
				RequireDeployAction: true,
				IsFirstAction:       true,
				NumTurn:             1,
				Actions: []LogAction{
					{PlayerId: "2", Turn: 0, At: 1, To: 1, From: 0, T: core.MoveAction, Direction: core.East},
				},
			}},
			// 自分に行動なし
			{func(b *BattleService) {
				b.setGame(newGame("4"))
				b.appendAction(newAction("4", "2"))
			}, newInput("4", "1"), nil, &GetLogsOutput{
				RequireAction:       true,
				RequireDeployAction: true,
				IsFirstAction:       false,
				NumTurn:             1,
				Actions: []LogAction{
					{PlayerId: "2", Turn: 0, At: -1, To: -1, From: -1, T: core.MoveAction, Direction: core.East},
				},
			}},
			// ゲーム終了
			{func(b *BattleService) {
				b.setGame(newGame("5"))
				b.appendAction(newAction("5", "2"))
				b.appendAction(newAction("5", "1"))
				b.timeoutDuration = time.Duration(time.Nanosecond * 1)
				time.Sleep(time.Duration(time.Nanosecond * 2))
			}, newInput("5", "2"), nil, &GetLogsOutput{
				RequireAction:       false,
				RequireDeployAction: false,
				IsFirstAction:       true,
				NumTurn:             1,
				GameOver: &GameOver{
					Winner: "1",
					Reason: core.Timeout,
				},
				Actions: []LogAction{
					// 相手の行動でも開示
					{PlayerId: "1", Turn: 0, At: 1, T: core.MoveAction, To: 1, Direction: core.East},
					{},
				},
			}},
		}
		for _, tt := range test {
			var battle = New()
			tt.preProcess(battle)
			output, err := battle.GetLogs(ctx, tt.input)
			assert.Equal(t, tt.expErr, err)
			if tt.expOutput == nil {
				assert.Nil(t, output)
			} else {
				assert.Equal(t, tt.expOutput.RequireAction, output.RequireAction)
				assert.Equal(t, tt.expOutput.RequireDeployAction, output.RequireDeployAction)
				assert.Equal(t, tt.expOutput.NumTurn, output.NumTurn)
				assert.Equal(t, tt.expOutput.GameOver, output.GameOver)
				assert.Len(t, output.Actions, len(tt.expOutput.Actions))
			}
			if len(tt.expOutput.Actions) > 0 {
				assert.Equal(t, tt.expOutput.Actions[0].PlayerId, output.Actions[0].PlayerId)
				assert.Equal(t, tt.expOutput.Actions[0].T, output.Actions[0].T)
				assert.Equal(t, tt.expOutput.Actions[0].ActionResult, output.Actions[0].ActionResult)
				assert.Equal(t, tt.expOutput.Actions[0].At, output.Actions[0].At)
				assert.Equal(t, tt.expOutput.Actions[0].To, output.Actions[0].To)
				assert.Equal(t, tt.expOutput.Actions[0].Turn, output.Actions[0].Turn)
				assert.Equal(t, tt.expOutput.Actions[0].Direction, output.Actions[0].Direction)
			}
		}

	})

	t.Run("gameOver", func(t *testing.T) {
		t.Parallel()

		test := []struct {
			input *core.Action
			exp   *GameOver
		}{
			{nil, nil},
			{&core.Action{
				PlayerId:     "a",
				ActionResult: core.Hit,
				T:            core.TorpedoFireAction,
				Timestamp:    time.Now(),
			}, &GameOver{Winner: "a", Reason: core.TorpedoHit}},
			{&core.Action{
				PlayerId:     "a",
				ActionResult: core.Hit,
				T:            core.MineTriggerAction,
				Timestamp:    time.Now(),
			}, &GameOver{Winner: "a", Reason: core.MineHit}},
			{&core.Action{
				PlayerId:  "a",
				Timestamp: time.Now().Add(-time.Duration(1 * time.Minute)),
			}, &GameOver{
				Winner: "a",
				Reason: core.Timeout,
			}},
		}
		for _, tt := range test {
			battle := BattleService{}
			battle.init()
			act := battle.gameOver(tt.input)
			if tt.exp == nil {
				assert.Nil(t, act)
			} else {
				assert.Equal(t, tt.exp, act)
			}
		}
	})

	t.Run("PlayGame", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancelCause(context.Background())
		defer cancel(nil)
		noError := func(err error) {
			assert.NoError(t, err)
			if err != nil {
				cancel(err)
			}
		}

		var (
			gameId  = "PlayGame"
			playerA = "playerA"
			playerB = "playerB"
			battle  = New()
		)
		battle.setGame(core.Game{
			GameId:     gameId,
			PlayerIds:  [2]string{playerA, playerB},
			OceanMap:   core.DefaultOceanMap,
			Islands:    []core.Sector{20, 35},
			Submarines: map[string]core.Submarine{playerA: core.DefaultSubmarine, playerB: core.DefaultSubmarine},
		})
		wait := func(playerId string) {
			waitInput := WaitTurnInput{
				GameId:   gameId,
				PlayerId: playerId,
			}
			done, err := battle.WaitTurn(ctx, &waitInput)
			noError(err)
			select {
			case <-ctx.Done():
			case <-done:
			}
		}

		// PlayerAが先攻
		err := battle.DeploySubmarineAndMines(ctx, &DeploySubmarineAndMinesInput{
			GameId:   gameId,
			PlayerId: playerA,
			At:       10,
			Mines:    []int8{13, 28},
		})
		noError(err)
		go func() {
			me := playerA
			wait(me)
			// 2
			err = battle.TriggerMine(ctx, &ActionInput{
				GameId:   gameId,
				PlayerId: me,
				At:       13,
			})
			noError(err)
			wait(me)
			// 3
			err = battle.TriggerMine(ctx, &ActionInput{
				GameId:   gameId,
				PlayerId: me,
				At:       28,
			})
			noError(err)
			wait(me)
			// 4
			err = battle.Move(ctx, &ActionInput{
				GameId:   gameId,
				PlayerId: me,
				At:       16,
			})
			noError(err)
			wait(me)
			// 5
			err = battle.FireTorpedo(ctx, &ActionInput{
				GameId:   gameId,
				PlayerId: me,
				At:       15,
			})
			noError(err)
			wait(me)
			// end
			logs, err := battle.GetLogs(ctx, &GetLogsInput{
				GameId:   gameId,
				PlayerId: me,
			})
			assert.NoError(t, err)
			assert.NotNil(t, logs.GameOver)
			cancel(nil)
		}()
		go func() {
			me := playerB
			wait(me)
			// 1
			err := battle.DeploySubmarineAndMines(ctx, &DeploySubmarineAndMinesInput{
				GameId:   gameId,
				PlayerId: me,
				At:       27,
				Mines:    []int8{5, 19},
			})
			noError(err)
			wait(me)
			// 2
			err = battle.FireTorpedo(ctx, &ActionInput{
				GameId:   gameId,
				PlayerId: me,
				At:       21,
			})
			noError(err)
			wait(me)
			// 3
			err = battle.Move(ctx, &ActionInput{
				GameId:   gameId,
				PlayerId: me,
				At:       21,
			})
			noError(err)
			wait(me)
			// 4
			err = battle.Move(ctx, &ActionInput{
				GameId:   gameId,
				PlayerId: me,
				At:       15,
			})
			noError(err)
			wait(me)
			// end
			logs, err := battle.GetLogs(ctx, &GetLogsInput{
				GameId:   gameId,
				PlayerId: me,
			})
			assert.NoError(t, err)
			assert.NotNil(t, logs.GameOver)
			cancel(nil)
		}()

		<-ctx.Done()
		err = context.Cause(ctx)
		assert.ErrorIs(t, err, context.Canceled)
		logs, err := battle.GetLogs(ctx, &GetLogsInput{
			GameId:   gameId,
			PlayerId: playerA,
		})
		assert.NoError(t, err)
		assert.NotNil(t, logs.GameOver)
		assert.Equal(t, playerA, logs.GameOver.Winner)
		assert.Equal(t, core.TorpedoHit, logs.GameOver.Reason)
	})
}

func TestRepository(t *testing.T) {
	t.Run("Game", func(t *testing.T) {
		var battle = New()

		want1 := core.Game{
			GameId:    "1",
			OceanMap:  core.DefaultOceanMap,
			PlayerIds: [2]string{"a", "b"},
			Islands:   []core.Sector{1, 2},
			Submarines: map[string]core.Submarine{
				"a": core.DefaultSubmarine,
				"b": core.DefaultSubmarine,
			}}
		want2 := core.Game{
			GameId:    "2",
			OceanMap:  core.DefaultOceanMap,
			PlayerIds: [2]string{"c", "d"},
			Islands:   []core.Sector{3, 4},
			Submarines: map[string]core.Submarine{
				"c": core.DefaultSubmarine,
				"d": core.DefaultSubmarine,
			}}
		err := battle.setGame(want1)
		assert.NoError(t, err)
		err = battle.setGame(want2)
		assert.NoError(t, err)

		// empty
		dist, err := battle.getGame("not found")
		assert.ErrorIs(t, err, ErrGameNotFound)
		assert.Nil(t, dist)

		// get
		dist, err = battle.getGame("1")
		assert.NoError(t, err)
		testEqualGame(t, want1, *dist)
		dist, err = battle.getGame("2")
		assert.NoError(t, err)
		testEqualGame(t, want2, *dist)

		// delete 1
		err = battle.deleteGame("1")
		assert.NoError(t, err)
		_, err = battle.getGame("1")
		assert.ErrorIs(t, err, ErrGameNotFound)
		dist, err = battle.getGame("2")
		assert.NoError(t, err)
		testEqualGame(t, want2, *dist)

		// delete 2
		err = battle.deleteGame("2")
		assert.NoError(t, err)
		_, err = battle.getGame("1")
		assert.ErrorIs(t, err, ErrGameNotFound)
		_, err = battle.getGame("2")
		assert.ErrorIs(t, err, ErrGameNotFound)
	})

	t.Run("Action", func(t *testing.T) {
		var battle = New()

		dist, err := battle.getLatestAction(battle.db.DB, "a")
		assert.NoError(t, err)
		assert.Nil(t, dist)
		prevs, err := battle.getPrevActions("a")
		assert.Nil(t, err)
		assert.Empty(t, prevs)

		want1 := core.Action{
			GameId:   "a",
			PlayerId: "1",
			At:       core.Sector(1),
			Mines:    []core.Sector{10, 20},
		}
		want2 := core.Action{
			GameId:   "a",
			PlayerId: "2",
			At:       core.Sector(8),
			Mines:    []core.Sector{1, 9},
		}
		err = battle.appendAction(want1)
		assert.NoError(t, err)
		err = battle.appendAction(want2)
		assert.NoError(t, err)

		dist, err = battle.getLatestAction(battle.db.DB, "a")
		assert.NoError(t, err)
		testEqualAction(t, want2, *dist)

		dist, err = battle.getPrevAction("a", "2")
		assert.NoError(t, err)
		testEqualAction(t, want2, *dist)

		dist, err = battle.getPrevAction("a", "1")
		assert.NoError(t, err)
		testEqualAction(t, want1, *dist)

		prevs, err = battle.getPrevActions("a")
		assert.Nil(t, err)
		testEqualAction(t, want1, *prevs["1"])
		testEqualAction(t, want2, *prevs["2"])

		// 新しい順
		actions, err := battle.getAllAction("a")
		assert.NoError(t, err)
		assert.Len(t, actions, 2)
		testEqualAction(t, want2, actions[0])
		testEqualAction(t, want1, actions[1])

		want3 := core.Action{
			GameId:       "a",
			PlayerId:     "1",
			At:           core.Sector(2),
			T:            core.MineTriggerAction,
			From:         want1.At,
			To:           core.Sector(10),
			ActionResult: core.FullSpeedAhead,
			Mines:        []core.Sector{20},
		}
		err = battle.appendAction(want3)
		assert.NoError(t, err)
		dist, err = battle.getLatestAction(battle.db.DB, "a")
		assert.NoError(t, err)
		testEqualAction(t, want3, *dist)
		dist, err = battle.getPrevAction("a", "1")
		assert.NoError(t, err)
		testEqualAction(t, want3, *dist)
		dist, err = battle.getPrevAction("a", "2")
		assert.NoError(t, err)
		testEqualAction(t, want2, *dist)
		prevs, err = battle.getPrevActions("a")
		assert.Nil(t, err)
		testEqualAction(t, want3, *prevs["1"])
		testEqualAction(t, want2, *prevs["2"])

		// 新しい順
		actions, err = battle.getAllAction("a")
		assert.NoError(t, err)
		assert.Len(t, actions, 3)
		testEqualAction(t, want3, actions[0])
		testEqualAction(t, want2, actions[1])
		testEqualAction(t, want1, actions[2])
	})
}

func testEqualGame(t *testing.T, exp, act core.Game) {
	t.Helper()
	assert.Equal(t, exp.GameId, act.GameId)
	assert.Equal(t, exp.Islands, act.Islands)
	assert.Equal(t, exp.OceanMap, act.OceanMap)
	assert.Equal(t, exp.PlayerIds, act.PlayerIds)
	assert.Equal(t, exp.Submarines, act.Submarines)
	assert.NotZero(t, act.Timestamp)
}

func testEqualAction(t *testing.T, exp, act core.Action) {
	t.Helper()
	assert.Equal(t, exp.GameId, act.GameId)
	assert.Equal(t, exp.PlayerId, act.PlayerId)
	assert.Equal(t, exp.At, act.At)
	assert.Equal(t, exp.Mines, act.Mines)
	assert.Equal(t, exp.T, act.T)
	assert.Equal(t, exp.From, act.From)
	assert.Equal(t, exp.To, act.To)
	assert.Equal(t, exp.ActionResult, act.ActionResult)
	assert.NotZero(t, act.Timestamp)
}
