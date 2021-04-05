// Package resource implements virtual resources, which simulate some
// real-world resources without actually consuming them.
package resource

import (
	"context"
	"fmt"
	"time"
)

type Ctx = context.Context

// CPU implements app.CPU interface.
type CPU struct {
	freq time.Duration
	tick chan struct{}
}

// NewCPU creates and returns new CPU.
func NewCPU(hz int) *CPU {
	const compensate = 2 * time.Millisecond // Compensate slow consumer for up to 0.002s.
	if hz < 1 {
		panic(fmt.Sprintf("hz must be a positive number: %d", hz))
	}
	return &CPU{
		freq: time.Second / time.Duration(hz),
		tick: make(chan struct{}, hz/int(time.Second/compensate)),
	}
}

// Provide should be started in background once per each CPU instance.
// It'll return nil when ctx is done.
//
// When Provide is not running Consume won't be able to consume resources
// and will return only when their ctx is done.
func (c *CPU) Provide(ctx Ctx) error {
	prev := time.Now().Round(c.freq)
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case now := <-ticker.C:
			now = now.Round(c.freq)
			for i := 0; i < int(now.Sub(prev)/c.freq); i++ {
				select {
				case c.tick <- struct{}{}:
				default:
				}
			}
			prev = now
		}
	}
}

// Consume implements app.CPU interface.
func (c *CPU) Consume(ctx Ctx, t time.Duration) error {
	for ticks := int(t.Round(c.freq) / c.freq); ticks > 0; ticks-- {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.tick:
		}
	}
	return nil
}
