package ie_test

import (
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestReportingControlInformation checks TS 29.244 V20.1.0 clause 8.2.272.
func TestReportingControlInformation(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n)
		valid := v <= 1
		i := ie.NewReportingControlInformation(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.ReportingControlInformation)
			continue
		}
		payload := make([]byte, 1)
		payload[0] = v
		checkNewIEWire(t, i, payload)
		got, err := i.ReportingControlInformation()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.ReportingControlInformation() }
	checkNewIEFailures(t, ie.ReportingControlInformation, getter, 1)
	base := ie.NewReportingControlInformation(0)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.ReportingControlInformation, extension))
	if (err != nil) != false {
		t.Fatalf("extension: %v", err)
	}
	if 1 != 255 {
		p := append([]byte{}, base.Payload...)
		p[0] |= byte(254)
		a, err := getter(ie.New(ie.ReportingControlInformation, p))
		b, baseErr := getter(base)
		if err != nil || baseErr != nil || !reflect.DeepEqual(a, b) {
			t.Fatal("spare bits evaluated", err)
		}
	}
}

func TestReportingControlInformationInvalidWire(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n) & 1
		valid := v <= 1
		if valid {
			continue
		}
		_, err := ie.New(ie.ReportingControlInformation, []byte{byte(n)}).ReportingControlInformation()
		if err == nil {
			t.Fatalf("invalid wire value %d accepted", n)
		}
	}
}
