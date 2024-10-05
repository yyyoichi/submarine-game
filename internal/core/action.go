package core

import (
	"math"
	"time"
)

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
	// NOTE badgerのReverseイテレーションができないので応急処置
	// MaxInt64から現在時刻を引いた行動時刻
	ReverseUnixNano int64
	Timestamp       time.Time
}

func (a *Action) SetTimestamp() {
	a.Timestamp = time.Now()
	// 最大値から引いて反転
	a.ReverseUnixNano = math.MaxInt64 - a.Timestamp.UnixNano()
}

func (a *Action) RestoreTime() {
	nano := math.MaxInt64 - a.ReverseUnixNano
	a.Timestamp = time.Unix(0, nano)
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
