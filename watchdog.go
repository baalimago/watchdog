package watchdog

import (
	"errors"
	"time"
)

// Watchdog which will execute a callback after d duration. Initiate using watchdog.New
type Watchdog struct {
	*time.Timer
	d time.Duration
}

// New Watchdog which will execute the cb after d duration, unless peted
func New(d time.Duration, cb func()) *Watchdog {
	if d <= 0 {
		panic(errors.New("non-positive interval for New"))
	}
	return &Watchdog{
		time.AfterFunc(d, cb),
		d,
	}
}

// Pet the watchdog to reset the timer duration. Will return an error if attempting to pet an inactive
// watchdog, i.e., a watchdog which has already executed it's cb or has been stopped.
func (w *Watchdog) Pet() error {
	ok := w.Reset(w.d)
	if !ok {
		return errors.New("can't pet inactive watchdog")
	}
	return nil
}

// PetWithUpdate resets the watchdog timer to a new duration and restarts the countdown.
// Will return an error if attempting to pet an inactive watchdog,
// i.e., a watchdog which has already executed its callback or has been stopped.
// `newDuration` specifies the new time duration to use for the watchdog.
func (w *Watchdog) PetWithUpdate(newDuration time.Duration) error {
	w.d = newDuration
	return w.Pet()
}
