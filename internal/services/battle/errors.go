package battle

import "errors"

var (
	ErrGameNotFound      = errors.New("game not found")
	ErrInvalidMineCount  = errors.New("invalid mine count")
	ErrInvalidActionType = errors.New("invalid action type")
	ErrTimeout           = errors.New("timeout")
	ErrInvalidSector     = errors.New("invalid sector")
)
