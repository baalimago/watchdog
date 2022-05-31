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

// New Watchdog which will execute the cb after d duration, unless kicked
func New(d time.Duration, cb func()) *Watchdog {
	if d <= 0 {
		panic(errors.New("non-positive interval for New"))
	}
	return &Watchdog{
		time.AfterFunc(d, cb),
		d,
	}
}

// Kick the watchdog to reset the timer duration. Will return an error if attempting to kick an inactive
// watchdog, i.e., a watchdog which has already executed it's cb or has been stopped.
func (w *Watchdog) Kick() error {
	ok := w.Reset(w.d)
	if !ok {
		return errors.New("can't kick inactive watchdog")
	}
	return nil
}
