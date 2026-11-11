package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestReportingEndpointID checks TS 29.244 V20.1.0 clause 8.2.253.
func TestReportingEndpointID(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n)
		valid := true
		i := ie.NewReportingEndpointID(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.ReportingEndpointID)
			continue
		}
		payload := make([]byte, 1)
		payload[0] = v
		checkNewIEWire(t, i, payload)
		got, err := i.ReportingEndpointID()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.ReportingEndpointID() }
	checkNewIEFailures(t, ie.ReportingEndpointID, getter, 1)
	base := ie.NewReportingEndpointID(0)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.ReportingEndpointID, extension))
	if (err != nil) != true {
		t.Fatalf("extension: %v", err)
	}
}
