package ie_test

import (
	"net/netip"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestDNSServerAddressCurrent(t *testing.T) {
	for _, v := range []string{"0.0.0.0", "192.0.2.1", "255.255.255.255"} {
		a := netip.MustParseAddr(v)
		assertAddedVector(t, ie.NewDNSServerAddress(a), a.AsSlice(), "DNSServerAddress", false)
		got, err := ie.NewDNSServerAddress(a).DNSServerAddress()
		if err != nil || got != a {
			t.Fatalf("address %v: %v %v", a, got, err)
		}
	}
	for _, a := range []netip.Addr{{}, netip.MustParseAddr("2001:db8::1"), netip.MustParseAddr("::ffff:192.0.2.1")} {
		assertAddedInvalid(t, ie.NewDNSServerAddress(a))
	}
}
