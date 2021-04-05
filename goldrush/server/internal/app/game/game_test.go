package game_test

import (
	"sort"
	"testing"

	"github.com/powerman/check"

	"github.com/Djarvur/allcups-itrally-2020-task/internal/app/game"
)

func TestSmoke(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()

	g, err := game.Factory{}.New(ctx, cfg3x2x2s666)
	t.Nil(err)

	count, err := g.CountTreasures(game.Area{X: 0, Y: 0, SizeX: 3, SizeY: 2}, 1)
	t.Nil(err)
	t.Equal(count, 1)
	count, err = g.CountTreasures(game.Area{X: 0, Y: 0, SizeX: 3, SizeY: 2}, 2)
	t.Nil(err)
	t.Equal(count, 2)

	count, err = g.CountTreasures(game.Area{X: 2, Y: 0, SizeX: 1, SizeY: 1}, 1)
	t.Nil(err)
	t.Equal(count, 1)
	count, err = g.CountTreasures(game.Area{X: 0, Y: 0, SizeX: 1, SizeY: 2}, 2)
	t.Nil(err)
	t.Equal(count, 2)

	count, err = g.CountTreasures(game.Area{X: 0, Y: 0, SizeX: 2, SizeY: 2}, 1)
	t.Nil(err)
	t.Equal(count, 0)
	count, err = g.CountTreasures(game.Area{X: 0, Y: 1, SizeX: 2, SizeY: 2}, 1)
	t.Err(err, game.ErrWrongCoord)
	t.Equal(count, 0)

	lic1, err := g.IssueLicense(nil)
	t.Nil(err)
	t.DeepEqual(lic1, game.License{ID: 0, DigAllowed: 3})
	lic2, err := g.IssueLicense(nil)
	t.Nil(err)
	t.DeepEqual(lic2, game.License{ID: 1, DigAllowed: 3})
	lic3, err := g.IssueLicense(nil)
	t.Err(err, game.ErrActiveLicenseLimit)
	t.Zero(lic3)

	found, err := g.Dig(lic1.ID, game.Coord{X: 2, Y: 2, Depth: 1})
	t.Err(err, game.ErrWrongCoord)
	t.False(found)

	found, err = g.Dig(lic1.ID, game.Coord{X: 0, Y: 0, Depth: 1})
	t.Nil(err)
	t.False(found)
	found, err = g.Dig(lic1.ID, game.Coord{X: 0, Y: 0, Depth: 2})
	t.Nil(err)
	t.True(found)
	wallet, err := g.Cash(game.Coord{X: 0, Y: 0, Depth: 2})
	t.Nil(err)
	t.Len(wallet, 3)
	t.Len(g.Licenses(), 1)

	found, err = g.Dig(lic2.ID, game.Coord{X: 0, Y: 1, Depth: 2})
	t.Err(err, game.ErrWrongDepth)
	t.False(found)
	found, err = g.Dig(lic2.ID, game.Coord{X: 0, Y: 1, Depth: 1})
	t.Nil(err)
	t.False(found)
	found, err = g.Dig(lic2.ID, game.Coord{X: 0, Y: 1, Depth: 2})
	t.Nil(err)
	t.True(found)
	wallet, err = g.Cash(game.Coord{X: 0, Y: 1, Depth: 2})
	t.Nil(err)
	t.Len(wallet, 4)
	t.Len(g.Licenses(), 0)
	wallet, err = g.Cash(game.Coord{X: 0, Y: 1, Depth: 2})
	t.Err(err, game.ErrNoThreasure)
	t.Nil(wallet)

	found, err = g.Dig(lic2.ID, game.Coord{X: 2, Y: 0, Depth: 1})
	t.Err(err, game.ErrNoSuchLicense)
	t.False(found)
	wallet, err = g.Cash(game.Coord{X: 2, Y: 0, Depth: 1})
	t.Err(err, game.ErrNotDigged)
	t.Nil(wallet)

	lic3, err = g.IssueLicense([]int{0, 2})
	t.Nil(err)
	t.DeepEqual(lic3, game.License{ID: 2, DigAllowed: 5})
	balance, wallet := g.Balance()
	t.Equal(balance, 5)
	t.DeepEqual(wallet, []int{6, 5, 4, 3, 1})

	lic3, err = g.IssueLicense([]int{1, 3, 4, 5, 6, 6})
	t.Err(err, game.ErrBogusCoin)
	t.Zero(lic3)
	lic3, err = g.IssueLicense([]int{1, 2, 3})
	t.Err(err, game.ErrBogusCoin)
	t.Zero(lic3)
	lic3, err = g.IssueLicense([]int{1, 3, 3})
	t.Err(err, game.ErrBogusCoin)
	t.Zero(lic3)
}

func TestNew(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()

	tests := []struct {
		cfg  C
		want int
	}{
		{C{Density: 1, SizeX: 1, SizeY: 1, Depth: 1}, 1},
		{C{Density: 3500, SizeX: 3500, SizeY: 3500, Depth: 10}, 3506},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := game.Factory{}.New(ctx, addTreasureValues(tc.cfg))
			t.Nil(err)
			count, err := res.CountTreasures(game.Area{
				SizeX: tc.cfg.SizeX,
				SizeY: tc.cfg.SizeY,
			}, 1)
			t.Nil(err)
			t.Equal(count, tc.want)
		})
	}
}

func TestNewErr(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()

	ints := func(v []int) *[]int { return &v }

	tests := []struct {
		cfg       C
		wantErr   error
		wantPanic string
	}{
		{C{Density: 1, SizeX: 1, SizeY: 1, Depth: 1}, nil, ``},
		{C{Density: 6000 * 6000 * 10, SizeX: 6000, SizeY: 6000, Depth: 10}, nil, ``},
		{C{Density: 0, SizeX: 1, SizeY: 1, Depth: 1}, errOutOfBounds, ``},
		{C{Density: 1, SizeX: 0, SizeY: 1, Depth: 1}, errOutOfBounds, ``},
		{C{Density: 1, SizeX: 1, SizeY: 0, Depth: 1}, errOutOfBounds, ``},
		{C{Density: 1, SizeX: 1, SizeY: 1, Depth: 0}, errOutOfBounds, ``},
		{C{Density: 6000*6000*10 + 1, SizeX: 6000, SizeY: 6000, Depth: 10}, errOutOfBounds, ``},
		{C{Density: 6000 * 6000 * 10, SizeX: 6001, SizeY: 6000, Depth: 10}, errOutOfBounds, ``},
		{C{Density: 6000 * 6000 * 10, SizeX: 6000, SizeY: 6001, Depth: 10}, errOutOfBounds, ``},
		{C{Density: 6000 * 6000 * 10, SizeX: 6000, SizeY: 6000, Depth: 11}, errOutOfBounds, ``},
		{C{Density: 1, SizeX: 1, SizeY: 1, Depth: 1, TreasureValueAlg: game.AlgDoubleMax}, nil, `TreasureValue length`},
		{C{Density: 1, SizeX: 1, SizeY: 1, Depth: 1, TreasureValueAlg: game.AlgDoubleMax, TreasureValue: ints([]int{})}, nil, `TreasureValue length`},
		{C{Density: 1, SizeX: 1, SizeY: 1, Depth: 1, TreasureValueAlg: game.AlgDoubleMax, TreasureValue: ints([]int{0, 1})}, nil, `TreasureValue length`},
		{C{Density: 1, SizeX: 1, SizeY: 1, Depth: 1, TreasureValueAlg: -1, TreasureValue: ints([]int{1})}, nil, `unknown TreasureValueAlg:`},
		{C{}, errOutOfBounds, ``},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			if tc.wantPanic != `` {
				t.PanicMatch(func() {
					game.Factory{}.New(ctx, addTreasureValues(tc.cfg))
				}, tc.wantPanic)
				return
			}
			res, err := game.Factory{}.New(ctx, addTreasureValues(tc.cfg))
			if tc.wantErr == nil {
				t.Nil(err)
				t.NotNil(res)
			} else {
				t.Err(err, tc.wantErr)
				t.Nil(res)
			}
		})
	}
}

func TestBalance(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()

	g, err := game.Factory{}.New(ctx, cfg2x2x1)
	t.Nil(err)

	balance, wallet := g.Balance()
	t.Equal(balance, 0)
	t.Len(wallet, 0)

	g.IssueLicense(nil)
	g.Dig(0, game.Coord{X: 0, Y: 0, Depth: 1})
	g.Cash(game.Coord{X: 0, Y: 0, Depth: 1})
	g.IssueLicense([]int{1})
	g.Dig(1, game.Coord{X: 0, Y: 1, Depth: 1})
	g.Cash(game.Coord{X: 0, Y: 1, Depth: 1})
	balance, wallet = g.Balance()
	t.Equal(balance, 2)
	t.DeepEqual(wallet, []int{2, 0})

	g, err = game.Factory{}.New(ctx, cfg40x40x1)
	t.Nil(err)
	for x := 0; x < 40; x++ {
		for y := 0; y < 40; y++ {
			g.IssueLicense(nil)
			g.Dig(x*40+y, game.Coord{X: x, Y: y, Depth: 1})
			g.Cash(game.Coord{X: x, Y: y, Depth: 1})
		}
	}
	balance, wallet = g.Balance()
	t.Equal(balance, 1538)
	t.Len(wallet, 1000)

	l, err := g.IssueLicense(wallet)
	t.Nil(err)
	t.Equal(l, game.License{ID: 40 * 40, DigAllowed: 43})
}

func TestLicenses(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()

	g, err := game.Factory{}.New(ctx, cfg1x1x1s1)
	t.Nil(err)

	t.Len(g.Licenses(), 0)
	l0, err := g.IssueLicense(nil)
	t.Nil(err)
	t.DeepEqual(l0, game.License{ID: 0, DigAllowed: 3, DigUsed: 0})
	t.DeepEqual(g.Licenses(), []game.License{l0})

	g.Dig(l0.ID, game.Coord{X: 0, Y: 0, Depth: 1})
	g.Cash(game.Coord{X: 0, Y: 0, Depth: 1})
	l0.DigUsed++
	t.DeepEqual(g.Licenses(), []game.License{l0})

	l1, err := g.IssueLicense([]int{0, 2})
	t.Err(err, game.ErrBogusCoin)
	t.Zero(l1)
	l2, err := g.IssueLicense([]int{0})
	t.Nil(err)
	t.DeepEqual(l2, game.License{ID: 2, DigAllowed: 5, DigUsed: 0})

	l, err := g.IssueLicense(nil)
	t.Err(err, game.ErrActiveLicenseLimit)
	t.Zero(l)

	ls := g.Licenses()
	sort.Slice(ls, func(i, j int) bool { return ls[i].ID < ls[j].ID })
	t.DeepEqual(ls, []game.License{l0, l2})

	g.Dig(l0.ID, game.Coord{X: 0, Y: 0, Depth: 1})
	g.Dig(l0.ID, game.Coord{X: 0, Y: 0, Depth: 1})
	t.DeepEqual(g.Licenses(), []game.License{l2})

	l, err = g.IssueLicense([]int{0})
	t.Err(err, game.ErrBogusCoin)
	t.DeepEqual(g.Licenses(), []game.License{l2})
	t.Zero(l)
}

func TestCountTreasures(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()

	g1, err := game.Factory{}.New(ctx, cfg1x1x1)
	t.Nil(err)
	g5, err := game.Factory{}.New(ctx, cfg5x5x5)
	t.Nil(err)

	tests := []struct {
		g       game.Game
		area    game.Area
		depth   uint8
		want    int
		wantErr error
	}{
		{g1, game.Area{X: 0, Y: 0, SizeX: 1, SizeY: 1}, 1, 1, nil},
		{g1, game.Area{X: 1, Y: 0, SizeX: 1, SizeY: 1}, 1, 0, game.ErrWrongCoord},
		{g1, game.Area{X: 0, Y: 1, SizeX: 1, SizeY: 1}, 1, 0, game.ErrWrongCoord},
		{g1, game.Area{X: -1, Y: 0, SizeX: 1, SizeY: 1}, 1, 0, game.ErrWrongCoord},
		{g1, game.Area{X: 0, Y: -1, SizeX: 1, SizeY: 1}, 1, 0, game.ErrWrongCoord},
		{g1, game.Area{X: 0, Y: 0, SizeX: 0, SizeY: 1}, 1, 0, game.ErrWrongCoord},
		{g1, game.Area{X: 0, Y: 0, SizeX: 2, SizeY: 1}, 1, 0, game.ErrWrongCoord},
		{g1, game.Area{X: 0, Y: 0, SizeX: 1, SizeY: 0}, 1, 0, game.ErrWrongCoord},
		{g1, game.Area{X: 0, Y: 0, SizeX: 1, SizeY: 2}, 1, 0, game.ErrWrongCoord},
		{g1, game.Area{X: 0, Y: 0, SizeX: 1, SizeY: 1}, 0, 0, game.ErrWrongDepth},
		{g1, game.Area{X: 0, Y: 0, SizeX: 1, SizeY: 1}, 2, 0, game.ErrWrongDepth},
		{g5, game.Area{X: 0, Y: 0, SizeX: 5, SizeY: 5}, 3, 7, nil},
		{g5, game.Area{X: 1, Y: 2, SizeX: 4, SizeY: 3}, 5, 2, nil},
		{g5, game.Area{X: 1, Y: 0, SizeX: 5, SizeY: 5}, 1, 0, game.ErrWrongCoord},
		{g5, game.Area{X: 0, Y: 1, SizeX: 5, SizeY: 5}, 1, 0, game.ErrWrongCoord},
		{g5, game.Area{X: 5, Y: 0, SizeX: 1, SizeY: 1}, 1, 0, game.ErrWrongCoord},
		{g5, game.Area{X: 0, Y: 0, SizeX: 1, SizeY: 1}, 0, 0, game.ErrWrongDepth},
		{g5, game.Area{X: 0, Y: 0, SizeX: 1, SizeY: 1}, 6, 0, game.ErrWrongDepth},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := tc.g.CountTreasures(tc.area, tc.depth)
			t.Err(err, tc.wantErr)
			t.Equal(res, tc.want)
		})
	}
}

func TestDig(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()

	g, err := game.Factory{}.New(ctx, cfg2x2x2)
	t.Nil(err)

	for i := 0; i < 10; i++ {
		g.IssueLicense(nil)
	}

	tests := []struct {
		licenseID int
		coord     game.Coord
		want      bool
		wantErr   error
	}{
		{10, game.Coord{X: 0, Y: 0, Depth: 1}, false, game.ErrNoSuchLicense},
		{-1, game.Coord{X: 0, Y: 0, Depth: 1}, false, game.ErrNoSuchLicense},
		{0, game.Coord{X: 0, Y: 0, Depth: 1}, true, nil},
		{0, game.Coord{X: 0, Y: 0, Depth: 2}, true, nil},
		{0, game.Coord{X: 1, Y: 1, Depth: 1}, true, nil},
		{0, game.Coord{X: 1, Y: 1, Depth: 2}, false, game.ErrNoSuchLicense},
		{1, game.Coord{X: 1, Y: 1, Depth: 2}, false, nil},
		{1, game.Coord{X: 1, Y: 1, Depth: 1}, false, game.ErrWrongDepth},
		{1, game.Coord{X: 1, Y: 1, Depth: 2}, false, game.ErrWrongDepth},
		{1, game.Coord{X: 1, Y: 1, Depth: 2}, false, game.ErrNoSuchLicense},
		{2, game.Coord{X: -1, Y: 1, Depth: 1}, false, game.ErrWrongCoord},
		{2, game.Coord{X: 1, Y: -1, Depth: 1}, false, game.ErrWrongCoord},
		{3, game.Coord{X: 2, Y: 1, Depth: 1}, false, game.ErrWrongCoord},
		{3, game.Coord{X: 1, Y: 2, Depth: 1}, false, game.ErrWrongCoord},
		{4, game.Coord{X: 1, Y: 1, Depth: 0}, false, game.ErrWrongCoord},
		{4, game.Coord{X: 1, Y: 1, Depth: 3}, false, game.ErrWrongCoord},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := g.Dig(tc.licenseID, tc.coord)
			t.Err(err, tc.wantErr)
			t.Equal(res, tc.want)
		})
	}
}

func TestCash(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()

	g, err := game.Factory{}.New(ctx, cfg2x2x2)
	t.Nil(err)

	g.IssueLicense(nil)
	g.Dig(0, game.Coord{X: 0, Y: 0, Depth: 1})
	g.Dig(0, game.Coord{X: 0, Y: 0, Depth: 2})

	tests := []struct {
		coord   game.Coord
		want    []int
		wantErr error
	}{
		{game.Coord{X: 0, Y: 0, Depth: 2}, []int{0, 1, 2, 3}, nil},
		{game.Coord{X: 0, Y: 0, Depth: 1}, []int{4}, nil},
		{game.Coord{X: 0, Y: 0, Depth: 1}, nil, game.ErrNoThreasure},
		{game.Coord{X: 1, Y: 1, Depth: 1}, nil, game.ErrNotDigged},
		{game.Coord{X: 1, Y: 1, Depth: 2}, nil, game.ErrNoThreasure},
		{game.Coord{X: 2, Y: 1, Depth: 2}, nil, game.ErrWrongCoord},
		{game.Coord{X: 1, Y: 2, Depth: 2}, nil, game.ErrWrongCoord},
		{game.Coord{X: 1, Y: 1, Depth: 3}, nil, game.ErrWrongCoord},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt) //nolint:govet // Shadow t? Sure! Why not…
			res, err := g.Cash(tc.coord)
			t.Err(err, tc.wantErr)
			t.DeepEqual(res, tc.want)
		})
	}

	balance, wallet := g.Balance()
	t.Equal(balance, 5)
	t.DeepEqual(wallet, []int{4, 3, 2, 1, 0})
}
