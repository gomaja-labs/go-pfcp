package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestOffendingIEInformationCurrent(t *testing.T) {
	for _, v := range []ie.OffendingIEInformationFields{{Type: ie.Cause}, {Type: ie.IEType(65535), Value: []byte{0, 255, 7}}} {
		b := []byte{byte(v.Type >> 8), byte(v.Type)}
		b = append(b, v.Value...)
		assertAddedFields(t, ie.NewOffendingIEInformation(v), &v, b, "OffendingIEInformation", false)
	}

	var nilFields *ie.OffendingIEInformationFields
	assertAddedNilFields(t, nilFields)
}

func addedGetter(i *ie.IE, name string) (any, error) {
	out := reflect.ValueOf(i).MethodByName(name).Call(nil)
	if !out[1].IsNil() {
		return nil, out[1].Interface().(error)
	}
	return out[0].Interface(), nil
}
func assertAddedInvalid(t *testing.T, i *ie.IE) {
	t.Helper()
	if i == nil {
		t.Fatal("constructor returned nil")
	}
	for _, marshal := range []func() error{func() error { _, e := i.Marshal(); return e }, func() error { return i.MarshalTo(make([]byte, i.MarshalLen())) }, func() error { _, e := ie.NewGroupedIE(ie.CreatePDR, i).Marshal(); return e }} {
		err := marshal()
		var enc *ie.EncodeError
		if !errors.Is(err, ie.ErrInvalidValue) || !errors.As(err, &enc) || enc.Type != i.Type {
			t.Fatalf("invalid argument did not propagate EncodeError: %v", err)
		}
	}
}
func assertAddedNull(t *testing.T, i *ie.IE, name string) {
	t.Helper()
	b, e := i.Marshal()
	if e != nil {
		t.Fatal(e)
	}
	p, e := ie.Parse(b)
	if e != nil || len(p.Payload) != 0 {
		t.Fatalf("null IE: %v", e)
	}
	if _, e = addedGetter(p, name); !errors.Is(e, io.ErrUnexpectedEOF) {
		t.Fatalf("null helper: %v", e)
	}
}
func assertAddedVector(t *testing.T, i *ie.IE, payload []byte, name string, fixed bool) {
	t.Helper()
	if i == nil {
		t.Fatal("nil constructor")
	}
	if !bytes.Equal(i.Payload, payload) {
		t.Fatalf("%s: payload differs from clause vector", name)
	}
	want, err := addedGetter(i, name)
	if err != nil {
		t.Fatalf("%s: constructor getter: %v", name, err)
	}
	wire, err := i.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if len(wire) != len(payload)+4 || !bytes.Equal(wire[4:], payload) {
		t.Fatal("wire differs from vector")
	}
	p, err := ie.Parse(wire)
	if err != nil {
		t.Fatal(err)
	}
	got, err := addedGetter(p, name)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("helper round trip: %v", err)
	}
	again, err := p.Marshal()
	if err != nil || !bytes.Equal(again, wire) {
		t.Fatal("wire round trip")
	}
	for _, n := range []int{0, 1, 3, i.MarshalLen() - 1} {
		if !errors.Is(i.MarshalTo(make([]byte, n)), ie.ErrInvalidLength) {
			t.Fatalf("short buffer %d", n)
		}
	}
	if _, err = addedGetter(nil, name); !errors.Is(err, ie.ErrIENotFound) {
		t.Fatalf("nil helper: %v", err)
	}
	if _, err = addedGetter(ie.New(0, nil), name); err == nil {
		t.Fatal("wrong type accepted")
	}
	assertAddedNull(t, ie.New(i.Type, nil), name)
	minimum := len(payload)
	variable := false
	switch i.Type {
	case ie.CallingNumber, ie.CalledNumber, ie.GroupID, ie.EventNotificationURI, ie.NotificationCorrelationID, ie.PredefinedRulesName, ie.TunnelPassword, ie.Metadata, ie.TLContainer:
		minimum = 1
		variable = true
	case ie.OffendingIEInformation:
		minimum = 2
		variable = true
	case ie.VendorSpecificNodeReportType:
		minimum = 3
		variable = true
	}
	for n := 0; n < minimum; n++ {
		if _, err = addedGetter(ie.New(i.Type, payload[:n]), name); err == nil {
			t.Fatalf("%s: truncated payload %d accepted", name, n)
		}
	}
	ext := ie.New(i.Type, append(append([]byte(nil), payload...), 0x65))
	extra, err := addedGetter(ext, name)
	if fixed || i.Type == ie.CPIPAddress {
		if !errors.Is(err, ie.ErrInvalidLength) {
			t.Fatalf("fixed IE accepts trailing octets: %v", err)
		}
	} else if !variable {
		if err != nil || !reflect.DeepEqual(extra, want) {
			t.Fatalf("extension changed known value: %v", err)
		}
	}
}

type addedFields interface {
	Marshal() ([]byte, error)
	MarshalTo([]byte) error
	MarshalLen() int
	UnmarshalBinary([]byte) error
}

func assertAddedNilFields(t *testing.T, f addedFields) {
	t.Helper()
	if f.MarshalLen() != 0 {
		t.Fatal("nil fields length")
	}
	if _, e := f.Marshal(); !errors.Is(e, ie.ErrIENotFound) {
		t.Fatalf("nil Marshal: %v", e)
	}
	if e := f.MarshalTo(nil); !errors.Is(e, ie.ErrIENotFound) {
		t.Fatalf("nil MarshalTo: %v", e)
	}
	if e := f.UnmarshalBinary(nil); !errors.Is(e, ie.ErrIENotFound) {
		t.Fatalf("nil UnmarshalBinary: %v", e)
	}
}
func assertAddedFields(t *testing.T, i *ie.IE, f addedFields, payload []byte, name string, fixed bool) {
	t.Helper()
	assertAddedVector(t, i, payload, name, fixed)
	var spare uint8
	switch i.Type {
	case ie.Thresholds, ie.CPIPAddress, ie.MPQUICAddressInformation:
		spare = 0xfc
	case ie.MBSSessionIdentifier, ie.LocalIngressTunnel:
		spare = 0xf8
	case ie.TrafficParameterThreshold, ie.N6JitterMeasurement:
		spare = 0xfe
	case ie.ReportingSuggestionInfo:
		spare = 0xf0
	case ie.MulticastTransportInformation:
		spare = 0xff
	}
	for bits := 1; bits <= 255; bits++ {
		if uint8(bits)&^spare != 0 {
			continue
		}
		variant := bytes.Clone(payload)
		variant[0] |= uint8(bits)
		got, err := addedGetter(ie.New(i.Type, variant), name)
		if err != nil || !reflect.DeepEqual(got, f) {
			t.Fatalf("%s: spare bits %d changed value: %v", name, bits, err)
		}
		encoded, err := ie.New(i.Type, variant).Marshal()
		if err != nil || !bytes.Equal(encoded[4:], variant) {
			t.Fatal("spare bits not retained on IE")
		}
	}

	b, e := f.Marshal()
	if e != nil || !bytes.Equal(b, payload) || f.MarshalLen() != len(payload) {
		t.Fatalf("fields encoding %s: %v", name, e)
	}
	got, e := addedGetter(i, name)
	if e != nil || !reflect.DeepEqual(got, f) {
		t.Fatalf("%s: constructor inputs changed: %v", name, e)
	}
	for n := 0; n < min(len(payload), 64); n++ {
		if !errors.Is(f.MarshalTo(make([]byte, n)), io.ErrUnexpectedEOF) {
			t.Fatalf("fields short buffer %d", n)
		}
	}
	dst := bytes.Repeat([]byte{0xa5}, len(payload)+3)
	if e = f.MarshalTo(dst); e != nil || !bytes.Equal(dst[:len(payload)], payload) || !bytes.Equal(dst[len(payload):], []byte{0xa5, 0xa5, 0xa5}) {
		t.Fatal("fields MarshalTo")
	}
	reused := reflect.New(reflect.TypeOf(f).Elem()).Interface().(addedFields)
	if e = reused.UnmarshalBinary(payload); e != nil {
		t.Fatal(e)
	}
	if e = reused.UnmarshalBinary(payload); e != nil || !reflect.DeepEqual(reused, f) {
		t.Fatal("reused fields")
	}
}

func TestOffendingIEInformationLengthBoundary(t *testing.T) {
	v := ie.OffendingIEInformationFields{Type: ie.Cause, Value: bytes.Repeat([]byte{1}, 65533)}
	payload := append([]byte{0, 19}, v.Value...)
	assertAddedFields(t, ie.NewOffendingIEInformation(v), &v, payload, "OffendingIEInformation", false)
	v.Value = append(v.Value, 1)
	assertAddedInvalid(t, ie.NewOffendingIEInformation(v))
}
