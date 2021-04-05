package game_test

import (
	"bytes"
	"sort"
	"testing"

	"github.com/powerman/check"

	"github.com/Djarvur/allcups-itrally-2020-task/internal/app/game"
)

func licenses(g game.Game) []game.License {
	ls := g.Licenses()
	sort.Slice(ls, func(i, j int) bool { return ls[i].ID < ls[j].ID })
	return ls
}

func TestSaveLoad(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()

	g, err := game.Factory{}.New(ctx, cfg3x2x2s666)
	t.Nil(err)

	var buf bytes.Buffer
	n, err := g.WriteTo(&buf)
	t.Nil(err)
	t.Equal(n, int64(166))
	_, err = game.Factory{}.Continue(ctx, bytes.NewReader(buf.Bytes()))
	t.Nil(err)

	lic1, _ := g.IssueLicense(nil)
	g.Dig(lic1.ID, game.Coord{X: 0, Y: 0, Depth: 1})
	g.Dig(lic1.ID, game.Coord{X: 0, Y: 0, Depth: 2})
	g.Cash(game.Coord{X: 0, Y: 0, Depth: 2})
	lic2, _ := g.IssueLicense([]int{1})
	g.Dig(lic2.ID, game.Coord{X: 2, Y: 0, Depth: 1})

	balance, wallet := g.Balance()
	t.Equal(balance, 2)
	t.DeepEqual(wallet, []int{2, 0})
	t.DeepEqual(licenses(g), []game.License{
		{ID: 0, DigAllowed: 3, DigUsed: 2},
		{ID: 1, DigAllowed: 5, DigUsed: 1},
	})

	buf.Reset()
	n, err = g.WriteTo(&buf)
	t.Nil(err)
	t.Equal(n, int64(237))
	g, err = game.Factory{}.Continue(ctx, bytes.NewReader(buf.Bytes()))
	t.Nil(err)

	balance, wallet = g.Balance()
	t.Equal(balance, 2)
	t.DeepEqual(wallet, []int{1, 0})
	t.DeepEqual(licenses(g), []game.License{
		{ID: 0, DigAllowed: 3, DigUsed: 2},
		{ID: 1, DigAllowed: 5, DigUsed: 1},
	})

	_, err = g.Cash(game.Coord{X: 0, Y: 0, Depth: 2}) // Already digged and cashed.
	t.Err(err, game.ErrNoThreasure)
	_, err = g.Cash(game.Coord{X: 2, Y: 0, Depth: 1}) // Already digged but not cashed.
	t.Err(err, game.ErrNoThreasure)

	_, err = g.Dig(lic1.ID, game.Coord{X: 2, Y: 0, Depth: 1}) // Already digged.
	t.Err(err, game.ErrWrongDepth)
	_, err = g.Dig(lic2.ID, game.Coord{X: 2, Y: 0, Depth: 2}) // Not digged.
	t.Nil(err)
	_, err = g.Dig(lic2.ID, game.Coord{X: 0, Y: 1, Depth: 1})
	t.Nil(err)
	_, err = g.Dig(lic2.ID, game.Coord{X: 0, Y: 1, Depth: 2})
	t.Nil(err)
	wallet, err = g.Cash(game.Coord{X: 0, Y: 1, Depth: 2}) // Not digged and not cashed.
	t.Nil(err)
	t.DeepEqual(wallet, []int{2, 3, 4})
}
