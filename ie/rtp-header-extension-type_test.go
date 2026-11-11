package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.234: every defined value and every received spare-bit pattern.
func TestRTPHeaderExtensionTypeCurrent(t *testing.T) {
	for v := 1; v <= 3; v++ {
		assertAddedVector(t, ie.NewRTPHeaderExtensionType(ie.RTPHeaderExtensionTypeValue(v)), []byte{byte(v)}, "RTPHeaderExtensionType", true)
	}
	for v := 0; v < 256; v++ {
		got, err := ie.New(ie.RTPHeaderExtensionType, []byte{byte(v)}).RTPHeaderExtensionType()
		if err != nil || uint8(got) != uint8(v&255) {
			t.Fatalf("received %d: %v, %v", v, got, err)
		}
	}
	for _, v := range []int{0, 255} {
		if v >= 1 && v <= 3 {
			continue
		}
		assertAddedInvalid(t, ie.NewRTPHeaderExtensionType(ie.RTPHeaderExtensionTypeValue(v)))
	}
}

var _ = registerConstructorFixture(ie.NewRTPHeaderExtensionType, ie.RTPHeaderExtensionTypeValue(1))
