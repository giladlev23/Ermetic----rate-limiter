package rate_limiter

import (
	"sync"
	"time"
)

type Window interface {
	Start() time.Time
	Count() int64
	AddCount()
	Set(s time.Time, c int64)
}

func NewLimiter(size time.Duration, limit int64) *Limiter {
	now := time.Now()

	lim := &Limiter{
		size:       size,
		limit:      limit,
		currWindow: NewWindow(now, 0),
		prevWindow: NewWindow(now, 0),
	}

	return lim
}

type Limiter struct {
	size  time.Duration
	limit int64

	mu sync.Mutex

	currWindow Window
	prevWindow Window
}

func (lim *Limiter) Allow() bool {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	now := time.Now()

	lim.advance(now)

	elapsed := now.Sub(lim.currWindow.Start())
	weight := float64(lim.size-elapsed) / float64(lim.size)
	count := int64(weight*float64(lim.prevWindow.Count())) + lim.currWindow.Count()

	if count+1 > lim.limit {
		return false
	}

	lim.currWindow.AddCount()
	return true
}

func (lim *Limiter) advance(now time.Time) {
	// Calculate the start boundary of the expected current-window.
	newCurrStart := now.Truncate(lim.size)

	diffSize := newCurrStart.Sub(lim.currWindow.Start()) / lim.size
	if diffSize >= 1 {
		// The current-window is at least one-window-size behind the expected one.

		newPrevCount := int64(0)
		if diffSize == 1 {
			// The new previous-window will overlap with the old current-window,
			// so it inherits the count.
			//
			// Note that the count here may be not accurate, since it is only a
			// SNAPSHOT of the current-window's count, which in itself tends to
			// be inaccurate due to the asynchronous nature of the sync behaviour.
			newPrevCount = lim.currWindow.Count()
		}
		lim.prevWindow.Set(newCurrStart.Add(-lim.size), newPrevCount)

		// The new current-window always has zero count.
		lim.currWindow.Set(newCurrStart, 0)
	}
}
