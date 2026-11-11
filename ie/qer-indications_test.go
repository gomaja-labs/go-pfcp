package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.216: every defined value and every received spare-bit pattern.
func TestQERIndicationsCurrent(t *testing.T) {
	for v := 0; v <= 127; v++ {
		assertAddedVector(t, ie.NewQERIndications(ie.QERIndicationFlags(v)), []byte{byte(v)}, "QERIndications", false)
	}
	for v := 0; v < 256; v++ {
		got, err := ie.New(ie.QERIndications, []byte{byte(v)}).QERIndications()
		if err != nil || uint8(got) != uint8(v&127) {
			t.Fatalf("received %d: %v, %v", v, got, err)
		}
	}
	for _, v := range []int{128, 255} {
		if v >= 0 && v <= 127 {
			continue
		}
		assertAddedInvalid(t, ie.NewQERIndications(ie.QERIndicationFlags(v)))
	}
}

var _ = registerConstructorFixture(ie.NewQERIndications, ie.QERIndicationFlags(1))
