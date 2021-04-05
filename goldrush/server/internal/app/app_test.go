package app_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/powerman/check"

	"github.com/Djarvur/allcups-itrally-2020-task/internal/app"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/app/game"
)

var (
	start = time.Now()
	dump  = nopCloser{bytes.NewReader([]byte("save"))}
)

func TestNew(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	cleanup, mockRepo, mockCPU, mockLicenseSvc, mockGame, mockGameFactory, wantErr := testPrepare(t)
	defer cleanup()

	mockRepo.EXPECT().LoadStartTime().Return(nil, io.EOF)
	mockRepo.EXPECT().LoadStartTime().Return(&time.Time{}, nil).Times(7)
	wantErr(app.New(ctx, mockRepo, mockCPU, mockLicenseSvc, mockGameFactory, cfg))

	// Enforce random seed if difficulty is not "test" and seed is 0.
	cfgNormal := app.Difficulty["normal"]
	cfgNormal.TreasureValue = app.Difficulty["test"].TreasureValue
	cfgNormal.TreasureValueAlg = app.Difficulty["test"].TreasureValueAlg
	cfgNormal7 := cfgNormal
	cfgNormal7.Seed = 7
	cfgNormalNoAlg := app.Difficulty["normal"]
	cfgNormalNoAlg.Seed = 7
	cfgNormalDefaultAlg := cfgNormalNoAlg
	cfgNormalDefaultAlg.TreasureValue = func(v []int) *[]int { return &v }([]int{4, 9, 15, 23, 27, 30, 33, 35, 37, 38})
	cfgNormalDefaultAlg.TreasureValueAlg = game.AlgQuarterAround
	cfgTest := app.Difficulty["test"]
	mockGameFactory.EXPECT().New(gomock.Any(), matchRandomSeed(cfgNormal)).Return(nil, io.EOF)
	mockGameFactory.EXPECT().New(gomock.Any(), cfgNormal7).Return(nil, io.EOF)
	mockGameFactory.EXPECT().New(gomock.Any(), cfgNormalDefaultAlg).Return(nil, io.EOF)
	mockGameFactory.EXPECT().New(gomock.Any(), cfgTest).Return(nil, io.EOF)
	mockGameFactory.EXPECT().New(gomock.Any(), cfgTest).Return(mockGame, nil).Times(3)
	cfg := cfg
	cfg.Game = cfgNormal
	wantErr(app.New(ctx, mockRepo, mockCPU, mockLicenseSvc, mockGameFactory, cfg))
	cfg.Game = cfgNormal7
	wantErr(app.New(ctx, mockRepo, mockCPU, mockLicenseSvc, mockGameFactory, cfg))
	cfg.Game = cfgNormalNoAlg
	wantErr(app.New(ctx, mockRepo, mockCPU, mockLicenseSvc, mockGameFactory, cfg))
	cfg.Game = cfgTest
	wantErr(app.New(ctx, mockRepo, mockCPU, mockLicenseSvc, mockGameFactory, cfg))

	mockRepo.EXPECT().SaveTreasureKey(gomock.Len(32)).Return(io.EOF)
	mockRepo.EXPECT().SaveTreasureKey(gomock.Len(32)).Return(nil).Times(2)
	wantErr(app.New(ctx, mockRepo, mockCPU, mockLicenseSvc, mockGameFactory, cfg))

	mockRepo.EXPECT().SaveGame(mockGame).Return(io.EOF)
	mockRepo.EXPECT().SaveGame(mockGame).Return(nil).MinTimes(1)
	wantErr(app.New(ctx, mockRepo, mockCPU, mockLicenseSvc, mockGameFactory, cfg))

	a, err := app.New(ctx, mockRepo, mockCPU, mockLicenseSvc, mockGameFactory, cfg)
	t.Nil(err)
	t.NotNil(a)
}

func TestContinue(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	cleanup, mockRepo, mockCPU, mockLicenseSvc, mockGame, mockGameFactory, wantErr := testPrepare(t)
	defer cleanup()

	mockRepo.EXPECT().LoadStartTime().Return(&start, nil).AnyTimes()

	mockRepo.EXPECT().LoadTreasureKey().Return(nil, io.EOF)
	mockRepo.EXPECT().LoadTreasureKey().Return(make([]byte, 33), nil)
	mockRepo.EXPECT().LoadTreasureKey().Return(make([]byte, 32), nil).Times(4)

	wantErr(app.New(ctx, mockRepo, mockCPU, mockLicenseSvc, mockGameFactory, cfg))

	a, err := app.New(ctx, mockRepo, mockCPU, mockLicenseSvc, mockGameFactory, cfg)
	t.Match(err, `bad PASETO key size`)
	t.Nil(a)

	mockRepo.EXPECT().LoadGame().Return(nil, io.EOF)
	mockRepo.EXPECT().LoadGame().Return(dump, nil).Times(3)
	wantErr(app.New(ctx, mockRepo, mockCPU, mockLicenseSvc, mockGameFactory, cfg))

	mockGameFactory.EXPECT().Continue(gomock.Any(), dump).Return(nil, io.EOF)
	mockGameFactory.EXPECT().Continue(gomock.Any(), dump).Return(mockGame, nil).Times(2)
	wantErr(app.New(ctx, mockRepo, mockCPU, mockLicenseSvc, mockGameFactory, cfg))

	mockRepo.EXPECT().SaveStartTime(start).Return(io.EOF)
	wantErr(app.New(ctx, mockRepo, mockCPU, mockLicenseSvc, mockGameFactory, cfg))

	mockRepo.EXPECT().SaveStartTime(start).Return(nil)
	a, err = app.New(ctx, mockRepo, mockCPU, mockLicenseSvc, mockGameFactory, cfg)
	t.Nil(err)
	t.NotNil(a)
}

func TestRestoreKey(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	cleanup, mockRepo, mockCPU, mockLicenseSvc, mockGame, mockGameFactory, _ := testPrepare(t)
	defer cleanup()
	var key []byte

	mockRepo.EXPECT().LoadStartTime().Return(&time.Time{}, nil)
	mockGameFactory.EXPECT().New(gomock.Any(), cfg.Game).Return(mockGame, nil)
	mockRepo.EXPECT().SaveTreasureKey(gomock.Len(32)).DoAndReturn(func(k []byte) error {
		key = k
		return nil
	})
	mockRepo.EXPECT().SaveGame(mockGame).Return(nil).MinTimes(1)
	a, err := app.New(ctx, mockRepo, mockCPU, mockLicenseSvc, mockGameFactory, cfg)
	t.Nil(err)
	t.NotNil(a)

	mockGame.EXPECT().Dig(1, game.Coord{X: 0, Y: 0, Depth: 1}).Return(true, nil)
	treasure, _ := a.Dig(ctx, 1, game.Coord{X: 0, Y: 0, Depth: 1})

	mockRepo.EXPECT().LoadStartTime().Return(&start, nil)
	mockRepo.EXPECT().LoadTreasureKey().Return(key, nil)
	mockRepo.EXPECT().LoadGame().Return(dump, nil)
	mockGameFactory.EXPECT().Continue(gomock.Any(), dump).Return(mockGame, nil)
	mockRepo.EXPECT().SaveStartTime(start).Return(nil)
	a, err = app.New(ctx, mockRepo, mockCPU, mockLicenseSvc, mockGameFactory, cfg)
	t.Nil(err)
	t.NotNil(a)

	mockGame.EXPECT().Cash(game.Coord{X: 0, Y: 0, Depth: 1}).Return([]int{42}, nil)
	res, err := a.Cash(ctx, treasure)
	t.Nil(err)
	t.DeepEqual(res, []int{42})
}

type matchRandomSeed game.Config

func (m matchRandomSeed) String() string { return "has random Seed" }
func (m matchRandomSeed) Matches(x interface{}) bool {
	cfg, ok := x.(game.Config)
	if !ok {
		return false
	}
	if m.Seed == 0 && cfg.Seed > 0 {
		cfg.Seed = 0
		return game.Config(m) == cfg
	}
	return false
}

type nopCloser struct{ io.ReadSeeker }

func (nopCloser) Close() error { return nil }
