package core

type (
	Sector int8
	// 相対位置[x,y]
	RelativeSectors [][2]int8
)

// 相対位置の方向を転換する
func (rss RelativeSectors) Trun(d Direction) RelativeSectors {
	var f func([2]int8) [2]int8
	switch d {
	case North:
		f = func(i [2]int8) [2]int8 { return i }
	case East:
		f = func(i [2]int8) [2]int8 {
			return [2]int8{i[1] * -1, i[0]}
		}
	case South:
		f = func(i [2]int8) [2]int8 {
			return [2]int8{i[0] * -1, i[1] * -1}
		}
	case West:
		f = func(i [2]int8) [2]int8 {
			return [2]int8{i[1], i[0] * -1}
		}
	}
	var resp = make([][2]int8, len(rss))
	for i, s := range rss {
		resp[i] = f(s)
	}
	return RelativeSectors(resp)
}

// 海域の状態
type SectorStatus int8

const (
	SelfOccupied   SectorStatus = 1
	CanMove        SectorStatus = 2
	CanFireTorpedo SectorStatus = 3
	CanTriggerMine SectorStatus = 4
	IslandSector   SectorStatus = 99
)
