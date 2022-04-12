package xsf

import (
	"time"
)

type authFilter struct {
	valid bool

	win    time.Duration
	baseTs time.Time
}

func (a *authFilter) init(win time.Duration) *authFilter {
	a.win = win
	a.valid = true
	a.baseTs = time.Now()
	return a
}
func (a *authFilter) filter(auth int32) int32 {
	if !a.valid {
		return auth
	}
	elapsed := time.Since(a.baseTs)
	if elapsed > a.win {
		a.valid = false
		return auth
	}
	return int32(float64(time.Since(a.baseTs)) / float64(a.win) * float64(auth))
}
