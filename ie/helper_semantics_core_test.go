package ie_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestHelperSemanticHeaderOperations(t *testing.T) {
	// TS 29.244 V20.1.0 clauses 8.1.1–8.1.2. These operations
	// interpret headers and opaque payloads; they accept every IE type and
	// null values. Payload field validation belongs to typed value helpers.
	for _, tc := range []struct {
		name   string
		input  *ie.IE
		wire   []byte
		vendor bool
	}{
		{"value", ie.NewCause(1), []byte{0, 19, 0, 1, 1}, false},
		{"unknown-type", ie.New(0x7fff, []byte{1}), []byte{0x7f, 0xff, 0, 1, 1}, false},
		{"null", ie.New(ie.Cause, nil), []byte{0, 19, 0, 0}, false},
		{"opaque-prefix", ie.New(ie.PDRID, []byte{1}), []byte{0, 56, 0, 1, 1}, false},
		{"vendor", ie.NewVendorSpecificIE(0x8001, 0x1234, []byte{1}), []byte{0x80, 1, 0, 3, 0x12, 0x34, 1}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.input.Marshal()
			if err != nil || !bytes.Equal(got, tc.wire) {
				t.Fatalf("Marshal mismatch error=%v", err)
			}
			if got := tc.input.MarshalLen(); got != len(tc.wire) {
				t.Fatalf("MarshalLen=%d; want %d", got, len(tc.wire))
			}
			if got := tc.input.IsVendorSpecific(); got != tc.vendor {
				t.Fatalf("IsVendorSpecific=%v; want %v", got, tc.vendor)
			}
		})
	}
	var absent *ie.IE
	if _, err := absent.Marshal(); !errors.Is(err, ie.ErrIENotFound) {
		t.Fatal(err)
	}
	if absent.MarshalLen() != 0 || absent.IsVendorSpecific() {
		t.Fatal("nil header operations returned a value")
	}
}

func TestHelperSemanticValueAsGrouped(t *testing.T) {
	// The explicit interpretation is intentionally independent of automatic
	// grouped classification (TS 29.244 V20.1.0 clause 8.1.2).
	for _, typ := range []ie.IEType{ie.CreatePDR, ie.ApplicationID, 0x8001} {
		group := ie.New(typ, []byte{0, 19, 0, 1, 7})
		got, err := group.ValueAsGrouped()
		if err != nil || len(got) != 1 || !got[0].Equal(ie.NewCause(7)) {
			t.Fatalf("type=%d interpreted group error=%v", typ, err)
		}
		null := ie.New(typ, nil)
		got, err = null.ValueAsGrouped()
		if err != nil || len(got) != 0 {
			t.Fatalf("type=%d empty group error=%v", typ, err)
		}
		for n := 1; n < 5; n++ {
			if _, err := ie.New(typ, group.Payload[:n]).ValueAsGrouped(); err == nil {
				t.Fatalf("type=%d accepted truncated group length %d", typ, n)
			}
		}
	}
	var absent *ie.IE
	if _, err := absent.ValueAsGrouped(); !errors.Is(err, ie.ErrIENotFound) {
		t.Fatal(err)
	}
}
