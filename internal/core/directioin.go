package core

var (
	North Direction = 1
	South Direction = 2
	East  Direction = 3
	West  Direction = 4
)

var (
	Top   Rotation = 0
	Right Rotation = 1
	Down  Rotation = 2
	Left  Rotation = 3
)

type (
	// 方向
	Direction int8
	// 右回転度
	Rotation int8
)
