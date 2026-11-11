package ie_test

import (
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestHeaderHandlingAction checks TS 29.244 V20.1.0 clause 8.2.255.
func TestHeaderHandlingAction(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n)
		valid := v <= 3
		i := ie.NewHeaderHandlingAction(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.HeaderHandlingAction)
			continue
		}
		payload := make([]byte, 1)
		payload[0] = v
		checkNewIEWire(t, i, payload)
		got, err := i.HeaderHandlingAction()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.HeaderHandlingAction() }
	checkNewIEFailures(t, ie.HeaderHandlingAction, getter, 1)
	base := ie.NewHeaderHandlingAction(0)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.HeaderHandlingAction, extension))
	if (err != nil) != true {
		t.Fatalf("extension: %v", err)
	}
	if 15 != 255 {
		p := append([]byte{}, base.Payload...)
		p[0] |= byte(240)
		a, err := getter(ie.New(ie.HeaderHandlingAction, p))
		b, baseErr := getter(base)
		if err != nil || baseErr != nil || !reflect.DeepEqual(a, b) {
			t.Fatal("spare bits evaluated", err)
		}
	}
}

func TestHeaderHandlingActionInvalidWire(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n) & 15
		valid := v <= 3
		if valid {
			continue
		}
		_, err := ie.New(ie.HeaderHandlingAction, []byte{byte(n)}).HeaderHandlingAction()
		if err == nil {
			t.Fatalf("invalid wire value %d accepted", n)
		}
	}
}
