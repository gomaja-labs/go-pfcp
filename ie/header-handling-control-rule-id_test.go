package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestHeaderHandlingControlRuleID checks TS 29.244 V20.1.0 clause 8.2.260.
func TestHeaderHandlingControlRuleID(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n)
		valid := true
		i := ie.NewHeaderHandlingControlRuleID(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.HeaderHandlingControlRuleID)
			continue
		}
		payload := make([]byte, 1)
		payload[0] = v
		checkNewIEWire(t, i, payload)
		got, err := i.HeaderHandlingControlRuleID()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.HeaderHandlingControlRuleID() }
	checkNewIEFailures(t, ie.HeaderHandlingControlRuleID, getter, 1)
	base := ie.NewHeaderHandlingControlRuleID(0)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.HeaderHandlingControlRuleID, extension))
	if (err != nil) != true {
		t.Fatalf("extension: %v", err)
	}
}
