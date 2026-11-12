package ie_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestTruncatedIEUnknownType(t *testing.T) {
	for _, b := range [][]byte{nil, {0}, {0, 0}, {0, 19}, {0, 19, 0}} {
		_, err := ie.Parse(b)
		var de *ie.DecodeError
		if !errors.Is(err, ie.ErrInvalidLength) || !errors.Is(err, io.ErrUnexpectedEOF) || !errors.As(err, &de) {
			t.Fatalf("length=%d: %v", len(b), err)
		}
		if de.TypeUnknown != (len(b) < 2) {
			t.Fatalf("length=%d: TypeUnknown=%v", len(b), de.TypeUnknown)
		}
		if len(b) < 2 {
			if !strings.Contains(err.Error(), "type unknown") || strings.Contains(err.Error(), "decode IE 0") {
				t.Fatalf("fabricated type: %v", err)
			}
		} else if strings.Contains(err.Error(), "type unknown") {
			t.Fatalf("lost known type: %v", err)
		}
	}
	// A partial child's type remains unknown even inside a known group.
	_, err := ie.Parse([]byte{0, 1, 0, 1, 0})
	var de *ie.DecodeError
	if !errors.As(err, &de) || !de.TypeUnknown || len(de.Path) != 1 || de.Path[0] != ie.CreatePDR || !strings.Contains(err.Error(), "type unknown") {
		t.Fatalf("group context: %v", err)
	}
}
