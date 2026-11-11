package ie_test

import (
	"encoding/binary"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestN6DelayMeasurementControlInformationID checks TS 29.244 V20.1.0 clause 8.2.264.
func TestN6DelayMeasurementControlInformationID(t *testing.T) {
	for n := 0; n <= 65535; n++ {
		v := uint16(n)
		valid := true
		i := ie.NewN6DelayMeasurementControlInformationID(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.N6DelayMeasurementControlInformationID)
			continue
		}
		payload := make([]byte, 2)
		binary.BigEndian.PutUint16(payload, v)
		checkNewIEWire(t, i, payload)
		got, err := i.N6DelayMeasurementControlInformationID()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.N6DelayMeasurementControlInformationID() }
	checkNewIEFailures(t, ie.N6DelayMeasurementControlInformationID, getter, 2)
	base := ie.NewN6DelayMeasurementControlInformationID(0)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.N6DelayMeasurementControlInformationID, extension))
	if (err != nil) != true {
		t.Fatalf("extension: %v", err)
	}
}
