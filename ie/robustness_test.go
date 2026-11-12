package ie_test

import (
	"bytes"
	"errors"
	"math/rand/v2"
	"net/netip"
	"reflect"
	"strings"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func payloadFamily() [][]byte {
	out := [][]byte{nil, {}}
	r := rand.New(rand.NewPCG(17, 91))
	for n := 1; n <= 64; n++ {
		b := make([]byte, n)
		for j := range b {
			b[j] = byte(r.Uint32())
		}
		out = append(out, make([]byte, n), bytes.Repeat([]byte{255}, n), b)
	}
	return out
}

func namedTypes() []ie.IEType {
	out := []ie.IEType{0x8001}
	for n := 1; n <= 65535; n++ {
		typ := ie.IEType(n)
		if !strings.HasPrefix(typ.String(), "IEType(") {
			out = append(out, typ)
		}
	}
	return out
}

func panicOf(fn func()) (p any) { defer func() { p = recover() }(); fn(); return nil }

// TestAllHelpersNoPanic calls every exported no-argument method of *ie.IE on
// IEs of every named type with every payload of the family. A method that
// rejects an IE type with InvalidTypeError does so before reading the payload,
// so the remaining payloads are skipped for that type and method; every other
// combination sees the whole family.
func TestAllHelpersNoPanic(t *testing.T) {
	methods := reflect.TypeFor[*ie.IE]()
	payloads := payloadFamily()
	calls := 0
	for _, typ := range namedTypes() {
		for m := 0; m < methods.NumMethod(); m++ {
			method := methods.Method(m)
			if method.Type.NumIn() != 1 {
				continue
			}
			for k, b := range payloads {
				var out []reflect.Value
				v := reflect.ValueOf(ie.New(typ, b))
				calls++
				if p := panicOf(func() { out = v.Method(m).Call(nil) }); p != nil {
					t.Errorf("%s/%s length %d: panic %v", typ, method.Name, len(b), p)
					break
				}
				if k == 0 && rejectsType(out) {
					break
				}
			}
		}
	}
	t.Logf("%d named/vendor types; %d exported methods; %d payloads; %d calls", len(namedTypes()), methods.NumMethod(), len(payloads), calls)
}

// rejectsType reports whether a method's results end in an *ie.InvalidTypeError.
func rejectsType(out []reflect.Value) bool {
	if len(out) == 0 {
		return false
	}
	last := out[len(out)-1]
	if last.Type() != reflect.TypeFor[error]() || last.IsNil() {
		return false
	}
	var invalid *ie.InvalidTypeError
	return errors.As(last.Interface().(error), &invalid)
}

func TestNilIEAPIs(t *testing.T) {
	var i *ie.IE
	v := reflect.ValueOf(i)
	for m := 0; m < v.NumMethod(); m++ {
		name := v.Type().Method(m).Name
		if name == "Equal" {
			// A nil IE equals another nil IE; TestIEEqual covers it.
			continue
		}
		t.Run(name, func(t *testing.T) {
			method := v.Method(m)
			args := make([]reflect.Value, method.Type().NumIn())
			for k := range args {
				args[k] = reflect.Zero(method.Type().In(k))
			}
			if p := panicOf(func() {
				var results []reflect.Value
				if method.Type().IsVariadic() {
					results = method.CallSlice(args)
				} else {
					results = method.Call(args)
				}
				for _, result := range results {
					if result.Type() == reflect.TypeFor[error]() {
						if result.IsNil() || !errors.Is(result.Interface().(error), ie.ErrIENotFound) {
							t.Errorf("nil receiver did not report ErrIENotFound")
						}
					} else if !result.IsZero() {
						t.Errorf("nil receiver returned a nonzero value")
					}
				}

			}); p != nil {
				t.Fatalf("panic %v", p)
			}
		})
	}
}

func TestFieldsNoPanic(t *testing.T) {
	for name, parse := range fieldParsers {
		t.Run(name, func(t *testing.T) {
			for _, b := range payloadFamily() {
				if p := panicOf(func() { reflect.ValueOf(parse).Call([]reflect.Value{reflect.ValueOf(b)}) }); p != nil {
					t.Fatalf("length %d: panic %v", len(b), p)
				}
			}
		})
	}
}

func TestNilGroupedChildren(t *testing.T) {
	for _, typ := range namedTypes() {
		i := ie.New(typ, nil)
		if !i.IsGrouped() {
			continue
		}
		i.ChildIEs = []*ie.IE{nil, ie.NewCause(1)}
		t.Run(typ.String(), func(t *testing.T) {
			v := reflect.ValueOf(i)
			for m := 0; m < v.NumMethod(); m++ {
				method := v.Method(m)
				if method.Type().NumIn() != 0 {
					continue
				}
				if p := panicOf(func() { method.Call(nil) }); p != nil {
					t.Errorf("%s: panic %v", v.Type().Method(m).Name, p)
				}
			}
		})
	}
}

func TestNilFieldsReceivers(t *testing.T) {
	for name, parse := range fieldParsers {
		t.Run(name, func(t *testing.T) {
			v := reflect.Zero(reflect.TypeOf(parse).Out(0))
			for m := 0; m < v.NumMethod(); m++ {
				method := v.Method(m)
				args := make([]reflect.Value, method.Type().NumIn())
				for n := range args {
					args[n] = reflect.Zero(method.Type().In(n))
				}
				if p := panicOf(func() { method.Call(args) }); p != nil {
					t.Errorf("%s panic %v", v.Type().Method(m).Name, p)
				}
			}
		})
	}
}

func TestFieldsReceiverReuse(t *testing.T) {
	for name, parse := range fieldParsers {
		t.Run(name, func(t *testing.T) {
			v := reflect.New(reflect.TypeOf(parse).Out(0).Elem())
			for _, payload := range payloadFamily() {
				result := reflect.ValueOf(parse).Call([]reflect.Value{reflect.ValueOf(payload)})
				err := v.Interface().(interface{ UnmarshalBinary([]byte) error }).UnmarshalBinary(payload)
				if !result[1].IsNil() {
					if err == nil {
						t.Fatal("reuse accepted a payload rejected by fresh decoding")
					}
					continue
				}
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(v.Interface(), result[0].Interface()) {
					t.Fatal("reused receiver retained fields from the previous input")
				}
			}
		})
	}
}

func TestIEEqual(t *testing.T) {
	var none *ie.IE
	decoded, err := ie.Parse([]byte{0, 1, 0, 6, 0, 56, 0, 2, 0, 1})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		name string
		a, b *ie.IE
		want bool
	}{
		{"nil nil", none, none, true},
		{"nil value", none, ie.NewPDRID(1), false},
		{"same value", ie.NewPDRID(1), ie.NewPDRID(1), true},
		{"different value", ie.NewPDRID(1), ie.NewPDRID(2), false},
		{"different type", ie.NewPDRID(1), ie.New(ie.Precedence, []byte{0, 1}), false},
		{"decoded and constructed group", decoded, ie.NewCreatePDR(ie.NewPDRID(1)), true},
		{"same invalid value", ie.NewNodeIDFQDN(""), ie.NewNodeIDFQDN(""), true},
		{"invalid and valid", ie.NewNodeIDFQDN(""), ie.NewNodeIDAddress(netip.MustParseAddr("192.0.2.1")), false},
	} {
		if got := c.a.Equal(c.b); got != c.want {
			t.Errorf("%s: Equal = %v, want %v", c.name, got, c.want)
		}
	}
}
