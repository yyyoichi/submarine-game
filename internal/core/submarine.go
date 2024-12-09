package core

var (
	DefaultSubmarine = Submarine{
		Id: 101,
		Movable: [][2]int8{
			/*     */ {00, -1},
			{-1, 00} /**    */, {01, 00},
			/*     */ {00, 01},
		},
		TorpedoTargetable: [][2]int8{
			{-1, -1}, {00, -1}, {01, -1},
			{-1, 00} /**    */, {01, 00},
			{-1, 01}, {00, 01}, {01, 01},
		},
		MineCount: 2,
	}
)

// 潜水艦
type Submarine struct {
	Id                uint8
	Movable           RelativeSectors
	TorpedoTargetable RelativeSectors
	MineCount         uint8
}
