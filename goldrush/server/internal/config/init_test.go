package config

import (
	"os"
	"testing"

	"github.com/powerman/check"
	_ "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/pflag"

	"github.com/Djarvur/allcups-itrally-2020-task/pkg/def"
)

var (
	testAll      = all
	testFlagsets = FlagSets{
		Serve: pflag.NewFlagSet("", 0),
	}
)

func TestMain(m *testing.M) {
	def.Init()
	os.Clearenv()
	check.TestMain(m)
}

func testGetServe(flags ...string) (*ServeConfig, error) {
	all = testAll
	err := Init(testFlagsets)
	if err != nil {
		return nil, err
	}
	if len(flags) > 0 {
		testFlagsets.Serve.Parse(flags)
	}
	return GetServe()
}

// Require helps testing for missing env var (required to set
// configuration value which don't have default value).
func require(t *check.C, field string) {
	t.Helper()
	c, err := testGetServe()
	t.Match(err, `^`+field+` .* required`)
	t.Nil(c)
}

// Constraint helps testing for invalid env var value.
func constraint(t *check.C, name, val, match string) {
	t.Helper()
	old, ok := os.LookupEnv(name)

	t.Nil(os.Setenv(name, val))
	c, err := testGetServe()
	t.Match(err, match)
	t.Nil(c)

	if ok {
		os.Setenv(name, old)
	} else {
		os.Unsetenv(name)
	}
}
