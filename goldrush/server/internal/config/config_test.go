package config

import (
	"os"
	"testing"
	"time"

	"github.com/powerman/check"

	"github.com/Djarvur/allcups-itrally-2020-task/internal/app"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/def"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/netx"
)

func Test(t *testing.T) {
	want := &ServeConfig{
		AccessLog:             true,
		Addr:                  netx.NewAddr(def.Hostname, 8000),
		AutosavePeriod:        time.Second,
		DepthProfitChange:     0.1,
		DigBaseDelay:          time.Millisecond,
		DigExtraDelay:         time.Millisecond / 10,
		Duration:              10 * time.Minute,
		Game:                  app.Difficulty["test"],
		LicenseMaxDelay:       time.Second / 10,
		LicenseMinDelay:       time.Second / 100,
		LicensePercentFail:    60,
		LicensePercentTimeout: 10,
		LicenseTimeoutDelay:   time.Second,
		MetricsAddr:           netx.NewAddr(def.Hostname, 9000),
		OpCashPercentFail:     5,
		OpCashRate:            300,
		OpDigRate:             1000,
		OpDigTimeout:          2 * time.Second,
		OpExploreAreaRate:     1000,
		OpExploreAreaTimeout:  time.Second,
		OpGetBalanceRate:      100,
		OpIssueLicenseRate:    350,
		OpListLicensesRate:    100,
		Pprof:                 true,
		ResultDir:             "var/data",
		StartTimeout:          2 * time.Minute,
		WorkDir:               "var",
	}

	t.Run("required", func(tt *testing.T) {
		t := check.T(tt)
		require(t, "Difficulty")
		os.Setenv("HLCUP2020_DIFFICULTY", "test")
	})
	t.Run("default", func(tt *testing.T) {
		t := check.T(tt)
		c, err := testGetServe()
		t.Nil(err)
		t.DeepEqual(c, want)
	})
	t.Run("constraint", func(tt *testing.T) {
		t := check.T(tt)
		constraint(t, "HLCUP2020_ACCESS_LOG", "x", `^AccessLog .* invalid syntax`)
		constraint(t, "HLCUP2020_ADDR_PORT", "x", `^AddrPort .* invalid syntax`)
		constraint(t, "HLCUP2020_DURATION", "x", `^Duration .* invalid duration`)
		constraint(t, "HLCUP2020_DIFFICULTY", "x", `^Difficulty .* not one of`)
		constraint(t, "HLCUP2020_METRICS_ADDR_PORT", "x", `^MetricsAddrPort .* invalid syntax`)
		constraint(t, "HLCUP2020_PPROF", "x", `^Pprof .* invalid syntax`)
		constraint(t, "HLCUP2020_RESULT_DIR", "", `^ResultDir .* empty`)
		constraint(t, "HLCUP2020_START_TIMEOUT", "x", `^StartTimeout .* invalid duration`)
		constraint(t, "HLCUP2020_WORK_DIR", "", `^WorkDir .* empty`)
	})
	t.Run("env", func(tt *testing.T) {
		t := check.T(tt)
		os.Setenv("HLCUP2020_ACCESS_LOG", "false")
		os.Setenv("HLCUP2020_ADDR_HOST", "localhost3")
		os.Setenv("HLCUP2020_ADDR_PORT", "8003")
		os.Setenv("HLCUP2020_DIFFICULTY", "normal")
		os.Setenv("HLCUP2020_DURATION", "3s")
		os.Setenv("HLCUP2020_METRICS_ADDR_PORT", "9003")
		os.Setenv("HLCUP2020_PPROF", "false")
		os.Setenv("HLCUP2020_RESULT_DIR", "/data/3")
		os.Setenv("HLCUP2020_START_TIMEOUT", "3m")
		os.Setenv("HLCUP2020_WORK_DIR", "/work/3")
		c, err := testGetServe()
		t.Nil(err)
		want.AccessLog = false
		want.Addr = netx.NewAddr("localhost3", 8003)
		want.Duration = 3 * time.Second
		want.Game = app.Difficulty["normal"]
		want.MetricsAddr = netx.NewAddr("localhost3", 9003)
		want.Pprof = false
		want.ResultDir = "/data/3"
		want.StartTimeout = 3 * time.Minute
		want.WorkDir = "/work/3"
		t.DeepEqual(c, want)
	})
	t.Run("flag", func(tt *testing.T) {
		t := check.T(tt)
		c, err := testGetServe(
			"--host=localhost4",
			"--port=8004",
			"--duration=4ms",
			"--metrics.port=9004",
			"--accesslog=true",
			"--start-timeout=4s",
		)
		t.Nil(err)
		want.AccessLog = true
		want.Addr = netx.NewAddr("localhost4", 8004)
		want.Duration = 4 * time.Millisecond
		want.MetricsAddr = netx.NewAddr("localhost4", 9004)
		want.StartTimeout = 4 * time.Second
		t.DeepEqual(c, want)
	})
	t.Run("cleanup", func(tt *testing.T) {
		t := check.T(tt)
		t.Panic(func() { GetServe() })
	})
}
