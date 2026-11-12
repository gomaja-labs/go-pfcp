package ie_test

import (
	"errors"
	"io"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestDecodeErrorContext(t *testing.T) {
	for _, tc := range []struct {
		name string
		wire []byte
		typ  ie.IEType
		want error
	}{
		{"short_header", []byte{0, 19, 0}, 0, io.ErrUnexpectedEOF},
		{"short_value", []byte{0, 19, 0, 1}, ie.Cause, io.ErrUnexpectedEOF},
		{"vendor_length", []byte{128, 1, 0, 1, 0}, 0x8001, ie.ErrInvalidLength},
		{"grouped_child", []byte{0, 1, 0, 4, 0, 19, 0, 1}, ie.Cause, io.ErrUnexpectedEOF},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ie.Parse(tc.wire)
			if !errors.Is(err, tc.want) {
				t.Fatal("lost underlying error", err)
			}
			// The interface allows this regression to compile on revisions
			// without typed decoding context.
			var context interface{ IEType() ie.IEType }
			found := errors.As(err, &context)
			if tc.typ == 0 {
				if found {
					t.Fatal("incomplete header reported an IE type")
				}
				return
			}
			if !found || context.IEType() != tc.typ {
				t.Fatal("missing offending IE type", err)
			}
		})
	}
}
