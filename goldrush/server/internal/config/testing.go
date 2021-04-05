package config

import (
	"github.com/powerman/must"
	"github.com/spf13/pflag"

	"github.com/Djarvur/allcups-itrally-2020-task/internal/app"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/netx"
)

// MustGetServeTest returns config suitable for use in tests.
func MustGetServeTest() *ServeConfig {
	err := Init(FlagSets{
		Serve: pflag.NewFlagSet("", pflag.ContinueOnError),
	})
	must.NoErr(err)
	cfg, err := GetServe()
	must.NoErr(err)

	const host = "localhost"
	cfg.Addr = netx.NewAddr(host, netx.UnusedTCPPort(host))
	cfg.MetricsAddr = netx.NewAddr(host, 0)
	cfg.Game = app.Difficulty["test"]
	cfg.Game.Seed = 3

	return cfg
}
