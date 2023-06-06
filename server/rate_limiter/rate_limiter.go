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

	durationSinceCurrWindowStart := now.Sub(lim.currWindow.Start())
	prevWindowPart := float64(lim.size - durationSinceCurrWindowStart)
	prevWindowWeight := prevWindowPart / float64(lim.size)
	prevWindowWeightedCount := int64(prevWindowWeight * float64(lim.prevWindow.Count()))
	count := prevWindowWeightedCount + lim.currWindow.Count()

	if count+1 > lim.limit {
		return false
	}

	lim.currWindow.AddCount()
	return true
}

func (lim *Limiter) advance(now time.Time) {
	newCurrStart := now.Truncate(lim.size)
	timeSinceLastWindow := newCurrStart.Sub(lim.currWindow.Start())
	diff := timeSinceLastWindow - lim.size

	if diff >= 0 {
		newPrevCount := int64(0)
		if diff == 0 {
			// Exactly overlapping windows
			newPrevCount = lim.currWindow.Count()
		}

		lim.prevWindow.Set(newCurrStart.Add(-lim.size), newPrevCount)
		lim.currWindow.Set(newCurrStart, 0)
	}
}
