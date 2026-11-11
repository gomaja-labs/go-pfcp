package ie_test

import (
	"net/netip"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestMBSSessionIdentifierCurrent(t *testing.T) {
	for flags := 0; flags < 8; flags++ {
		v := ie.MBSSessionIdentifierFields{Flags: ie.MBSSessionIdentifierFlags(flags)}
		b := []byte{byte(flags)}
		if flags&1 != 0 {
			v.TMGI = [6]byte{1, 2, 3, 4, 5, 6}
			b = append(b, 1, 2, 3, 4, 5, 6)
		}
		if flags&2 != 0 {
			v.MulticastAddress = netip.MustParseAddr("232.0.2.1")
			v.SourceAddress = netip.MustParseAddr("192.0.2.1")
			b = append(b, 4, 232, 0, 2, 1, 4, 192, 0, 2, 1)
		}
		if flags&4 != 0 {
			v.NID = 0x123456789ab
			b = append(b, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xb0)
		}
		assertAddedFields(t, ie.NewMBSSessionIdentifier(v), &v, b, "MBSSessionIdentifier", false)
	}

	var nilFields *ie.MBSSessionIdentifierFields
	assertAddedNilFields(t, nilFields)
}

func TestMBSSessionIdentifierInvalidFields(t *testing.T) {
	for _, v := range []ie.MBSSessionIdentifierFields{{Flags: 8}, {Flags: ie.MBSSessionNIDI, NID: 0x100000000000}, {TMGI: [6]byte{1}}, {NID: 1}, {Flags: ie.MBSSessionSSMI}, {MulticastAddress: netip.MustParseAddr("232.0.2.1")}} {
		assertAddedInvalid(t, ie.NewMBSSessionIdentifier(v))
		if _, err := v.Marshal(); err == nil {
			t.Fatal("invalid fields accepted")
		}
	}
}
