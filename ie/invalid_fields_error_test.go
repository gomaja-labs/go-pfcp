package ie_test

import (
	"errors"
	"net"
	"strings"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// A constructor whose fields fail to encode reports ErrInvalidValue once, and
// keeps any other cause reachable.
func TestInvalidFieldsErrorReportsInvalidValueOnce(t *testing.T) {
	for _, tc := range []struct {
		name  string
		ie    *ie.IE
		cause error
	}{
		{"FSEID_IPv6_as_IPv4", ie.NewFSEID(1, net.ParseIP("2001:db8::1"), nil), nil},
		{"FTEID_no_family", ie.NewFTEID(0, 1, nil, nil, 0), nil},
		{"UEIPAddress_IPv6D_with_IP6PL", ie.NewUEIPAddress(0x49, "", "2001:db8::1", 4, 72), nil},
		{"OuterHeaderCreation_no_description", ie.NewOuterHeaderCreation(0, 0, "", "", 0, 0, 0), nil},
		{"SourceIPAddress_IPv4_as_IPv6", ie.NewSourceIPAddress(nil, net.ParseIP("192.0.2.1"), 0), nil},
		{"HeaderEnrichment_too_long", ie.NewHeaderEnrichment(0, strings.Repeat("a", 70000), "v"), ie.ErrInvalidLength},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.ie.Marshal()
			var encodeErr *ie.EncodeError
			if !errors.As(err, &encodeErr) || !errors.Is(err, ie.ErrInvalidValue) {
				t.Fatalf("Marshal error = %v, want EncodeError wrapping ErrInvalidValue", err)
			}
			if n := strings.Count(err.Error(), ie.ErrInvalidValue.Error()); n != 1 {
				t.Fatalf("%q states %q %d times", err, ie.ErrInvalidValue, n)
			}
			if tc.cause != nil && !errors.Is(err, tc.cause) {
				t.Fatalf("Marshal error = %v, want it to wrap %v", err, tc.cause)
			}
		})
	}
}
