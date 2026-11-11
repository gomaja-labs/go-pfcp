package ie

import (
	"encoding/binary"
	"time"
)

// NewDLPeriodicity creates a DL Periodicity IE with an exact millisecond duration.
// TS 29.244 V20.1.0 clause 8.2.221.
func NewDLPeriodicity(value time.Duration) *IE {
	if !addedDurationFits(value, time.Millisecond, 0xffffffff) {
		return newInvalidIE(DLPeriodicity, ErrInvalidValue)
	}
	b := make([]byte, 5)
	binary.BigEndian.PutUint32(b[1:], uint32(value/time.Millisecond))
	return New(DLPeriodicity, b)
}

// DLPeriodicity returns the periodicity, ignoring the spare octet and extensions.
// TS 29.244 V20.1.0 clauses 8.2.221 and 8.1.2.
func (i *IE) DLPeriodicity() (time.Duration, error) {
	b, err := addedValuePayload(i, DLPeriodicity, 5, false)
	if err != nil {
		return 0, err
	}
	return time.Duration(binary.BigEndian.Uint32(b[1:5])) * time.Millisecond, nil
}
