package animation

import (
	"fmt"
	"time"
)

type Animation struct {
	itime time.Time // 初期化時間
	// レンダリング設定。[0]がパーセンテージ、[1]が値。0と1パーセンテージは必須。
	// パーセンテージの昇順で設定すること。
	ByPercentage [][2]float32
	Duration     time.Duration // アニメーション時間
}

func (a *Animation) Clear() {
	a.itime = time.Now()
}

func (a *Animation) Value() float32 {
	for i := range len(a.ByPercentage) - 1 {
		current := a.ByPercentage[i]
		next := a.ByPercentage[i+1]
		p := a.nowPercentage()
		if next[0] <= p {
			continue
		}
		// current <= p < next
		// 増加分に比例して値を計算
		return current[1] + (next[1]-current[1])*(p-current[0])/(next[0]-current[0])
	}
	return a.ByPercentage[len(a.ByPercentage)-1][1]
}

// 現在のパーセンテージを取得する
func (a *Animation) nowPercentage() float32 {
	sub := time.Since(a.itime)
	return float32(sub) / float32(a.Duration)
}

func (a *Animation) String() string {
	sub := time.Since(a.itime)
	return fmt.Sprintf("Animation{ByPercentage: %v, Duration: %v, sub: %v, Value: %f}", a.ByPercentage, a.Duration, sub, float32(sub)/float32(a.Duration))
}
