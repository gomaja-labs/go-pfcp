package ie_test

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand/v2"
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

func TestAllHelpersNoPanic(t *testing.T) {
	methods := reflect.TypeFor[*ie.IE]()
	failures := map[string]bool{}
	for _, typ := range namedTypes() {
		for _, b := range payloadFamily() {
			v := reflect.ValueOf(ie.New(typ, b))
			for m := 0; m < methods.NumMethod(); m++ {
				method := methods.Method(m)
				if method.Type.NumIn() != 1 {
					continue
				}
				key := fmt.Sprintf("%s/%s", typ, method.Name)
				if failures[key] {
					continue
				}
				if p := panicOf(func() { v.Method(m).Call(nil) }); p != nil {
					failures[key] = true
					t.Errorf("%s length %d: panic %v", key, len(b), p)
				}
			}
		}
	}
	t.Logf("%d named/vendor types; %d exported methods; %d payloads", len(namedTypes()), methods.NumMethod(), len(payloadFamily()))
}

func TestNilIEAPIs(t *testing.T) {
	var i *ie.IE
	v := reflect.ValueOf(i)
	for m := 0; m < v.NumMethod(); m++ {
		name := v.Type().Method(m).Name
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
