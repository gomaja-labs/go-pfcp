package ie_test

import (
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestSessionReflectorModeSTAMP checks TS 29.244 V20.1.0 clause 8.2.279.
func TestSessionReflectorModeSTAMP(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n)
		valid := v <= 1
		i := ie.NewSessionReflectorModeSTAMP(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.SessionReflectorModeSTAMP)
			continue
		}
		payload := make([]byte, 1)
		payload[0] = v
		checkNewIEWire(t, i, payload)
		got, err := i.SessionReflectorModeSTAMP()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.SessionReflectorModeSTAMP() }
	checkNewIEFailures(t, ie.SessionReflectorModeSTAMP, getter, 1)
	base := ie.NewSessionReflectorModeSTAMP(0)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.SessionReflectorModeSTAMP, extension))
	if (err != nil) != true {
		t.Fatalf("extension: %v", err)
	}
	if 1 != 255 {
		p := append([]byte{}, base.Payload...)
		p[0] |= byte(254)
		a, err := getter(ie.New(ie.SessionReflectorModeSTAMP, p))
		b, baseErr := getter(base)
		if err != nil || baseErr != nil || !reflect.DeepEqual(a, b) {
			t.Fatal("spare bits evaluated", err)
		}
	}
}

func TestSessionReflectorModeSTAMPInvalidWire(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n) & 1
		valid := v <= 1
		if valid {
			continue
		}
		_, err := ie.New(ie.SessionReflectorModeSTAMP, []byte{byte(n)}).SessionReflectorModeSTAMP()
		if err == nil {
			t.Fatalf("invalid wire value %d accepted", n)
		}
	}
}
