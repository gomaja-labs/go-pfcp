package ie_test

import (
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestExtendedDLBufferingNotificationPolicy checks TS 29.244 V20.1.0 clause 8.2.238.
func TestExtendedDLBufferingNotificationPolicy(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n)
		valid := v <= 1
		i := ie.NewExtendedDLBufferingNotificationPolicy(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.ExtendedDLBufferingNotificationPolicy)
			continue
		}
		payload := make([]byte, 1)
		payload[0] = v
		checkNewIEWire(t, i, payload)
		got, err := i.ExtendedDLBufferingNotificationPolicy()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.ExtendedDLBufferingNotificationPolicy() }
	checkNewIEFailures(t, ie.ExtendedDLBufferingNotificationPolicy, getter, 1)
	base := ie.NewExtendedDLBufferingNotificationPolicy(0)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.ExtendedDLBufferingNotificationPolicy, extension))
	if (err != nil) != false {
		t.Fatalf("extension: %v", err)
	}
	if 1 != 255 {
		p := append([]byte{}, base.Payload...)
		p[0] |= byte(254)
		a, err := getter(ie.New(ie.ExtendedDLBufferingNotificationPolicy, p))
		b, baseErr := getter(base)
		if err != nil || baseErr != nil || !reflect.DeepEqual(a, b) {
			t.Fatal("spare bits evaluated", err)
		}
	}
}

func TestExtendedDLBufferingNotificationPolicyInvalidWire(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n) & 1
		valid := v <= 1
		if valid {
			continue
		}
		_, err := ie.New(ie.ExtendedDLBufferingNotificationPolicy, []byte{byte(n)}).ExtendedDLBufferingNotificationPolicy()
		if err == nil {
			t.Fatalf("invalid wire value %d accepted", n)
		}
	}
}
