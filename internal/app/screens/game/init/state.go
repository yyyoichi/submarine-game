package gameinit

type (
	pstate  string
	trigger string
)

const (
	pstate_sector  pstate = "sector"
	pstate_mines   pstate = "mines"
	pstate_pending pstate = "pending"
)

const (
	trigger_sector  trigger = "sector"
	trigger_mines   trigger = "mines"
	trigger_pending trigger = "pending"
)

type istate struct {
	selectedSector *int
	selectedMines  []int
}
