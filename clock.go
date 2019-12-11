package cache

import "time"

type Clock interface {
	Now() time.Time
	Reset()
}

type _ClockNone struct{}

func (cn _ClockNone) Now() time.Time {
	return time.Time{}
}
func (cn _ClockNone) Reset() {}

var ClockNone = _ClockNone{}

type _ClockSimple struct{}

func (cs _ClockSimple) Now() time.Time {
	return time.Now()
}
func (cs _ClockSimple) Reset() {}

var ClockSimple = _ClockSimple{}
