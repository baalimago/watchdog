package watchdog_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/baalimago/watchdog"
)

func Test_noPet(t *testing.T) {
	t.Run("it should call cb after duration", func(t0 *testing.T) {
		timeout := time.Millisecond * 10
		testTimeoutCtx, cancel := context.WithTimeout(context.Background(), timeout*2)
		defer cancel()

		cbCalled := make(chan struct{})

		watchdog.New(timeout, func() {
			close(cbCalled)
		})

		select {
		case <-testTimeoutCtx.Done():
			t0.Error("failed to call cb in time")
		case <-cbCalled:
		}
	})
}

func Test_Pet(t *testing.T) {
	t.Run("it shouldn't call cb if peted", func(t0 *testing.T) {
		timeout := time.Millisecond * 10
		testTimeoutCtx, cancel := context.WithTimeout(context.Background(), timeout*2)
		defer cancel()
		var w *watchdog.Watchdog

		cbCalled := make(chan struct{})

		awaitPetRoutine := sync.WaitGroup{}
		awaitPetRoutine.Add(1)
		// Anonymous routine which pets the watchdog 4 times more often than timeout duration
		go func() {
			awaitPetRoutine.Done()
			t := time.NewTicker(timeout / 4)
			for {
				select {
				case <-testTimeoutCtx.Done():
					// Non related, just to cleanup
					t.Stop()
				case <-t.C:
					err := w.Pet()
					if err != nil {
						t0.Errorf("failed to pet: %v", err)
					}
				}
			}
		}()

		awaitPetRoutine.Wait()
		w = watchdog.New(timeout, func() {
			close(cbCalled)
		})

		select {
		case <-cbCalled:
			t0.Error("failed to reset the Done timer")
		case <-testTimeoutCtx.Done():
		}
	})
}

func Test_Stop(t *testing.T) {
	t.Run("it shouldn't call cb when Stop()", func(t0 *testing.T) {
		timeout := time.Millisecond * 10
		testTimeoutCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		var w *watchdog.Watchdog

		cbCalled := make(chan struct{})

		awaitStopRoutine := sync.WaitGroup{}
		awaitWatchdogCreateRoutine := sync.WaitGroup{}
		awaitStopRoutine.Add(1)
		awaitWatchdogCreateRoutine.Add(1)
		go func() {
			awaitStopRoutine.Done()
			awaitWatchdogCreateRoutine.Wait()
			w.Stop()
		}()

		awaitStopRoutine.Wait()
		w = watchdog.New(timeout/2, func() {
			close(cbCalled)
		})
		awaitWatchdogCreateRoutine.Done()

		select {
		case <-cbCalled:
			t0.Error("failed to stop watchdog cb firing")
		case <-testTimeoutCtx.Done():
		}
	})
}

func Test_PetWithUpdate(t *testing.T) {
	t.Run("it should update duration and not call cb if peted", func(t0 *testing.T) {
		initialTimeout := time.Millisecond * 10
		updatedTimeout := time.Millisecond * 20
		testTimeoutCtx, cancel := context.WithTimeout(context.Background(), updatedTimeout*2)
		defer cancel()
		var w *watchdog.Watchdog

		cbCalled := make(chan struct{})

		awaitPetRoutine := sync.WaitGroup{}
		awaitPetRoutine.Add(1)
		// Anonymous routine which pets the watchdog 4 times more often than updatedTimeout duration
		go func() {
			awaitPetRoutine.Done()
			t := time.NewTicker(updatedTimeout / 4)
			for {
				select {
				case <-testTimeoutCtx.Done():
					// Non related, just to cleanup
					t.Stop()
				case <-t.C:
					err := w.PetWithUpdate(updatedTimeout)
					if err != nil {
						t0.Errorf("failed to pet: %v", err)
					}
				}
			}
		}()

		awaitPetRoutine.Wait()
		w = watchdog.New(initialTimeout, func() {
			close(cbCalled)
		})

		select {
		case <-cbCalled:
			t0.Error("failed to reset the Done timer")
		case <-testTimeoutCtx.Done():
		}
	})
}
