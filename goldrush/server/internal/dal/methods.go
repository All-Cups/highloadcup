package dal

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/Djarvur/allcups-itrally-2020-task/internal/app"
)

const (
	fnStartTime   = "start.time"
	fnTreasureKey = "treasure.key"
	fnGame        = "game.data"
	fnResult      = "result.json"
)

func (r *Repo) LoadStartTime() (*time.Time, error) {
	t := new(time.Time)
	path := filepath.Join(r.cfg.WorkDir, fnStartTime)
	buf, err := ioutil.ReadFile(path) //nolint:gosec // False positive.
	switch {
	case errors.Is(err, os.ErrNotExist):
		err = nil
	case err == nil:
		err = t.UnmarshalText(buf)
	}
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *Repo) SaveStartTime(t time.Time) error {
	buf, err := t.MarshalText()
	if err != nil {
		return err
	}
	return save(r.cfg.WorkDir, fnStartTime, buf)
}

func (r *Repo) LoadTreasureKey() ([]byte, error) {
	path := filepath.Join(r.cfg.WorkDir, fnTreasureKey)
	return ioutil.ReadFile(path) //nolint:gosec // False positive.
}

func (r *Repo) SaveTreasureKey(buf []byte) error {
	return save(r.cfg.WorkDir, fnTreasureKey, buf)
}

func (r *Repo) LoadGame() (app.ReadSeekCloser, error) {
	path := filepath.Join(r.cfg.WorkDir, fnGame)
	return os.Open(path) //nolint:gosec // False positive.
}

func (r *Repo) SaveGame(from io.WriterTo) error {
	return saveFrom(r.cfg.WorkDir, fnGame, from)
}

func (r *Repo) SaveResult(result int) error {
	path := filepath.Join(r.cfg.ResultDir, fnResult)
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return os.ErrExist
	}
	data := struct {
		Status string `json:"status"`
		Score  int    `json:"score"`
	}{
		Status: "OK",
		Score:  result,
	}
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return save(r.cfg.ResultDir, fnResult, buf)
}

func (r *Repo) SaveError(msg string) error {
	path := filepath.Join(r.cfg.ResultDir, fnResult)
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return os.ErrExist
	}
	data := struct {
		Status string   `json:"status"`
		Errors []string `json:"errors"`
	}{
		Status: "ERR",
		Errors: []string{msg},
	}
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return save(r.cfg.ResultDir, fnResult, buf)
}
