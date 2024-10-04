package core

var DefaultOceanMap = OceanMap{
	Id:          201,
	H:           6,
	W:           6,
	IsLandCount: 2,
}

// 対戦ボード
type OceanMap struct {
	Id   uint8
	H, W uint8
	// 島数
	IsLandCount uint8
}
