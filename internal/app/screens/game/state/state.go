package gamestate

import "github.com/qmuntal/stateless"

type (
	State   string
	Trigger string
)

const (
	State_Init State = "init"
	State_Play State = "play"
	State_Over State = "over"
)

const (
	Trigger_Play Trigger = "play"
	Trigger_Over Trigger = "over"
)

func New() stateless.State {
	sm := stateless.NewStateMachine(State_Init)
	sm.Configure(State_Init).Permit(Trigger_Play, State_Play)
	sm.Configure(State_Play).Permit(Trigger_Over, State_Over)
	return sm
}
