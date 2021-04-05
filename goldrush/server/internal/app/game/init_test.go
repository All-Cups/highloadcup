package game_test

import (
	"context"
	"errors"
	"testing"

	"github.com/powerman/check"

	"github.com/Djarvur/allcups-itrally-2020-task/internal/app/game"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/def"
)

type C = game.Config

func TestMain(m *testing.M) {
	def.Init()
	check.TestMain(m)
}

var (
	ctx            = context.Background()
	errOutOfBounds = errors.New("out of bounds")
)

func addTreasureValues(cfg C) C {
	treasureValues := [11]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	treasureValue := [][]int{
		treasureValues[:0],
		treasureValues[:1],
		treasureValues[:2],
		5:  treasureValues[:5],
		10: treasureValues[:10],
		11: treasureValues[:11],
	}
	if cfg.TreasureValueAlg == 0 {
		cfg.TreasureValue = &treasureValue[cfg.Depth]
		cfg.TreasureValueAlg = game.AlgDoubleMax
	}
	return cfg
}

var (
	cfg3x2x2s666 = addTreasureValues(C{
		Seed:              666, // {2 0 1}), {0 0 2}, {0 1 2} and 1 duplicate
		MaxActiveLicenses: 2,
		Density:           3,
		SizeX:             3,
		SizeY:             2,
		Depth:             2,
	})
	cfg2x2x1 = addTreasureValues(C{ // {0 0 1}), {0 1 1}, {1 1 1}
		MaxActiveLicenses: 2,
		Density:           1,
		SizeX:             2,
		SizeY:             2,
		Depth:             1,
	})
	cfg40x40x1 = addTreasureValues(C{
		MaxActiveLicenses: 40*40 + 1,
		Density:           1,
		SizeX:             40,
		SizeY:             40,
		Depth:             1,
	})
	cfg1x1x1s1 = addTreasureValues(C{
		Seed:              1,
		MaxActiveLicenses: 2,
		Density:           1,
		SizeX:             1,
		SizeY:             1,
		Depth:             1,
	})
	cfg1x1x1 = addTreasureValues(C{
		Density: 1,
		SizeX:   1,
		SizeY:   1,
		Depth:   1,
	})
	cfg5x5x5 = addTreasureValues(C{
		Density: 5,
		SizeX:   5,
		SizeY:   5,
		Depth:   5,
	})
	cfg2x2x2 = addTreasureValues(C{ // {0 0 1}), {0 1 1}, {1 0 1}, {1 1 1}, {0 0 2}
		MaxActiveLicenses: 10,
		Density:           1,
		SizeX:             2,
		SizeY:             2,
		Depth:             2,
	})
)
