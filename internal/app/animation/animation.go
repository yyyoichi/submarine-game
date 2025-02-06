package animation

import (
	"fmt"
	"time"
)

var (
	DefaultFade = func() *Animation {
		return &Animation{
			ByPercentage: [][2]float32{{0, 0}, {0.1, 1}, {0.8, 1}, {1, 0}},
			Duration:     time.Duration(time.Millisecond * 1200),
		}
	}
	DefaultFadeOut = func() *Animation {
		return &Animation{
			ByPercentage: [][2]float32{{0, 1}, {0.8, 1}, {1, 0}},
			Duration:     time.Duration(time.Millisecond * 1200),
		}
	}
	DefaultLoop = func() *LoopAnimation {
		return &LoopAnimation{
			Animation: Animation{
				ByPercentage: [][2]float32{{0, 0.1}, {0.3, 0.3}, {1, 0.1}},
				Duration:     time.Duration(time.Millisecond * 900),
			},
		}
	}
)

type LoopAnimation struct {
	Animation
}

func (a *LoopAnimation) Value() float32 {
	if i := a.Animation.step(); i == len(a.Animation.ByPercentage)-1 {
		a.Clear()
	}
	return a.Animation.Value()
}

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

func (a *Animation) Zero() {
	a.itime = time.Now().Add(-a.Duration)
}

func (a *Animation) Value() float32 {
	if len(a.ByPercentage) == 0 {
		return 0
	}
	p := a.nowPercentage()
	i := a.step()
	if i == len(a.ByPercentage)-1 {
		return a.ByPercentage[i][1]
	}
	current := a.ByPercentage[i]
	next := a.ByPercentage[i+1]
	// current <= p < next
	// 増加分に比例して値を計算
	return current[1] + (next[1]-current[1])*(p-current[0])/(next[0]-current[0])
}

func (a *Animation) step() int {
	p := a.nowPercentage()
	for i := range len(a.ByPercentage) - 1 {
		next := a.ByPercentage[i+1]
		if next[0] <= p {
			continue
		}
		return i
	}
	return len(a.ByPercentage) - 1
}

// 現在のパーセンテージを取得する
func (a *Animation) nowPercentage() float32 {
	if a.Duration == 0 {
		return 0
	}
	sub := time.Since(a.itime)
	return float32(sub) / float32(a.Duration)
}

func (a *Animation) String() string {
	sub := time.Since(a.itime)
	return fmt.Sprintf("Animation{ByPercentage: %v, Duration: %v, sub: %v, Value: %f}", a.ByPercentage, a.Duration, sub, float32(sub)/float32(a.Duration))
}
