// Package config provides configurations for subcommands.
//
// Default values can be obtained from various sources (constants,
// environment variables, etc.) and then overridden by flags.
//
// As configuration is global you can get it only once for safety:
// you can call only one of Get… functions and call it just once.
package config

import (
	"errors"
	"time"

	"github.com/powerman/appcfg"
	"github.com/spf13/pflag"

	"github.com/Djarvur/allcups-itrally-2020-task/internal/app"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/app/game"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/def"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/netx"
)

// EnvPrefix defines common prefix for environment variables.
const envPrefix = "HLCUP2020_"

var errLicenseMinDelay = errors.New("LicenseMinDelay should be <= LicenseMaxDelay")

// All configurable values of the microservice.
//
// If microservice may runs in different ways (e.g. using CLI subcommands)
// then these subcommands may use subset of these values.
var all = &struct { //nolint:gochecknoglobals // Config is global anyway.
	AccessLog             appcfg.Bool           `env:"ACCESS_LOG"`
	AddrHost              appcfg.NotEmptyString `env:"ADDR_HOST"`
	AddrPort              appcfg.Port           `env:"ADDR_PORT"`
	AutosavePeriod        appcfg.Duration
	DepthProfitChange     appcfg.Float64
	Difficulty            appcfg.OneOfString `env:"DIFFICULTY"`
	DigBaseDelay          appcfg.Duration
	DigExtraDelay         appcfg.Duration
	Duration              appcfg.Duration `env:"DURATION"`
	LicenseMaxDelay       appcfg.Duration
	LicenseMinDelay       appcfg.Duration
	LicensePercentFail    appcfg.IntBetween
	LicensePercentTimeout appcfg.IntBetween
	LicenseTimeoutDelay   appcfg.Duration
	MetricsAddrPort       appcfg.Port `env:"METRICS_ADDR_PORT"`
	OpCashPercentFail     appcfg.IntBetween
	OpCashRate            appcfg.Uint
	OpDigRate             appcfg.Uint
	OpDigTimeout          appcfg.Duration
	OpExploreAreaRate     appcfg.Uint
	OpExploreAreaTimeout  appcfg.Duration
	OpGetBalanceRate      appcfg.Uint
	OpIssueLicenseRate    appcfg.Uint
	OpListLicensesRate    appcfg.Uint
	Pprof                 appcfg.Bool           `env:"PPROF"`
	ResultDir             appcfg.NotEmptyString `env:"RESULT_DIR"`
	Seed                  appcfg.Int64          `env:"SEED"`
	StartTimeout          appcfg.Duration       `env:"START_TIMEOUT"`
	WorkDir               appcfg.NotEmptyString `env:"WORK_DIR"`
}{ // Defaults, if any:
	AccessLog:             appcfg.MustBool("true"),
	AddrHost:              appcfg.MustNotEmptyString(def.Hostname),
	AddrPort:              appcfg.MustPort("8000"),
	AutosavePeriod:        appcfg.MustDuration("1s"),
	DepthProfitChange:     appcfg.MustFloat64("0.1"), // 10% of max profit per depth distance from "best" depth.
	Difficulty:            appcfg.NewOneOfString(difficulties()),
	DigBaseDelay:          appcfg.MustDuration("1ms"),   // Real dig rate limit (1000 RPS).
	DigExtraDelay:         appcfg.MustDuration("0.1ms"), // (Optional) To make dig a bit slower with increased depth.
	Duration:              appcfg.MustDuration("10m"),
	LicenseMaxDelay:       appcfg.MustDuration("0.1s"),
	LicenseMinDelay:       appcfg.MustDuration("0.01s"),
	LicensePercentFail:    appcfg.MustIntBetween("60", 0, 100), // Only for ListLicenses and issue free license.
	LicensePercentTimeout: appcfg.MustIntBetween("10", 0, 100),
	LicenseTimeoutDelay:   appcfg.MustDuration("1s"),
	MetricsAddrPort:       appcfg.MustPort("9000"),
	OpCashPercentFail:     appcfg.MustIntBetween("5", 0, 100),
	OpCashRate:            appcfg.MustUint("300"),
	OpDigRate:             appcfg.MustUint("1000"),
	OpDigTimeout:          appcfg.MustDuration("2s"),
	OpExploreAreaRate:     appcfg.MustUint("1000"),
	OpExploreAreaTimeout:  appcfg.MustDuration("1s"),
	OpGetBalanceRate:      appcfg.MustUint("100"),
	OpIssueLicenseRate:    appcfg.MustUint("350"),
	OpListLicensesRate:    appcfg.MustUint("100"),
	Pprof:                 appcfg.MustBool("true"),
	ResultDir:             appcfg.MustNotEmptyString("var/data"),
	Seed:                  appcfg.MustInt64("0"),
	StartTimeout:          appcfg.MustDuration("2m"),
	WorkDir:               appcfg.MustNotEmptyString("var"),
}

// FlagSets for all CLI subcommands which use flags to set config values.
type FlagSets struct {
	Serve *pflag.FlagSet
}

var fs FlagSets //nolint:gochecknoglobals // Flags are global anyway.

// Init updates config defaults (from env) and setup subcommands flags.
//
// Init must be called once before using this package.
func Init(flagsets FlagSets) error {
	fs = flagsets

	fromEnv := appcfg.NewFromEnv(envPrefix)
	err := appcfg.ProvideStruct(all, fromEnv)
	if err != nil {
		return err
	}

	appcfg.AddPFlag(fs.Serve, &all.AccessLog, "accesslog", "enable/disable accesslog")
	appcfg.AddPFlag(fs.Serve, &all.AddrHost, "host", "host to serve OpenAPI")
	appcfg.AddPFlag(fs.Serve, &all.AddrPort, "port", "port to serve OpenAPI")
	appcfg.AddPFlag(fs.Serve, &all.Duration, "duration", "overall task duration")
	appcfg.AddPFlag(fs.Serve, &all.MetricsAddrPort, "metrics.port", "port to serve Prometheus metrics")
	appcfg.AddPFlag(fs.Serve, &all.Seed, "seed", "PRNG seed (0 means random)")
	appcfg.AddPFlag(fs.Serve, &all.StartTimeout, "start-timeout", "task start timeout")

	return nil
}

// ServeConfig contains configuration for subcommand.
type ServeConfig struct {
	AccessLog             bool
	Addr                  netx.Addr
	AutosavePeriod        time.Duration
	DepthProfitChange     float64
	DigBaseDelay          time.Duration
	DigExtraDelay         time.Duration
	Duration              time.Duration
	Game                  game.Config
	LicenseMaxDelay       time.Duration
	LicenseMinDelay       time.Duration
	LicensePercentFail    int
	LicensePercentTimeout int
	LicenseTimeoutDelay   time.Duration
	MetricsAddr           netx.Addr
	OpCashPercentFail     int
	OpCashRate            int
	OpDigRate             int
	OpDigTimeout          time.Duration
	OpExploreAreaRate     int
	OpExploreAreaTimeout  time.Duration
	OpGetBalanceRate      int
	OpIssueLicenseRate    int
	OpListLicensesRate    int
	Pprof                 bool
	ResultDir             string
	StartTimeout          time.Duration
	WorkDir               string
}

// GetServe validates and returns configuration for subcommand.
func GetServe() (c *ServeConfig, err error) {
	defer cleanup()

	c = &ServeConfig{
		AccessLog:             all.AccessLog.Value(&err),
		Addr:                  netx.NewAddr(all.AddrHost.Value(&err), all.AddrPort.Value(&err)),
		AutosavePeriod:        all.AutosavePeriod.Value(&err),
		DepthProfitChange:     all.DepthProfitChange.Value(&err),
		DigBaseDelay:          all.DigBaseDelay.Value(&err),
		DigExtraDelay:         all.DigExtraDelay.Value(&err),
		Duration:              all.Duration.Value(&err),
		Game:                  app.Difficulty[all.Difficulty.Value(&err)],
		LicenseMaxDelay:       all.LicenseMaxDelay.Value(&err),
		LicenseMinDelay:       all.LicenseMinDelay.Value(&err),
		LicensePercentFail:    all.LicensePercentFail.Value(&err),
		LicensePercentTimeout: all.LicensePercentTimeout.Value(&err),
		LicenseTimeoutDelay:   all.LicenseTimeoutDelay.Value(&err),
		MetricsAddr:           netx.NewAddr(all.AddrHost.Value(&err), all.MetricsAddrPort.Value(&err)),
		OpCashPercentFail:     all.OpCashPercentFail.Value(&err),
		OpCashRate:            int(all.OpCashRate.Value(&err)),
		OpDigRate:             int(all.OpDigRate.Value(&err)),
		OpDigTimeout:          all.OpDigTimeout.Value(&err),
		OpExploreAreaRate:     int(all.OpExploreAreaRate.Value(&err)),
		OpExploreAreaTimeout:  all.OpExploreAreaTimeout.Value(&err),
		OpGetBalanceRate:      int(all.OpGetBalanceRate.Value(&err)),
		OpIssueLicenseRate:    int(all.OpIssueLicenseRate.Value(&err)),
		OpListLicensesRate:    int(all.OpListLicensesRate.Value(&err)),
		Pprof:                 all.Pprof.Value(&err),
		ResultDir:             all.ResultDir.Value(&err),
		StartTimeout:          all.StartTimeout.Value(&err),
		WorkDir:               all.WorkDir.Value(&err),
	}
	c.Game.Seed = all.Seed.Value(&err)
	if err == nil && c.LicenseMinDelay > c.LicenseMaxDelay {
		err = errLicenseMinDelay
	}
	if err != nil {
		return nil, appcfg.WrapPErr(err, fs.Serve, all)
	}
	return c, nil
}

// Cleanup must be called by all Get* functions to ensure second call to
// any of them will panic.
func cleanup() {
	all = nil
}

func difficulties() []string {
	levels := make([]string, 0, len(app.Difficulty))
	for level := range app.Difficulty {
		levels = append(levels, level)
	}
	return levels
}
