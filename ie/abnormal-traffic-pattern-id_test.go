package ie_test

import (
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestAbnormalTrafficPatternID checks TS 29.244 V20.1.0 clause 8.2.285.
func TestAbnormalTrafficPatternID(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n)
		valid := v <= 3 || (v >= 32 && v <= 63)
		i := ie.NewAbnormalTrafficPatternID(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.AbnormalTrafficPatternID)
			continue
		}
		payload := make([]byte, 1)
		payload[0] = v
		checkNewIEWire(t, i, payload)
		got, err := i.AbnormalTrafficPatternID()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.AbnormalTrafficPatternID() }
	checkNewIEFailures(t, ie.AbnormalTrafficPatternID, getter, 1)
	base := ie.NewAbnormalTrafficPatternID(0)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.AbnormalTrafficPatternID, extension))
	if (err != nil) != true {
		t.Fatalf("extension: %v", err)
	}
	if 63 != 255 {
		p := append([]byte{}, base.Payload...)
		p[0] |= byte(192)
		a, err := getter(ie.New(ie.AbnormalTrafficPatternID, p))
		b, baseErr := getter(base)
		if err != nil || baseErr != nil || !reflect.DeepEqual(a, b) {
			t.Fatal("spare bits evaluated", err)
		}
	}
}

func TestAbnormalTrafficPatternIDInvalidWire(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n) & 63
		valid := v <= 3 || (v >= 32 && v <= 63)
		if valid {
			continue
		}
		_, err := ie.New(ie.AbnormalTrafficPatternID, []byte{byte(n)}).AbnormalTrafficPatternID()
		if err == nil {
			t.Fatalf("invalid wire value %d accepted", n)
		}
	}
}
