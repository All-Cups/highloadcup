package resource_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/powerman/check"

	"github.com/Djarvur/allcups-itrally-2020-task/internal/app/resource"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/def"
)

func TestCPUSmoke(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	if os.Getenv("CI") != "" {
		t.Skip("skipping test on CI")
	}

	ctxShutdown, shutdown := context.WithCancel(ctx)
	defer shutdown()
	cpu := resource.NewCPU(10000)
	go func() { cpu.Provide(ctxShutdown) }()
	errc := make(chan error)
	errc2 := make(chan error)

	for _, d := range []time.Duration{time.Second / 10, time.Second / 100, time.Second / 1000} {
		go func() { errc <- cpu.Consume(ctx, d) }()
		waitErr(t, errc, d, nil)
		go func() { errc2 <- cpu.Consume(ctx, d) }()
		go func() { errc2 <- cpu.Consume(ctx, d) }()
		go func() { <-errc2; errc <- <-errc2 }()
		waitErr(t, errc, d*2, nil)
	}
}

func TestCPUProduceCtx(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	errc := make(chan error)
	ctx, cancel := context.WithCancel(ctx)

	cpu := resource.NewCPU(1000)
	go func() { time.Sleep(def.TestSecond / 10); cancel() }()
	go func() { errc <- cpu.Provide(ctx) }()
	waitErr(t, errc, def.TestSecond/10, nil)
	go func() { time.Sleep(def.TestSecond / 10); errc <- cpu.Provide(ctx) }()
	waitErr(t, errc, def.TestSecond/10, nil)
}

func TestCPUConsumeCtx(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	errc := make(chan error)
	ctxShutdown, shutdown := context.WithCancel(ctx)
	defer shutdown()
	ctx, cancel := context.WithCancel(ctx)

	cpu := resource.NewCPU(1000)
	go func() { time.Sleep(def.TestSecond / 10); cancel() }()
	go func() { errc <- cpu.Provide(ctxShutdown) }()
	go func() { errc <- cpu.Consume(ctx, def.TestSecond) }()
	waitErr(t, errc, def.TestSecond/10, context.Canceled)
	go func() { time.Sleep(def.TestSecond / 10); errc <- cpu.Consume(ctx, def.TestSecond) }()
	waitErr(t, errc, def.TestSecond/10, context.Canceled)
}
