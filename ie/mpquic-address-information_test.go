package ie_test

import (
	"net/netip"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestMPQUICAddressInformationCurrent(t *testing.T) {
	for flags := 0; flags < 4; flags++ {
		for typ := 2; typ <= 5; typ++ {
			for _, port := range []uint16{0, 65535} {
				v := ie.MPQUICAddressInformationFields{Flags: ie.MPQUICAddressFlags(flags), ProxyType: ie.MPQUICProxyType(typ), Port: port}
				b := []byte{byte(flags), byte(typ), byte(port >> 8), byte(port)}
				if flags&1 != 0 {
					v.IPv4Address = netip.MustParseAddr("192.0.2.1")
					b = append(b, 192, 0, 2, 1)
				}
				if flags&2 != 0 {
					v.IPv6Address = netip.MustParseAddr("2001:db8::1")
					b = append(b, v.IPv6Address.AsSlice()...)
				}
				assertAddedFields(t, ie.NewMPQUICAddressInformation(v), &v, b, "MPQUICAddressInformation", false)
			}
		}
	}

	var nilFields *ie.MPQUICAddressInformationFields
	assertAddedNilFields(t, nilFields)
}

func TestMPQUICAddressInformationInvalidFields(t *testing.T) {
	for _, v := range []ie.MPQUICAddressInformationFields{{ProxyType: 0}, {ProxyType: 1}, {ProxyType: 6}, {ProxyType: 2, Flags: 4}, {ProxyType: 2, Flags: ie.MPQUICAddressV4}, {ProxyType: 2, Flags: ie.MPQUICAddressV6, IPv6Address: netip.MustParseAddr("192.0.2.1")}} {
		assertAddedInvalid(t, ie.NewMPQUICAddressInformation(v))
		if _, err := v.Marshal(); err == nil {
			t.Fatal("invalid fields accepted")
		}
	}
}
