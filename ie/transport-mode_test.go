package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.227: every defined value and every received spare-bit pattern.
func TestTransportModeCurrent(t *testing.T) {
	for v := 0; v <= 2; v++ {
		assertAddedVector(t, ie.NewTransportMode(ie.TransportModeValue(v)), []byte{byte(v)}, "TransportMode", false)
	}
	for v := 0; v < 256; v++ {
		got, err := ie.New(ie.TransportMode, []byte{byte(v)}).TransportMode()
		if err != nil || uint8(got) != uint8(v&15) {
			t.Fatalf("received %d: %v, %v", v, got, err)
		}
	}
	for _, v := range []int{3, 255} {
		if v >= 0 && v <= 2 {
			continue
		}
		assertAddedInvalid(t, ie.NewTransportMode(ie.TransportModeValue(v)))
	}
}

var _ = registerConstructorFixture(ie.NewTransportMode, ie.TransportModeValue(1))
