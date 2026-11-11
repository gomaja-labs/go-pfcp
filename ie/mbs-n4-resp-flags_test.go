package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.211: every defined value and every received spare-bit pattern.
func TestMBSN4RespFlagsCurrent(t *testing.T) {
	for v := 0; v <= 7; v++ {
		assertAddedVector(t, ie.NewMBSN4RespFlags(ie.MBSN4ResponseFlags(v)), []byte{byte(v)}, "MBSN4RespFlags", false)
	}
	for v := 0; v < 256; v++ {
		got, err := ie.New(ie.MBSN4RespFlags, []byte{byte(v)}).MBSN4RespFlags()
		if err != nil || uint8(got) != uint8(v&7) {
			t.Fatalf("received %d: %v, %v", v, got, err)
		}
	}
	for _, v := range []int{8, 255} {
		if v >= 0 && v <= 7 {
			continue
		}
		assertAddedInvalid(t, ie.NewMBSN4RespFlags(ie.MBSN4ResponseFlags(v)))
	}
}

var _ = registerConstructorFixture(ie.NewMBSN4RespFlags, ie.MBSN4ResponseFlags(1))
