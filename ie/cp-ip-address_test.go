package ie_test

import (
	"net/netip"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestCPIPAddressCurrent(t *testing.T) {
	for flags := 1; flags < 4; flags++ {
		v := ie.CPIPAddressFields{Flags: ie.CPIPAddressFlags(flags)}
		b := []byte{byte(flags)}
		if flags&2 != 0 {
			v.IPv4Address = netip.MustParseAddr("192.0.2.1")
			b = append(b, 192, 0, 2, 1)
		}
		if flags&1 != 0 {
			v.IPv6Address = netip.MustParseAddr("2001:db8::1")
			b = append(b, v.IPv6Address.AsSlice()...)
		}
		assertAddedFields(t, ie.NewCPIPAddress(v), &v, b, "CPIPAddress", false)
	}

	var nilFields *ie.CPIPAddressFields
	assertAddedNilFields(t, nilFields)
}

func TestCPIPAddressInvalidFields(t *testing.T) {
	for _, v := range []ie.CPIPAddressFields{{}, {Flags: 4}, {Flags: ie.CPIPAddressV4}, {Flags: ie.CPIPAddressV6, IPv6Address: netip.MustParseAddr("192.0.2.1")}, {IPv4Address: netip.MustParseAddr("192.0.2.1")}} {
		assertAddedInvalid(t, ie.NewCPIPAddress(v))
		if _, err := v.Marshal(); err == nil {
			t.Fatal("invalid fields accepted")
		}
	}
}

var _ = registerConstructorFixture(ie.NewCPIPAddress, ie.CPIPAddressFields{Flags: 2, IPv4Address: netip.MustParseAddr("192.0.2.1")})
