package ie_test

import (
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestHeaderHandlingReportingIndication checks TS 29.244 V20.1.0 clause 8.2.266.
func TestHeaderHandlingReportingIndication(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n)
		valid := v <= 1
		i := ie.NewHeaderHandlingReportingIndication(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.HeaderHandlingReportingIndication)
			continue
		}
		payload := make([]byte, 1)
		payload[0] = v
		checkNewIEWire(t, i, payload)
		got, err := i.HeaderHandlingReportingIndication()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.HeaderHandlingReportingIndication() }
	checkNewIEFailures(t, ie.HeaderHandlingReportingIndication, getter, 1)
	base := ie.NewHeaderHandlingReportingIndication(0)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.HeaderHandlingReportingIndication, extension))
	if (err != nil) != false {
		t.Fatalf("extension: %v", err)
	}
	if 1 != 255 {
		p := append([]byte{}, base.Payload...)
		p[0] |= byte(254)
		a, err := getter(ie.New(ie.HeaderHandlingReportingIndication, p))
		b, baseErr := getter(base)
		if err != nil || baseErr != nil || !reflect.DeepEqual(a, b) {
			t.Fatal("spare bits evaluated", err)
		}
	}
}

func TestHeaderHandlingReportingIndicationInvalidWire(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n) & 1
		valid := v <= 1
		if valid {
			continue
		}
		_, err := ie.New(ie.HeaderHandlingReportingIndication, []byte{byte(n)}).HeaderHandlingReportingIndication()
		if err == nil {
			t.Fatalf("invalid wire value %d accepted", n)
		}
	}
}
