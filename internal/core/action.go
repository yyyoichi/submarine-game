package core

import "time"

type Action struct {
	GameId   string
	PlayerId string
	// 現在位置
	At Sector
	// 行動内容
	T ActionType
	// 行動対象位置
	To Sector
	// 行動結果
	ActionResult ActionResult
	// 残機雷位置
	Mines []Sector
	// 行動時刻
	Timestamp time.Time
}

// 行動タイプ
type ActionType int8

const (
	MoveAction        ActionType = 1
	TorpedoFireAction ActionType = 2
	MineTriggerAction ActionType = 3
)

// 行動結果
type ActionResult int8

const (
	// なし
	NoResult ActionResult = 0
	// ヨーソロー
	FullSpeedAhead ActionResult = 1
	// 面舵一杯
	HardToStarboard ActionResult = 2
	// 命中
	Hit ActionResult = 3
)
