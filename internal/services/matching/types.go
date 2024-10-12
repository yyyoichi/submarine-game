package matching

import "context"

type JoinOutput struct {
	PlayerId     string
	GameId       string
	EnemyId      string
	Matched      bool
	WaitMatching func(context.Context) error
}
