package ie_test

import (
	"net/netip"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestLocalIngressTunnelCurrent(t *testing.T) {
	for _, flags := range []uint8{1, 2, 5, 6} {
		for _, port := range []uint16{0, 65535} {
			v := ie.LocalIngressTunnelFields{Flags: ie.LocalIngressTunnelFlags(flags)}
			b := []byte{flags}
			if flags&4 == 0 {
				v.Port = port
				b = append(b, byte(port>>8), byte(port))
				s := "192.0.2.1"
				if flags&2 != 0 {
					s = "2001:db8::1"
				}
				v.Address = netip.MustParseAddr(s)
				b = append(b, v.Address.AsSlice()...)
			}
			assertAddedFields(t, ie.NewLocalIngressTunnel(v), &v, b, "LocalIngressTunnel", false)
		}
	}

	var nilFields *ie.LocalIngressTunnelFields
	assertAddedNilFields(t, nilFields)
}

func TestLocalIngressTunnelInvalidFields(t *testing.T) {
	for _, v := range []ie.LocalIngressTunnelFields{{}, {Flags: 3}, {Flags: 4}, {Flags: 7}, {Flags: 8}, {Flags: ie.LocalIngressV4}, {Flags: ie.LocalIngressV4 | ie.LocalIngressCH, Port: 1}, {Flags: ie.LocalIngressV4 | ie.LocalIngressCH, Address: netip.MustParseAddr("192.0.2.1")}} {
		assertAddedInvalid(t, ie.NewLocalIngressTunnel(v))
		if _, err := v.Marshal(); err == nil {
			t.Fatal("invalid fields accepted")
		}
	}
}
