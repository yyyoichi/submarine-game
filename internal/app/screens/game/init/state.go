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
		selectedSector *uint8
		selectedMines  []uint8
		config         IStateConfig
	}
	IStateConfig struct {
		MineCount uint8
	}
)

func newState(config IStateConfig) istate {
	is := istate{config: config}
	is.selectedMines = make([]uint8, 0, config.MineCount)
	return is
}
