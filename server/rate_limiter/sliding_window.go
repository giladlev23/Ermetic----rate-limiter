package rate_limiter

import "time"

type SlidingWindow struct {
	start time.Time

	count int64
}

func (w *SlidingWindow) Start() time.Time {
	return w.start
}

func (w *SlidingWindow) Count() int64 {
	return w.count
}

func (w *SlidingWindow) AddCount() {
	w.count += 1
}

func (w *SlidingWindow) Set(s time.Time, c int64) {
	w.start = s
	w.count = c
}

func NewWindow(start time.Time, count int64) *SlidingWindow {
	return &SlidingWindow{start: start,
		count: count}
}
