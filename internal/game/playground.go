package game

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

// 6x6
const lineSize = 6
const campSize = 36

type Playground struct {
	gameById   map[string]*Game
	gameByUser map[string]*Game
}

func NewPlayground() *Playground {
	return &Playground{
		gameById:   map[string]*Game{},
		gameByUser: map[string]*Game{},
	}
}

func (pg *Playground) NewGame(users [2]string) *Game {
	island1 := rand.Int31n(campSize)
	island2 := rand.Int31n(campSize)
	g := &Game{
		Id:     uuid.NewString(),
		Users:  users,
		Island: [2]uint32{uint32(island1), uint32(island2)},
		clock: clock{
			createdAt: time.Now(),
		},

		mu:       sync.RWMutex{},
		NextUser: users[1],
		mines: map[string][]uint32{
			users[0]: {},
			users[1]: {},
		},
	}
	pg.gameById[g.Id] = g
	pg.gameByUser[users[0]] = g
	pg.gameByUser[users[1]] = g
	return g
}

func (pg *Playground) Use(id string) (*Game, error) {
	gm, found := pg.gameById[id]
	if !found {
		return nil, ErrInvalidGameId
	}
	return gm, nil
}

func (pg *Playground) Found(user string) *Game {
	return pg.gameByUser[user]
}

func (pg *Playground) DeleteRunner(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(1) * time.Minute):
			pg.deleteTimeoutGame()
		}
	}
}

func (pg *Playground) deleteTimeoutGame() {
	// 開始後1時間経過のゲームは削除する。
	for id, g := range pg.gameById {
		hours := time.Since(g.clock.createdAt).Hours()
		if hours > 1 {
			user1 := g.Users[0]
			user2 := g.Users[1]
			delete(pg.gameById, id)
			delete(pg.gameByUser, user1)
			delete(pg.gameByUser, user2)
		}
	}
}
