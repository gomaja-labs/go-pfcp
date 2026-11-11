package ie_test

import (
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestN6JitterMeasurementCurrent(t *testing.T) {
	assertAddedFields(t, ie.NewN6JitterMeasurement(ie.N6JitterMeasurementFields{}), &ie.N6JitterMeasurementFields{}, []byte{0}, "N6JitterMeasurement", false)
	for _, lo := range []time.Duration{-0x80000000 * time.Millisecond, -time.Millisecond, 0, 0x7fffffff * time.Millisecond} {
		v := ie.N6JitterMeasurementFields{Flags: ie.N6JitterDL, DLPeriodicity: 0xffffffff * time.Millisecond, LowerDLJitter: lo, HigherDLJitter: lo}
		b := []byte{1, 255, 255, 255, 255}
		for k := 0; k < 2; k++ {
			u := uint32(int32(lo / time.Millisecond))
			b = append(b, byte(u>>24), byte(u>>16), byte(u>>8), byte(u))
		}
		assertAddedFields(t, ie.NewN6JitterMeasurement(v), &v, b, "N6JitterMeasurement", false)
	}

	var nilFields *ie.N6JitterMeasurementFields
	assertAddedNilFields(t, nilFields)
}

func TestN6JitterMeasurementInvalidFields(t *testing.T) {
	for _, v := range []ie.N6JitterMeasurementFields{{Flags: 2}, {DLPeriodicity: time.Millisecond}, {Flags: ie.N6JitterDL, DLPeriodicity: -1}, {Flags: ie.N6JitterDL, LowerDLJitter: time.Nanosecond}, {Flags: ie.N6JitterDL, LowerDLJitter: -0x80000001 * time.Millisecond}, {Flags: ie.N6JitterDL, HigherDLJitter: 0x80000000 * time.Millisecond}} {
		assertAddedInvalid(t, ie.NewN6JitterMeasurement(v))
		if _, err := v.Marshal(); err == nil {
			t.Fatal("invalid fields accepted")
		}
	}
}
