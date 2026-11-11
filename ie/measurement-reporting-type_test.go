package ie_test

import (
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestMeasurementReportingType checks TS 29.244 V20.1.0 clause 8.2.262.
func TestMeasurementReportingType(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n)
		valid := v > 0 && v <= 7 && v&3 != 3
		i := ie.NewMeasurementReportingType(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.MeasurementReportingType)
			continue
		}
		payload := make([]byte, 1)
		payload[0] = v
		checkNewIEWire(t, i, payload)
		got, err := i.MeasurementReportingType()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.MeasurementReportingType() }
	checkNewIEFailures(t, ie.MeasurementReportingType, getter, 1)
	base := ie.NewMeasurementReportingType(1)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.MeasurementReportingType, extension))
	if (err != nil) != false {
		t.Fatalf("extension: %v", err)
	}
	if 7 != 255 {
		p := append([]byte{}, base.Payload...)
		p[0] |= byte(248)
		a, err := getter(ie.New(ie.MeasurementReportingType, p))
		b, baseErr := getter(base)
		if err != nil || baseErr != nil || !reflect.DeepEqual(a, b) {
			t.Fatal("spare bits evaluated", err)
		}
	}
}

func TestMeasurementReportingTypeInvalidWire(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n) & 7
		valid := v > 0 && v <= 7 && v&3 != 3
		if valid {
			continue
		}
		_, err := ie.New(ie.MeasurementReportingType, []byte{byte(n)}).MeasurementReportingType()
		if err == nil {
			t.Fatalf("invalid wire value %d accepted", n)
		}
	}
}
