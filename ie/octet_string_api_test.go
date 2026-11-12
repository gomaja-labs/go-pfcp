package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// These identifiers have no specified character repertoire. OctetString does
// not imply text: preserve arbitrary octets, including NUL and invalid UTF-8.
// TS 29.244 V20.1.0 clauses 8.2.4, 8.2.6, 8.2.23, 8.2.60,
// 8.2.72, 8.2.73, 8.2.128 and 8.2.161.
func TestOpaqueIdentifierAPIs(t *testing.T) {
	for _, tc := range []struct {
		name      string
		ctor, get any
	}{
		{"NetworkInstance", ie.NewNetworkInstance, (*ie.IE).NetworkInstance},
		{"ApplicationID", ie.NewApplicationID, (*ie.IE).ApplicationID},
		{"ApplicationInstanceID", ie.NewApplicationInstanceID, (*ie.IE).ApplicationInstanceID},
		{"ForwardingPolicyIdentifier", ie.NewForwardingPolicy, (*ie.IE).ForwardingPolicyIdentifier},
		{"ActivatePredefinedRules", ie.NewActivatePredefinedRules, (*ie.IE).ActivatePredefinedRules},
		{"DeactivatePredefinedRules", ie.NewDeactivatePredefinedRules, (*ie.IE).DeactivatePredefinedRules},
		{"UEIPAddressPoolIdentity", ie.NewUEIPAddressPoolIdentity, (*ie.IE).UEIPAddressPoolIdentity},
		{"DataNetworkAccessIdentifier", ie.NewDataNetworkAccessIdentifier, (*ie.IE).DataNetworkAccessIdentifier},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctor, ok := tc.ctor.(func([]byte) *ie.IE)
			if !ok {
				t.Fatalf("opaque constructor takes %T, require []byte", tc.ctor)
			}
			get, ok := tc.get.(func(*ie.IE) ([]byte, error))
			if !ok {
				t.Fatalf("opaque getter has type %T, require []byte", tc.get)
			}
			for _, b := range [][]byte{{0}, {0xff, 0, 0x80, 0xc0, 0xaf}, {1, 2, 3}} {
				item := ctor(b)
				wire, err := item.Marshal()
				if err != nil {
					t.Fatal(err)
				}
				item, err = ie.Parse(wire)
				if err != nil {
					t.Fatal(err)
				}
				got, err := get(item)
				if err != nil {
					t.Fatal(err)
				}
				// The pool's existing value helper includes its explicit two-octet length.
				if tc.name == "UEIPAddressPoolIdentity" {
					got = got[2:]
				}
				if !bytes.Equal(got, b) {
					t.Error("opaque identifier changed")
				}
			}
			if _, err := get(nil); !errors.Is(err, ie.ErrIENotFound) {
				t.Fatal(err)
			}
			if tc.name != "ForwardingPolicyIdentifier" && tc.name != "UEIPAddressPoolIdentity" {
				item := ctor(nil)
				wire, err := item.Marshal()
				if err != nil || len(wire) != 4 {
					t.Fatalf("null: %v", err)
				}
				if _, err := get(item); !errors.Is(err, io.ErrUnexpectedEOF) {
					t.Fatalf("null value: %v", err)
				}
			}
		})
	}
}

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
