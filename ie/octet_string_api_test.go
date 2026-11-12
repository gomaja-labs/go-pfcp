package ie_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// These identifiers have no specified character repertoire. OctetString does
// not imply text: preserve arbitrary octets, including NUL and invalid UTF-8.
// TS 29.244 V20.1.0 clauses 8.2.4, 8.2.6, 8.2.23, 8.2.60,
// 8.2.72, 8.2.73, 8.2.128 and 8.2.161.
// TS 29.244 V20.1.0 clauses 8.2.12 and 8.2.66 encode markings
// as two octets in wire order, rather than an unsigned integer.
func TestBinaryMarkingAPIs(t *testing.T) {
	ctor, ok := any(ie.NewTransportLevelMarking).(func([]byte) *ie.IE)
	if !ok {
		t.Fatal("Transport Level Marking requires a []byte API")
	}
	get, ok := any((*ie.IE).TransportLevelMarking).(func(*ie.IE) ([]byte, error))
	if !ok {
		t.Fatal("Transport Level Marking getter requires []byte")
	}
	for _, b := range [][]byte{{0, 255}, {255, 0}, {128, 252}} {
		got, err := get(ctor(b))
		if err != nil || !bytes.Equal(got, b) {
			t.Fatalf("marking: %v", err)
		}
	}
	for _, n := range []int{1, 3} {
		assertInvalidConstructor(t, ctor(make([]byte, n)))
	}
	typ := reflect.TypeFor[ie.DLFlowLevelMarkingFields]()
	for _, name := range []string{"ToSTrafficClass", "ServiceClassIdentifier"} {
		f, _ := typ.FieldByName(name)
		if f.Type != reflect.TypeFor[[]byte]() {
			t.Errorf("%s requires []byte", name)
		}
	}
}

func TestDLBinaryMarkings(t *testing.T) {
	ctor, ok := any(ie.NewDLFlowLevelMarking).(func(uint8, []byte, []byte) *ie.IE)
	if !ok {
		t.Fatal("DL markings require binary constructor arguments")
	}
	fields, ok := any(ie.NewDLFlowLevelMarkingFields).(func(uint8, []byte, []byte) *ie.DLFlowLevelMarkingFields)
	if !ok {
		t.Fatal("DL marking Fields require binary constructor arguments")
	}
	for _, tc := range []struct {
		flags          uint8
		ttc, sci, want []byte
	}{
		{0, nil, nil, []byte{0}}, {1, []byte{0x12, 0x34}, nil, []byte{1, 0x12, 0x34}},
		{2, nil, []byte{0x56, 0x78}, []byte{2, 0x56, 0x78}},
		{3, []byte{0x12, 0x34}, []byte{0x56, 0x78}, []byte{3, 0x12, 0x34, 0x56, 0x78}},
	} {
		item := ctor(tc.flags, tc.ttc, tc.sci)
		if !bytes.Equal(item.Payload, tc.want) {
			t.Error("binary marking octet order changed")
		}
		f, err := item.DLFlowLevelMarking()
		if err != nil {
			t.Fatal(err)
		}
		v := reflect.ValueOf(f).Elem()
		for name, want := range map[string][]byte{"ToSTrafficClass": tc.ttc, "ServiceClassIdentifier": tc.sci} {
			if !reflect.DeepEqual(v.FieldByName(name).Interface(), want) {
				t.Errorf("%s changed", name)
			}
		}
	}
	for _, tc := range []struct {
		flags    uint8
		ttc, sci []byte
	}{
		{1, nil, nil}, {2, nil, nil}, {1, []byte{1}, nil}, {2, nil, []byte{1}},
		{1, []byte{1, 2, 3}, nil}, {2, nil, []byte{1, 2, 3}},
		{0, []byte{1, 2}, nil}, {0, nil, []byte{1, 2}}, {0x80, nil, nil},
	} {
		assertInvalidConstructor(t, ctor(tc.flags, tc.ttc, tc.sci))
		f := fields(tc.flags, tc.ttc, tc.sci)
		if err := f.MarshalTo(make([]byte, f.MarshalLen())); !errors.Is(err, ie.ErrInvalidValue) {
			t.Fatalf("invalid Fields accepted: %v", err)
		}
	}
}
