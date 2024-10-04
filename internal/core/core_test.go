package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCore(t *testing.T) {
	t.Run("Direction", func(t *testing.T) {
		assert.Equal(t, East, North.Rotate(Right))
		assert.Equal(t, South, North.Rotate(Back))
		assert.Equal(t, West, North.Rotate(Left))
		assert.Equal(t, North, North.Rotate(Front))

		assert.Equal(t, East, West.Rotate(Back))
		assert.Equal(t, South, West.Rotate(Left))
		assert.Equal(t, West, West.Rotate(Front))
		assert.Equal(t, North, West.Rotate(Right))
	})

	t.Run("Sector", func(t *testing.T) {
		src := RelativeSectors{{4, 3}, {-1, 2}}
		test := []struct {
			d   Direction
			exp RelativeSectors
		}{
			{
				South,
				RelativeSectors{{-4, -3}, {1, -2}}},
			{
				West,
				RelativeSectors{{3, -4}, {2, 1}}},
			{
				East,
				RelativeSectors{{-3, 4}, {-2, -1}}},
			{
				North,
				RelativeSectors{{4, 3}, {-1, 2}}},
		}

		for _, tt := range test {
			assert.Equal(t, tt.exp, src.Trun(tt.d))
		}
	})

	t.Run("OceanMap", func(t *testing.T) {
		src := OceanMap{
			H: 7,
			W: 6,
		}
		rss := RelativeSectors{{1, -1}, {-1, 1}, {2, 3}, {-4, -5}}
		test := []struct {
			at  Sector
			exp []Sector
		}{
			{0, []Sector{20}},
			{20, []Sector{15, 25, 40}},
			{41, []Sector{7}},
			{5, []Sector{10}},
		}
		for _, tt := range test {
			act := src.Sectors(tt.at, rss)
			assert.Equal(t, tt.exp, act)
		}

		assert.Equal(t, [][]Sector{
			{00, 01, 02, 03, 04, 05},
			{06, 07, 8, 9, 10, 11},
			{12, 13, 14, 15, 16, 17},
			{18, 19, 20, 21, 22, 23},
			{24, 25, 26, 27, 28, 29},
			{30, 31, 32, 33, 34, 35},
			{36, 37, 38, 39, 40, 41},
		}, src.Lines())
	})
}
