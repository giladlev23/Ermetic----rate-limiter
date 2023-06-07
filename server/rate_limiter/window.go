package rate_limiter

import "time"

func NewWindow(start time.Time, count int64, prevCount int64) *Window {
	return &Window{start: start,
		count:     count,
		prevCount: prevCount}
}

type Window struct {
	start time.Time

	count     int64
	prevCount int64
}

func (w *Window) Start() time.Time {
	return w.start
}

func (w *Window) PrevCount() int64 {
	return w.prevCount
}

func (w *Window) Count() int64 {
	return w.count
}

func (w *Window) AddCount() {
	w.count += 1
}

func (w *Window) Set(start time.Time, count int64, prevCount int64) {
	w.start = start
	w.count = count
	w.prevCount = prevCount
}
