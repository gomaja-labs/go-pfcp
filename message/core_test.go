package message_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/gomaja-labs/go-pfcp/message"
)

func TestVersionBeforeFraming(t *testing.T) {
	parsers := []func([]byte) error{
		func(b []byte) error { _, e := message.Parse(b); return e },
		func(b []byte) error { _, e := message.ParseDatagram(b); return e },
		func(b []byte) error { _, e := message.ParseHeader(b); return e },
		func(b []byte) error { _, e := message.ParseGeneric(b); return e },
	}
	for _, m := range allMessages() {
		parsers = append(parsers, m.UnmarshalBinary)
	}
	for version := 0; version < 8; version++ {
		if version == 1 {
			continue
		}
		for _, s := range []byte{0, 1} {
			wire := wireHeader(byte(version<<5)|s, 54, nil)
			for _, length := range []uint16{0, 4, 12, 200} {
				wire[2], wire[3] = byte(length>>8), byte(length)
				for n := 1; n <= len(wire); n++ {
					for p, parse := range parsers {
						err := parse(wire[:n])
						var v *message.UnsupportedVersionError
						if !errors.As(err, &v) || v.Header == nil || v.Header.Version() != version {
							t.Fatalf("version=%d S=%d length=%d prefix=%d parser=%d: %v", version, s, length, n, p, err)
						}
						if n >= 2 && v.Header.Type != 54 {
							t.Fatal("lost received type")
						}
					}
				}
			}
		}
	}
	for _, parse := range parsers {
		if !errors.Is(parse(nil), io.ErrUnexpectedEOF) {
			t.Fatal("empty input must be EOF")
		}
	}
}

func TestSingleEncodingClearsFO(t *testing.T) {
	for _, typ := range []uint8{1, 54, 200} {
		h := message.NewHeader(typ, 42, 1, nil)
		h.Flags |= 4
		original := h.Flags
		for _, encode := range []func() ([]byte, error){h.Marshal, func() ([]byte, error) { b := make([]byte, h.MarshalLen()); e := h.MarshalTo(b); return b, e }} {
			b, e := encode()
			if e != nil || b[0]&4 != 0 {
				t.Fatal("standalone header FO", e)
			}
			if _, e = message.Parse(b); e != nil {
				t.Fatal(e)
			}
			if h.Flags != original {
				t.Fatal("header encoding mutated flags")
			}
		}
	}

	for _, m := range allMessages() {
		if _, ok := reflect.TypeOf(m).MethodByName("SetFO"); ok {
			t.Fatalf("%T still exposes SetFO", m)
		}
		h := reflect.ValueOf(m).Elem().FieldByName("Header").Interface().(*message.Header)
		h.Flags |= 4
		original := h.Flags
		for _, encode := range []func() ([]byte, error){
			func() ([]byte, error) { b := make([]byte, m.MarshalLen()); e := m.MarshalTo(b); return b, e },
			m.(interface{ Marshal() ([]byte, error) }).Marshal,
		} {
			b, e := encode()
			if e != nil {
				t.Fatalf("%T: %v", m, e)
			}
			if b[0]&4 != 0 {
				t.Fatalf("%T: standalone FO set", m)
			}
			if _, e = message.Parse(b); e != nil {
				t.Fatal(e)
			}
			if h.Flags != original {
				t.Fatal("encoding changed source flags")
			}
		}
		b, e := message.MarshalDatagram(m)
		if e != nil || b[0]&4 != 0 {
			t.Fatalf("single datagram %T: %v", m, e)
		}
	}
}

func TestDatagramBodyFailuresContinue(t *testing.T) {
	for count := 1; count <= 8; count++ {
		for mask := 0; mask < 1<<count; mask++ {
			var wire []byte
			for n := 0; n < count; n++ {
				body := []byte(nil)
				if mask&(1<<n) != 0 {
					body = []byte{0, 19, 0, 2, 1}
				}
				frame := wireHeader(0x21, 54, body)
				frame[14] = byte(n)
				if n+1 < count {
					frame[0] |= 4
				}
				wire = append(wire, frame...)
			}
			got, err := message.ParseDatagram(wire)
			want := 0
			for n := 0; n < count; n++ {
				if mask&(1<<n) == 0 {
					if want >= len(got) || got[want].Sequence() != uint32(n) {
						t.Fatalf("count %d mask %d: lost frame %d (%d returned), err %v", count, mask, n, len(got), err)
					}
					want++
				}
			}
			if len(got) != want || (err == nil) != (mask == 0) {
				t.Fatalf("count=%d mask=%d: messages=%d err=%v", count, mask, len(got), err)
			}
			var failures []*message.DatagramError
			if err != nil {
				joined, ok := err.(interface{ Unwrap() []error })
				if !ok {
					t.Fatal("errors were not joined")
				}
				for _, cause := range joined.Unwrap() {
					var f *message.DatagramError
					if !errors.As(cause, &f) {
						t.Fatal("missing frame context")
					}
					failures = append(failures, f)
				}
			}
			offset, failed := 0, 0
			for n := 0; n < count; n++ {
				if mask&(1<<n) != 0 {
					if failed >= len(failures) {
						t.Fatal("missing frame failure")
					}
					f := failures[failed]
					failed++
					if f.Index != n || f.Offset != offset || f.Header == nil || f.Header.Sequence() != uint32(n) {
						t.Fatalf("bad frame context: %+v", f)
					}
					offset += 21
				} else {
					offset += 16
				}
			}
			if len(failures) != failed {
				t.Fatal("extra frame failures")
			}
			if mask != 0 && (!errors.Is(err, ie.ErrInvalidLength) || !errors.Is(err, io.ErrUnexpectedEOF)) {
				t.Fatalf("missing IE framing causes: %v", err)
			}
		}
	}
}

func TestDatagramFramingFailureKeepsPrefix(t *testing.T) {
	good := wireHeader(0x25, 54, nil)
	for _, bad := range [][]byte{wireHeader(0x25, 54, nil), {0x21}, wireHeader(0x21, 54, []byte{0})} {
		wire := append(bytes.Clone(good), bad...)
		if len(bad) > 4 && bad[0]&4 == 0 {
			wire[len(good)+3]++
		}
		got, err := message.ParseDatagram(wire)
		if err == nil || len(got) != 1 || !got[0].(interface{ HasFO() bool }).HasFO() {
			t.Fatalf("lost valid prefix: %d %v", len(got), err)
		}
	}
}

func TestGenericCopiesAndFiltersIEs(t *testing.T) {
	a, b, c := ie.NewCause(1), ie.NewCause(2), ie.NewCause(3)
	input := make([]*ie.IE, 2, 4)
	input[0] = a
	m := message.NewGeneric(200, 1, 1, input...)
	if len(m.OtherIEs) != 1 || m.OtherIEs[0] != a {
		t.Fatal("nil IE retained")
	}
	input[0] = b
	if m.OtherIEs[0] != a {
		t.Fatal("constructor retained caller slice")
	}
	m.OtherIEs = make([]*ie.IE, 1, 4)
	m.OtherIEs[0] = a
	cp := *m
	m.AddIE(b, nil)
	cp.AddIE(c)
	if len(m.OtherIEs) != 2 || m.OtherIEs[1] != b || len(cp.OtherIEs) != 2 || cp.OtherIEs[1] != c {
		t.Fatal("AddIE mutated a shared backing array or kept nil")
	}
}

func TestDatagramUnsupportedVersionStopsAtPrefix(t *testing.T) {
	wire := wireHeader(0x25, 54, nil)
	wire = append(wire, 0x40)
	got, err := message.ParseDatagram(wire)
	var f *message.DatagramError
	var v *message.UnsupportedVersionError
	if len(got) != 1 || !errors.As(err, &f) || f.Index != 1 || f.Offset != 16 || !errors.As(err, &v) || v.Header.Version() != 2 || f.Header != v.Header {
		t.Fatalf("unsupported-version context: %d %v", len(got), err)
	}
}
