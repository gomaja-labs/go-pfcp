package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestSecurityModeOWAMPOrTWAMP(t *testing.T) {
	for v := uint32(0); v < 2048; v++ {
		mode := v & 15
		valid := v < 1024 && (mode == 1 || mode == 2 || mode == 4 || mode == 8 || v == 0)
		i := ie.NewSecurityModeOWAMPOrTWAMP(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.SecurityModeOWAMPOrTWAMP)
			continue
		}
		checkNewIEWire(t, i, []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)})
		got, err := i.SecurityModeOWAMPOrTWAMP()
		if err != nil || got != v {
			t.Fatal("mode", err)
		}
		payload := append([]byte{}, i.Payload...)
		payload[0] |= 0xff
		got, err = ie.New(ie.SecurityModeOWAMPOrTWAMP, payload).SecurityModeOWAMPOrTWAMP()
		if err != nil || got != v {
			t.Fatal("unassigned bits evaluated", err)
		}
	}
	checkNewIEFailures(t, ie.SecurityModeOWAMPOrTWAMP, func(i *ie.IE) (any, error) { return i.SecurityModeOWAMPOrTWAMP() }, 4)
	if _, err := ie.New(ie.SecurityModeOWAMPOrTWAMP, []byte{0, 0, 0, 1, 0}).SecurityModeOWAMPOrTWAMP(); err == nil {
		t.Fatal("fixed extension accepted")
	}
}
