package ie_test

import (
	"errors"
	"io"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.27 and RFC 2474 section 3: PPI carries six DSCP bits.
func TestPPISixBitValue(t *testing.T) {
	for n := 0; n < 256; n++ {
		item := ie.New(ie.DownlinkDataServiceInformation, []byte{0xf9, byte(n), 0xaa})
		want := byte(n) & 63
		ppi, err := item.PPI()
		if err != nil || ppi != want {
			t.Errorf("octet %d: PPI=%d, %v; want %d", n, ppi, err, want)
		}
		fields, err := item.DownlinkDataServiceInformation()
		if err != nil || fields.PagingPolicyIndication != want {
			t.Errorf("octet %d: fields=%v, %v", n, fields, err)
		}
		wire, err := item.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		parsed, err := ie.Parse(wire)
		if err != nil {
			t.Fatal(err)
		}
		again, err := parsed.PPI()
		if err != nil || again != want {
			t.Errorf("octet %d: parsed %d, %v", n, again, err)
		}
		if n <= 63 {
			constructed := ie.NewDownlinkDataServiceInformation(true, true, uint8(n), 63)
			b, err := constructed.Marshal()
			if err != nil {
				t.Errorf("value %d: %v", n, err)
				continue
			}
			if b[5] != uint8(n) || b[6] != 63 {
				t.Error("constructor changed PPI or QFI")
			}
			f, err := constructed.DownlinkDataServiceInformation()
			if err != nil || f.PagingPolicyIndication != uint8(n) || f.QFI != 63 {
				t.Errorf("constructor round trip: %v", err)
			}
		} else {
			assertInvalidConstructor(t, ie.NewDownlinkDataServiceInformation(true, false, uint8(n), 0))
			f := &ie.DownlinkDataServiceInformationFields{Flags: 1, PagingPolicyIndication: uint8(n)}
			if _, err := f.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
				t.Errorf("fields value %d: %v", n, err)
			}
			if err := f.MarshalTo(make([]byte, 8)); !errors.Is(err, ie.ErrInvalidValue) {
				t.Errorf("fields MarshalTo value %d: %v", n, err)
			}
		}
	}
	for _, b := range [][]byte{nil, {1}} {
		if _, err := ie.New(ie.DownlinkDataServiceInformation, b).PPI(); !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatal(err)
		}
	}
	if _, err := ie.NewDownlinkDataServiceInformation(false, true, 0, 9).PPI(); !errors.Is(err, ie.ErrIENotFound) {
		t.Fatal(err)
	}
}
