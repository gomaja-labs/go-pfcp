package ie_test

import (
	"encoding/binary"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestMeasurementEndpointPortNumber checks TS 29.244 V20.1.0 clause 8.2.265.
func TestMeasurementEndpointPortNumber(t *testing.T) {
	for n := 0; n <= 65535; n++ {
		v := uint16(n)
		valid := true
		i := ie.NewMeasurementEndpointPortNumber(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.MeasurementEndpointPortNumber)
			continue
		}
		payload := make([]byte, 2)
		binary.BigEndian.PutUint16(payload, v)
		checkNewIEWire(t, i, payload)
		got, err := i.MeasurementEndpointPortNumber()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.MeasurementEndpointPortNumber() }
	checkNewIEFailures(t, ie.MeasurementEndpointPortNumber, getter, 2)
	base := ie.NewMeasurementEndpointPortNumber(0)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.MeasurementEndpointPortNumber, extension))
	if (err != nil) != true {
		t.Fatalf("extension: %v", err)
	}
}
