package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.233: every defined value and every received spare-bit pattern.
func TestMediaTransportProtocolCurrent(t *testing.T) {
	for v := 0; v <= 3; v++ {
		assertAddedVector(t, ie.NewMediaTransportProtocol(ie.MediaTransportProtocolValue(v)), []byte{byte(v)}, "MediaTransportProtocol", false)
	}
	for v := 0; v < 256; v++ {
		got, err := ie.New(ie.MediaTransportProtocol, []byte{byte(v)}).MediaTransportProtocol()
		if err != nil || uint8(got) != uint8(v&15) {
			t.Fatalf("received %d: %v, %v", v, got, err)
		}
	}
	for _, v := range []int{4, 255} {
		if v >= 0 && v <= 3 {
			continue
		}
		assertAddedInvalid(t, ie.NewMediaTransportProtocol(ie.MediaTransportProtocolValue(v)))
	}
}

var _ = registerConstructorFixture(ie.NewMediaTransportProtocol, ie.MediaTransportProtocolValue(1))
