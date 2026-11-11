package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestNodeLevelRuleID checks TS 29.244 V20.1.0 clause 8.2.284.
func TestNodeLevelRuleID(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n)
		valid := v >= 1
		i := ie.NewNodeLevelRuleID(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.NodeLevelRuleID)
			continue
		}
		payload := make([]byte, 1)
		payload[0] = v
		checkNewIEWire(t, i, payload)
		got, err := i.NodeLevelRuleID()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.NodeLevelRuleID() }
	checkNewIEFailures(t, ie.NodeLevelRuleID, getter, 1)
	base := ie.NewNodeLevelRuleID(1)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.NodeLevelRuleID, extension))
	if (err != nil) != true {
		t.Fatalf("extension: %v", err)
	}
}

func TestNodeLevelRuleIDInvalidWire(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n) & 255
		valid := v >= 1
		if valid {
			continue
		}
		_, err := ie.New(ie.NodeLevelRuleID, []byte{byte(n)}).NodeLevelRuleID()
		if err == nil {
			t.Fatalf("invalid wire value %d accepted", n)
		}
	}
}
