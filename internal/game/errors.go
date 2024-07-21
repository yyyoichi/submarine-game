package game

import "errors"

var (
	ErrIsnotYourTurn = errors.New("is not your trun")
	ErrOutOfCampSize = errors.New("out of camp size")
	ErrInvalidAction = errors.New("invalid action")
	ErrInvalidCamp   = errors.New("invalid camp")
	ErrGameIsOver    = errors.New("game is over")
)
