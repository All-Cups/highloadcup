package game

type rr struct {
	step    int // 1 or -1
	length  int // 1..
	current int // 0..length-1
	count   int // 0..length
}

func newRR(start int, forward bool, length int) *rr {
	if start < 0 || start >= length || length <= 0 {
		panic("never here")
	}
	step := 1
	if !forward {
		step = -1
	}
	return &rr{
		step:    step,
		length:  length,
		current: start,
		count:   0,
	}
}

// Next returns next value (moving by one from start value) until it reach
// start value again (then it panics).
func (r *rr) next() int {
	r.current += r.step
	if r.current >= r.length {
		r.current = 0
	} else if r.current < 0 {
		r.current = r.length - 1
	}
	r.count++
	if r.count > r.length {
		panic("never here")
	}
	return r.current
}
