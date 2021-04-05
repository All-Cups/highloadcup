package main

import (
	"context"
	"regexp"

	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"

	"github.com/Djarvur/allcups-itrally-2020-task/api/openapi/restapi"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/app"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/app/game"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/app/resource"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/config"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/dal"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/srv/openapi"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/concurrent"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/def"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/serve"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

var reg = prometheus.NewPedanticRegistry() //nolint:gochecknoglobals // Metrics are global anyway.

type service struct {
	cfg        *config.ServeConfig
	repo       *dal.Repo
	cpu        *resource.CPU
	svcLicense *resource.LicenseSvc
	appl       *app.App
	srv        *restapi.Server
}

func initService(_, serveCmd *cobra.Command) error {
	namespace := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(def.ProgName, "_")
	initMetrics(reg, namespace)
	app.InitMetrics(reg)
	openapi.InitMetrics(reg, namespace)

	return config.Init(config.FlagSets{
		Serve: serveCmd.Flags(),
	})
}

func (s *service) runServe(ctxStartup, ctxShutdown Ctx, shutdown func()) (err error) {
	log := structlog.FromContext(ctxShutdown, nil)
	if s.cfg == nil {
		s.cfg, err = config.GetServe()
	}
	if err != nil {
		return log.Err("failed to get config", "err", err)
	}

	err = concurrent.Setup(ctxStartup, map[interface{}]concurrent.SetupFunc{
		&s.repo: s.connectRepo,
	})
	if err != nil {
		return log.Err("failed to connect", "err", err)
	}

	const hz = 10000
	s.cpu = resource.NewCPU(hz)
	s.svcLicense = resource.NewLicenseSvc(resource.LicenseSvcConfig{
		Seed:           s.cfg.Game.Seed,
		PercentTimeout: s.cfg.LicensePercentTimeout,
		MinDelay:       s.cfg.LicenseMinDelay,
		MaxDelay:       s.cfg.LicenseMaxDelay,
		TimeoutDelay:   s.cfg.LicenseTimeoutDelay,
	})

	if s.appl == nil {
		s.appl, err = app.New(ctxStartup, s.repo, s.cpu, s.svcLicense, game.Factory{}, app.Config{
			AutosavePeriod:     s.cfg.AutosavePeriod,
			DepthProfitChange:  s.cfg.DepthProfitChange,
			DigBaseDelay:       s.cfg.DigBaseDelay,
			DigExtraDelay:      s.cfg.DigExtraDelay,
			Duration:           s.cfg.Duration,
			Game:               s.cfg.Game,
			LicensePercentFail: s.cfg.LicensePercentFail,
			StartTimeout:       s.cfg.StartTimeout,
		})
	}
	if err != nil {
		return log.Err("failed to app.New", "err", err)
	}
	s.srv, err = openapi.NewServer(s.appl, openapi.Config{
		Seed:                 s.cfg.Game.Seed,
		DisableAccessLog:     !s.cfg.AccessLog,
		Addr:                 s.cfg.Addr,
		Pprof:                s.cfg.Pprof,
		OpCashPercentFail:    s.cfg.OpCashPercentFail,
		OpCashRate:           s.cfg.OpCashRate,
		OpDigRate:            s.cfg.OpDigRate,
		OpDigTimeout:         s.cfg.OpDigTimeout,
		OpExploreAreaRate:    s.cfg.OpExploreAreaRate,
		OpExploreAreaTimeout: s.cfg.OpExploreAreaTimeout,
		OpGetBalanceRate:     s.cfg.OpGetBalanceRate,
		OpIssueLicenseRate:   s.cfg.OpIssueLicenseRate,
		OpListLicensesRate:   s.cfg.OpListLicensesRate,
	})
	if err != nil {
		return log.Err("failed to openapi.NewServer", "err", err)
	}

	err = concurrent.Serve(ctxShutdown, shutdown,
		s.cpu.Provide,
		s.appl.Wait,
		s.serveOpenAPI,
		s.serveMetrics,
	)
	if err != nil {
		return log.Err("failed to serve", "err", err)
	}
	return nil
}

func (s *service) connectRepo(ctx Ctx) (interface{}, error) {
	return dal.New(ctx, dal.Config{
		ResultDir: s.cfg.ResultDir,
		WorkDir:   s.cfg.WorkDir,
	})
}

func (s *service) serveMetrics(ctx Ctx) error {
	return serve.Metrics(ctx, s.cfg.MetricsAddr, reg)
}

func (s *service) serveOpenAPI(ctx Ctx) error {
	return serve.OpenAPI(ctx, s.srv, "OpenAPI")
}
