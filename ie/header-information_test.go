package ie_test

import (
	"bytes"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestHeaderInformation(t *testing.T) {
	for _, n := range []int{0, 1, 255, 256, 65535, 65536} {
		payload := bytes.Repeat([]byte{0x61}, n)
		i := ie.NewHeaderInformation(string(payload))
		if n > 65535 {
			assertNewIEInvalid(t, i, ie.HeaderInformation)
			continue
		}
		checkNewIEWire(t, i, payload)
		got, err := i.HeaderInformation()
		if n == 0 {
			if err == nil {
				t.Fatal("null helper accepted")
			}
			continue
		}
		if err != nil || got != string(payload) {
			t.Fatal("value", err)
		}
	}
	checkNewIEFailures(t, ie.HeaderInformation, func(i *ie.IE) (any, error) { return i.HeaderInformation() }, 1)
}

func TestHeaderInformationBinary(t *testing.T) {
	p := []byte{0, 1, 255, 128}
	i := ie.NewHeaderInformation(string(p))
	checkNewIEWire(t, i, p)
}
