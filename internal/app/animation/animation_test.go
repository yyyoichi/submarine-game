package animation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAnimation(t *testing.T) {
	type exp struct {
		times int
		value float32
	}
	test := []struct {
		a                Animation
		expOpacityByCall []exp
	}{
		{Animation{
			TPS:          60,
			Duration:     time.Duration(time.Second * 1),
			ByPercentage: [][2]float32{{0, 0}, {0.5, 0.5}, {1, 1}},
		},
			[]exp{{0, 0}, {30, 0.5}, {30, 1}, {1, 1}},
		},
		{Animation{
			TPS:          100, // 1回のUpdateで10ms
			Duration:     time.Duration(float64(time.Millisecond * 1200)),
			ByPercentage: [][2]float32{{0, 0}, {0.4, 0.5}, {0.5, 0.8}, {1, 1}},
		},
			// 480msで0.5, 600ms(+120ms)で0.8, 900ms(+300ms)で0.9, 1200ms(+300ms)で1
			[]exp{{48, 0.5}, {12, 0.8}, {30, 0.9}, {30, 1}, {1, 1}},
		},
		{Animation{
			TPS:          100, // 1回のUpdateで10ms
			Duration:     time.Duration(float64(time.Millisecond * 1200)),
			ByPercentage: [][2]float32{{0, 0}, {0.4, 1}, {0.6, 1}, {1, 0}},
		},
			// 240msで0.5, 480ms(+240ms)で1, 720ms(+240ms)で1, 960ms(+240ms)で0.5, 1200ms(+240ms)で0
			[]exp{{24, 0.5}, {24, 1}, {24, 1}, {24, 0.5}, {24, 0}, {1, 0}},
		},
	}
	for i, tt := range test {
		for _, exp := range tt.expOpacityByCall {
			for range exp.times {
				tt.a.Update()
			}
			val := tt.a.Value()
			assert.Equalf(t, exp.value, val, "test[%d]: exp is '%v' but got='%v'", i, exp.value, val)
		}
	}
}
