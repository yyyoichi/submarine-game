package matching

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
					errCh := output.WaitMatching(ctx)
					err := <-errCh
					assert.NoError(t, err)
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
