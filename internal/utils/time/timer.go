package time

import "time"

type Time interface {
	Now() time.Time
	NewTicker(d time.Duration) *time.Ticker
}

type timeImpl struct{}

func NewTime() Time {
	return &timeImpl{}
}

func (t *timeImpl) Now() time.Time {
	return time.Now()
}

func (t *timeImpl) NewTicker(d time.Duration) *time.Ticker {
	return time.NewTicker(d)
}
