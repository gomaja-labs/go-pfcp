package ie_test

import (
	"bytes"
	"errors"
	"io"
	"net"
	"reflect"
	"sync"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/gomaja-labs/go-pfcp/internal/spec"
)

type groupedFields interface {
	IE() *ie.IE
	MarshalLen() int
	MarshalTo([]byte) error
	UnmarshalBinary([]byte) error
}
type groupedCase struct {
	typ    ie.IEType
	newIE  func(...*ie.IE) *ie.IE
	empty  func() groupedFields
	access func(*ie.IE) (groupedFields, error)
	parse  func([]byte) (groupedFields, error)
}
type childField struct {
	name  string
	typ   ie.IEType
	slice bool
}

func groupedSchema(t *testing.T) map[ie.IEType][]childField {
	t.Helper()
	c, err := spec.Load()
	if err != nil {
		t.Fatal(err)
	}
	schemas := map[ie.IEType][]childField{}
	for _, g := range c.Grouped {
		// Count repetitions within a table, not between contexts of a union.
		counts := map[int]int{}
		for _, r := range g.IEs {
			counts[r.IEType]++
		}
		typ := ie.IEType(g.Type)
		for _, r := range g.IEs {
			list := r.Multiple || r.SharedTypeGroup != nil || counts[r.IEType] > 1
			found := false
			for j := range schemas[typ] {
				f := &schemas[typ][j]
				if int(f.typ) == r.IEType {
					if f.name != r.GoField {
						t.Fatalf("inconsistent union field %d", r.IEType)
					}
					f.slice = f.slice || list
					found = true
					break
				}
			}
			if !found {
				schemas[typ] = append(schemas[typ], childField{r.GoField, ie.IEType(r.IEType), list})
			}
		}
	}
	// The two Session Report tables impose compatible constraints; their union
	// is SRR ID, Access Availability Report, QoS Monitoring Report, traffic report.
	fields := schemas[ie.SessionReport]
	schemas[ie.SessionReport] = []childField{fields[0], fields[2], fields[1], fields[3]}
	return schemas
}
func groupedChild(typ ie.IEType, n byte) *ie.IE {
	if (&ie.IE{Type: typ}).IsGrouped() {
		return ie.NewGroupedIE(typ, ie.New(409, []byte{n}))
	}
	return ie.New(typ, []byte{n})
}
func fieldValue(t *testing.T, f groupedFields, field childField) reflect.Value {
	t.Helper()
	v := reflect.ValueOf(f).Elem().FieldByName(field.name)
	if !v.IsValid() {
		t.Fatalf("missing field %s", field.name)
	}
	want := reflect.TypeFor[*ie.IE]()
	if field.slice {
		want = reflect.SliceOf(want)
	}
	if v.Type() != want {
		t.Fatalf("%s type %v want %v", field.name, v.Type(), want)
	}
	return v
}
func valueBytes(t *testing.T, f groupedFields) []byte {
	t.Helper()
	b := make([]byte, f.MarshalLen())
	if err := f.MarshalTo(b); err != nil {
		t.Fatal(err)
	}
	return b
}
func TestGroupedFieldsCatalogue(t *testing.T) {
	schemas := groupedSchema(t)
	cases := catalogueGrouped()
	if len(cases) != 121 || len(schemas) != 121 {
		t.Fatal("grouped census", len(cases), len(schemas))
	}
	seen := map[ie.IEType]bool{}
	for _, tc := range cases {
		t.Run(tc.typ.String(), func(t *testing.T) {
			if seen[tc.typ] {
				t.Fatal("duplicate grouped API")
			}
			seen[tc.typ] = true
			schema := schemas[tc.typ]
			if schema == nil {
				t.Fatal("missing schema")
			}
			if reflect.ValueOf(tc.empty()).Elem().NumField() != len(schema)+1 {
				t.Fatal("wrong field count")
			}
			var children, canonical, others []*ie.IE
			expected := map[string][]*ie.IE{}
			for _, field := range schema {
				first, second := groupedChild(field.typ, 1), groupedChild(field.typ, 2)
				children = append(children, first, nil, second)
				expected[field.name] = []*ie.IE{first}
				canonical = append(canonical, first)
				if field.slice {
					children = append(children, groupedChild(field.typ, 3))
					expected[field.name] = append(expected[field.name], second, children[len(children)-1])
					canonical = append(canonical, expected[field.name][1:]...)
				} else {
					others = append(others, second)
				}
			}
			unknown := ie.New(409, []byte{7})
			children = append(children, nil, unknown)
			others = append(others, unknown)
			canonical = append(canonical, others...)
			parent := tc.newIE(children...)
			if !parent.IsGrouped() || parent.Type != tc.typ || len(parent.ChildIEs) != len(children)-len(schema)-1 {
				t.Fatal("constructor lost children or retained nil")
			}
			f, err := tc.access(parent)
			if err != nil {
				t.Fatal(err)
			}
			for _, field := range schema {
				v := fieldValue(t, f, field)
				if field.slice {
					if !reflect.DeepEqual(v.Interface(), expected[field.name]) {
						t.Fatal("list lost a child", field.name)
					}
				} else if v.Interface() != expected[field.name][0] {
					t.Fatal("first occurrence lost", field.name)
				}
			}
			if !reflect.DeepEqual(reflect.ValueOf(f).Elem().FieldByName("OtherIEs").Interface(), others) {
				t.Fatal("OtherIEs lost duplicates or unknowns")
			}
			if !reflect.DeepEqual(f.IE().ChildIEs, canonical) {
				t.Fatal("table order or IE symmetry")
			}
			// Direct struct construction, including nil list entries, exercises IE().
			direct := tc.empty()
			for _, field := range schema {
				v := fieldValue(t, direct, field)
				if field.slice {
					v.Set(reflect.ValueOf(append([]*ie.IE{nil}, expected[field.name]...)))
				} else {
					v.Set(reflect.ValueOf(expected[field.name][0]))
				}
			}
			reflect.ValueOf(direct).Elem().FieldByName("OtherIEs").Set(reflect.ValueOf(append([]*ie.IE{nil}, others...)))
			if !reflect.DeepEqual(direct.IE().ChildIEs, canonical) {
				t.Fatal("struct fields not encoded")
			}
			b := valueBytes(t, f)
			if !bytes.Equal(b, valueBytes(t, direct)) {
				t.Fatal("struct and accessor differ")
			}
			parsed, err := tc.parse(b)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(b, valueBytes(t, parsed)) {
				t.Fatal("value round trip")
			}
			// Receiver reuse clears fields, lists and OtherIEs.
			if err := parsed.UnmarshalBinary(nil); err != nil {
				t.Fatal(err)
			}
			if parsed.MarshalLen() != 0 || len(parsed.IE().ChildIEs) != 0 {
				t.Fatal("receiver reuse retained children")
			}
			// Raw IE decoding retains non-table order byte for byte; typed views use table order.
			reversed := append([]*ie.IE{}, canonical...)
			for a, z := 0, len(reversed)-1; a < z; a, z = a+1, z-1 {
				reversed[a], reversed[z] = reversed[z], reversed[a]
			}
			raw := tc.newIE(reversed...)
			wire, err := raw.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := ie.Parse(wire)
			if err != nil {
				t.Fatal(err)
			}
			again, err := decoded.Marshal()
			if err != nil || !bytes.Equal(wire, again) {
				t.Fatal("raw group changed order", err)
			}
			view, err := tc.access(decoded)
			if err != nil {
				t.Fatal(err)
			}
			fromValue, err := tc.parse(wire[4:])
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(valueBytes(t, view), valueBytes(t, fromValue)) {
				t.Fatal("non-table order decode differs")
			}
			for _, field := range schema {
				var want *ie.IE
				for _, child := range decoded.ChildIEs {
					if child.Type == field.typ {
						want = child
						break
					}
				}
				v := fieldValue(t, view, field)
				if field.slice {
					if v.Len() == 0 || v.Index(0).Interface() != want {
						t.Fatal("list receive order")
					}
				} else if v.Interface() != want {
					t.Fatal("non-table first occurrence")
				}
			}
			// Every short buffer is an error, with no panic.
			for n := 0; n < len(b); n++ {
				short := bytes.Repeat([]byte{0xee}, n)
				if err := f.MarshalTo(short); !errors.Is(err, ie.ErrInvalidLength) {
					t.Fatal("short buffer", n, err)
				}
				if !bytes.Equal(short, bytes.Repeat([]byte{0xee}, n)) {
					t.Fatal("short buffer modified", n)
				}
			}
			extra := bytes.Repeat([]byte{0xee}, len(b)+5)
			if err := f.MarshalTo(extra); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(extra[len(b):], bytes.Repeat([]byte{0xee}, 5)) {
				t.Fatal("overwrote trailing buffer")
			}
			// Preserve a constructor error through every grouped builder/encoder.
			invalid := ie.NewNetworkInstanceFQDN("bad..name")
			bad := tc.empty()
			reflect.ValueOf(bad).Elem().FieldByName("OtherIEs").Set(reflect.ValueOf([]*ie.IE{invalid}))
			for _, encode := range []func() error{func() error { _, err := tc.newIE(invalid).Marshal(); return err }, func() error { _, err := bad.IE().Marshal(); return err }, func() error { return bad.MarshalTo(make([]byte, bad.MarshalLen())) }} {
				var ee *ie.EncodeError
				if err := encode(); !errors.As(err, &ee) || ee.Type != ie.NetworkInstance {
					t.Fatal("lost child EncodeError", err)
				}
			}
		})
	}
}

func TestGroupedFieldsErrors(t *testing.T) {
	for _, tc := range catalogueGrouped() {
		t.Run(tc.typ.String(), func(t *testing.T) {
			if _, err := tc.access(nil); !errors.Is(err, ie.ErrIENotFound) {
				t.Fatal("nil IE", err)
			}
			var invalid *ie.InvalidTypeError
			if _, err := tc.access(ie.NewPDRID(1)); !errors.As(err, &invalid) || invalid.Type != ie.PDRID {
				t.Fatal("wrong type", err)
			}
			// Null-length groups and empty fields encode an empty value.
			for _, parent := range []*ie.IE{tc.newIE(nil), ie.New(tc.typ, nil), {Type: tc.typ, ChildIEs: []*ie.IE{nil}}} {
				f, err := tc.access(parent)
				if err != nil || f.MarshalLen() != 0 {
					t.Fatal("null group", err)
				}
			}
			empty := tc.empty()
			nilFields := reflect.Zero(reflect.TypeOf(empty)).Interface().(groupedFields)
			if nilFields.MarshalLen() != 0 || !errors.Is(nilFields.MarshalTo(nil), ie.ErrIENotFound) || !errors.Is(nilFields.UnmarshalBinary(nil), ie.ErrIENotFound) {
				t.Fatal("nil fields")
			}
			var ee *ie.EncodeError
			if _, err := nilFields.IE().Marshal(); !errors.As(err, &ee) || !errors.Is(err, ie.ErrIENotFound) {
				t.Fatal("nil fields IE", err)
			}
			for _, b := range [][]byte{{0}, {0, 1, 0}, {0, 19, 0, 1}, make([]byte, 65536)} {
				if _, err := tc.parse(b); err == nil {
					t.Fatal("malformed value accepted")
				}
				f := tc.empty()
				reflect.ValueOf(f).Elem().FieldByName("OtherIEs").Set(reflect.ValueOf([]*ie.IE{ie.NewPDRID(1)}))
				before := valueBytes(t, f)
				if err := f.UnmarshalBinary(b); err == nil || !bytes.Equal(before, valueBytes(t, f)) {
					t.Fatal("failed decode mutated receiver", err)
				}
			}
			// Truncated child reports the offending type through DecodeError.
			if _, err := tc.parse([]byte{0, 19, 0, 1}); !errors.Is(err, io.ErrUnexpectedEOF) {
				t.Fatal("truncation error", err)
			}
			raw := ie.New(tc.typ, []byte{0, 19, 0, 1})
			if _, err := tc.access(raw); err == nil {
				t.Fatal("raw malformed payload accepted")
			}
			tooBig := tc.empty()
			reflect.ValueOf(tooBig).Elem().FieldByName("OtherIEs").Set(reflect.ValueOf([]*ie.IE{ie.New(409, make([]byte, 65532))}))
			if err := tooBig.MarshalTo(make([]byte, tooBig.MarshalLen())); !errors.Is(err, ie.ErrInvalidLength) {
				t.Fatal("oversize grouped value", err)
			}
		})
	}
}

func TestGroupedFieldsConcurrentMarshal(t *testing.T) {
	for _, tc := range catalogueGrouped() {
		f, err := tc.access(tc.newIE(ie.NewPDRID(42), ie.New(409, []byte{1})))
		if err != nil {
			t.Fatal(err)
		}
		want := valueBytes(t, f)
		var wg sync.WaitGroup
		for range 4 {
			wg.Go(func() {
				for range 5 {
					b := make([]byte, f.MarshalLen())
					if err := f.MarshalTo(b); err != nil || !bytes.Equal(b, want) {
						t.Error("concurrent encoding", err)
					}
				}
			})
		}
		wg.Wait()
	}
}

func TestGroupedFieldsNesting(t *testing.T) {
	// The parent counts toward the core's sixteen enclosing grouped IEs.
	for _, tc := range catalogueGrouped() {
		for _, depth := range []int{15, 16} {
			child := ie.NewPDRID(1)
			for range depth {
				child = ie.NewCreatePDR(child)
			}
			b, err := child.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			for _, parse := range []func() error{
				func() error { _, err := tc.parse(b); return err },
				func() error { _, err := tc.access(ie.New(tc.typ, b)); return err },
			} {
				err := parse()
				if depth == 15 && err != nil {
					t.Fatal("valid nesting", tc.typ, err)
				}
				if depth == 16 && !errors.Is(err, ie.ErrNestingTooDeep) {
					t.Fatal("unbounded nesting", tc.typ, err)
				}
			}
		}
	}
}

func TestCreatedTrafficEndpointLocalFTEID(t *testing.T) {
	var empty *ie.CreatedTrafficEndpointFields
	if _, err := empty.LocalFTEID(); !errors.Is(err, ie.ErrIENotFound) {
		t.Fatal(err)
	}
	first := ie.NewFTEID(1, 1, net.ParseIP("192.0.2.1"), nil, 0)
	second := ie.NewFTEID(1, 2, net.ParseIP("192.0.2.2"), nil, 0)
	f := &ie.CreatedTrafficEndpointFields{FTEID: []*ie.IE{nil, first, nil, second}}
	for _, n := range []int{-1, 0, 3} {
		if _, err := f.LocalFTEIDN(n); !errors.Is(err, ie.ErrIENotFound) {
			t.Fatal("missing occurrence", n, err)
		}
	}
	for _, n := range []int{1, 2} {
		if value, err := f.LocalFTEIDN(n); err != nil || value.TEID != uint32(n) {
			t.Fatal(n, err)
		}
	}
	if _, err := f.LocalFTEID(); err != nil {
		t.Fatal(err)
	}
	f.FTEID = []*ie.IE{ie.New(ie.FTEID, nil)}
	if _, err := f.LocalFTEIDN(1); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatal("malformed F-TEID", err)
	}
}

func FuzzGroupedFields(f *testing.F) {
	cases := catalogueGrouped()
	for n := range cases {
		f.Add(uint8(n), []byte{0, 56, 0, 2, 0, 42})
		f.Add(uint8(n), []byte{})
	}
	f.Fuzz(func(t *testing.T, n uint8, b []byte) {
		tc := cases[int(n)%len(cases)]
		fields, err := tc.parse(b)
		if err != nil {
			return
		}
		encoded := valueBytes(t, fields)
		next, err := tc.parse(encoded)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(encoded, valueBytes(t, next)) {
			t.Fatal("canonical value changed")
		}
		wire, err := fields.IE().Marshal()
		if err != nil {
			t.Fatal(err)
		}
		parent, err := ie.Parse(wire)
		if err != nil {
			t.Fatal(err)
		}
		view, err := tc.access(parent)
		if err != nil || !bytes.Equal(encoded, valueBytes(t, view)) {
			t.Fatal("accessor differs", err)
		}
	})
}
