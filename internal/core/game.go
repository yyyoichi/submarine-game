package core

import (
	"slices"
	"time"
)

// ゲーム開始設定
type Game struct {
	GameId     string
	PlayerIds  [2]string
	OceanMap   OceanMap
	Islands    []Sector
	Submarines map[string]Submarine
	// ゲーム開始時間
	Timestamp time.Time
}

// prevActionからplayerIdの海域状態を返す。
func (g Game) SectorStatus(playerId string, prev Action, ats ...Sector) map[Sector][]SectorStatus {
	if len(ats) == 0 {
		ats = g.OceanMap.Sectors()
	}
	submarine := g.Submarines[playerId]
	canMoves := g.OceanMap.EnableSectors(prev.At, submarine.Movable)
	canFireTorpedo := g.OceanMap.EnableSectors(prev.At, submarine.TorpedoTargetable)
	canTriggerMine := prev.Mines

	var resp = make(map[Sector][]SectorStatus, len(ats))
	for _, at := range ats {
		resp[at] = make([]SectorStatus, 0, 3)
		if slices.Contains(g.Islands, at) {
			resp[at] = append(resp[at], IslandSector)
			continue
		}
		if prev.At == at {
			resp[at] = append(resp[at], SelfOccupied)
			continue
		}
		if slices.Contains(canMoves, at) {
			resp[at] = append(resp[at], CanMove)
		}
		if slices.Contains(canFireTorpedo, at) {
			resp[at] = append(resp[at], CanFireTorpedo)
		}
		if slices.Contains(canTriggerMine, at) {
			resp[at] = append(resp[at], CanTriggerMine)
		}
	}
	return resp
}

// atの位置にいるときのtoに対する攻撃結果
func (g Game) ActionResult(at Sector, to Sector) ActionResult {
	if at == to {
		return Hit
	}
	sectors := g.OceanMap.EnableSectors(at, [][2]int8{
		{-1, -1}, {00, -1}, {01, -1},
		{-1, 00} /**    */, {01, 00},
		{-1, 01}, {00, 01}, {01, 01},
	})
	if slices.Contains(sectors, to) {
		// atの周囲にtoが含まれる
		return HardToStarboard
	}
	return FullSpeedAhead
}

func (g Game) Enemy(playerId string) string {
	if g.PlayerIds[0] == playerId {
		return g.PlayerIds[1]
	}
	return g.PlayerIds[0]
}

func (g Game) Since() time.Duration {
	return time.Since(g.Timestamp)
}

type GameOverReason int8

const (
	Ongoing    GameOverReason = 0  // まだゲーム中
	TorpedoHit GameOverReason = 1  // 魚雷が命中
	MineHit    GameOverReason = 2  // 機雷が命中
	Timeout    GameOverReason = 10 // 時間切れ
)
