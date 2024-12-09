package matching

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yyyoichi/hookdb"
)

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
					out, err := m.Match(ctx)
					assert.NoError(t, err)
					select {
					case <-ctx.Done():
						return
					case v, ok := <-out:
						if !ok {
							return
						}
						assert.NotEmpty(t, v)
						gamePalyerCh <- [2]string{v.GameId, v.PlayerId}
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

	t.Run("Leave", func(t *testing.T) {
		t.Parallel()
		m := MatchingService{
			db: hookdb.New(),
		}
		timeCtx, timecancel := context.WithTimeout(context.Background(), time.Millisecond*1)
		defer timecancel()

		done := make(chan struct{})
		go func() {
			defer close(done)
			// タイムアウトでキャンセルされる
			ch, err := m.Match(timeCtx)
			assert.NoError(t, err)
			_, ok := <-ch
			assert.False(t, ok)
		}()
		<-done
		assert.ErrorIs(t, context.Cause(timeCtx), context.DeadlineExceeded)
		assert.Empty(t, m.waitPlayer)
	})
}
