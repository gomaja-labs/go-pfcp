package ie_test

import (
	"errors"
	"io"
	"net/netip"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestMoQRelayIPAddress(t *testing.T) {
	for flags := uint8(0); flags <= 3; flags++ {
		for boundary := 0; boundary < 2; boundary++ {
			f := &ie.MoQRelayIPAddressFields{Flags: flags}
			want := []byte{flags}
			if flags&2 != 0 {
				f.IPv4Address = netip.MustParseAddr("192.0.2.129")
				b := f.IPv4Address.AsSlice()
				want = append(want, b...)
			}
			if flags&1 != 0 {
				f.IPv6Address = netip.MustParseAddr("2001:db8::1234")
				b := f.IPv6Address.AsSlice()
				want = append(want, b...)
			}
			i := ie.NewMoQRelayIPAddress(flags, f.IPv4Address, f.IPv6Address)
			if !(true) {
				assertNewIEInvalid(t, i, ie.MoQRelayIPAddress)
				continue
			}
			checkNewIEWire(t, i, want)
			got, err := i.MoQRelayIPAddress()
			if err != nil || !reflect.DeepEqual(f, got) {
				t.Fatalf("flags %d: got %#v want %#v error %v", flags, got, f, err)
			}
			for n := 0; n < len(want); n++ {
				if _, err = ie.ParseMoQRelayIPAddressFields(want[:n]); err == nil {
					t.Fatalf("accepted truncation %d flags %d", n, flags)
				}
			}
			for n := 0; n < f.MarshalLen(); n++ {
				if err = f.MarshalTo(make([]byte, n)); !errors.Is(err, io.ErrUnexpectedEOF) {
					t.Fatal("short Fields buffer", err)
				}
			}
			extended := append(append([]byte{}, want...), 0xa5, 0x5a)
			extended[0] |= 252
			got, err = ie.ParseMoQRelayIPAddressFields(extended)
			if err != nil || !reflect.DeepEqual(f, got) {
				t.Fatal("spares/extensions", err)
			}
			*got = ie.MoQRelayIPAddressFields{Flags: 255}
			if err = got.UnmarshalBinary(want); err != nil || !reflect.DeepEqual(f, got) {
				t.Fatal("reused receiver", err)
			}
		}
	}
	checkNewIEFailures(t, ie.MoQRelayIPAddress, func(i *ie.IE) (any, error) { return i.MoQRelayIPAddress() }, 1)
	var f *ie.MoQRelayIPAddressFields
	if _, err := f.Marshal(); !errors.Is(err, ie.ErrIENotFound) {
		t.Fatal(err)
	}
	if err := f.MarshalTo(nil); !errors.Is(err, ie.ErrIENotFound) {
		t.Fatal(err)
	}
	if err := f.UnmarshalBinary(nil); !errors.Is(err, ie.ErrIENotFound) {
		t.Fatal(err)
	}
	if f.MarshalLen() != 0 {
		t.Fatal("nil length")
	}
}

func TestMoQRelayIPAddressInvalidFields(t *testing.T) {
	valid := func() *ie.MoQRelayIPAddressFields {
		return &ie.MoQRelayIPAddressFields{Flags: 3, IPv4Address: netip.MustParseAddr("192.0.2.1"), IPv6Address: netip.MustParseAddr("2001:db8::1")}
	}
	{
		f := valid()
		f.Flags |= 4
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags |= 4`)
		}
		assertNewIEInvalid(t, ie.NewMoQRelayIPAddress(f.Flags, f.IPv4Address, f.IPv6Address), ie.MoQRelayIPAddress)
	}
	{
		f := valid()
		f.Flags &^= 2
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 2`)
		}
		assertNewIEInvalid(t, ie.NewMoQRelayIPAddress(f.Flags, f.IPv4Address, f.IPv6Address), ie.MoQRelayIPAddress)
	}
	{
		f := valid()
		f.IPv4Address = netip.Addr{}
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.IPv4Address=netip.Addr{}`)
		}
		assertNewIEInvalid(t, ie.NewMoQRelayIPAddress(f.Flags, f.IPv4Address, f.IPv6Address), ie.MoQRelayIPAddress)
	}
	{
		f := valid()
		f.IPv4Address = netip.MustParseAddr("2001:db8::1")
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.IPv4Address=netip.MustParseAddr("2001:db8::1")`)
		}
		assertNewIEInvalid(t, ie.NewMoQRelayIPAddress(f.Flags, f.IPv4Address, f.IPv6Address), ie.MoQRelayIPAddress)
	}
	{
		f := valid()
		f.Flags &^= 1
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 1`)
		}
		assertNewIEInvalid(t, ie.NewMoQRelayIPAddress(f.Flags, f.IPv4Address, f.IPv6Address), ie.MoQRelayIPAddress)
	}
	{
		f := valid()
		f.IPv6Address = netip.Addr{}
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.IPv6Address=netip.Addr{}`)
		}
		assertNewIEInvalid(t, ie.NewMoQRelayIPAddress(f.Flags, f.IPv4Address, f.IPv6Address), ie.MoQRelayIPAddress)
	}
	{
		f := valid()
		f.IPv6Address = netip.MustParseAddr("192.0.2.1")
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.IPv6Address=netip.MustParseAddr("192.0.2.1")`)
		}
		assertNewIEInvalid(t, ie.NewMoQRelayIPAddress(f.Flags, f.IPv4Address, f.IPv6Address), ie.MoQRelayIPAddress)
	}
	{
		f := valid()
		f.IPv6Address = netip.MustParseAddr("fe80::1%zone")
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.IPv6Address=netip.MustParseAddr("fe80::1%zone")`)
		}
		assertNewIEInvalid(t, ie.NewMoQRelayIPAddress(f.Flags, f.IPv4Address, f.IPv6Address), ie.MoQRelayIPAddress)
	}
}
