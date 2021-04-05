package game

import (
	"encoding/json"
	"io"
)

type state struct {
	Cfg            Config
	ActiveLicenses []License
	Balance        int
}

// Continue implements app.GameFactory interface.
func (factory Factory) Continue(ctx Ctx, r io.ReadSeeker) (g Game, err error) {
	var previous state
	dec := json.NewDecoder(r)
	err = dec.Decode(&previous)
	if err == nil {
		g, err = factory.New(ctx, previous.Cfg)
	}
	if err == nil {
		_, err = r.Seek(dec.InputOffset()+1, io.SeekStart)
	}
	if err == nil {
		_, err = io.ReadFull(r, g.(*game).field.depth)
	}
	if err == nil && previous.Balance > 0 {
		_, err = g.(*game).bank.earn(previous.Balance)
	}
	if err != nil {
		return nil, err
	}
	g.(*game).licenses.init(previous.ActiveLicenses)
	g.(*game).field.init()
	return g, nil
}

func (g *game) WriteTo(w io.Writer) (n int64, err error) {
	g.muModify.Lock()
	current := state{
		Cfg:            g.cfg,
		ActiveLicenses: g.licenses.active(),
		Balance:        g.bank.balance,
	}
	g.muModify.Unlock()

	depth := make([]byte, len(g.field.depth))
	g.field.mu.RLock()
	copy(depth, g.field.depth)
	g.field.mu.RUnlock()

	cw := &countingWriter{w: w}
	err = json.NewEncoder(cw).Encode(current)
	if err != nil {
		return cw.n, err
	}
	_, err = cw.Write(depth)
	return cw.n, err
}

type countingWriter struct {
	w io.Writer
	n int64
}

func (c *countingWriter) Write(buf []byte) (int, error) {
	n, err := c.w.Write(buf)
	c.n += int64(n)
	return n, err
}
