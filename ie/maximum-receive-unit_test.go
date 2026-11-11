package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestMaximumReceiveUnitCurrent(t *testing.T) {
	for _, v := range []uint16{0, 1, 0x1234, 65535} {
		assertAddedVector(t, ie.NewMaximumReceiveUnit(v), []byte{byte(v >> 8), byte(v)}, "MaximumReceiveUnit", true)
		got, err := ie.NewMaximumReceiveUnit(v).MaximumReceiveUnit()
		if err != nil || got != v {
			t.Fatalf("value %d: %d %v", v, got, err)
		}
	}

}
