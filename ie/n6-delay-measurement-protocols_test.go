package ie_test

import (
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestN6DelayMeasurementProtocols checks TS 29.244 V20.1.0 clause 8.2.246.
func TestN6DelayMeasurementProtocols(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n)
		valid := v <= 15
		i := ie.NewN6DelayMeasurementProtocols(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.N6DelayMeasurementProtocols)
			continue
		}
		payload := make([]byte, 1)
		payload[0] = v
		checkNewIEWire(t, i, payload)
		got, err := i.N6DelayMeasurementProtocols()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.N6DelayMeasurementProtocols() }
	checkNewIEFailures(t, ie.N6DelayMeasurementProtocols, getter, 1)
	base := ie.NewN6DelayMeasurementProtocols(0)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.N6DelayMeasurementProtocols, extension))
	if (err != nil) != false {
		t.Fatalf("extension: %v", err)
	}
	if 15 != 255 {
		p := append([]byte{}, base.Payload...)
		p[0] |= byte(240)
		a, err := getter(ie.New(ie.N6DelayMeasurementProtocols, p))
		b, baseErr := getter(base)
		if err != nil || baseErr != nil || !reflect.DeepEqual(a, b) {
			t.Fatal("spare bits evaluated", err)
		}
	}
}

func TestN6DelayMeasurementProtocolsInvalidWire(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n) & 15
		valid := v <= 15
		if valid {
			continue
		}
		_, err := ie.New(ie.N6DelayMeasurementProtocols, []byte{byte(n)}).N6DelayMeasurementProtocols()
		if err == nil {
			t.Fatalf("invalid wire value %d accepted", n)
		}
	}
}
