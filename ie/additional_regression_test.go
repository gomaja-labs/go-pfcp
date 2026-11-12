package ie_test

import (
	"net"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestNFInstanceIDShortInput(t *testing.T) {
	for n := 0; n < 16; n++ {
		if p := panicOf(func() {
			if ie.NewNFInstanceID(make([]byte, n)) != nil {
				t.Error("accepted short UUID")
			}
		}); p != nil {
			t.Errorf("length %d panic: %v", n, p)
		}
	}
}

func TestUserIDEvenDigits(t *testing.T) {
	for _, digits := range []string{"12", "123456", "1234567"} {
		i := ie.NewUserID(7, digits, digits, digits, "")
		got, err := i.UserID()
		if err != nil {
			t.Fatal(err)
		}
		if got.IMSI != digits || got.IMEI != digits || got.MSISDN != digits {
			t.Error("lost subscriber identity digit")
		}
	}
}

func TestRedirectAddresses(t *testing.T) {
	for _, typ := range []uint8{ie.RedirectAddrIPv4, ie.RedirectAddrIPv4AndIPv6} {
		addrs := []string{"192.0.2.1"}
		if typ == ie.RedirectAddrIPv4AndIPv6 {
			addrs = append(addrs, "2001:db8::1")
		}
		i := ie.NewRedirectInformation(typ, addrs...)
		got, err := i.RedirectInformation()
		if err != nil {
			t.Fatal(err)
		}
		if got.RedirectServerAddress != "192.0.2.1" {
			t.Error("first address changed")
		}
		if typ == ie.RedirectAddrIPv4AndIPv6 && got.OtherRedirectServerAddress != "2001:db8::1" {
			t.Error("second address changed")
		}
		if typ == ie.RedirectAddrIPv4 && got.OtherRedirectServerAddress != "" {
			t.Error("unexpected second address")
		}
	}
}

func TestPMFOptionalAddresses(t *testing.T) {
	mac := net.HardwareAddr{2, 0, 0, 0, 0, 1}
	for flags := 0; flags < 8; flags++ {
		var v4, v6 net.IP
		var m1, m2 net.HardwareAddr
		if flags&1 != 0 {
			v4 = net.ParseIP("192.0.2.1")
		}
		if flags&2 != 0 {
			v6 = net.ParseIP("2001:db8::1")
		}
		if flags&4 != 0 {
			m1 = mac
			m2 = mac
		}
		i := ie.NewPMFAddressInformation(v4, v6, 8805, 8806, m1, m2)
		f, err := i.PMFAddressInformation()
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(f.PMFMACAddressFor3GPPAccess, m1) || !reflect.DeepEqual(f.PMFMACAddressForNon3GPPAccess, m2) {
			t.Error("MAC addresses lost")
		}
		if v4 != nil && !f.PMFIPv4Address.Equal(v4) {
			t.Error("IPv4 address lost")
		}
		if v6 != nil && !f.PMFIPv6Address.Equal(v6) {
			t.Error("IPv6 address lost")
		}
	}
}

func TestRemoteZeroDestinationLength(t *testing.T) {
	f := ie.NewRemoteGTPUPeerFields(6, "192.0.2.1", "", 0, "")
	f.DILength = 0
	if p := panicOf(func() {
		if _, err := f.Marshal(); err == nil {
			t.Error("accepted zero destination-interface length")
		}
	}); p != nil {
		t.Fatal("panic", p)
	}
}

func TestRemoteInterfaceExtension(t *testing.T) {
	// TS 29.244 V20.1.0 clauses 8.2.70 and 8.2.24 permit the extendable
	// Destination Interface value. Its first octet must still be present.
	i := ie.New(ie.RemoteGTPUPeer, []byte{6, 192, 0, 2, 1, 0, 2, 0, 99})
	fields, err := i.RemoteGTPUPeer()
	if err != nil {
		t.Fatal(err)
	}
	if fields.DestinationInterface != 0 || fields.DILength != 2 {
		t.Error("destination interface extension changed known fields")
	}
	wire, err := i.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ie.Parse(wire)
	if err != nil {
		t.Fatal(err)
	}
	again, err := parsed.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(wire, again) {
		t.Error("destination interface extension changed wire bytes")
	}
}
