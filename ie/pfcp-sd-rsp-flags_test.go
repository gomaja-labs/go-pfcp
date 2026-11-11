package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.215: every defined value and every received spare-bit pattern.
func TestPFCPSDRspFlagsCurrent(t *testing.T) {
	for v := 0; v <= 1; v++ {
		assertAddedVector(t, ie.NewPFCPSDRspFlags(ie.PFCPSDResponseFlags(v)), []byte{byte(v)}, "PFCPSDRspFlags", false)
	}
	for v := 0; v < 256; v++ {
		got, err := ie.New(ie.PFCPSDRspFlags, []byte{byte(v)}).PFCPSDRspFlags()
		if err != nil || uint8(got) != uint8(v&1) {
			t.Fatalf("received %d: %v, %v", v, got, err)
		}
	}
	for _, v := range []int{2, 255} {
		if v >= 0 && v <= 1 {
			continue
		}
		assertAddedInvalid(t, ie.NewPFCPSDRspFlags(ie.PFCPSDResponseFlags(v)))
	}
}
