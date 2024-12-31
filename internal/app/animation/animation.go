package animation

import "time"

type Animation struct {
	TPS float32 // 1 秒あたりの描画回数
	// レンダリング設定。[0]がパーセンテージ、[1]が値。0と1パーセンテージは必須。
	// パーセンテージの昇順で設定すること。
	ByPercentage [][2]float32
	Duration     time.Duration // アニメーション時間
	count        float32       // カウント
}

func (a *Animation) Clear() {
	a.count = 0
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

// 描画
func (a *Animation) Update() {
	a.count++
}

// 現在のパーセンテージを取得する
func (a *Animation) nowPercentage() float32 {
	// 1. 経過時間を計算
	// 1countあたりの経過時間
	t := float32(time.Duration(time.Second*1)) / a.TPS
	// 現在の経過時間
	epi := t * a.count
	// 2. 経過時間に対する全体の時間に対する割合を計算
	return epi / float32(a.Duration)
}
