package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
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
	b := []byte{0, 19, 0, 1, 7, 99}
	i, err := ie.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	_ = append(i.Payload, 8)
	if b[5] != 99 {
		t.Fatal("append overwrote following IE")
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
	v, err := decoded.PDRID()
	if err != nil || v != 7 {
		t.Errorf("child lookup %d %v", v, err)
	}
	i = ie.NewCreatePDR(ie.NewPDRID(7), nil)
	if i == nil || i.Payload != nil || len(i.ChildIEs) != 1 {
		t.Fatal("constructor must collect children only")
	}
	i.ChildIEs[0].Payload[1] = 9
	v, err = i.PDRID()
	if err != nil || v != 9 {
		t.Errorf("lookup used stale payload: %d %v", v, err)
	}
	i = ie.New(ie.CreatePDR, raw)
	vIE, err := i.FindByType(ie.PDRID)
	if err != nil || vIE == nil {
		t.Fatal("raw group lookup", err)
	}
	i.Add(nil, ie.NewCause(4))
	i.Remove(ie.PDRID)
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
	grouped := map[ie.IEType]bool{
		1:   true,
		2:   true,
		3:   true,
		4:   true,
		5:   true,
		6:   true,
		7:   true,
		8:   true,
		9:   true,
		10:  true,
		11:  true,
		12:  true,
		13:  true,
		14:  true,
		15:  true,
		16:  true,
		17:  true,
		18:  true,
		51:  true,
		54:  true,
		58:  true,
		59:  true,
		68:  true,
		77:  true,
		78:  true,
		79:  true,
		80:  true,
		83:  true,
		85:  true,
		86:  true,
		87:  true,
		99:  true,
		102: true,
		105: true,
		118: true,
		127: true,
		128: true,
		129: true,
		130: true,
		132: true,
		143: true,
		147: true,
		165: true,
		166: true,
		167: true,
		168: true,
		169: true,
		175: true,
		176: true,
		183: true,
		187: true,
		188: true,
		189: true,
		190: true,
		195: true,
		199: true,
		200: true,
		201: true,
		203: true,
		205: true,
		211: true,
		212: true,
		213: true,
		214: true,
		216: true,
		218: true,
		220: true,
		221: true,
		225: true,
		226: true,
		227: true,
		233: true,
		238: true,
		239: true,
		240: true,
		242: true,
		247: true,
		252: true,
		254: true,
		255: true,
		256: true,
		261: true,
		263: true,
		264: true,
		267: true,
		270: true,
		271: true,
	}
	for n := 1; n <= 65535; n++ {
		typ := ie.IEType(n)
		if ie.New(typ, nil).IsGrouped() != grouped[typ] {
			t.Errorf("wrong grouped status for type %d", n)
		}
	}
}
