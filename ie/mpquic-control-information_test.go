package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.225: every defined value and every received spare-bit pattern.
func TestMPQUICControlInformationCurrent(t *testing.T) {
	for v := 1; v <= 7; v++ {
		assertAddedVector(t, ie.NewMPQUICControlInformation(ie.MPQUICControlFlags(v)), []byte{byte(v)}, "MPQUICControlInformation", false)
	}
	for v := 0; v < 256; v++ {
		got, err := ie.New(ie.MPQUICControlInformation, []byte{byte(v)}).MPQUICControlInformation()
		if v&7 == 0 {
			if err == nil {
				t.Fatal("prohibited value accepted")
			}
			continue
		}
		if err != nil || uint8(got) != uint8(v&7) {
			t.Fatalf("received %d: %v, %v", v, got, err)
		}
	}
	for _, v := range []int{0, 255} {
		if v >= 1 && v <= 7 {
			continue
		}
		assertAddedInvalid(t, ie.NewMPQUICControlInformation(ie.MPQUICControlFlags(v)))
	}
}
