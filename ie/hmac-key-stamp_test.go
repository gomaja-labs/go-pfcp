package ie_test

import (
	"bytes"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestHMACKeySTAMP(t *testing.T) {
	for _, n := range []int{0, 1, 255, 256, 65535, 65536} {
		payload := bytes.Repeat([]byte{0x61}, n)
		i := ie.NewHMACKeySTAMP(payload)
		if n > 65535 {
			assertNewIEInvalid(t, i, ie.HMACKeySTAMP)
			continue
		}
		checkNewIEWire(t, i, payload)
		got, err := i.HMACKeySTAMP()
		if n == 0 {
			if err == nil {
				t.Fatal("null helper accepted")
			}
			continue
		}
		if err != nil || !bytes.Equal(got, payload) {
			t.Fatal("value", err)
		}
	}
	checkNewIEFailures(t, ie.HMACKeySTAMP, func(i *ie.IE) (any, error) { return i.HMACKeySTAMP() }, 1)
}

func TestHMACKeySTAMPBinary(t *testing.T) {
	p := []byte{0, 1, 255, 128}
	i := ie.NewHMACKeySTAMP(p)
	checkNewIEWire(t, i, p)
}
