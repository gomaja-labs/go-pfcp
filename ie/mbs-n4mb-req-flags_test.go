package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.208: every defined value and every received spare-bit pattern.
func TestMBSN4mbReqFlagsCurrent(t *testing.T) {
	for v := 0; v <= 7; v++ {
		assertAddedVector(t, ie.NewMBSN4mbReqFlags(ie.MBSN4mbRequestFlags(v)), []byte{byte(v)}, "MBSN4mbReqFlags", false)
	}
	for v := 0; v < 256; v++ {
		got, err := ie.New(ie.MBSN4mbReqFlags, []byte{byte(v)}).MBSN4mbReqFlags()
		if err != nil || uint8(got) != uint8(v&7) {
			t.Fatalf("received %d: %v, %v", v, got, err)
		}
	}
	for _, v := range []int{8, 255} {
		if v >= 0 && v <= 7 {
			continue
		}
		assertAddedInvalid(t, ie.NewMBSN4mbReqFlags(ie.MBSN4mbRequestFlags(v)))
	}
}

var _ = registerConstructorFixture(ie.NewMBSN4mbReqFlags, ie.MBSN4mbRequestFlags(1))
