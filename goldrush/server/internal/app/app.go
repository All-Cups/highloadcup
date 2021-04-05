//go:generate gobin -m -run github.com/golang/mock/mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mock.$GOFILE Appl,Repo,GameFactory,CPU,LicenseSvc

// Package app provides business logic.
package app

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	prng "math/rand"
	"sync"
	"time"

	"github.com/powerman/must"
	"github.com/powerman/structlog"

	"github.com/Djarvur/allcups-itrally-2020-task/internal/app/game"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

// Errors.
var (
	ErrContactExists    = errors.New("contact already exists")
	errBadPASETOKeySize = errors.New("bad PASETO key size")
)

// Appl provides application features (use cases) service.
type Appl interface {
	// HealthCheck returns error if service is unhealthy or current
	// status otherwise.
	// Errors: none.
	HealthCheck(Ctx) (interface{}, error)
	// Start must be called before any other method to ensure task
	// will be available for cfg.Duration since given time. Second and
	// following calls will have no effect, so it's safe to call Start
	// on every API call.
	// Errors: none.
	Start(time.Time) error
	// Balance returns current balance and up to 1000 issued coins.
	// Errors: none.
	Balance(Ctx) (balance int, wallet []int, err error)
	// Licenses returns all active licenses.
	// Errors: resource.ErrRPCInternal, resource.ErrRPCTimeout.
	Licenses(Ctx) ([]game.License, error)
	// IssueLicense creates and returns a new license with given digAllowed.
	// Errors: game.ErrActiveLicenseLimit, game.ErrBogusCoin,
	// resource.ErrRPCInternal, resource.ErrRPCTimeout.
	IssueLicense(_ Ctx, wallet []int) (game.License, error)
	// ExploreArea returns amount of not-digged-yet treasures in the
	// area at depth.
	// Errors: game.ErrWrongCoord.
	ExploreArea(_ Ctx, area game.Area) (int, error)
	// Dig tries to dig at pos and returns if any treasure was found.
	// The pos depth must be next to current (already digged) one.
	// Also it increment amount of used dig calls in given active license.
	// If amount of used dig calls became equal to amount of allowed dig calls
	// then license will became inactive after the call.
	// Errors: game.ErrNoSuchLicense, game.ErrWrongCoord, game.ErrWrongDepth.
	Dig(_ Ctx, licenseID int, pos game.Coord) (treasure string, _ error)
	// Cash returns coins earned for treasure as given pos.
	// Errors: game.ErrWrongCoord, game.ErrNotDigged, game.ErrAlreadyCached.
	Cash(_ Ctx, treasure string) (wallet []int, err error)
}

// Repo provides data storage.
type Repo interface {
	// LoadStartTime returns start time or zero time if not started.
	// Errors: none.
	LoadStartTime() (*time.Time, error)
	// SaveStartTime stores start time.
	// Errors: none.
	SaveStartTime(t time.Time) error
	// LoadTreasureKey returns treasure key.
	// Errors: none.
	LoadTreasureKey() ([]byte, error)
	// SaveTreasureKey stores treasure key.
	// Errors: none.
	SaveTreasureKey([]byte) error
	// LoadGame returns game state.
	// Errors: none.
	LoadGame() (ReadSeekCloser, error)
	// SaveGame stores game state.
	// Errors: none.
	SaveGame(io.WriterTo) error
	// SaveResult stores final game result.
	// Errors: none.
	SaveResult(int) error
	// SaveError stores final game error.
	// Errors: none.
	SaveError(msg string) error
}

// GameFactory provides different ways to create a new game.
type GameFactory interface {
	// New creates and returns new game.
	New(Ctx, game.Config) (game.Game, error)
	// Continue creates and returns new game restored from given reader, which
	// should contain data written by Game.WriteTo.
	Continue(Ctx, io.ReadSeeker) (game.Game, error)
}

// CPU is a resource which can be consumed for up to time.Second per
// real-time second (i.e. it's a single-core CPU).
type CPU interface {
	// Consume t resources of this CPU instance.
	// It returns nil if consumed successfully or ctx.Err() if ctx is done
	// earlier than t resources will be consumed.
	Consume(Ctx, time.Duration) error
}

// LicenseSvc is a virtual resource which pretends to be an RPC client.
type LicenseSvc interface {
	// Call will use percentFail and percentTimeout to decide call result:
	//   - delay 0.01…0.1 sec without error
	//   - delay 0.01 sec with ErrRPCInternal
	//   - delay 1 sec with ErrRPCTimeout
	Call(ctx Ctx, percentFail int) error
}

type (
	// Contact describes record in address book.
	Contact struct {
		ID   int
		Name string
	}
	// ReadSeekCloser is the interface that groups the basic Read,
	// Seek and Close methods.
	ReadSeekCloser interface {
		io.ReadSeeker
		io.Closer
	}
)

const pasetoKeySize = 32

// Difficulty contains predefined game difficulty levels.
//nolint:gochecknoglobals,gomnd // Const.
var Difficulty = map[string]game.Config{
	"test": {
		MaxActiveLicenses: 3,
		Density:           4,
		SizeX:             5,
		SizeY:             5,
		Depth:             10,
		TreasureValue: func() *[]int {
			treasureValues := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			return &treasureValues
		}(),
		TreasureValueAlg: game.AlgDoubleMax,
	},
	"normal": {
		MaxActiveLicenses: 10,
		Density:           250,
		SizeX:             3500,
		SizeY:             3500,
		Depth:             10,
	},
}

type Config struct {
	StartTimeout       time.Duration
	Duration           time.Duration
	Game               game.Config
	AutosavePeriod     time.Duration
	DigBaseDelay       time.Duration
	DigExtraDelay      time.Duration
	DepthProfitChange  float64
	LicensePercentFail int
}

// App implements interface Appl.
type App struct {
	repo       Repo
	cpu        CPU
	svcLicense LicenseSvc
	cfg        Config
	game       game.Game
	started    chan time.Time
	startOnce  sync.Once
	key        []byte
}

func New(ctx Ctx, repo Repo, cpu CPU, svcLicense LicenseSvc, factory GameFactory, cfg Config) (*App, error) {
	a := &App{
		repo:       repo,
		cpu:        cpu,
		svcLicense: svcLicense,
		cfg:        cfg,
		started:    make(chan time.Time, 1),
		key:        make([]byte, pasetoKeySize),
	}
	t, err := a.repo.LoadStartTime()
	if err != nil {
		return nil, fmt.Errorf("LoadStartTime: %w", err)
	}
	if t.IsZero() {
		err = a.newGame(ctx, factory)
	} else {
		err = a.continueGame(ctx, factory, *t)
	}
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) newGame(ctx Ctx, factory GameFactory) (err error) {
	log := structlog.FromContext(ctx, nil)
	if a.cfg.Game != Difficulty["test"] && a.cfg.Game.Seed == 0 {
		a.cfg.Game.Seed = time.Now().UnixNano()
	}
	if a.cfg.Game.TreasureValue == nil {
		a.cfg.Game.TreasureValueAlg = game.AlgQuarterAround
		a.calcTreasureValue(ctx)
	}

	_, err = io.ReadFull(rand.Reader, a.key)
	must.NoErr(err)

	a.game, err = factory.New(ctx, a.cfg.Game)
	if err != nil {
		return fmt.Errorf("newGame: %w", err)
	}

	err = a.repo.SaveTreasureKey(a.key)
	if err != nil {
		return fmt.Errorf("SaveTreasureKey: %w", err)
	}
	err = a.repo.SaveGame(a.game)
	if err != nil {
		return fmt.Errorf("SaveGame: %w", err)
	}

	log.Info("new game")
	return nil
}

func (a *App) continueGame(ctx Ctx, factory GameFactory, t time.Time) (err error) {
	log := structlog.FromContext(ctx, nil)
	a.key, err = a.repo.LoadTreasureKey()
	if err != nil {
		return fmt.Errorf("LoadTreasureKey: %w", err)
	}
	if len(a.key) != pasetoKeySize {
		return fmt.Errorf("%w: %d", errBadPASETOKeySize, len(a.key))
	}

	f, err := a.repo.LoadGame()
	if err != nil {
		return fmt.Errorf("LoadGame: %w", err)
	}
	a.game, err = factory.Continue(ctx, f)
	if err != nil {
		return fmt.Errorf("factory.Continue: %w", err)
	}
	err = f.Close()
	if err != nil {
		return fmt.Errorf("LoadGame.Close: %w", err)
	}

	log.Info("continue game")
	err = a.Start(t)
	if err != nil {
		return fmt.Errorf("SaveStartTime: %w", err)
	}
	return nil
}

func (a *App) HealthCheck(_ Ctx) (interface{}, error) {
	return "OK", nil
}

func (a *App) calcTreasureValue(ctx Ctx) {
	log := structlog.FromContext(ctx, nil)
	value := make([]int, a.cfg.Game.Depth)
	a.cfg.Game.TreasureValue = &value

	const minValue = 4 // Should be >=4, so ±25% in game.treasureValueAt will result in random, non-constant value.

	bestDepth := uint8(3 + prng.New(prng.NewSource(a.cfg.Game.Seed)).Intn(5)) //nolint:gomnd,gosec // Balance: 5±2.

	delay := make([]float64, a.cfg.Game.Depth)
	timeToDig := make([]float64, a.cfg.Game.Depth)
	delay[0] = a.cfg.DigBaseDelay.Seconds()
	timeToDig[0] = a.cfg.DigBaseDelay.Seconds()
	value[0] = minValue
	log.Debug("calcTreasureValue", "depth", 1, "value", value[0], "maxProfitPerSecond", value[0]*int(1/timeToDig[0]))
	for i := 1; i < int(a.cfg.Game.Depth); i++ {
		depth := uint8(i + 1)
		delay[i] = delay[i-1] + a.cfg.DigExtraDelay.Seconds()
		timeToDig[i] = timeToDig[i-1] + delay[i]
		coeff := 1 + a.cfg.DepthProfitChange
		if bestDepth < depth {
			coeff = 1 / (1 + a.cfg.DepthProfitChange)
		}
		profitPerSecond := float64(value[i-1]) * (1 / timeToDig[i-1]) * coeff
		foundPerSecond := 1 / timeToDig[i]
		value[i] = int(profitPerSecond / foundPerSecond)
		log.Debug("calcTreasureValue", "depth", depth, "value", value[i], "maxProfitPerSecond", int(profitPerSecond))
	}
}
