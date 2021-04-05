package game

import (
	"fmt"
	"sync"

	"github.com/powerman/must"
)

type field struct {
	cfg       Config
	mu        sync.RWMutex
	depth     []uint8 // 0 means not digged yet.
	treasure  []bool
	notCashed []bool
}

func newField(cfg Config) *field {
	return &field{
		cfg:       cfg,
		depth:     make([]uint8, cfg.area()),
		treasure:  make([]bool, cfg.volume()),
		notCashed: make([]bool, cfg.volume()),
	}
}

func (f *field) init() {
	for x := 0; x < f.cfg.SizeX; x++ {
		for y := 0; y < f.cfg.SizeY; y++ {
			areaOffset, err := f.areaOffset(x, y)
			must.NoErr(err)
			for depth := uint8(1); depth <= f.depth[areaOffset]; depth++ {
				offset, err := f.offset(Coord{X: x, Y: y, Depth: depth})
				must.NoErr(err)
				f.treasure[offset] = false
				f.notCashed[offset] = false
			}
		}
	}
}

func (f *field) areaOffset(x, y int) (int, error) {
	if x < 0 || x >= f.cfg.SizeX {
		return 0, fmt.Errorf("%w: x=%d >= %d", ErrWrongCoord, x, f.cfg.SizeX)
	} else if y < 0 || y >= f.cfg.SizeY {
		return 0, fmt.Errorf("%w: y=%d >= %d", ErrWrongCoord, y, f.cfg.SizeY)
	}
	return x + f.cfg.SizeX*y, nil
}

func (f *field) offset(pos Coord) (int, error) {
	if pos.Depth < 1 || pos.Depth > f.cfg.Depth {
		return 0, fmt.Errorf("%w: depth=%d > %d", ErrWrongCoord, pos.Depth, f.cfg.Depth)
	}
	offset, err := f.areaOffset(pos.X, pos.Y)
	if err != nil {
		return 0, err
	}
	return offset + f.cfg.SizeX*f.cfg.SizeY*int(pos.Depth-1), nil
}

// AddTreasure is not safe for concurrent use by multiple goroutines.
func (f *field) addTreasure(pos Coord) (ok bool) {
	offset, err := f.offset(pos)
	if err != nil {
		panic(err)
	}
	if f.treasure[offset] {
		return false
	}
	f.treasure[offset] = true
	f.notCashed[offset] = true
	return true
}

func (f *field) countTreasures(area Area, depth uint8) (int, error) { //nolint:gocognit // False positive.
	f.mu.RLock()
	defer f.mu.RUnlock()

	lastX, lastY := area.X+area.SizeX-1, area.Y+area.SizeY-1
	switch {
	case area.X < 0 || area.SizeX < 1 || lastX < area.X || lastX >= f.cfg.SizeX:
		return 0, fmt.Errorf("%w: X %d-%d is outside the allowed range 0-%d", ErrWrongCoord, area.X, lastX, f.cfg.SizeX-1)
	case area.Y < 0 || area.SizeY < 1 || lastY < area.Y || lastY >= f.cfg.SizeY:
		return 0, fmt.Errorf("%w: Y %d-%d is outside the allowed range 0-%d", ErrWrongCoord, area.Y, lastY, f.cfg.SizeY-1)
	case depth < 1 || depth > f.cfg.Depth:
		return 0, fmt.Errorf("%w: depth %d is outside the allowed range 1-%d", ErrWrongDepth, depth, f.cfg.Depth)
	}

	found := 0
	for x := area.X; x <= lastX; x++ {
		for y := area.Y; y <= lastY; y++ {
			offset, err := f.offset(Coord{X: x, Y: y, Depth: depth})
			if err != nil {
				panic(err)
			}
			if f.treasure[offset] {
				found++
			}
		}
	}
	return found, nil
}

func (f *field) dig(pos Coord) (found bool, _ error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	offset, err := f.offset(pos)
	if err != nil {
		return false, err
	}
	areaOffset, err := f.areaOffset(pos.X, pos.Y)
	if err != nil {
		panic(err)
	}

	if f.depth[areaOffset] != pos.Depth-1 {
		return false, fmt.Errorf("%w: %d (should be %d)", ErrWrongDepth, pos.Depth, f.depth[areaOffset]+1)
	}
	f.depth[areaOffset]++

	if !f.treasure[offset] {
		return false, nil
	}
	f.treasure[offset] = false
	return true, nil
}

func (f *field) cash(pos Coord) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	offset, err := f.offset(pos)
	if err != nil {
		return err
	}
	if !f.notCashed[offset] {
		return ErrNoThreasure
	}
	if f.treasure[offset] {
		return ErrNotDigged
	}
	f.notCashed[offset] = false
	return nil
}
