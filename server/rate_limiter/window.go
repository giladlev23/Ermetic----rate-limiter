package rate_limiter

import "time"

type Window struct {
	start time.Time

	count int64
}

func (w *Window) Start() time.Time {
	return w.start
}

func (w *Window) Count() int64 {
	return w.count
}

func (w *Window) AddCount() {
	w.count += 1
}

func (w *Window) Set(s time.Time, c int64) {
	w.start = s
	w.count = c
}

func NewWindow(start time.Time, count int64) *Window {
	return &Window{start: start,
		count: count}
}
