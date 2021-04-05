package app_test

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/powerman/check"

	"github.com/Djarvur/allcups-itrally-2020-task/pkg/def"
)

func TestWait_TimeoutBeforeStart(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	cleanup, a, mockRepo, _ := testNew(t)
	defer cleanup()

	mockRepo.EXPECT().SaveError("waiting for first API request: timeout").Return(nil)
	errc := make(chan error, 1)
	go func() { errc <- a.Wait(ctx) }()
	waitErr(t, errc, cfg.StartTimeout, nil)
}

func TestWait_SaveErrorErr(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	cleanup, a, mockRepo, _ := testNew(t)
	defer cleanup()

	mockRepo.EXPECT().SaveError("waiting for first API request: timeout").Return(io.EOF)
	errc := make(chan error, 1)
	go func() { errc <- a.Wait(ctx) }()
	waitErr(t, errc, cfg.StartTimeout, io.EOF)
}

func TestWait_ShutdownBeforeStart(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	cleanup, a, _, _ := testNew(t)
	defer cleanup()

	ctx, shutdown := context.WithCancel(ctx)
	errc := make(chan error, 1)
	go func() { errc <- a.Wait(ctx) }()
	go func() { time.Sleep(def.TestSecond / 10); shutdown() }()
	waitErr(t, errc, def.TestSecond/10, nil)
}

func TestWait_ShutdownAfterStart(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	cleanup, a, mockRepo, _ := testNew(t)
	defer cleanup()

	ctx, shutdown := context.WithCancel(ctx)
	errc := make(chan error, 1)
	go func() { errc <- a.Wait(ctx) }()
	mockRepo.EXPECT().SaveStartTime(gomock.Any()).Return(nil)
	a.Start(time.Now())
	go func() { time.Sleep(def.TestSecond / 10); shutdown() }()
	waitErr(t, errc, def.TestSecond/10, nil)
}

func TestWait_AutosaveErr(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	cleanup, a, mockRepo, mockGame := testNew(t)
	defer cleanup()

	errc := make(chan error, 1)
	go func() { errc <- a.Wait(ctx) }()
	mockRepo.EXPECT().SaveStartTime(gomock.Any()).Return(nil)
	a.Start(time.Now())
	mockRepo.EXPECT().SaveGame(mockGame).Return(io.EOF)
	waitErr(t, errc, cfg.AutosavePeriod, io.EOF)
}

func TestWait_SaveResultErr(t *testing.T) {
	tests := []struct {
		err     error
		wantErr error
	}{
		{io.EOF, io.EOF},
		{os.ErrExist, nil},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			t.Parallel()
			cleanup, a, mockRepo, mockGame := testNew(t)
			defer cleanup()

			mockRepo.EXPECT().SaveGame(mockGame).Return(nil).AnyTimes()
			mockGame.EXPECT().Balance().Return(0, nil)
			mockRepo.EXPECT().SaveResult(0).Return(tc.err)

			errc := make(chan error, 1)
			go func() { errc <- a.Wait(ctx) }()

			now := time.Now().Add(def.TestSecond - cfg.Duration)
			mockRepo.EXPECT().SaveStartTime(now).Return(nil)
			a.Start(now)
			waitErr(t, errc, def.TestSecond, tc.wantErr)
		})
	}
}

func TestWait(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	cleanup, a, mockRepo, mockGame := testNew(t)
	defer cleanup()

	mockRepo.EXPECT().SaveGame(mockGame).Return(nil).AnyTimes()
	mockGame.EXPECT().Balance().Return(2, []int{0, 1})

	ctx, shutdown := context.WithCancel(ctx)
	defer shutdown()
	errc := make(chan error, 1)
	go func() { errc <- a.Wait(ctx) }()
	// Waiting for a.Start().
	select {
	case <-errc:
		t.FailNow()
	case <-time.After(def.TestSecond + def.TestSecond/4):
	}
	// Finish in cfg.Duration after first a.Start().
	// Second Start() will be ignored.
	now := time.Now().Add(def.TestSecond - cfg.Duration)
	mockRepo.EXPECT().SaveStartTime(now).Return(nil)
	a.Start(now)
	time.Sleep(def.TestSecond / 2)
	a.Start(now.Add(def.TestSecond / 2))
	mockRepo.EXPECT().SaveResult(2).Return(nil)
	waitErr(t, errc, def.TestSecond/2, nil)
}
