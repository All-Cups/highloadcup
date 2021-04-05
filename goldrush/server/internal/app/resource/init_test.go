package resource_test

import (
	"context"
	"testing"
	"time"

	"github.com/powerman/check"

	"github.com/Djarvur/allcups-itrally-2020-task/pkg/def"
)

func TestMain(m *testing.M) {
	def.Init()
	check.TestMain(m)
}

var ctx = context.Background()

func waitErr(t *check.C, errc <-chan error, wait time.Duration, wantErr error) {
	t.Helper()
	now := time.Now()
	select {
	case err := <-errc:
		t.Between(time.Since(now), wait-wait/2, wait+wait/2)
		t.Err(err, wantErr)
	case <-time.After(def.TestTimeout):
		t.FailNow()
	}
}
