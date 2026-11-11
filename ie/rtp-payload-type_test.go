package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TestRTPPayloadType checks TS 29.244 V20.1.0 clause 8.2.236.
func TestRTPPayloadType(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n)
		valid := v <= 127
		i := ie.NewRTPPayloadType(v)
		if !valid {
			assertNewIEInvalid(t, i, ie.RTPPayloadType)
			continue
		}
		payload := make([]byte, 1)
		payload[0] = v
		checkNewIEWire(t, i, payload)
		got, err := i.RTPPayloadType()
		if err != nil || got != v {
			t.Fatalf("value %d: got %d, error %v", v, got, err)
		}
	}
	getter := func(i *ie.IE) (any, error) { return i.RTPPayloadType() }
	checkNewIEFailures(t, ie.RTPPayloadType, getter, 1)
	base := ie.NewRTPPayloadType(0)
	extension := append(append([]byte{}, base.Payload...), 0xa5)
	_, err := getter(ie.New(ie.RTPPayloadType, extension))
	if (err != nil) != true {
		t.Fatalf("extension: %v", err)
	}
}

// Shared assertions for the IEs defined in clauses 8.2.236–286.
func checkNewIEWire(t *testing.T, i *ie.IE, payload []byte) {
	t.Helper()
	if i == nil {
		t.Fatal("nil constructor")
	}
	wire, err := i.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	want := append([]byte{byte(i.Type >> 8), byte(i.Type), byte(len(payload) >> 8), byte(len(payload))}, payload...)
	if !bytes.Equal(wire, want) {
		t.Fatal("wire differs from clause figure")
	}
	parsed, err := ie.Parse(want)
	if err != nil {
		t.Fatal(err)
	}
	again, err := parsed.Marshal()
	if err != nil || !bytes.Equal(again, want) {
		t.Fatal("wire round trip", err)
	}

	method := i.Type.String()
	before := reflect.ValueOf(i).MethodByName(method).Call(nil)
	after := reflect.ValueOf(parsed).MethodByName(method).Call(nil)
	if !before[1].IsNil() && len(payload) > 0 {
		t.Fatalf("constructor value helper failed: %v", before[1].Interface())
	}
	if before[1].IsNil() != after[1].IsNil() || !reflect.DeepEqual(before[0].Interface(), after[0].Interface()) {
		t.Fatal("decoded constructor inputs changed")
	}
	if err = i.MarshalTo(make([]byte, len(wire)-1)); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatal("short buffer", err)
	}
}
func assertNewIEInvalid(t *testing.T, i *ie.IE, typ ie.IEType) {
	t.Helper()
	if i == nil {
		t.Fatal("invalid constructor returned nil")
	}
	_, err := i.Marshal()
	var enc *ie.EncodeError
	if !errors.Is(err, ie.ErrInvalidValue) || !errors.As(err, &enc) || enc.Type != typ {
		t.Fatal("invalid value not reported", err)
	}
}
func checkNewIEFailures(t *testing.T, typ ie.IEType, get func(*ie.IE) (any, error), minimum int) {
	t.Helper()
	if _, err := get(nil); !errors.Is(err, ie.ErrIENotFound) {
		t.Fatal("nil receiver", err)
	}
	var wrong *ie.InvalidTypeError
	if _, err := get(ie.New(ie.Cause, nil)); !errors.As(err, &wrong) {
		t.Fatal("wrong type", err)
	}
	for n := 0; n < minimum; n++ {
		i := ie.New(typ, make([]byte, n))
		if _, err := get(i); err == nil {
			t.Fatalf("accepted length %d", n)
		} else {
			var dec *ie.DecodeError
			if !errors.As(err, &dec) || dec.Type != typ {
				t.Fatal("missing typed decode context", err)
			}
		}
		if n == 0 {
			b, err := i.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if _, err = ie.Parse(b); err != nil {
				t.Fatal("null IE decoding", err)
			}
		}
	}
}

type newFields interface {
	Marshal() ([]byte, error)
	MarshalTo([]byte) error
	MarshalLen() int
	UnmarshalBinary([]byte) error
}

func checkNewFields(t *testing.T, f newFields, parse func([]byte) (newFields, error)) {
	t.Helper()
	b, err := f.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	got, err := parse(b)
	if err != nil || !reflect.DeepEqual(f, got) {
		t.Fatal("Fields round trip", err)
	}
	for n := 0; n < len(b); n++ {
		if _, err = parse(b[:n]); err == nil {
			t.Fatalf("accepted Fields truncation at %d", n)
		}
		if err = f.MarshalTo(make([]byte, n)); !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatal("Fields short buffer", err)
		}
	}
	got, err = parse(append(append([]byte{}, b...), 0xa5))
	if err != nil || !reflect.DeepEqual(f, got) {
		t.Fatal("Fields extension", err)
	}
	if err = got.UnmarshalBinary(b); err != nil || !reflect.DeepEqual(f, got) {
		t.Fatal("Fields reuse", err)
	}
}
func checkNewFieldsNil(t *testing.T, f newFields) {
	t.Helper()
	if _, err := f.Marshal(); !errors.Is(err, ie.ErrIENotFound) {
		t.Fatal(err)
	}
	if err := f.MarshalTo(nil); !errors.Is(err, ie.ErrIENotFound) {
		t.Fatal(err)
	}
	if err := f.UnmarshalBinary(nil); !errors.Is(err, ie.ErrIENotFound) {
		t.Fatal(err)
	}
	if f.MarshalLen() != 0 {
		t.Fatal("nil Fields length")
	}
}

func TestRTPPayloadTypeInvalidWire(t *testing.T) {
	for n := 0; n <= 255; n++ {
		v := uint8(n) & 255
		valid := v <= 127
		if valid {
			continue
		}
		_, err := ie.New(ie.RTPPayloadType, []byte{byte(n)}).RTPPayloadType()
		if err == nil {
			t.Fatalf("invalid wire value %d accepted", n)
		}
	}
}
