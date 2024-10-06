package battle

import "github.com/yyyoichi/submarine-game/internal/core"

type DeploySubmarineAndMinesInput struct {
	GameId   string
	PlayerId string
	At       int8
	Mines    []int8
}

type MoveInput struct {
	GameId   string
	PlayerId string
	At       int8
}

type GetValidPrevActionsInput struct {
	Game            core.Game
	PlayerId        string
	ExpSectorStatus core.SectorStatus
	At              core.Sector
}
