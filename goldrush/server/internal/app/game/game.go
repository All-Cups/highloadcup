//go:generate gobin -m -run github.com/golang/mock/mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mock.$GOFILE Game

// Package game implements treasure hunting game.
package game

import (
	"context"
	"errors"
	"fmt"
	"io"
	prng "math/rand"
	"sync"

	"github.com/powerman/structlog"
)

type Ctx = context.Context

const (
	maxSizeX, maxSizeY, maxDepth = 6000, 6000, 10 // About 1GB RAM without bit-optimization.
	maxWalletSize                = 1000           // Needed to keep Balance fast.
	maxDigAllowed                = 100            // Not required, just for fun.
)

// Errors.
var (
	ErrActiveLicenseLimit = errors.New("no more active licenses allowed")
	ErrNoThreasure        = errors.New("no treasure")
	ErrBogusCoin          = errors.New("bogus coin")
	ErrNoSuchLicense      = errors.New("no such license")
	ErrNotDigged          = errors.New("treasure is not digged")
	ErrWrongCoord         = errors.New("wrong coordinates")
	ErrWrongDepth         = errors.New("wrong depth")

	errOutOfBounds = errors.New("out of bounds")
	errWrongAmount = errors.New("wrong amount of coins")
)

// Game implements treasure hunting game.
type Game interface {
	// WriteTo saves current game state.
	WriteTo(w io.Writer) (n int64, err error)
	// Balance returns current balance and up to 1000 issued coins.
	Balance() (balance int, wallet []int)
	// Licenses returns all active licenses.
	Licenses() []License
	// IssueLicense creates and returns a new license with given digAllowed.
	// Errors: ErrActiveLicenseLimit, ErrBogusCoin.
	IssueLicense(wallet []int) (License, error)
	// CountTreasures returns amount of not-digged-yet treasures in the area
	// at depth.
	// Errors: ErrWrongCoord, ErrWrongDepth.
	CountTreasures(area Area, depth uint8) (int, error)
	// Dig tries to dig at pos and returns if any treasure was found.
	// The pos depth must be next to current (already digged) one.
	// Also it increment amount of used dig calls in given active license.
	// If amount of used dig calls became equal to amount of allowed dig calls
	// then license will became inactive after the call.
	// Errors: ErrNoSuchLicense, ErrWrongCoord, ErrWrongDepth.
	Dig(licenseID int, pos Coord) (found bool, _ error)
	// Cash returns coins earned for treasure as given pos.
	// Errors: ErrWrongCoord, ErrNotDigged, ErrAlreadyCached.
	Cash(pos Coord) (wallet []int, err error)
}

type (
	// License defines amount of allowed dig calls.
	License struct {
		ID         int
		DigAllowed int
		DigUsed    int
	}
	// Area describes rectangle.
	Area struct {
		X     int // From 0.
		Y     int // From 0.
		SizeX int // From 1.
		SizeY int // From 1.
	}
	// Coord describes single cell.
	Coord struct {
		X     int   // From 0.
		Y     int   // From 0.
		Depth uint8 // From 1.
	}
	// TreasureValueAlg define which algorithm to use for calculating
	// min/max treasure value range at given depth.
	TreasureValueAlg int
)

const (
	// AlgDoubleMax generates treasure values in range x…x*2
	// where x is Config.TreasureValue[depth].
	AlgDoubleMax TreasureValueAlg = iota + 1
	// AlgQuarterAround generates treasure values in range x*0.75…x*1.25
	// where x is Config.TreasureValue[depth].
	AlgQuarterAround
)

// Config contains game configuration.
type Config struct {
	Seed              int64
	MaxActiveLicenses int
	Density           int // About one treasure per Density cells.
	SizeX             int
	SizeY             int
	Depth             uint8
	TreasureValue     *[]int // Pointer to keep this struct comparable.
	TreasureValueAlg  TreasureValueAlg
}

type game struct {
	cfg      Config
	muModify sync.RWMutex
	licenses *licenses
	bank     *bank
	field    *field
	muPRNG   sync.Mutex
	prng     *prng.Rand
}

// Factory implements app.GameFactory interface.
type Factory struct{}

// New implements app.GameFactory interface.
func (Factory) New(ctx Ctx, cfg Config) (Game, error) {
	log := structlog.FromContext(ctx, nil)
	switch {
	case cfg.Density <= 0, cfg.Density > cfg.volume(): // Min 1 treasure.
		return nil, fmt.Errorf("%w: Density", errOutOfBounds)
	case cfg.SizeX <= 0, cfg.SizeX > maxSizeX:
		return nil, fmt.Errorf("%w: SizeX", errOutOfBounds)
	case cfg.SizeY <= 0, cfg.SizeY > maxSizeY:
		return nil, fmt.Errorf("%w: SizeY", errOutOfBounds)
	case cfg.Depth <= 0, cfg.Depth > maxDepth:
		return nil, fmt.Errorf("%w: Depth", errOutOfBounds)
	case cfg.TreasureValue == nil || len(*cfg.TreasureValue) != int(cfg.Depth):
		panic(fmt.Sprintf("TreasureValue length must be %d", cfg.Depth))
	}

	g := &game{
		cfg:      cfg,
		licenses: newLicenses(cfg.MaxActiveLicenses),
		bank:     newBank(ctx, cfg.totalCash()),
		field:    newField(cfg),
		prng:     prng.New(prng.NewSource(cfg.Seed)), //nolint:gosec // We need repeatable results.
	}

	skipped := 0
	for i := 0; i < cfg.treasures(); i++ {
		pos := Coord{
			X:     g.prng.Intn(cfg.SizeX),
			Y:     g.prng.Intn(cfg.SizeY),
			Depth: uint8(g.prng.Intn(int(cfg.Depth)) + 1),
		}
		if !g.field.addTreasure(pos) {
			skipped++
		} else if i < 10 { //nolint:gomnd // Debug.
			log.Debug("buried one of first 10 treasures", "pos", pos)
		}
	}
	log.Info("the treasures were buried", "all", cfg.treasures(), "skipped", skipped)
	return g, nil
}

func (g *game) Balance() (balance int, wallet []int) {
	return g.bank.getBalance()
}

func (g *game) Licenses() []License {
	return g.licenses.active()
}

func (g *game) IssueLicense(wallet []int) (l License, err error) {
	g.muModify.RLock()
	defer g.muModify.RUnlock()

	digAllowed := g.licensePrice(len(wallet))
	license, err := g.licenses.beginIssue(digAllowed)
	if err != nil {
		return l, err
	}
	err = g.bank.spend(wallet)
	if err != nil {
		g.licenses.rollbackIssue(license.ID)
		return l, err
	}
	g.licenses.commitIssue(license.ID)
	return license, nil
}

func (g *game) CountTreasures(area Area, depth uint8) (int, error) {
	return g.field.countTreasures(area, depth)
}

func (g *game) Dig(licenseID int, pos Coord) (found bool, _ error) {
	g.muModify.RLock()
	defer g.muModify.RUnlock()

	err := g.licenses.use(licenseID)
	if err != nil {
		return false, err
	}
	return g.field.dig(pos)
}

func (g *game) Cash(pos Coord) (wallet []int, err error) {
	g.muModify.RLock()
	defer g.muModify.RUnlock()

	err = g.field.cash(pos)
	if err != nil {
		return nil, err
	}

	g.muPRNG.Lock()
	defer g.muPRNG.Unlock()
	min, max := g.cfg.treasureValueAt(pos.Depth)
	amount := min + g.prng.Intn(max-min+1)

	if amount <= 0 {
		return nil, nil
	}
	return g.bank.earn(amount)
}

func (g *game) licensePrice(coins int) (digAllowed int) {
	//nolint:gomnd // TODO Balance?
	switch {
	case coins == 0:
		return 3
	case coins <= 5:
		return 5
	case coins <= 10:
		return 10
	case coins <= 20:
		return 20 + g.prng.Intn(10)
	default:
		return 40 + g.prng.Intn(10)
	}
}
