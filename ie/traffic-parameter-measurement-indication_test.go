package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.223: every defined value and every received spare-bit pattern.
func TestTrafficParameterMeasurementIndicationCurrent(t *testing.T) {
	for v := 0; v <= 7; v++ {
		assertAddedVector(t, ie.NewTrafficParameterMeasurementIndication(ie.TrafficParameterMeasurementFlags(v)), []byte{byte(v)}, "TrafficParameterMeasurementIndication", false)
	}
	for v := 0; v < 256; v++ {
		got, err := ie.New(ie.TrafficParameterMeasurementIndication, []byte{byte(v)}).TrafficParameterMeasurementIndication()
		if err != nil || uint8(got) != uint8(v&7) {
			t.Fatalf("received %d: %v, %v", v, got, err)
		}
	}
	for _, v := range []int{8, 255} {
		if v >= 0 && v <= 7 {
			continue
		}
		assertAddedInvalid(t, ie.NewTrafficParameterMeasurementIndication(ie.TrafficParameterMeasurementFlags(v)))
	}
}

var _ = registerConstructorFixture(ie.NewTrafficParameterMeasurementIndication, ie.TrafficParameterMeasurementFlags(1))
