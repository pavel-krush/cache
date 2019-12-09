package cache

import "time"

type Clock interface {
	Now() time.Time
	Reset()
}

type ClockNone struct {}
func (cn ClockNone) Now() time.Time {
	return time.Time{}
}
func (cn ClockNone) Reset() {}

type ClockSimple struct {}
func (cs ClockSimple) Now() time.Time {
	return time.Now()
}

func (cs ClockSimple) Reset() {}
