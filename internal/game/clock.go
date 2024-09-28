package game

import "time"

type clock struct {
	// ゲーム開始
	createdAt time.Time
	// 最近の行動
	latestAt *time.Time
}

var timeout = time.Duration(31) * time.Second
var timeoutDiff = time.Duration(500) * time.Millisecond

// 初回行動期限か
func (c clock) isFirstActionTimeout() bool {
	return c.getFirstActionTimeout().Add(timeoutDiff).Before(time.Now())
}

// 行動期限か
func (c clock) isActionTimeout() bool {
	return c.getActionTimeout().Add(timeoutDiff).Before(time.Now())
}

// 初回行動期限（ゲーム開始から30秒）を取得する。
func (c clock) getFirstActionTimeout() time.Time {
	return c.createdAt.Add(timeout)
}

// 行動期限（最後の行動から30秒）を取得する。
func (c clock) getActionTimeout() time.Time {
	if c.latestAt == nil {
		return c.createdAt.Add(timeout)
	}
	return c.latestAt.Add(timeout)
}
