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

type (
	istate struct {
		selectedSector *int
		selectedMines  []int
		config         IStateConfig
	}
	IStateConfig struct {
		MineCount     int
		IslandSectors []int
	}
)

func newState(config IStateConfig) istate {
	is := istate{config: config}
	is.selectedMines = make([]int, 0, config.MineCount)
	return is
}
