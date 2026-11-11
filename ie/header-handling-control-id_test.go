package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestHeaderHandlingControlID checks TS 29.244 V20.1.0 clause 8.2.259.
func TestHeaderHandlingControlID(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n)
		valid := true
		i := ie.NewHeaderHandlingControlID(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.HeaderHandlingControlID)
			continue
		}
		payload := make([]byte, 1)
		payload[0] = v
		checkNewIEWire(t, i, payload)
		got, err := i.HeaderHandlingControlID()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.HeaderHandlingControlID() }
	checkNewIEFailures(t, ie.HeaderHandlingControlID, getter, 1)
	base := ie.NewHeaderHandlingControlID(0)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.HeaderHandlingControlID, extension))
	if (err != nil) != true {
		t.Fatalf("extension: %v", err)
	}
}
