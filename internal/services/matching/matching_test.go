package matching

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yyyoichi/hookdb"
	"github.com/yyyoichi/submarine-game/internal/store"
)

func TestMatching(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	store, err := store.New()
	require.NoError(t, err)

	matching := &MatchingService{
		store:              store,
		waitTickerDuration: time.Millisecond * 2,
	}

	const C = 10_000

	var outputCh = make(chan *JoinOutput)
	go func() {
		defer close(outputCh)
		var wg sync.WaitGroup
		for range C {
			wg.Add(1)
			go func() {
				defer wg.Done()
				output, err := matching.Join()
				assert.NoError(t, err)
				if !output.Matched {
					gameCh := matching.Wait(ctx, output.PlayerId)
					game := <-gameCh
					assert.NotEmpty(t, game)
					output.GameId = game[0]
					output.EnemyId = game[1]
					output.Matched = true
				}
				outputCh <- output
			}()
		}
		wg.Wait()
	}()

	go func() {
		defer cancel(nil)
		var outputs = make(map[string][]JoinOutput)
		for output := range outputCh {
			assert.NotNil(t, output)
			if _, found := outputs[output.GameId]; !found {
				outputs[output.GameId] = make([]JoinOutput, 0, 2)
			}
			outputs[output.GameId] = append(outputs[output.GameId], *output)
		}

		assert.Len(t, outputs, C/2)
		for _, o := range outputs {
			assert.Len(t, o, 2)
			assert.Equal(t, o[0].GameId, o[1].GameId)
			assert.Equal(t, o[0].PlayerId, o[1].EnemyId)
			assert.Equal(t, o[0].EnemyId, o[1].PlayerId)
			assert.True(t, o[0].Matched)
			assert.True(t, o[1].Matched)
		}
	}()
	<-ctx.Done()
	assert.ErrorIs(t, context.Cause(ctx), context.Canceled)
}

func TestMatchingMatch(t *testing.T) {
	t.Run("Match", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancelCause(context.Background())
		defer cancel(nil)

		m := MatchingService{
			db: hookdb.New(),
		}

		gamePalyerCh := make(chan [2]string)

		var wg sync.WaitGroup
		go func() {
			defer close(gamePalyerCh)
			for range 10_000 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					playerId, gameIdCh := m.Match(ctx, cancel)
					select {
					case <-ctx.Done():
						return
					case gameId, ok := <-gameIdCh:
						if !ok {
							return
						}
						assert.NotEmpty(t, gameId)
						gamePalyerCh <- [2]string{gameId, playerId}
					}
				}()
			}
			wg.Wait()
		}()

		gamePlayers := make(map[string][]string)
		for gamePlayer := range gamePalyerCh {
			if _, found := gamePlayers[gamePlayer[0]]; !found {
				gamePlayers[gamePlayer[0]] = make([]string, 0, 2)
			}
			gamePlayers[gamePlayer[0]] = append(gamePlayers[gamePlayer[0]], gamePlayer[1])
		}

		assert.Len(t, gamePlayers, 5_000)
		for _, players := range gamePlayers {
			assert.Len(t, players, 2)
			assert.NotEqual(t, players[0], players[1])
			assert.NotEmpty(t, players[0])
			assert.NotEmpty(t, players[1])
		}

		assert.NoError(t, context.Cause(ctx))
	})

	t.Run("Reave", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancelCause(context.Background())
		defer cancel(nil)

		m := MatchingService{
			db: hookdb.New(),
		}
		timeCtx, timecancel := context.WithTimeout(ctx, time.Millisecond*1)
		defer timecancel()

		done := make(chan struct{})
		go func() {
			defer close(done)
			// タイムアウトでキャンセルされる
			_, ch := m.Match(timeCtx, cancel)
			_, ok := <-ch
			assert.False(t, ok)
		}()
		<-done
		assert.ErrorIs(t, context.Cause(timeCtx), context.DeadlineExceeded)
		assert.Empty(t, m.waitPlayer)
	})
}
