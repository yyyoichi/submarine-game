package core

var DefaultOceanMap = OceanMap{
	Id:          201,
	H:           6,
	W:           6,
	IslandCount: 2,
}

// 対戦ボード
type OceanMap struct {
	Id   int32
	H, W int8
	// 島数
	IslandCount int8
}

// atからの相対位置rssを絶対位置で返す。
func (om OceanMap) EnableSectors(at Sector, rss RelativeSectors) []Sector {
	var resp = make([]Sector, 0, len(rss))
	atX, atY := int8(at)%om.W, int8(at)/om.W
	for _, s := range rss {
		x := atX + s[0]
		if x < 0 || om.W <= x {
			continue
		}
		y := atY + s[1]
		if y < 0 || om.H <= y {
			continue
		}
		resp = append(resp, Sector(om.W*y+x))
	}
	return resp
}

func (om OceanMap) Sectors() []Sector {
	var resp = make([]Sector, om.H*om.W)
	for i := range resp {
		resp[i] = Sector(i)
	}
	return resp
}

// 2次元スライスで区域を返す。
// ex) [[0,1,2], [3,4,5], [6,7,8]]
func (om OceanMap) Lines() [][]Sector {
	var resp = make([][]Sector, om.H)
	for y := range om.H {
		line := make([]Sector, om.W)
		for x := range om.W {
			line[x] = Sector(om.W*y + x)
		}
		resp[y] = line
	}
	return resp
}
