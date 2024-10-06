package battle

type DeploySubmarineAndMinesInput struct {
	GameId   string
	PlayerId string
	At       int8
	Mines    []int8
}
