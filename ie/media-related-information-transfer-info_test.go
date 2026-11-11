package ie_test

import (
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestMediaRelatedInformationTransferInfo checks TS 29.244 V20.1.0 clause 8.2.271.
func TestMediaRelatedInformationTransferInfo(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n)
		valid := v <= 1
		i := ie.NewMediaRelatedInformationTransferInfo(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.MediaRelatedInformationTransferInfo)
			continue
		}
		payload := make([]byte, 1)
		payload[0] = v
		checkNewIEWire(t, i, payload)
		got, err := i.MediaRelatedInformationTransferInfo()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.MediaRelatedInformationTransferInfo() }
	checkNewIEFailures(t, ie.MediaRelatedInformationTransferInfo, getter, 1)
	base := ie.NewMediaRelatedInformationTransferInfo(0)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.MediaRelatedInformationTransferInfo, extension))
	if (err != nil) != false {
		t.Fatalf("extension: %v", err)
	}
	if 15 != 255 {
		p := append([]byte{}, base.Payload...)
		p[0] |= byte(240)
		a, err := getter(ie.New(ie.MediaRelatedInformationTransferInfo, p))
		b, baseErr := getter(base)
		if err != nil || baseErr != nil || !reflect.DeepEqual(a, b) {
			t.Fatal("spare bits evaluated", err)
		}
	}
}

func TestMediaRelatedInformationTransferInfoInvalidWire(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n) & 15
		valid := v <= 1
		if valid {
			continue
		}
		_, err := ie.New(ie.MediaRelatedInformationTransferInfo, []byte{byte(n)}).MediaRelatedInformationTransferInfo()
		if err == nil {
			t.Fatalf("invalid wire value %d accepted", n)
		}
	}
}
