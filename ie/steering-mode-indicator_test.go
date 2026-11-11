package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.197: every defined value and every received spare-bit pattern.
func TestSteeringModeIndicatorCurrent(t *testing.T) {
	for v := 0; v <= 2; v++ {
		assertAddedVector(t, ie.NewSteeringModeIndicator(ie.SteeringModeFlags(v)), []byte{byte(v)}, "SteeringModeIndicator", false)
	}
	for v := 0; v < 256; v++ {
		got, err := ie.New(ie.SteeringModeIndicator, []byte{byte(v)}).SteeringModeIndicator()
		if v&3 == 3 {
			if err == nil {
				t.Fatal("prohibited value accepted")
			}
			continue
		}
		if err != nil || uint8(got) != uint8(v&3) {
			t.Fatalf("received %d: %v, %v", v, got, err)
		}
	}
	for _, v := range []int{3, 255} {
		if v >= 0 && v <= 2 {
			continue
		}
		assertAddedInvalid(t, ie.NewSteeringModeIndicator(ie.SteeringModeFlags(v)))
	}
}

var _ = registerConstructorFixture(ie.NewSteeringModeIndicator, ie.SteeringModeFlags(1))
