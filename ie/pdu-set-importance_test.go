package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestPDUSetImportance(t *testing.T) {
	for n := 1; n <= 65535; n++ {
		v := uint16(n)
		i := ie.NewPDUSetImportance(v)
		checkNewIEWire(t, i, []byte{byte(v), byte(v >> 8)})
		got, err := i.PDUSetImportance()
		if err != nil || got != v {
			t.Fatal("PSI bitmap", err)
		}
	}
	assertNewIEInvalid(t, ie.NewPDUSetImportance(0), ie.PDUSetImportance)
	checkNewIEFailures(t, ie.PDUSetImportance, func(i *ie.IE) (any, error) { return i.PDUSetImportance() }, 2)
	if _, err := ie.New(ie.PDUSetImportance, []byte{1, 0, 0}).PDUSetImportance(); err == nil {
		t.Fatal("fixed-length extension accepted")
	}
}
