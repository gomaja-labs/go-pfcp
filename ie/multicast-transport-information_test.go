package ie_test

import (
	"net/netip"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestMulticastTransportInformationCurrent(t *testing.T) {
	for _, dst := range []string{"232.0.2.1", "ff3e::1"} {
		for _, src := range []string{"192.0.2.1", "2001:db8::1"} {
			v := ie.MulticastTransportInformationFields{CommonTEID: 0x12345678, MulticastAddress: netip.MustParseAddr(dst), SourceAddress: netip.MustParseAddr(src)}
			b := []byte{0, 0x12, 0x34, 0x56, 0x78}
			for _, a := range []netip.Addr{v.MulticastAddress, v.SourceAddress} {
				tag := byte(4)
				if a.Is6() {
					tag = 0x50
				}
				b = append(b, tag)
				b = append(b, a.AsSlice()...)
			}
			assertAddedFields(t, ie.NewMulticastTransportInformation(v), &v, b, "MulticastTransportInformation", false)
		}
	}

	var nilFields *ie.MulticastTransportInformationFields
	assertAddedNilFields(t, nilFields)
}

func TestMulticastTransportInformationInvalidFields(t *testing.T) {
	for _, v := range []ie.MulticastTransportInformationFields{{}, {MulticastAddress: netip.MustParseAddr("192.0.2.2"), SourceAddress: netip.MustParseAddr("192.0.2.1")}, {MulticastAddress: netip.MustParseAddr("232.0.2.2"), SourceAddress: netip.MustParseAddr("232.0.2.1")}} {
		assertAddedInvalid(t, ie.NewMulticastTransportInformation(v))
		if _, err := v.Marshal(); err == nil {
			t.Fatal("invalid fields accepted")
		}
	}
}
