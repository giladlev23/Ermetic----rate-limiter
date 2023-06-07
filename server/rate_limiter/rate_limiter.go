package rate_limiter

import (
	"sync"
	"time"
)

type SlidingWindow interface {
	Start() time.Time
	Count() int64
	PrevCount() int64
	AddCount()
	Set(start time.Time, count int64, prevCount int64)
}

func NewLimiter(size time.Duration, limit int64) *Limiter {
	now := time.Now()

	return &Limiter{
		size:   size,
		limit:  limit,
		window: newWindow(now, 0, 0),
	}
}

type Limiter struct {
	limit int64
	size  time.Duration

	mu sync.Mutex

	window SlidingWindow
}

func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	l.advanceWindow(now)

	timeIntoCurrentWindow := now.Sub(l.window.Start())
	prevWindowPart := l.size - timeIntoCurrentWindow
	prevWindowWeight := float64(prevWindowPart) / float64(l.size)
	prevWindowWeightedCount := int64(prevWindowWeight * float64(l.window.PrevCount()))
	count := prevWindowWeightedCount + l.window.Count()

	if count+1 > l.limit {
		return false
	}

	l.window.AddCount()
	return true
}

func (l *Limiter) advanceWindow(now time.Time) {
	if now.Sub(l.window.Start()) > l.size {
		currCount := l.window.Count()
		l.window.Set(now, 0, currCount)
	}
}
