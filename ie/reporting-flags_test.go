package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.204: every defined value and every received spare-bit pattern.
func TestReportingFlagsCurrent(t *testing.T) {
	for v := 0; v <= 1; v++ {
		assertAddedVector(t, ie.NewReportingFlags(ie.ReportingFlagBits(v)), []byte{byte(v)}, "ReportingFlags", false)
	}
	for v := 0; v < 256; v++ {
		got, err := ie.New(ie.ReportingFlags, []byte{byte(v)}).ReportingFlags()
		if err != nil || uint8(got) != uint8(v&1) {
			t.Fatalf("received %d: %v, %v", v, got, err)
		}
	}
	for _, v := range []int{2, 255} {
		if v >= 0 && v <= 1 {
			continue
		}
		assertAddedInvalid(t, ie.NewReportingFlags(ie.ReportingFlagBits(v)))
	}
}

var _ = registerConstructorFixture(ie.NewReportingFlags, ie.ReportingFlagBits(1))
