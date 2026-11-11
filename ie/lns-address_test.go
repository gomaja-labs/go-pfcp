package ie_test

import (
	"bytes"
	"errors"
	"net/netip"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestLNSAddressCurrent(t *testing.T) {
	for _, s := range []string{"0.0.0.0", "192.0.2.1", "::", "2001:db8::1"} {
		a := netip.MustParseAddr(s)
		i := ie.NewLNSAddress(a)
		for n := 0; n < i.MarshalLen(); n++ {
			if e := i.MarshalTo(make([]byte, n)); !errors.Is(e, ie.ErrInvalidLength) {
				t.Fatal("short IE buffer", e)
			}
		}
		got, e := i.LNSAddress()
		if e != nil || got != a {
			t.Fatal("constructor value", e)
		}
		wire, e := i.Marshal()
		if e != nil || !bytes.Equal(wire[4:], a.AsSlice()) {
			t.Fatal("address encoding", e)
		}
		p, e := ie.Parse(wire)
		if e != nil {
			t.Fatal(e)
		}
		got, e = p.LNSAddress()
		if e != nil || got != a {
			t.Fatal(e)
		}
		again, e := p.Marshal()
		if e != nil || !bytes.Equal(again, wire) {
			t.Fatal("roundtrip", e)
		}
	}
	var p *ie.IE
	if _, e := p.LNSAddress(); !errors.Is(e, ie.ErrIENotFound) {
		t.Fatal(e)
	}
	if _, e := ie.New(0, nil).LNSAddress(); e == nil {
		t.Fatal("wrong type")
	}
	for n := 0; n <= 32; n++ {
		if n == 4 || n == 16 {
			continue
		}
		i := ie.New(ie.LNSAddress, make([]byte, n))
		wire, e := i.Marshal()
		if e != nil {
			t.Fatal(e)
		}
		p, e := ie.Parse(wire)
		if e != nil {
			t.Fatal(e)
		}
		_, e = p.LNSAddress()
		if !errors.Is(e, ie.ErrMalformed) {
			t.Fatalf("length %d: %v", n, e)
		}
		again, e := p.Marshal()
		if e != nil || !bytes.Equal(wire, again) {
			t.Fatal("raw malformed length lost")
		}
	}
	for _, a := range []netip.Addr{{}, netip.MustParseAddr("fe80::1%eth0")} {
		assertAddedInvalid(t, ie.NewLNSAddress(a))
	}
}

var _ = registerConstructorFixture(ie.NewLNSAddress, netip.MustParseAddr("192.0.2.1"))
