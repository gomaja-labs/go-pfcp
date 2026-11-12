package ie_test

import (
	"errors"
	"net/netip"
	"strings"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestNodeIDBounds(t *testing.T) {
	for _, b := range [][]byte{{0, 1}, {1, 1}, {2, 1, 'a', 5, 'b'}, {2, 64, 'a'}} {
		if _, err := ie.New(ie.NodeID, b).NodeID(); err == nil {
			t.Error("accepted malformed node ID")
		}
	}
	long := strings.Repeat("x", 64)
	for name, i := range map[string]*ie.IE{
		"overlong label":  ie.NewNodeIDFQDN(long),
		"no value":        ie.NewNodeIDFQDN(""),
		"invalid address": ie.NewNodeIDAddress(netip.Addr{}),
	} {
		requireInvalidValue(t, name, ie.NodeID, i)
	}
}

// requireInvalidValue checks that a constructor given an argument it cannot
// encode returns an IE of the right type whose Marshal reports ErrInvalidValue,
// and that a grouped IE containing it fails to encode too instead of silently
// dropping it.
func requireInvalidValue(t *testing.T, name string, typ ie.IEType, i *ie.IE) {
	t.Helper()
	if i == nil || i.Type != typ {
		t.Fatalf("%s: got %v, want an IE of type %v", name, i, typ)
	}
	_, err := i.Marshal()
	var ee *ie.EncodeError
	if !errors.Is(err, ie.ErrInvalidValue) || !errors.As(err, &ee) || ee.IEType() != typ {
		t.Errorf("%s: Marshal returned %v, want an EncodeError for %v wrapping ErrInvalidValue", name, err, typ)
	}
	if _, err := ie.NewCreatePDR(ie.NewPDRID(1), i).Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Errorf("%s: grouped parent Marshal returned %v, want ErrInvalidValue", name, err)
	}
}
