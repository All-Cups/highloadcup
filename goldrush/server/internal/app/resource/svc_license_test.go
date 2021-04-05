package resource_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/powerman/check"

	"github.com/Djarvur/allcups-itrally-2020-task/internal/app/resource"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/def"
)

func TestSvcLicense(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	s := resource.NewLicenseSvc(resource.LicenseSvcConfig{
		Seed:           11,
		PercentTimeout: 10,
		MinDelay:       def.TestSecond / 100,
		MaxDelay:       def.TestSecond / 10,
		TimeoutDelay:   def.TestSecond,
	})
	errc := make(chan error)
	go func() { errc <- s.Call(ctx, 50) }()
	waitErr(t, errc, def.TestSecond/35, nil)
	go func() { errc <- s.Call(ctx, 50) }()
	waitErr(t, errc, def.TestSecond/100, resource.ErrRPCInternal)
	go func() { errc <- s.Call(ctx, 50) }()
	waitErr(t, errc, def.TestSecond, resource.ErrRPCTimeout)
	go func() { errc <- s.Call(ctx, 50) }()
	waitErr(t, errc, def.TestSecond/100, resource.ErrRPCInternal)
	go func() { errc <- s.Call(ctx, 50) }()
	waitErr(t, errc, def.TestSecond/26, nil)
	ctx, cancel := context.WithCancel(ctx)
	go func() { errc <- s.Call(ctx, 0) }()
	go func() { time.Sleep(def.TestSecond / 100); cancel() }()
	waitErr(t, errc, def.TestSecond/100, context.Canceled)
}

func TestSvcLicensePercent(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	s := resource.NewLicenseSvc(resource.LicenseSvcConfig{
		Seed:           11,
		PercentTimeout: 10,
	})
	var timeout, internal int
	for i := 0; i < 1000; i++ {
		switch err := s.Call(ctx, 50); true {
		case errors.Is(err, resource.ErrRPCTimeout):
			timeout++
		case errors.Is(err, resource.ErrRPCInternal):
			internal++
		}
	}
	t.Equal(timeout, 46)
	t.Equal(internal, 421)
}
