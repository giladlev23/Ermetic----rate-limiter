package rate_limiter

import "time"

func newWindow(start time.Time, count int64, prevCount int64) *window {
	return &window{start: start,
		count:     count,
		prevCount: prevCount}
}

type window struct {
	start time.Time

	count     int64
	prevCount int64
}

func (w *window) Start() time.Time {
	return w.start
}

func (w *window) PrevCount() int64 {
	return w.prevCount
}

func (w *window) Count() int64 {
	return w.count
}

func (w *window) AddCount() {
	w.count += 1
}

func (w *window) Set(start time.Time, count int64, prevCount int64) {
	w.start = start
	w.count = count
	w.prevCount = prevCount
}
