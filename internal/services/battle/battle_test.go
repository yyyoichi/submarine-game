package battle

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yyyoichi/submarine-game/internal/core"
	"github.com/yyyoichi/submarine-game/internal/store"
)

func TestBattleService(t *testing.T) {
	ctx := context.Background()
	store, err := store.New()
	require.NoError(t, err)

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
			battle := BattleService{Store: store}
			tt.preProcess(&battle)
			err := battle.DeploySubmarineAndMines(ctx, tt.input)
			if tt.exp == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.exp)
			}
		}
	})
}

func TestRepository(t *testing.T) {
	t.Run("Game", func(t *testing.T) {
		var err error
		battle := BattleService{}
		battle.Store, err = store.New()
		require.NoError(t, err)

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
		err = battle.setGame(want1)
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
		testEqualGame(t, want1, dist.Game)
		dist, err = battle.getGame("2")
		assert.NoError(t, err)
		testEqualGame(t, want2, dist.Game)

		// delete 1
		err = battle.deleteGame("1")
		assert.NoError(t, err)
		_, err = battle.getGame("1")
		assert.ErrorIs(t, err, ErrGameNotFound)
		dist, err = battle.getGame("2")
		assert.NoError(t, err)
		testEqualGame(t, want2, dist.Game)

		// delete 2
		err = battle.deleteGame("2")
		assert.NoError(t, err)
		_, err = battle.getGame("1")
		assert.ErrorIs(t, err, ErrGameNotFound)
		_, err = battle.getGame("2")
		assert.ErrorIs(t, err, ErrGameNotFound)
	})

	t.Run("Action", func(t *testing.T) {
		var err error
		battle := BattleService{}
		battle.Store, err = store.New()
		require.NoError(t, err)

		dist, err := battle.getLatestAction("a")
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

		dist, err = battle.getLatestAction("a")
		assert.NoError(t, err)
		testEqualAction(t, want2, dist.Action)

		dist, err = battle.getPrevAction("a", "2")
		assert.NoError(t, err)
		testEqualAction(t, want2, dist.Action)

		dist, err = battle.getPrevAction("a", "1")
		assert.NoError(t, err)
		testEqualAction(t, want1, dist.Action)

		prevs, err = battle.getPrevActions("a")
		assert.Nil(t, err)
		testEqualAction(t, want1, prevs["1"].Action)
		testEqualAction(t, want2, prevs["2"].Action)

		want3 := core.Action{
			GameId:       "a",
			PlayerId:     "1",
			At:           core.Sector(2),
			T:            core.MineTriggerAction,
			To:           core.Sector(10),
			ActionResult: core.FullSpeedAhead,
			Mines:        []core.Sector{20},
		}
		err = battle.appendAction(want3)
		assert.NoError(t, err)
		dist, err = battle.getLatestAction("a")
		assert.NoError(t, err)
		testEqualAction(t, want3, dist.Action)
		dist, err = battle.getPrevAction("a", "1")
		assert.NoError(t, err)
		testEqualAction(t, want3, dist.Action)
		dist, err = battle.getPrevAction("a", "2")
		assert.NoError(t, err)
		testEqualAction(t, want2, dist.Action)
		prevs, err = battle.getPrevActions("a")
		assert.Nil(t, err)
		testEqualAction(t, want3, prevs["1"].Action)
		testEqualAction(t, want2, prevs["2"].Action)
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
	assert.Equal(t, exp.To, act.To)
	assert.Equal(t, exp.ActionResult, act.ActionResult)
	assert.NotZero(t, act.Timestamp)
}
