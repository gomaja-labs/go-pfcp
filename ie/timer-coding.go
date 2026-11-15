package ie

import (
	"math"
	"time"
)

// InfiniteTimer denotes the infinite timer in TS 29.244 V20.1.0
// clauses 8.2.29, 8.2.35 and 8.2.78. Other duration IEs do not use this sentinel.
const InfiniteTimer = time.Duration(math.MaxInt64)

// TS 29.244 V20.1.0 tables 8.2.29.1, 8.2.35.1 and 8.2.78.1:
// choose the largest unit that represents the duration exactly in five bits.
func newTimerIE(typ IEType, duration time.Duration) *IE {
	if duration == 0 {
		return newUint8ValIE(typ, 0)
	}
	if duration == InfiniteTimer {
		return newUint8ValIE(typ, 0xe0)
	}
	units := [...]time.Duration{2 * time.Second, time.Minute, 10 * time.Minute, time.Hour, 10 * time.Hour}
	for n := len(units) - 1; n >= 0; n-- {
		if duration > 0 && durationFits(duration, units[n], 31) {
			return newUint8ValIE(typ, uint8(n)<<5|uint8(duration/units[n]))
		}
	}
	return newInvalidIE(typ, ErrInvalidValue)
}

func decodeTimerOctet(value byte) time.Duration {
	unit := time.Minute // Unassigned units are interpreted as minutes.
	switch value >> 5 {
	case 0:
		unit = 2 * time.Second
	case 2:
		unit = 10 * time.Minute
	case 3:
		unit = time.Hour
	case 4:
		unit = 10 * time.Hour
	case 7:
		return InfiniteTimer
	}
	return time.Duration(value&31) * unit
}
