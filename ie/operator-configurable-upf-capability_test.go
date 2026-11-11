package ie_test

import (
	"bytes"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestOperatorConfigurableUPFCapability(t *testing.T) {
	for _, n := range []int{0, 1, 255, 256, 65535, 65536} {
		payload := bytes.Repeat([]byte{0x61}, n)
		i := ie.NewOperatorConfigurableUPFCapability(string(payload))
		if n > 65535 {
			assertNewIEInvalid(t, i, ie.OperatorConfigurableUPFCapability)
			continue
		}
		checkNewIEWire(t, i, payload)
		got, err := i.OperatorConfigurableUPFCapability()
		if n == 0 {
			if err == nil {
				t.Fatal("null helper accepted")
			}
			continue
		}
		if err != nil || !bytes.Equal([]byte(got), payload) {
			t.Fatal("value", err)
		}
	}
	checkNewIEFailures(t, ie.OperatorConfigurableUPFCapability, func(i *ie.IE) (any, error) { return i.OperatorConfigurableUPFCapability() }, 1)
}
func TestOperatorConfigurableUPFCapabilityUTF8(t *testing.T) {
	assertNewIEInvalid(t, ie.NewOperatorConfigurableUPFCapability(string([]byte{255})), ie.OperatorConfigurableUPFCapability)
}

func TestOperatorConfigurableUPFCapabilityText(t *testing.T) {
	i := ie.NewOperatorConfigurableUPFCapability("capacité")
	checkNewIEWire(t, i, []byte("capacité"))
}
