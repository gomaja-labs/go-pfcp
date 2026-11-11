package ie_test

import (
	"net/netip"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestIPAddressAndPortNumberReplacementCurrent(t *testing.T) {
	for flags := 0; flags < 128; flags++ {
		v := replacementFields(flags)
		if flags&64 != 0 && flags != 64 || flags&64 == 0 && flags&27 == 0 {
			assertAddedInvalid(t, ie.NewIPAddressAndPortNumberReplacement(v))
			continue
		}
		b := []byte{byte(flags)}
		if flags&1 != 0 {
			b = append(b, v.DestinationIPv4.AsSlice()...)
		}
		if flags&2 != 0 {
			b = append(b, v.DestinationIPv6.AsSlice()...)
		}
		if flags&4 != 0 {
			b = append(b, 0xff, 0xff)
		}
		if flags&8 != 0 {
			b = append(b, v.SourceIPv4.AsSlice()...)
		}
		if flags&16 != 0 {
			b = append(b, v.SourceIPv6.AsSlice()...)
		}
		if flags&32 != 0 {
			b = append(b, 0, 0)
		}
		assertAddedFields(t, ie.NewIPAddressAndPortNumberReplacement(v), &v, b, "IPAddressAndPortNumberReplacement", false)
		spare := append([]byte(nil), b...)
		spare[0] |= 128
		got, e := ie.New(ie.IPAddressAndPortNumberReplacement, spare).IPAddressAndPortNumberReplacement()
		if e != nil || got.Flags != v.Flags {
			t.Fatal("spare flag changed value", e)
		}

	}
	var f *ie.IPAddressAndPortNumberReplacementFields
	assertAddedNilFields(t, f)
}
func replacementFields(flags int) ie.IPAddressAndPortNumberReplacementFields {
	v := ie.IPAddressAndPortNumberReplacementFields{Flags: ie.IPReplacementFlags(flags)}
	if flags&1 != 0 {
		v.DestinationIPv4 = netip.MustParseAddr("192.0.2.1")
	}
	if flags&2 != 0 {
		v.DestinationIPv6 = netip.MustParseAddr("2001:db8::1")
	}
	if flags&4 != 0 {
		v.DestinationPort = 65535
	}
	if flags&8 != 0 {
		v.SourceIPv4 = netip.MustParseAddr("198.51.100.2")
	}
	if flags&16 != 0 {
		v.SourceIPv6 = netip.MustParseAddr("2001:db8::2")
	}
	return v
}
func TestIPAddressAndPortNumberReplacementValidation(t *testing.T) {
	for _, v := range []ie.IPAddressAndPortNumberReplacementFields{{Flags: 128}, {Flags: 1}, {Flags: 2, DestinationIPv6: netip.MustParseAddr("192.0.2.1")}, {Flags: 64, SourcePort: 1}, {Flags: 64, DestinationIPv4: netip.MustParseAddr("192.0.2.1")}} {
		assertAddedInvalid(t, ie.NewIPAddressAndPortNumberReplacement(v))
	}
	// Receiver interprets present fields rather than applying sender policy.
	for _, b := range [][]byte{{0}, {4, 0, 0}, {65, 192, 0, 2, 1}} {
		if _, e := ie.New(ie.IPAddressAndPortNumberReplacement, b).IPAddressAndPortNumberReplacement(); e != nil {
			t.Fatal(e)
		}
	}
}

var _ = registerConstructorFixture(ie.NewIPAddressAndPortNumberReplacement, ie.IPAddressAndPortNumberReplacementFields{Flags: ie.IPReplacementUseMappedN6})
