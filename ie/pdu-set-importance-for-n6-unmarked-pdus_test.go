package ie_test

import (
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestPDUSetImportanceForN6UnmarkedPDUs checks TS 29.244 V20.1.0 clause 8.2.282.
func TestPDUSetImportanceForN6UnmarkedPDUs(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n)
		valid := v >= 1 && v <= 15
		i := ie.NewPDUSetImportanceForN6UnmarkedPDUs(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.PDUSetImportanceForN6UnmarkedPDUs)
			continue
		}
		payload := make([]byte, 1)
		payload[0] = v
		checkNewIEWire(t, i, payload)
		got, err := i.PDUSetImportanceForN6UnmarkedPDUs()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.PDUSetImportanceForN6UnmarkedPDUs() }
	checkNewIEFailures(t, ie.PDUSetImportanceForN6UnmarkedPDUs, getter, 1)
	base := ie.NewPDUSetImportanceForN6UnmarkedPDUs(1)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.PDUSetImportanceForN6UnmarkedPDUs, extension))
	if (err != nil) != true {
		t.Fatalf("extension: %v", err)
	}
	if 15 != 255 {
		p := append([]byte{}, base.Payload...)
		p[0] |= byte(240)
		a, err := getter(ie.New(ie.PDUSetImportanceForN6UnmarkedPDUs, p))
		b, baseErr := getter(base)
		if err != nil || baseErr != nil || !reflect.DeepEqual(a, b) {
			t.Fatal("spare bits evaluated", err)
		}
	}
}

func TestPDUSetImportanceForN6UnmarkedPDUsInvalidWire(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n) & 15
		valid := v >= 1 && v <= 15
		if valid {
			continue
		}
		_, err := ie.New(ie.PDUSetImportanceForN6UnmarkedPDUs, []byte{byte(n)}).PDUSetImportanceForN6UnmarkedPDUs()
		if err == nil {
			t.Fatalf("invalid wire value %d accepted", n)
		}
	}
}
