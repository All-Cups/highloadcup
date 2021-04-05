// +build integration

package main

import (
	"context"
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/powerman/check"

	"github.com/Djarvur/allcups-itrally-2020-task/api/openapi/client"
	"github.com/Djarvur/allcups-itrally-2020-task/api/openapi/client/op"
	"github.com/Djarvur/allcups-itrally-2020-task/api/openapi/model"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/app"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/srv/openapi"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/def"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/netx"
)

func TestTaskDuration(tt *testing.T) {
	t := check.T(tt)

	s := &service{cfg: cfg}
	s.cfg.Duration = def.TestSecond
	s.cfg.WorkDir = t.TempDir()
	s.cfg.ResultDir = t.TempDir()

	ctxStartup, cancel := context.WithTimeout(ctx, def.TestTimeout)
	defer cancel()
	ctxShutdown, shutdown := context.WithCancel(ctx)
	errc := make(chan error, 1)
	go func() { errc <- s.runServe(ctxStartup, ctxShutdown, shutdown) }()
	t.Must(t.Nil(netx.WaitTCPPort(ctxStartup, cfg.Addr), "connect to service"))

	openapiClient := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Schemes:  []string{"http"},
		Host:     cfg.Addr.String(),
		BasePath: client.DefaultBasePath,
	})
	openapiClient.Op.GetBalance(op.NewGetBalanceParams())

	start := time.Now()
	select {
	case err := <-errc:
		t.Nil(err)
		t.Between(time.Since(start), def.TestSecond/2, def.TestSecond*2)
	case <-time.After(def.TestTimeout):
		t.Fail()
	}
	shutdown()
}

func TestSmoke(tt *testing.T) {
	t := check.T(tt)

	s := &service{cfg: cfg}
	s.cfg.WorkDir = t.TempDir()
	s.cfg.ResultDir = t.TempDir()
	s.cfg.Game = app.Difficulty["normal"]
	s.cfg.Game.Seed = 3

	ctxStartup, cancel := context.WithTimeout(ctx, def.TestTimeout)
	defer cancel()
	ctxShutdown, shutdown := context.WithCancel(ctx)
	errc := make(chan error)
	go func() { errc <- s.runServe(ctxStartup, ctxShutdown, shutdown) }()
	defer func() {
		shutdown()
		t.Nil(<-errc, "RunServe")
	}()
	t.Must(t.Nil(netx.WaitTCPPort(ctxStartup, cfg.Addr), "connect to service"))

	openapiClient := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Schemes:  []string{"http"},
		Host:     cfg.Addr.String(),
		BasePath: client.DefaultBasePath,
	})

	var (
		area = &model.Area{
			PosX:  swag.Int64(1247),
			PosY:  swag.Int64(1366),
			SizeX: 1,
			SizeY: 1,
		}
		treasure model.Treasure
	)

	{
		res, err := openapiClient.Op.GetBalance(op.NewGetBalanceParams())
		t.Nil(err)
		t.DeepEqual(res, &op.GetBalanceOK{Payload: &model.Balance{
			Balance: swag.Uint32(0),
			Wallet:  model.Wallet{},
		}})
	}
	{
		args := model.Wallet{}
		res, err := openapiClient.Op.IssueLicense(op.NewIssueLicenseParams().WithArgs(args))
		t.DeepEqual(openapi.ErrPayload(err), apiError502)
		t.Nil(res)
	}
	{
		args := model.Wallet{}
		res, err := openapiClient.Op.IssueLicense(op.NewIssueLicenseParams().WithArgs(args))
		t.Nil(err)
		t.DeepEqual(res, &op.IssueLicenseOK{Payload: &model.License{
			ID:         swag.Int64(0),
			DigAllowed: 3,
			DigUsed:    0,
		}})
	}
	{
		args := area
		res, err := openapiClient.Op.ExploreArea(op.NewExploreAreaParams().WithArgs(args))
		t.Nil(err)
		t.DeepEqual(res, &op.ExploreAreaOK{Payload: &model.Report{
			Area:           area,
			Amount:         1,
			AmountPerDepth: nil,
		}})
	}
	{
		args := &model.Dig{
			LicenseID: swag.Int64(0),
			PosX:      area.PosX,
			PosY:      area.PosY,
			Depth:     swag.Int64(1),
		}
		res, err := openapiClient.Op.Dig(op.NewDigParams().WithArgs(args))
		t.Nil(err)
		t.NotNil(res)
		t.Len(res.Payload, 1)
		treasure = res.Payload[0]
	}
	{
		args := treasure
		res, err := openapiClient.Op.Cash(op.NewCashParams().WithArgs(args))
		t.Nil(err)
		t.DeepEqual(res, &op.CashOK{Payload: model.Wallet{0, 1, 2}})
	}
	{
		res, err := openapiClient.Op.ListLicenses(op.NewListLicensesParams())
		t.Nil(err)
		t.DeepEqual(res, &op.ListLicensesOK{Payload: model.LicenseList{
			&model.License{
				ID:         swag.Int64(0),
				DigAllowed: 3,
				DigUsed:    1,
			},
		}})
	}
}
