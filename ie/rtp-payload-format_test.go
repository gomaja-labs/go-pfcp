package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestRTPPayloadFormat checks TS 29.244 V20.1.0 clause 8.2.237.
func TestRTPPayloadFormat(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n)
		valid := v == 1 || v == 2
		i := ie.NewRTPPayloadFormat(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.RTPPayloadFormat)
			continue
		}
		payload := make([]byte, 1)
		payload[0] = v
		checkNewIEWire(t, i, payload)
		got, err := i.RTPPayloadFormat()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.RTPPayloadFormat() }
	checkNewIEFailures(t, ie.RTPPayloadFormat, getter, 1)
	base := ie.NewRTPPayloadFormat(1)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.RTPPayloadFormat, extension))
	if (err != nil) != true {
		t.Fatalf("extension: %v", err)
	}
}

func TestRTPPayloadFormatInvalidWire(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n) & 255
		valid := v == 1 || v == 2
		if valid {
			continue
		}
		_, err := ie.New(ie.RTPPayloadFormat, []byte{byte(n)}).RTPPayloadFormat()
		if err == nil {
			t.Fatalf("invalid wire value %d accepted", n)
		}
	}
}
