package time

import "time"

type Time interface {
	Now() time.Time
}

type timeImpl struct{}

func NewTime() Time {
	return &timeImpl{}
}

func (t *timeImpl) Now() time.Time {
	return time.Now()
}
