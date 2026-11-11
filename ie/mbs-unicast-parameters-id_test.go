package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestMBSUnicastParametersIDCurrent(t *testing.T) {
	for _, v := range []uint16{0, 1, 0x1234, 65535} {
		assertAddedVector(t, ie.NewMBSUnicastParametersID(v), []byte{byte(v >> 8), byte(v)}, "MBSUnicastParametersID", false)
		got, err := ie.NewMBSUnicastParametersID(v).MBSUnicastParametersID()
		if err != nil || got != v {
			t.Fatalf("value %d: %d %v", v, got, err)
		}
	}

}

var _ = registerConstructorFixture(ie.NewMBSUnicastParametersID, uint16(1))
