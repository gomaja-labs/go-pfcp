package ie_test

import (
	"bytes"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestHeaderHandlingControlReference(t *testing.T) {
	for _, n := range []int{0, 1, 255, 256, 65535, 65536} {
		payload := bytes.Repeat([]byte{0x61}, n)
		i := ie.NewHeaderHandlingControlReference(payload)
		if n > 65535 {
			assertNewIEInvalid(t, i, ie.HeaderHandlingControlReference)
			continue
		}
		checkNewIEWire(t, i, payload)
		got, err := i.HeaderHandlingControlReference()
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
	checkNewIEFailures(t, ie.HeaderHandlingControlReference, func(i *ie.IE) (any, error) { return i.HeaderHandlingControlReference() }, 1)
}

func TestHeaderHandlingControlReferenceBinary(t *testing.T) {
	p := []byte{0, 1, 255, 128}
	i := ie.NewHeaderHandlingControlReference(p)
	checkNewIEWire(t, i, p)
}
