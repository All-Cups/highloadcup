package app_test

import (
	"io"
	"testing"

	"github.com/o1egl/paseto/v2"
	"github.com/powerman/check"

	"github.com/Djarvur/allcups-itrally-2020-task/internal/app/game"
)

func TestBalance(tt *testing.T) {
	t := check.T(tt)
	cleanup, a, _, mockGame := testNew(t)
	defer cleanup()

	mockGame.EXPECT().Balance().Return(0, nil)
	balance, wallet, err := a.Balance(ctx)
	t.Nil(err)
	t.Equal(balance, 0)
	t.Len(wallet, 0)
}

func TestExploreArea(tt *testing.T) {
	t := check.T(tt)
	cleanup, a, _, mockGame := testNew(t)
	defer cleanup()

	mockGame.EXPECT().CountTreasures(game.Area{X: 0, Y: 0, SizeX: 5, SizeY: 5}, uint8(1)).Return(5, nil)
	mockGame.EXPECT().CountTreasures(game.Area{X: 0, Y: 0, SizeX: 5, SizeY: 5}, uint8(2)).Return(4, nil)
	mockGame.EXPECT().CountTreasures(game.Area{X: 0, Y: 0, SizeX: 5, SizeY: 5}, uint8(3)).Return(3, nil)
	mockGame.EXPECT().CountTreasures(game.Area{X: 0, Y: 0, SizeX: 5, SizeY: 5}, uint8(4)).Return(2, nil)
	mockGame.EXPECT().CountTreasures(game.Area{X: 0, Y: 0, SizeX: 5, SizeY: 5}, uint8(5)).Return(1, nil)
	mockGame.EXPECT().CountTreasures(game.Area{X: 0, Y: 0, SizeX: 5, SizeY: 5}, uint8(6)).Return(0, nil)
	mockGame.EXPECT().CountTreasures(game.Area{X: 0, Y: 0, SizeX: 5, SizeY: 5}, uint8(7)).Return(0, nil)
	mockGame.EXPECT().CountTreasures(game.Area{X: 0, Y: 0, SizeX: 5, SizeY: 5}, uint8(8)).Return(1, nil)
	mockGame.EXPECT().CountTreasures(game.Area{X: 0, Y: 0, SizeX: 5, SizeY: 5}, uint8(9)).Return(2, nil)
	mockGame.EXPECT().CountTreasures(game.Area{X: 0, Y: 0, SizeX: 5, SizeY: 5}, uint8(10)).Return(3, nil)
	count, err := a.ExploreArea(ctx, game.Area{X: 0, Y: 0, SizeX: 5, SizeY: 5})
	t.Nil(err)
	t.Equal(count, 21)

	mockGame.EXPECT().CountTreasures(game.Area{X: 5, Y: 0, SizeX: 1, SizeY: 1}, uint8(1)).Return(0, io.EOF)
	count, err = a.ExploreArea(ctx, game.Area{X: 5, Y: 0, SizeX: 1, SizeY: 1})
	t.Err(err, io.EOF)
	t.Equal(count, 0)
}

func TestDig(tt *testing.T) {
	t := check.T(tt)
	cleanup, a, _, mockGame := testNew(t)
	defer cleanup()

	mockGame.EXPECT().Dig(1, game.Coord{X: 0, Y: 0, Depth: 0}).Return(false, game.ErrWrongDepth)
	mockGame.EXPECT().Dig(1, game.Coord{X: 0, Y: 0, Depth: 1}).Return(false, nil)
	mockGame.EXPECT().Dig(1, game.Coord{X: 0, Y: 0, Depth: 2}).Return(true, nil)

	tests := []struct {
		licenseID int
		pos       game.Coord
		want      string
		wantErr   error
	}{
		{1, game.Coord{X: 0, Y: 0, Depth: 0}, `^$`, game.ErrWrongDepth},
		{1, game.Coord{X: 0, Y: 0, Depth: 1}, `^$`, nil},
		{1, game.Coord{X: 0, Y: 0, Depth: 2}, `^v2[.]local[.]`, nil},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := a.Dig(ctx, tc.licenseID, tc.pos)
			t.Err(err, tc.wantErr)
			t.Match(res, tc.want)
		})
	}
}

func TestCash(tt *testing.T) {
	t := check.T(tt)
	cleanup, a, _, mockGame := testNew(t)
	defer cleanup()
	cleanup2, a2, _, _ := testNew(t)
	defer cleanup2()

	mockGame.EXPECT().Dig(1, game.Coord{X: 0, Y: 0, Depth: 1}).Return(true, nil)
	treasure, err := a.Dig(ctx, 1, game.Coord{X: 0, Y: 0, Depth: 1})
	t.Nil(err)
	t.HasPrefix(treasure, "v2.local.")

	mockGame.EXPECT().Cash(game.Coord{X: 0, Y: 0, Depth: 1}).Return([]int{0, 1}, nil)

	tests := []struct {
		treasure string
		want     []int
		wantErr  error
	}{
		{"", nil, paseto.ErrIncorrectTokenHeader},
		{"v2.local.AAAA", nil, paseto.ErrIncorrectTokenFormat},
		{treasure, []int{0, 1}, nil},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := a.Cash(ctx, tc.treasure)
			t.Err(err, tc.wantErr)
			t.DeepEqual(res, tc.want)
		})
	}

	res, err := a2.Cash(ctx, treasure)
	t.Err(err, paseto.ErrInvalidTokenAuth)
	t.Nil(res)
}
