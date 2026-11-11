package ie_test

import (
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestMTSDTControlInformation checks TS 29.244 V20.1.0 clause 8.2.239.
func TestMTSDTControlInformation(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n)
		valid := v <= 1
		i := ie.NewMTSDTControlInformation(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.MTSDTControlInformation)
			continue
		}
		payload := make([]byte, 1)
		payload[0] = v
		checkNewIEWire(t, i, payload)
		got, err := i.MTSDTControlInformation()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.MTSDTControlInformation() }
	checkNewIEFailures(t, ie.MTSDTControlInformation, getter, 1)
	base := ie.NewMTSDTControlInformation(0)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.MTSDTControlInformation, extension))
	if (err != nil) != false {
		t.Fatalf("extension: %v", err)
	}
	if 1 != 255 {
		p := append([]byte{}, base.Payload...)
		p[0] |= byte(254)
		a, err := getter(ie.New(ie.MTSDTControlInformation, p))
		b, baseErr := getter(base)
		if err != nil || baseErr != nil || !reflect.DeepEqual(a, b) {
			t.Fatal("spare bits evaluated", err)
		}
	}
}

func TestMTSDTControlInformationInvalidWire(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n) & 1
		valid := v <= 1
		if valid {
			continue
		}
		_, err := ie.New(ie.MTSDTControlInformation, []byte{byte(n)}).MTSDTControlInformation()
		if err == nil {
			t.Fatalf("invalid wire value %d accepted", n)
		}
	}
}
