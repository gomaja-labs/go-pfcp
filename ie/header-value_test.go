package ie_test

import (
	"bytes"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestHeaderValue(t *testing.T) {
	for _, n := range []int{0, 1, 255, 256, 65535, 65536} {
		payload := bytes.Repeat([]byte{0x61}, n)
		i := ie.NewHeaderValue(payload)
		if n > 65535 {
			assertNewIEInvalid(t, i, ie.HeaderValue)
			continue
		}
		checkNewIEWire(t, i, payload)
		got, err := i.HeaderValue()
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
	checkNewIEFailures(t, ie.HeaderValue, func(i *ie.IE) (any, error) { return i.HeaderValue() }, 1)
}

func TestHeaderValueBinary(t *testing.T) {
	p := []byte{0, 1, 255, 128}
	i := ie.NewHeaderValue(p)
	checkNewIEWire(t, i, p)
}
