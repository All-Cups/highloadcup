//nolint:testpackage // Testing internals.
package game

import (
	"testing"

	"github.com/powerman/check"
)

func TestTreasureValueAt(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()

	cfg := Config{TreasureValue: func(v []int) *[]int { return &v }([]int{10, 100})}
	tests := []struct {
		alg       TreasureValueAlg
		depth     uint8
		want      []int
		wantPanic string
	}{
		{AlgDoubleMax, 1, []int{10, 20}, ``},
		{AlgDoubleMax, 2, []int{100, 200}, ``},
		{AlgQuarterAround, 1, []int{8, 12}, ``},
		{AlgQuarterAround, 2, []int{75, 125}, ``},
		{0, 1, nil, `unknown TreasureValueAlg`},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			cfg.TreasureValueAlg = tc.alg
			if tc.wantPanic != `` {
				t.PanicMatch(func() { cfg.treasureValueAt(tc.depth) }, tc.wantPanic)
				return
			}
			min, max := cfg.treasureValueAt(tc.depth)
			t.DeepEqual([]int{min, max}, tc.want)
		})
	}
}
