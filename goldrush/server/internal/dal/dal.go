// Package dal implements Data Access Layer using in-memory DB.
package dal

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

type Config struct {
	WorkDir   string
	ResultDir string
}

// Repo provides access to storage.
type Repo struct {
	cfg Config
}

// New creates and returns new Repo.
func New(_ Ctx, cfg Config) (*Repo, error) {
	return &Repo{cfg: cfg}, nil
}

func save(dir, filename string, buf []byte) error {
	path := filepath.Join(dir, filename)
	err := ioutil.WriteFile(path+".tmp", buf, 0o600)
	if err == nil {
		err = os.Rename(path+".tmp", path)
	}
	return err
}

func saveFrom(dir, filename string, from io.WriterTo) error {
	path := filepath.Join(dir, filename)
	f, err := os.OpenFile(path+".tmp", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600) //nolint:gosec // False positive.
	if err == nil {
		_, err = from.WriteTo(f)
	}
	if err == nil {
		err = f.Close()
	}
	if err == nil {
		err = os.Rename(path+".tmp", path)
	}
	return err
}
