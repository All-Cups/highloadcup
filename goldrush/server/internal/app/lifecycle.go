package app

import (
	"os"
	"time"

	"github.com/powerman/structlog"
)

func (a *App) Wait(ctx Ctx) (err error) {
	log := structlog.FromContext(ctx, nil)
	select {
	case <-ctx.Done():
	case <-time.After(a.cfg.StartTimeout):
		err = a.repo.SaveError("waiting for first API request: timeout")
		if err != nil {
			log.PrintErr("SaveError", "err", err)
		}
		log.Info("task failed to start timely")
	case t := <-a.started:
		dur := time.Until(t.Add(a.cfg.Duration))
		log.Info("task started", "dur", dur)
		errc := make(chan error)
		go a.autosave(ctx, errc)
		select {
		case <-ctx.Done():
		case err = <-errc:
			log.PrintErr("autosave", "err", err)
		case <-time.After(dur):
			balance, _ := a.game.Balance()
			err = a.repo.SaveResult(balance)
			switch {
			case os.IsExist(err):
				log.Warn("SaveResult", "err", err)
				err = nil
			case err != nil:
				log.PrintErr("SaveResult", "err", err)
			}
			log.Info("task finished")
		}
	}
	return err
}

func (a *App) Start(t time.Time) (err error) {
	a.startOnce.Do(func() {
		a.started <- t
		err = a.repo.SaveStartTime(t)
	})
	return
}

func (a *App) autosave(ctx Ctx, errc chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(a.cfg.AutosavePeriod):
			t := time.Now()
			err := a.repo.SaveGame(a.game)
			if err != nil {
				errc <- err
				return
			}
			metric.autosaveDuration.Observe(time.Since(t).Seconds())
		}
	}
}
