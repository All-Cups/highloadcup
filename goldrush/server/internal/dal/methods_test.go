package dal_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/powerman/check"

	"github.com/Djarvur/allcups-itrally-2020-task/internal/dal"
)

func TestStartTime(tt *testing.T) {
	t := check.T(tt)
	cfg := dal.Config{
		ResultDir: t.TempDir(),
		WorkDir:   t.TempDir(),
	}
	r, err := dal.New(ctx, cfg)
	t.Nil(err)

	start, err := r.LoadStartTime()
	t.Nil(err)
	t.DeepEqual(start, new(time.Time))

	prev := time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC)
	t.Nil(r.SaveStartTime(prev))

	start, err = r.LoadStartTime()
	t.Nil(err)
	t.DeepEqual(start, &prev)

	t.Nil(os.Chmod(filepath.Join(cfg.WorkDir, "start.time"), 0o000))
	start, err = r.LoadStartTime()
	t.Match(err, "permission denied")
	t.Nil(start)

	t.Nil(os.Chmod(cfg.WorkDir, 0o500))
	defer os.Chmod(cfg.WorkDir, 0o700)
	err = r.SaveStartTime(prev)
	t.Match(err, "permission denied")
}

func TestTreasureKey(tt *testing.T) {
	t := check.T(tt)
	cfg := dal.Config{
		ResultDir: t.TempDir(),
		WorkDir:   t.TempDir(),
	}
	r, err := dal.New(ctx, cfg)
	t.Nil(err)

	key, err := r.LoadTreasureKey()
	t.True(os.IsNotExist(err))
	t.Nil(key)

	key1 := []byte{1, 2, 3, 30: 31, 32}
	t.Nil(r.SaveTreasureKey(key1))

	key, err = r.LoadTreasureKey()
	t.Nil(err)
	t.DeepEqual(key, key1)

	t.Nil(os.Chmod(filepath.Join(cfg.WorkDir, "treasure.key"), 0o000))
	key, err = r.LoadTreasureKey()
	t.Match(err, "permission denied")
	t.Nil(key)

	t.Nil(os.Chmod(cfg.WorkDir, 0o500))
	defer os.Chmod(cfg.WorkDir, 0o700)
	err = r.SaveTreasureKey(key1)
	t.Match(err, "permission denied")
}

func TestGame(tt *testing.T) {
	t := check.T(tt)
	cfg := dal.Config{
		ResultDir: t.TempDir(),
		WorkDir:   t.TempDir(),
	}
	r, err := dal.New(ctx, cfg)
	t.Nil(err)

	f, err := r.LoadGame()
	t.True(os.IsNotExist(err))
	t.Nil(f)

	game := bytes.NewBufferString("The game.")
	t.Nil(r.SaveGame(game))

	f, err = r.LoadGame()
	t.Nil(err)
	buf, err := ioutil.ReadAll(f)
	t.Nil(err)
	t.Equal(string(buf), "The game.")
	t.Nil(f.Close())

	t.Nil(os.Chmod(filepath.Join(cfg.WorkDir, "game.data"), 0o000))
	f, err = r.LoadGame()
	t.Match(err, "permission denied")
	t.Nil(f)

	t.Nil(os.Chmod(cfg.WorkDir, 0o500))
	defer os.Chmod(cfg.WorkDir, 0o700)
	err = r.SaveGame(game)
	t.Match(err, "permission denied")
}

func TestSaveResult(tt *testing.T) {
	t := check.T(tt)
	cfg := dal.Config{
		ResultDir: t.TempDir(),
		WorkDir:   t.TempDir(),
	}
	r, err := dal.New(ctx, cfg)
	t.Nil(err)

	t.Nil(r.SaveResult(42))
	t.Err(r.SaveResult(7), os.ErrExist)
	buf, err := ioutil.ReadFile(cfg.ResultDir + "/result.json")
	t.Nil(err)
	t.Equal(string(buf), `{"status":"OK","score":42}`)
}

func TestSaveError(tt *testing.T) {
	t := check.T(tt)
	cfg := dal.Config{
		ResultDir: t.TempDir(),
		WorkDir:   t.TempDir(),
	}
	r, err := dal.New(ctx, cfg)
	t.Nil(err)

	t.Nil(r.SaveError("Boo!"))
	t.Err(r.SaveError("Oops!"), os.ErrExist)
	buf, err := ioutil.ReadFile(cfg.ResultDir + "/result.json")
	t.Nil(err)
	t.Equal(string(buf), `{"status":"ERR","errors":["Boo!"]}`)
}
