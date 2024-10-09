package battle

import (
	"time"

	"github.com/yyyoichi/submarine-game/internal/core"
)

type DeploySubmarineAndMinesInput struct {
	GameId   string
	PlayerId string
	At       int8
	Mines    []int8
}

type ActionInput struct {
	GameId   string
	PlayerId string
	At       int8
}

type GetLogsInput struct {
	GameId   string
	PlayerId string
}

type GetLogsOutput struct {
	// なにかしらのアクションを要求するか
	RequireAction bool
	// 初回の行動要求
	RequireDeployAction bool
	// セクターごとに利用可能な行動タイプ
	SectorActionsMap map[core.Sector][]core.ActionType
	// ゲーム終了
	GameOver *GameOver

	// ターン数
	NumTurn int
	Actions []LogAction
	// 行動許容時間ミリ秒
	TimeoutDurationMSec int64
	// 行動期限
	Timeout time.Time
}

type LogAction struct {
	PlayerId string
	Turn     int
	// 現在位置
	At core.Sector
	// 行動内容
	T core.ActionType
	// 行動対象位置
	To core.Sector
	// 行動結果
	ActionResult core.ActionResult
}

type GameOver struct {
	Winner string
	Reason core.GameOverReason
}

type GetValidPrevActionsInput struct {
	Game            core.Game
	PlayerId        string
	ExpSectorStatus core.SectorStatus
	At              core.Sector
}
