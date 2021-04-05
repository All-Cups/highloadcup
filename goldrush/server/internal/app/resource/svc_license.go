package resource

import (
	"errors"
	prng "math/rand"
	"sync"
	"time"
)

// Errors.
var (
	ErrRPCTimeout  = errors.New("RPC timed out")
	ErrRPCInternal = errors.New("RPC failed")
)

// LicenseSvcConfig contains configuration for LicenseSvc.
type LicenseSvcConfig struct {
	Seed           int64
	PercentTimeout int
	MinDelay       time.Duration
	MaxDelay       time.Duration
	TimeoutDelay   time.Duration
}

// LicenseSvc implements app.LicenseSvc interface.
type LicenseSvc struct {
	cfg    LicenseSvcConfig
	muPRNG sync.Mutex
	prng   *prng.Rand
}

// NewLicenseSvc creates and returns new LicenseSvc.
func NewLicenseSvc(cfg LicenseSvcConfig) *LicenseSvc {
	if cfg.Seed == 0 {
		cfg.Seed = time.Now().UnixNano()
	}
	return &LicenseSvc{
		cfg:  cfg,
		prng: prng.New(prng.NewSource(cfg.Seed)), //nolint:gosec // We need repeatable results.
	}
}

// Call implements app.LicenseSvc interface.
func (s *LicenseSvc) Call(ctx Ctx, percentFail int) (err error) {
	delay, err := s.outcome(percentFail)
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case <-time.After(delay):
	}
	return err
}

func (s *LicenseSvc) outcome(percentFail int) (delay time.Duration, err error) {
	const percent100 = 100
	s.muPRNG.Lock()
	defer s.muPRNG.Unlock()
	if s.prng.Intn(percent100) < percentFail {
		if s.prng.Intn(percent100) < s.cfg.PercentTimeout {
			err = ErrRPCTimeout
			delay = s.cfg.TimeoutDelay
		} else {
			err = ErrRPCInternal
			delay = s.cfg.MinDelay
		}
	} else {
		delay = s.cfg.MinDelay + time.Duration(s.prng.Intn(int(s.cfg.MaxDelay-s.cfg.MinDelay+1))) // +1ns in case delays are equal.
	}
	return delay, err
}
