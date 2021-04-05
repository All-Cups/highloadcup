package openapi_test

import (
	"io"
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/golang/mock/gomock"
	"github.com/powerman/check"

	"github.com/Djarvur/allcups-itrally-2020-task/api/openapi/client/op"
	"github.com/Djarvur/allcups-itrally-2020-task/api/openapi/model"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/app/game"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/app/resource"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/srv/openapi"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/def"
)

func TestHealthCheck(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	cleanup, c, _, mockApp, _ := testNewServer(t, openapi.Config{})
	defer cleanup()

	mockApp.EXPECT().HealthCheck(gomock.Any()).Return(nil, io.EOF)
	mockApp.EXPECT().HealthCheck(gomock.Any()).Return(nil, nil)
	mockApp.EXPECT().HealthCheck(gomock.Any()).Return("OK", nil)
	mockApp.EXPECT().HealthCheck(gomock.Any()).Return(map[string]string{"main": "OK"}, nil)

	testCases := []struct {
		want    interface{}
		wantErr *model.Error
	}{
		{nil, apiError500},
		{nil, nil},
		{"OK", nil},
		{map[string]interface{}{"main": "OK"}, nil},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := c.Op.HealthCheck(op.NewHealthCheckParams())
			t.DeepEqual(openapi.ErrPayload(err), tc.wantErr)
			if res == nil {
				t.DeepEqual(nil, tc.want)
			} else {
				t.DeepEqual(res.Payload, tc.want)
			}
		})
	}
}

func TestGetBalance(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	cleanup, c, _, mockApp, _ := testNewServer(t, openapi.Config{
		OpGetBalanceRate: 3,
	})
	defer cleanup()

	mockApp.EXPECT().Balance(gomock.Any()).Return(0, nil, io.EOF)
	mockApp.EXPECT().Balance(gomock.Any()).Return(0, nil, nil)
	mockApp.EXPECT().Balance(gomock.Any()).Return(42, []int{1, 2}, nil)

	testCases := []struct {
		want    interface{}
		wantErr *model.Error
	}{
		{nil, apiError500},
		{&model.Balance{Balance: swag.Uint32(0), Wallet: model.Wallet{}}, nil},
		{&model.Balance{Balance: swag.Uint32(42), Wallet: model.Wallet{1, 2}}, nil},
		{nil, apiError429},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := c.Op.GetBalance(op.NewGetBalanceParams())
			t.DeepEqual(openapi.ErrPayload(err), tc.wantErr)
			if res == nil {
				t.DeepEqual(nil, tc.want)
			} else {
				t.DeepEqual(res.Payload, tc.want)
			}
		})
	}
}

func TestListLicenses(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	cleanup, c, _, mockApp, _ := testNewServer(t, openapi.Config{
		OpListLicensesRate: 5,
	})
	defer cleanup()

	mockApp.EXPECT().Licenses(gomock.Any()).Return(nil, io.EOF)
	mockApp.EXPECT().Licenses(gomock.Any()).Return(nil, resource.ErrRPCInternal)
	mockApp.EXPECT().Licenses(gomock.Any()).Return(nil, resource.ErrRPCTimeout)
	mockApp.EXPECT().Licenses(gomock.Any()).Return(nil, nil)
	mockApp.EXPECT().Licenses(gomock.Any()).Return([]game.License{
		{ID: 1, DigAllowed: 3, DigUsed: 0},
		{ID: 2, DigAllowed: 5, DigUsed: 1},
	}, nil)

	testCases := []struct {
		want    interface{}
		wantErr *model.Error
	}{
		{nil, apiError500},
		{nil, apiError502},
		{nil, apiError504},
		{model.LicenseList{}, nil},
		{model.LicenseList{
			&model.License{ID: swag.Int64(1), DigAllowed: 3, DigUsed: 0},
			&model.License{ID: swag.Int64(2), DigAllowed: 5, DigUsed: 1},
		}, nil},
		{nil, apiError429},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := c.Op.ListLicenses(op.NewListLicensesParams())
			t.DeepEqual(openapi.ErrPayload(err), tc.wantErr)
			if res == nil {
				t.DeepEqual(nil, tc.want)
			} else {
				t.DeepEqual(res.Payload, tc.want)
			}
		})
	}
}

func TestIssueLicense(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	cleanup, c, _, mockApp, _ := testNewServer(t, openapi.Config{
		OpIssueLicenseRate: 6,
	})
	defer cleanup()

	mockApp.EXPECT().IssueLicense(gomock.Any(), []int{}).Return(game.License{}, io.EOF)
	mockApp.EXPECT().IssueLicense(gomock.Any(), []int{}).Return(game.License{}, resource.ErrRPCInternal)
	mockApp.EXPECT().IssueLicense(gomock.Any(), []int{}).Return(game.License{}, resource.ErrRPCTimeout)
	mockApp.EXPECT().IssueLicense(gomock.Any(), []int{0}).Return(game.License{ID: 1, DigAllowed: 3, DigUsed: 2}, nil)
	mockApp.EXPECT().IssueLicense(gomock.Any(), []int{0}).Return(game.License{}, game.ErrBogusCoin)
	mockApp.EXPECT().IssueLicense(gomock.Any(), []int{1, 2}).Return(game.License{}, game.ErrActiveLicenseLimit)

	testCases := []struct {
		args    model.Wallet
		want    interface{}
		wantErr *model.Error
	}{
		{model.Wallet{}, nil, apiError500},
		{model.Wallet{}, nil, apiError502},
		{model.Wallet{}, nil, apiError504},
		{model.Wallet{0}, &model.License{ID: swag.Int64(1), DigAllowed: 3, DigUsed: 2}, nil},
		{model.Wallet{0}, nil, apiError402},
		{model.Wallet{1, 2}, nil, apiError1002},
		{model.Wallet{}, nil, apiError429},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := c.Op.IssueLicense(op.NewIssueLicenseParams().WithArgs(tc.args))
			t.DeepEqual(openapi.ErrPayload(err), tc.wantErr)
			if res == nil {
				t.DeepEqual(nil, tc.want)
			} else {
				t.DeepEqual(res.Payload, tc.want)
			}
		})
	}
}

func TestExploreArea(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	cleanup, c, _, mockApp, _ := testNewServer(t, openapi.Config{
		OpExploreAreaRate:    4,
		OpExploreAreaTimeout: def.TestSecond / 10,
	})
	defer cleanup()

	mockApp.EXPECT().ExploreArea(gomock.Any(), game.Area{X: 4, Y: 4, SizeX: 1, SizeY: 1}).DoAndReturn(
		func(_ Ctx, _ game.Area) (int, error) {
			time.Sleep(def.TestSecond / 10 * 2)
			return 3, nil
		})
	mockApp.EXPECT().ExploreArea(gomock.Any(), game.Area{X: 0, Y: 0, SizeX: 1, SizeY: 1}).Return(0, io.EOF)
	mockApp.EXPECT().ExploreArea(gomock.Any(), game.Area{X: 0, Y: 0, SizeX: 5, SizeY: 1}).Return(0, game.ErrWrongCoord)
	mockApp.EXPECT().ExploreArea(gomock.Any(), game.Area{X: 1, Y: 2, SizeX: 3, SizeY: 4}).Return(5, nil)

	testCases := []struct {
		args    *model.Area
		want    interface{}
		wantErr *model.Error
	}{
		{&model.Area{PosX: swag.Int64(4), PosY: swag.Int64(4), SizeX: 1, SizeY: 1}, nil, apiError503},
		{&model.Area{PosX: swag.Int64(0), PosY: swag.Int64(0), SizeX: 1, SizeY: 1}, nil, apiError500},
		{&model.Area{PosX: swag.Int64(0), PosY: swag.Int64(0), SizeX: 5, SizeY: 1}, nil, apiError1000},
		{&model.Area{PosX: swag.Int64(1), PosY: swag.Int64(2), SizeX: 3, SizeY: 4}, &model.Report{Amount: 5}, nil},
		{&model.Area{PosX: swag.Int64(4), PosY: swag.Int64(4), SizeX: 1, SizeY: 1}, nil, apiError429},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := c.Op.ExploreArea(op.NewExploreAreaParams().WithArgs(tc.args))
			t.DeepEqual(openapi.ErrPayload(err), tc.wantErr)
			if tc.want != nil {
				tc.want.(*model.Report).Area = tc.args
			}
			if res == nil {
				t.DeepEqual(nil, tc.want)
			} else {
				t.DeepEqual(res.Payload, tc.want)
			}
		})
	}
}

func TestDig(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	cleanup, c, _, mockApp, _ := testNewServer(t, openapi.Config{
		OpDigRate:    7 - 1,
		OpDigTimeout: def.TestSecond / 10,
	})
	defer cleanup()

	mockApp.EXPECT().Dig(gomock.Any(), 4, game.Coord{X: 4, Y: 4, Depth: 1}).DoAndReturn(
		func(_ Ctx, _ int, _ game.Coord) (string, error) {
			time.Sleep(def.TestSecond / 10 * 2)
			return "treasure4", nil
		})
	mockApp.EXPECT().Dig(gomock.Any(), 0, game.Coord{X: 0, Y: 0, Depth: 1}).Return("", io.EOF)
	mockApp.EXPECT().Dig(gomock.Any(), 9, game.Coord{X: 0, Y: 0, Depth: 1}).Return("", game.ErrNoSuchLicense)
	mockApp.EXPECT().Dig(gomock.Any(), 0, game.Coord{X: 9, Y: 0, Depth: 1}).Return("", game.ErrWrongCoord)
	mockApp.EXPECT().Dig(gomock.Any(), 0, game.Coord{X: 0, Y: 0, Depth: 9}).Return("", game.ErrWrongDepth)
	mockApp.EXPECT().Dig(gomock.Any(), 1, game.Coord{X: 1, Y: 1, Depth: 1}).Return("", nil)
	mockApp.EXPECT().Dig(gomock.Any(), 2, game.Coord{X: 1, Y: 1, Depth: 2}).Return("treasure1", nil)

	n := swag.Int64
	testCases := []struct {
		args    *model.Dig
		want    interface{}
		wantErr *model.Error
	}{
		{&model.Dig{LicenseID: n(4), PosX: n(4), PosY: n(4), Depth: n(1)}, nil, apiError503},
		{&model.Dig{LicenseID: n(0), PosX: n(0), PosY: n(0), Depth: n(1)}, nil, apiError500},
		{&model.Dig{LicenseID: n(9), PosX: n(0), PosY: n(0), Depth: n(1)}, nil, apiError403},
		{&model.Dig{LicenseID: n(0), PosX: n(9), PosY: n(0), Depth: n(1)}, nil, apiError1000},
		{&model.Dig{LicenseID: n(0), PosX: n(0), PosY: n(0), Depth: n(9)}, nil, apiError1001},
		{&model.Dig{LicenseID: n(1), PosX: n(1), PosY: n(1), Depth: n(1)}, nil, apiError404},
		{&model.Dig{LicenseID: n(2), PosX: n(1), PosY: n(1), Depth: n(2)}, model.TreasureList{"treasure1"}, nil},
		{&model.Dig{LicenseID: n(4), PosX: n(4), PosY: n(4), Depth: n(1)}, nil, apiError429},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := c.Op.Dig(op.NewDigParams().WithArgs(tc.args))
			t.DeepEqual(openapi.ErrPayload(err), tc.wantErr)
			if res == nil {
				t.DeepEqual(nil, tc.want)
			} else {
				t.DeepEqual(res.Payload, tc.want)
			}
		})
	}
}

func TestCash(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	cleanup, c, _, mockApp, _ := testNewServer(t, openapi.Config{
		Seed:              2,
		OpCashRate:        6,
		OpCashPercentFail: 5,
	})
	defer cleanup()

	mockApp.EXPECT().Cash(gomock.Any(), "").Return(nil, io.EOF)
	mockApp.EXPECT().Cash(gomock.Any(), "bad").Return(nil, game.ErrWrongCoord)
	mockApp.EXPECT().Cash(gomock.Any(), "treasure9").Return(nil, game.ErrNotDigged)
	mockApp.EXPECT().Cash(gomock.Any(), "empty").Return(nil, game.ErrNoThreasure)
	mockApp.EXPECT().Cash(gomock.Any(), "treasure1").Return([]int{0, 1}, nil)

	testCases := []struct {
		args    model.Treasure
		want    interface{}
		wantErr *model.Error
	}{
		{"", nil, apiError500},
		{"bad", nil, apiError1000},
		{"treasure9", nil, apiError1003},
		{"empty", nil, apiError404},
		{"", nil, apiError503},
		{"treasure1", model.Wallet{0, 1}, nil},
		{"", nil, apiError429},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := c.Op.Cash(op.NewCashParams().WithArgs(tc.args))
			t.DeepEqual(openapi.ErrPayload(err), tc.wantErr)
			if res == nil {
				t.DeepEqual(nil, tc.want)
			} else {
				t.DeepEqual(res.Payload, tc.want)
			}
		})
	}
}
