package ie_test

import (
	"errors"
	"fmt"
	"io"
	"net/netip"
	"reflect"
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type helperSemanticCase struct {
	name  string
	input *ie.IE
	want  any
	// minimum includes the fields selected by the fixture's flags. Zero means
	// no nonempty prefix is necessarily malformed (text/octet-string views).
	minimum int
	generic bool
}

func invokeSemanticHelper(i *ie.IE, name string) (any, error) {
	out := reflect.ValueOf(i).MethodByName(name).Call(nil)
	if len(out) < 2 {
		panic("unsupported semantic helper signature: " + name)
	}
	last := out[len(out)-1]
	value := out[0].Interface()
	if len(out) > 2 {
		values := make([]any, len(out)-1)
		for n := range values {
			values[n] = out[n].Interface()
		}
		value = values
	}
	if !last.IsNil() {
		return value, last.Interface().(error)
	}
	return value, nil
}

func equalSemanticValue(got, want any) bool {
	if g, ok := got.(time.Time); ok {
		w, ok := want.(time.Time)
		return ok && g.Equal(w)
	}
	g, w := reflect.ValueOf(got), reflect.ValueOf(want)
	if g.IsValid() && w.IsValid() && w.Type().ConvertibleTo(g.Type()) {
		want = w.Convert(g.Type()).Interface()
	}
	return cmp.Equal(got, want, cmpopts.IgnoreUnexported(ie.PacketRateStatusFields{}), cmpopts.EquateComparable(netip.Addr{}))
}

func TestHelperSemanticValues(t *testing.T) {
	for _, tc := range helperSemanticValues() {
		t.Run(tc.name, func(t *testing.T) {
			got, err := invokeSemanticHelper(tc.input, tc.name)
			if err != nil || !equalSemanticValue(got, tc.want) {
				t.Fatalf("value=%v error=%v; want %v", got, err, tc.want)
			}
		})
	}
}

// TS 29.244 V20.1.0 clause 8.1.2 permits null IEs for removal. A
// value helper must distinguish that absence from a provisioned zero value.
func TestHelperSemanticContracts(t *testing.T) {
	for _, tc := range helperSemanticValues() {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("nil", func(t *testing.T) {
				_, err := invokeSemanticHelper(nil, tc.name)
				if !errors.Is(err, ie.ErrIENotFound) {
					t.Fatalf("error=%v; want ErrIENotFound", err)
				}
			})
			t.Run("wrong-type", func(t *testing.T) {
				wrongInput := ie.New(0xffff, tc.input.Payload)
				if tc.generic {
					wrongInput = ie.NewCreatePDR()
				}
				_, err := invokeSemanticHelper(wrongInput, tc.name)
				var wrong *ie.InvalidTypeError
				if !errors.As(err, &wrong) {
					t.Fatalf("error=%v; want InvalidTypeError", err)
				}
			})
			t.Run("null", func(t *testing.T) {
				_, err := invokeSemanticHelper(ie.New(tc.input.Type, nil), tc.name)
				if err == nil {
					t.Fatal("accepted null value")
				}
			})
			for n := 1; n < tc.minimum && n < len(tc.input.Payload); n++ {
				t.Run(fmt.Sprintf("truncated-%d", n), func(t *testing.T) {
					_, err := invokeSemanticHelper(ie.New(tc.input.Type, tc.input.Payload[:n]), tc.name)
					if !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, ie.ErrInvalidLength) && !errors.Is(err, ie.ErrMalformed) {
						t.Fatalf("error=%v; want length error", err)
					}
				})
			}
		})
	}
}
