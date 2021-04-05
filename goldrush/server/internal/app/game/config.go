package game

import "fmt"

func (cfg Config) treasures() int {
	return cfg.area() * int(cfg.Depth) / cfg.Density
}

func (cfg Config) area() int {
	return cfg.SizeX * cfg.SizeY
}

func (cfg Config) volume() int {
	return cfg.area() * int(cfg.Depth)
}

// TotalCash returns amount of coins required to cash all treasures in
// worst case (all treasures cost as much as possible).
func (cfg Config) totalCash() (total int) {
	treasures := cfg.treasures()
	for depth := cfg.Depth; treasures > 0 && depth > 0; depth-- {
		var max int
		if treasures <= cfg.area() {
			max = treasures
		} else {
			max = cfg.area()
		}
		_, cost := cfg.treasureValueAt(depth)
		total += max * cost
		treasures -= max
	}
	if treasures > 0 {
		panic("not all treasures were counted")
	}
	return total
}

func (cfg Config) treasureValueAt(depth uint8) (min, max int) {
	avg := (*cfg.TreasureValue)[depth-1]
	//nolint:gomnd // Balance.
	switch cfg.TreasureValueAlg {
	case AlgDoubleMax:
		return avg, avg * 2
	case AlgQuarterAround:
		return avg - avg/4, avg + avg/4
	default:
		panic(fmt.Sprintf("unknown TreasureValueAlg: %d", cfg.TreasureValueAlg))
	}
}
