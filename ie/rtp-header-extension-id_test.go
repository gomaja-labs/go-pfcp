package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.235: every defined value and every received spare-bit pattern.
func TestRTPHeaderExtensionIDCurrent(t *testing.T) {
	for v := 1; v <= 255; v++ {
		assertAddedVector(t, ie.NewRTPHeaderExtensionID(uint8(v)), []byte{byte(v)}, "RTPHeaderExtensionID", true)
	}
	for v := 0; v < 256; v++ {
		got, err := ie.New(ie.RTPHeaderExtensionID, []byte{byte(v)}).RTPHeaderExtensionID()
		if v == 0 {
			if err == nil {
				t.Fatal("prohibited value accepted")
			}
			continue
		}
		if err != nil || got != uint8(v&255) {
			t.Fatalf("received %d: %v, %v", v, got, err)
		}
	}
	for _, v := range []int{0, 255} {
		if v >= 1 && v <= 255 {
			continue
		}
		assertAddedInvalid(t, ie.NewRTPHeaderExtensionID(uint8(v)))
	}
}

var _ = registerConstructorFixture(ie.NewRTPHeaderExtensionID, uint8(1))
