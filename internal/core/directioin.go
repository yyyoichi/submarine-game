package core

// 向方角
var (
	UnknownDirection Direction = -1
	North            Direction = 0
	East             Direction = 1
	South            Direction = 2
	West             Direction = 3
)

// 回転方向
var (
	Front Rotation = 0
	Right Rotation = 1
	Back  Rotation = 2
	Left  Rotation = 3
)

type (
	// 方角
	Direction int8
	// 右回転度
	Rotation int8
)

// 向き先を変える。
func (d Direction) Rotate(r Rotation) Direction {
	return Direction((int8(d) + int8(r)) % 4)
}

func (d Direction) Vector() RelativeSectors {
	switch d {
	case North:
		return [][2]int8{{0, -1}}
	case East:
		return [][2]int8{{1, 0}}
	case South:
		return [][2]int8{{0, 1}}
	case West:
		return [][2]int8{{-1, 0}}
	}
	return nil
}
