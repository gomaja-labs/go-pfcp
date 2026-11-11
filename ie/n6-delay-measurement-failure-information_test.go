package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestN6DelayMeasurementFailureInformation checks TS 29.244 V20.1.0 clause 8.2.263.
func TestN6DelayMeasurementFailureInformation(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n)
		valid := v <= 3
		i := ie.NewN6DelayMeasurementFailureInformation(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.N6DelayMeasurementFailureInformation)
			continue
		}
		payload := make([]byte, 1)
		payload[0] = v
		checkNewIEWire(t, i, payload)
		got, err := i.N6DelayMeasurementFailureInformation()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.N6DelayMeasurementFailureInformation() }
	checkNewIEFailures(t, ie.N6DelayMeasurementFailureInformation, getter, 1)
	base := ie.NewN6DelayMeasurementFailureInformation(0)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.N6DelayMeasurementFailureInformation, extension))
	if (err != nil) != false {
		t.Fatalf("extension: %v", err)
	}
}

func TestN6DelayMeasurementFailureInformationUnknownCause(t *testing.T) {
	for n := 4; n <= 255; n++ {
		i := ie.New(ie.N6DelayMeasurementFailureInformation, []byte{byte(n)})
		got, err := i.N6DelayMeasurementFailureInformation()
		if err != nil || got != ie.N6DelayMeasurementFailureInformationUnspecified {
			t.Fatalf("unknown cause %d: %d, %v", n, got, err)
		}
		checkNewIEWire(t, i, []byte{byte(n)})
	}
}
