package message

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"sync"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/gomaja-labs/go-pfcp/internal/spec"
)

// The test oracle reads the committed tables, independently of pfcpgen's field
// reduction and codec dispatch. Null-length IEs are valid (clause 8.1.2).
func TestMessagesCatalogue(t *testing.T) {
	c, err := spec.Load()
	if err != nil {
		t.Fatal(err)
	}
	kinds := map[int]string{}
	for _, e := range c.IETypes {
		kinds[e.Type] = e.Kind
	}
	factories := catalogueMessages()
	if len(factories) != len(c.Messages) {
		t.Fatal("message coverage")
	}
	for index, desc := range c.Messages {
		t.Run(desc.GoName, func(t *testing.T) {
			empty := factories[index]
			if int(empty.MessageType()) != desc.Type || reflect.TypeOf(empty).Elem().Name() != desc.GoName {
				t.Fatal("constructor/type mismatch")
			}
			if empty.IsRequest() != bytes.HasSuffix([]byte(desc.GoName), []byte("Request")) {
				t.Fatal("request classification")
			}
			if empty.MessageTypeName() != MessageTypeName(uint8(desc.Type)) {
				t.Fatal("message name")
			}
			type entry struct {
				r      spec.Row
				list   bool
				values []*ie.IE
			}
			entries := []entry{}
			byType := map[int]int{}
			for _, r := range desc.IEs {
				if n, ok := byType[r.IEType]; ok {
					entries[n].list = true
					continue
				}
				byType[r.IEType] = len(entries)
				entries = append(entries, entry{r: r, list: r.Multiple || r.SharedTypeGroup != nil})
			}
			value := reflect.ValueOf(empty).Elem()
			if value.NumField() != len(entries)+2 || value.Type().Field(value.NumField()-1).Name != "OtherIEs" {
				t.Fatal("wrong field count/order")
			}
			var input []*ie.IE
			expectedOthers := []*ie.IE{}
			for i := range entries {
				e := &entries[i]
				field := value.FieldByName(e.r.GoField)
				wantType := reflect.TypeOf((*ie.IE)(nil))
				if e.list {
					wantType = reflect.TypeOf([]*ie.IE{})
				}
				if !field.IsValid() || field.Type() != wantType {
					t.Fatalf("field %s has wrong cardinality", e.r.GoField)
				}
				count := 2
				if e.list {
					count = 3
				}
				for n := 0; n < count; n++ {
					marker := byte(i*3 + n + 1)
					var v *ie.IE
					if kinds[e.r.IEType] == "grouped" {
						v = ie.NewGroupedIE(ie.IEType(e.r.IEType), ie.NewApplicationID(string([]byte{'x', marker})))
					} else {
						v = ie.New(ie.IEType(e.r.IEType), []byte{marker})
					}
					e.values = append(e.values, v)
					input = append(input, nil, v, nil)
					if !e.list && n > 0 {
						expectedOthers = append(expectedOthers, v)
					}
				}
			}
			unknown := ie.New(0x7ffe, []byte("unknown"))
			input = append(input, unknown)
			expectedOthers = append(expectedOthers, unknown)
			m := catalogueMessages(input...)[index]
			v := reflect.ValueOf(m).Elem()
			for _, e := range entries {
				field := v.FieldByName(e.r.GoField)
				if e.list {
					got := field.Interface().([]*ie.IE)
					if !reflect.DeepEqual(got, e.values) {
						t.Fatal("lost list entries", e.r.GoField)
					}
				} else if field.Interface().(*ie.IE) != e.values[0] {
					t.Fatal("first occurrence not selected", e.r.GoField)
				}
			}
			if !reflect.DeepEqual(v.FieldByName("OtherIEs").Interface(), expectedOthers) {
				t.Fatal("duplicates/unknowns not retained")
			}
			out := encodeTestMessage(t, m)
			// New messages use unique table-field order, then OtherIEs.
			want := []*ie.IE{}
			for _, e := range entries {
				if e.list {
					want = append(want, e.values...)
				} else {
					want = append(want, e.values[0])
				}
			}
			want = append(want, expectedOthers...)
			gotOrder := testIEOrder(t, out)
			if !reflect.DeepEqual(gotOrder, ieTypes(want)) {
				t.Fatal("constructor ignored table order")
			}
			decoded, err := Parse(out)
			if err != nil {
				t.Fatal(err)
			}
			if reflect.TypeOf(decoded) != reflect.TypeOf(m) {
				t.Fatal("wrong concrete message dispatch")
			}
			if !bytes.Equal(encodeTestMessage(t, decoded), out) {
				t.Fatal("table-order round trip")
			}
			// Build an independently ordered frame. Sorting on encode must fail here.
			var payload []byte
			for n := len(input) - 1; n >= 0; n-- {
				if input[n] == nil {
					continue
				}
				b, err := input[n].Marshal()
				if err != nil {
					t.Fatal(err)
				}
				payload = append(payload, b...)
			}
			flags := byte(0x20)
			headerSize := 8
			if desc.Category == "session" {
				flags = 0x21
				headerSize = 16
			}
			wire := make([]byte, headerSize+len(payload))
			wire[0], wire[1] = flags, byte(desc.Type)
			binary.BigEndian.PutUint16(wire[2:4], uint16(len(wire)-4))
			wire[headerSize-2] = 2
			copy(wire[headerSize:], payload)
			decoded, err = Parse(wire)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(encodeTestMessage(t, decoded), wire) {
				t.Fatal("non-table-order bytes changed")
			}
			dv := reflect.ValueOf(decoded).Elem()
			for _, e := range entries {
				if e.list {
					if dv.FieldByName(e.r.GoField).Len() != len(e.values) {
						t.Fatal("wire list lost")
					}
				} else {
					got := dv.FieldByName(e.r.GoField).Interface().(*ie.IE)
					if !bytes.Equal(encodeTestIE(t, got), encodeTestIE(t, e.values[len(e.values)-1])) {
						t.Fatal("wire first occurrence lost", e.r.GoField)
					}
				}
			}
			// Reusing a receiver clears every old field/list/unknown IE.
			clean := encodeTestMessage(t, empty)
			if err := decoded.UnmarshalBinary(clean); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(encodeTestMessage(t, decoded), clean) {
				t.Fatal("receiver reuse retained fields")
			}
			// Every individual field accepts a null-length IE; nil fields/lists vanish.
			nulls := []*ie.IE{}
			for _, e := range entries {
				nulls = append(nulls, ie.New(ie.IEType(e.r.IEType), nil))
			}
			nullWire := encodeTestMessage(t, catalogueMessages(nulls...)[index])
			parsed, err := Parse(nullWire)
			if err != nil {
				t.Fatal("null IE rejected", err)
			}
			if !bytes.Equal(encodeTestMessage(t, parsed), nullWire) {
				t.Fatal("null IE round trip")
			}
			for size := 0; size < m.MarshalLen(); size++ {
				if err := m.MarshalTo(make([]byte, size)); !errors.Is(err, io.ErrShortBuffer) {
					t.Fatalf("short buffer %d: %v", size, err)
				}
			}
		})
	}
}
func encodeTestMessage(t *testing.T, m Message) []byte {
	t.Helper()
	b := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(b); err != nil {
		t.Fatal(err)
	}
	return b
}
func encodeTestIE(t *testing.T, i *ie.IE) []byte {
	t.Helper()
	b, err := i.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	return b
}
func ieTypes(ies []*ie.IE) []uint16 {
	out := []uint16{}
	for _, i := range ies {
		out = append(out, uint16(i.Type))
	}
	return out
}
func testIEOrder(t *testing.T, b []byte) []uint16 {
	t.Helper()
	h, err := ParseHeader(b)
	if err != nil {
		t.Fatal(err)
	}
	out := []uint16{}
	for p := h.Payload; len(p) > 0; {
		if len(p) < 4 {
			t.Fatal("short IE header")
		}
		n := int(binary.BigEndian.Uint16(p[2:4])) + 4
		if n > len(p) {
			t.Fatal("short IE payload")
		}
		out = append(out, binary.BigEndian.Uint16(p[:2]))
		p = p[n:]
	}
	return out
}

func TestEveryMessageConcurrentMarshal(t *testing.T) {
	c, err := spec.Load()
	if err != nil {
		t.Fatal(err)
	}
	for n, m := range catalogueMessages() {
		t.Run(c.Messages[n].GoName, func(t *testing.T) {
			// A complete opaque TLV is valid for null-length value semantics, and an
			// unknown IE exercises payload encoding even in the header-only response.
			ies := []*ie.IE{ie.New(0x7ffe, []byte("payload"))}
			for _, r := range c.Messages[n].IEs {
				ies = append(ies, ie.New(ie.IEType(r.IEType), nil))
			}
			m = catalogueMessages(ies...)[n]
			want := encodeTestMessage(t, m)
			parsed, err := Parse(want)
			if err != nil {
				t.Fatal(err)
			}
			for _, current := range []Message{m, parsed} {
				h := reflect.ValueOf(current).Elem().FieldByName("Header").Interface().(*Header)
				before := *h
				before.Payload = bytes.Clone(h.Payload)
				var wg sync.WaitGroup
				for range 8 {
					wg.Go(func() {
						for range 20 {
							b := make([]byte, current.MarshalLen())
							if err := current.MarshalTo(b); err != nil || !bytes.Equal(b, want) {
								t.Error("concurrent encoding differs", err)
							}
						}
					})
				}
				wg.Wait()
				if !reflect.DeepEqual(*h, before) {
					t.Fatal("encoding changed header")
				}
			}
		})
	}
}

// TS 29.244 V20.1.0 clause 8.2.82 is void; clause 7.6.9 preserves unknown IEs.
func TestReservedType116(t *testing.T) {
	for _, payload := range [][]byte{nil, {0xff, 0x01, 0x02}} {
		for _, m := range catalogueMessages(ie.New(116, payload)) {
			t.Run(m.MessageTypeName(), func(t *testing.T) {
				b := encodeTestMessage(t, m)
				got, err := Parse(b)
				if err != nil {
					t.Fatal(err)
				}
				field := reflect.ValueOf(got).Elem().FieldByName("OtherIEs")
				others, ok := field.Interface().([]*ie.IE)
				if !ok || len(others) != 1 || others[0].Type != 116 || !bytes.Equal(others[0].Payload, payload) {
					t.Fatal("reserved IE not retained as unknown")
				}
				after := encodeTestMessage(t, got)
				if !bytes.Equal(b, after) {
					t.Fatal("unknown IE bytes changed")
				}
			})
		}
	}
	if ie.IEType(116).String() != "IEType(116)" {
		t.Fatal("reserved type has an assigned name")
	}
}
