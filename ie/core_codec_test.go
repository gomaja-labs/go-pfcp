package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/gomaja-labs/go-pfcp/internal/spec"
)

func TestStrictIEBounds(t *testing.T) {
	for _, tc := range []struct {
		name string
		wire []byte
		want error
	}{
		{"header_only", []byte{0, 19, 0, 1}, io.ErrUnexpectedEOF},
		{"vendor_short_value", []byte{128, 1, 0, 1, 0, 0}, ie.ErrInvalidLength},
		{"vendor_missing_value", []byte{128, 1, 0, 2}, io.ErrUnexpectedEOF},
		{"grouped_truncated_child", []byte{0, 1, 0, 4, 0, 19, 0, 1}, io.ErrUnexpectedEOF},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ie.Parse(tc.wire)
			if !errors.Is(err, tc.want) {
				t.Fatalf("got %v, want %v", err, tc.want)
			}
		})
	}
}

func TestPayloadCapacity(t *testing.T) {
	b := []byte{0, 19, 0, 1, 7, 0, 19, 0, 1, 99}
	ies, err := ie.ParseMultiIEs(b)
	if err != nil || len(ies) != 2 {
		t.Fatal(ies, err)
	}
	_ = append(ies[0].Payload, 8)
	if b[5] != 0 || b[9] != 99 {
		t.Fatal("append overwrote the following IE")
	}
}

// TestParseRejectsTrailingOctets checks that Parse decodes exactly one IE:
// octets after the declared length are not silently ignored.
func TestParseRejectsTrailingOctets(t *testing.T) {
	_, err := ie.Parse([]byte{0, 19, 0, 1, 7, 99})
	var de *ie.DecodeError
	if !errors.Is(err, ie.ErrInvalidLength) || !errors.As(err, &de) || de.IEType() != ie.Cause {
		t.Fatalf("got %v, want a Cause DecodeError wrapping ErrInvalidLength", err)
	}
}

// nestedCreatePDR returns depth Create PDR IEs, each holding the next, around
// one PDR ID.
func nestedCreatePDR(depth int) []byte {
	b := make([]byte, 4*depth, 4*depth+6)
	b = append(b, 0, 56, 0, 2, 0, 1)
	for level := range depth {
		length := len(b) - 4*(level+1)
		copy(b[4*level:], []byte{0, 1, byte(length >> 8), byte(length)})
	}
	return b
}

// TestGroupedNestingBound checks the decoder's nesting bound: sixteen levels of
// grouped IEs decode, a seventeenth is rejected with ErrNestingTooDeep, and a
// 64 KiB datagram of nested headers is rejected without following thousands
// of levels.
func TestGroupedNestingBound(t *testing.T) {
	if _, err := ie.Parse(nestedCreatePDR(16)); err != nil {
		t.Fatalf("16 levels: %v", err)
	}
	if _, err := ie.Parse(nestedCreatePDR(17)); !errors.Is(err, ie.ErrNestingTooDeep) {
		t.Fatalf("17 levels: got %v, want ErrNestingTooDeep", err)
	}
	if _, err := ie.ParseMultiIEs(nestedCreatePDR(16380)); !errors.Is(err, ie.ErrNestingTooDeep) {
		t.Fatalf("16380 levels: got %v, want ErrNestingTooDeep", err)
	}
}

// TestParseMultiIEsAdvancesByDeclaredLength decodes many siblings that each
// hold the deepest accepted nesting; decoding work is linear in the input.
func TestParseMultiIEsAdvancesByDeclaredLength(t *testing.T) {
	one := nestedCreatePDR(16)
	var b []byte
	for len(b)+len(one) <= 65535 {
		b = append(b, one...)
	}
	ies, err := ie.ParseMultiIEs(b)
	if err != nil {
		t.Fatal(err)
	}
	if want := len(b) / len(one); len(ies) != want {
		t.Fatalf("decoded %d IEs, want %d", len(ies), want)
	}
	for _, i := range ies {
		if i.MarshalLen() != len(one) {
			t.Fatalf("decoded IE re-measures %d octets, want %d", i.MarshalLen(), len(one))
		}
	}
}

func TestComputedLength(t *testing.T) {
	if _, ok := reflect.TypeFor[ie.IE]().FieldByName("Length"); ok {
		t.Error("IE retains a second length source")
	}
	i := ie.NewCause(1)
	i.Payload = append(i.Payload, 2)
	b, err := i.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, []byte{0, 19, 0, 2, 1, 2}) {
		t.Fatal("header did not follow payload length")
	}
}

func TestMarshalBounds(t *testing.T) {
	for _, i := range []*ie.IE{ie.NewCause(1), ie.NewVendorSpecificIE(0x8001, 123, []byte{4}), ie.NewCreatePDR(ie.NewCause(1))} {
		for n := 0; n < i.MarshalLen(); n++ {
			t.Run(i.Type.String()+string(rune('a'+n)), func(t *testing.T) {
				if p := panicOf(func() {
					if err := i.MarshalTo(make([]byte, n)); err == nil {
						t.Error("accepted short buffer")
					}
				}); p != nil {
					t.Fatalf("panic %v", p)
				}
			})
		}
	}
	for _, i := range []*ie.IE{ie.New(19, make([]byte, 65536)), ie.NewVendorSpecificIE(0x8001, 1, make([]byte, 65534)), ie.NewCreatePDR(ie.New(19, make([]byte, 65532)))} {
		if _, err := i.Marshal(); err == nil {
			t.Error("accepted value length overflow")
		}
	}
}

func TestGroupedSourceOfTruth(t *testing.T) {
	raw, _ := ie.NewPDRID(7).Marshal()
	i := ie.New(ie.CreatePDR, raw)
	b, err := i.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 4+len(raw) {
		t.Fatal("raw grouped payload discarded")
	}
	decoded, err := ie.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Payload != nil || len(decoded.ChildIEs) != 1 {
		t.Error("decoded group must store children only")
	}
	fields, err := decoded.CreatePDR()
	if err != nil {
		t.Fatal(err)
	}
	v, err := fields.PDRID.PDRID()
	if err != nil || v != 7 {
		t.Errorf("child lookup %d %v", v, err)
	}
	i = ie.NewCreatePDR(ie.NewPDRID(7), nil)
	if i == nil || i.Payload != nil || len(i.ChildIEs) != 1 {
		t.Fatal("constructor must collect children only")
	}
	i.ChildIEs[0].Payload[1] = 9
	fields, err = i.CreatePDR()
	if err != nil {
		t.Fatal(err)
	}
	v, err = fields.PDRID.PDRID()
	if err != nil || v != 9 {
		t.Errorf("lookup used stale payload: %d %v", v, err)
	}
	i = ie.New(ie.CreatePDR, raw)
	vIE, err := i.FindByType(ie.PDRID)
	if err != nil || vIE == nil {
		t.Fatal("raw group lookup", err)
	}
	if err := i.Add(nil, ie.NewCause(4)); err != nil {
		t.Fatal(err)
	}
	if err := i.Remove(ie.PDRID); err != nil {
		t.Fatal(err)
	}
	if len(i.ChildIEs) != 1 || i.ChildIEs[0].Type != ie.Cause {
		t.Fatal("raw group edits failed")
	}
}

func TestGroupedEditErrors(t *testing.T) {
	i := ie.New(ie.CreatePDR, []byte{0, 19, 0, 1})
	before := append([]byte(nil), i.Payload...)
	for _, name := range []string{"Add", "Remove"} {
		method := reflect.ValueOf(i).MethodByName(name)
		args := []reflect.Value{reflect.ValueOf(ie.Cause)}
		if name == "Add" {
			args = []reflect.Value{reflect.ValueOf(ie.NewCause(1))}
		}
		results := method.Call(args)
		if len(results) != 1 || results[0].IsNil() {
			t.Errorf("%s did not report the malformed payload", name)
		}
		if !bytes.Equal(before, i.Payload) || i.ChildIEs != nil {
			t.Errorf("%s changed malformed raw group", name)
		}
	}
}

func TestGroupedCatalogue(t *testing.T) {
	// TS 29.244 V20.1.0 Table 8.1.2-1: definitions in clause 7 are grouped.
	c, err := spec.Load()
	if err != nil {
		t.Fatal(err)
	}
	grouped := map[ie.IEType]bool{}
	for _, entry := range c.IETypes {
		if entry.Kind == "grouped" {
			grouped[ie.IEType(entry.Type)] = true
		}
	}
	for n := 1; n <= 65535; n++ {
		typ := ie.IEType(n)
		if ie.New(typ, nil).IsGrouped() != grouped[typ] {
			t.Errorf("wrong grouped status for type %d", n)
		}
	}
}
