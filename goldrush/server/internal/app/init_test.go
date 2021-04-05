package app_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/smartystreets/goconvey/convey"

	"github.com/Djarvur/allcups-itrally-2020-task/internal/app"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/app/game"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/def"
)

func TestMain(m *testing.M) {
	def.Init()
	reg := prometheus.NewPedanticRegistry()
	app.InitMetrics(reg)
	check.TestMain(m)
}

type Ctx = context.Context

// Const shared by tests. Recommended naming scheme: <dataType><Variant>.
var (
	ctx = context.Background()
	cfg = app.Config{
		AutosavePeriod:    def.TestSecond,
		DepthProfitChange: 0.1,
		DigBaseDelay:      def.TestSecond / 1000,
		DigExtraDelay:     def.TestSecond / 10000,
		Duration:          60 * def.TestSecond,
		Game:              app.Difficulty["test"],
		StartTimeout:      3 * def.TestSecond,
	}
)

func testPrepare(t *check.C) (func(), *app.MockRepo, *app.MockCPU, *app.MockLicenseSvc, *game.MockGame, *app.MockGameFactory, func(a *app.App, err error)) {
	ctrl := gomock.NewController(t)
	mockRepo := app.NewMockRepo(ctrl)
	mockCPU := app.NewMockCPU(ctrl)
	mockLicenseSvc := app.NewMockLicenseSvc(ctrl)
	mockGame := game.NewMockGame(ctrl)
	mockGameFactory := app.NewMockGameFactory(ctrl)
	wantErr := func(a *app.App, err error) {
		t.Helper()
		t.Err(err, io.EOF)
		t.Nil(a)
	}
	mockCPU.EXPECT().Consume(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()     // TODO Test later.
	mockLicenseSvc.EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil).AnyTimes() // TODO Test later.
	return ctrl.Finish, mockRepo, mockCPU, mockLicenseSvc, mockGame, mockGameFactory, wantErr
}

func testNew(t *check.C) (func(), *app.App, *app.MockRepo, *game.MockGame) {
	t.Helper()
	cleanup, mockRepo, mockCPU, mockLicenseSvc, mockGame, mockGameFactory, _ := testPrepare(t)

	mockRepo.EXPECT().LoadStartTime().Return(&time.Time{}, nil)
	mockRepo.EXPECT().SaveTreasureKey(gomock.Any()).Return(nil)
	mockRepo.EXPECT().SaveGame(mockGame).Return(nil)
	mockGameFactory.EXPECT().New(gomock.Any(), cfg.Game).Return(mockGame, nil)

	a, err := app.New(ctx, mockRepo, mockCPU, mockLicenseSvc, mockGameFactory, cfg)
	t.Must(t.Nil(err))
	return cleanup, a, mockRepo, mockGame
}

func waitErr(t *check.C, errc <-chan error, wait time.Duration, wantErr error) {
	t.Helper()
	now := time.Now()
	select {
	case err := <-errc:
		t.Between(time.Since(now), wait-wait/4, wait+wait/4)
		t.Err(err, wantErr)
	case <-time.After(def.TestTimeout):
		t.FailNow()
	}
}
