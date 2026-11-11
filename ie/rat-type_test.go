package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.186: every defined value and every received spare-bit pattern.
func TestRATTypeCurrent(t *testing.T) {
	for v := 1; v <= 28; v++ {
		assertAddedVector(t, ie.NewRATType(ie.RATTypeValue(v)), []byte{byte(v)}, "RATType", false)
	}
	for v := 0; v < 256; v++ {
		got, err := ie.New(ie.RATType, []byte{byte(v)}).RATType()
		if err != nil || uint8(got) != uint8(v&255) {
			t.Fatalf("received %d: %v, %v", v, got, err)
		}
	}
	for _, v := range []int{0, 255} {
		if v >= 1 && v <= 28 {
			continue
		}
		assertAddedInvalid(t, ie.NewRATType(ie.RATTypeValue(v)))
	}
}

var _ = registerConstructorFixture(ie.NewRATType, ie.RATTypeValue(1))
