package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.192: every defined value and every received spare-bit pattern.
func TestL2TPSessionIndicationsCurrent(t *testing.T) {
	for v := 0; v <= 7; v++ {
		assertAddedVector(t, ie.NewL2TPSessionIndications(ie.L2TPSessionIndicationFlags(v)), []byte{byte(v)}, "L2TPSessionIndications", false)
	}
	for v := 0; v < 256; v++ {
		got, err := ie.New(ie.L2TPSessionIndications, []byte{byte(v)}).L2TPSessionIndications()
		if err != nil || uint8(got) != uint8(v&7) {
			t.Fatalf("received %d: %v, %v", v, got, err)
		}
	}
	for _, v := range []int{8, 255} {
		if v >= 0 && v <= 7 {
			continue
		}
		assertAddedInvalid(t, ie.NewL2TPSessionIndications(ie.L2TPSessionIndicationFlags(v)))
	}
}
