package ie_test

import (
	"errors"
	"io"
	"net/netip"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestOnPathN6ConnectionInformation(t *testing.T) {
	for flags := uint8(0); flags <= 63; flags++ {
		for boundary := 0; boundary < 2; boundary++ {
			f := &ie.OnPathN6ConnectionInformationFields{Flags: flags}
			want := []byte{flags}
			if flags&8 != 0 {
				f.IPv4Address = netip.MustParseAddr("192.0.2.129")
				b := f.IPv4Address.AsSlice()
				want = append(want, b...)
			}
			if flags&16 != 0 {
				f.IPv6Address = netip.MustParseAddr("2001:db8::1234")
				b := f.IPv6Address.AsSlice()
				want = append(want, b...)
			}
			if flags&32 != 0 {
				f.PortNumber = uint16(boundary) * 65535
				v := uint64(f.PortNumber)
				want = append(want, byte(v>>8))
				want = append(want, byte(v>>0))
			}
			i := ie.NewOnPathN6ConnectionInformation(flags, f.IPv4Address, f.IPv6Address, f.PortNumber)
			if flags&7 > 1 {
				assertNewIEInvalid(t, i, ie.OnPathN6ConnectionInformation)
				continue
			}
			checkNewIEWire(t, i, want)
			got, err := i.OnPathN6ConnectionInformation()
			if err != nil || !reflect.DeepEqual(f, got) {
				t.Fatalf("flags %d: got %#v want %#v error %v", flags, got, f, err)
			}
			for n := 0; n < len(want); n++ {
				if _, err = ie.ParseOnPathN6ConnectionInformationFields(want[:n]); err == nil {
					t.Fatalf("accepted truncation %d flags %d", n, flags)
				}
			}
			for n := 0; n < f.MarshalLen(); n++ {
				if err = f.MarshalTo(make([]byte, n)); !errors.Is(err, io.ErrUnexpectedEOF) {
					t.Fatal("short Fields buffer", err)
				}
			}
			extended := append(append([]byte{}, want...), 0xa5, 0x5a)
			extended[0] |= 192
			got, err = ie.ParseOnPathN6ConnectionInformationFields(extended)
			if err != nil || !reflect.DeepEqual(f, got) {
				t.Fatal("spares/extensions", err)
			}
			*got = ie.OnPathN6ConnectionInformationFields{Flags: 255}
			if err = got.UnmarshalBinary(want); err != nil || !reflect.DeepEqual(f, got) {
				t.Fatal("reused receiver", err)
			}
		}
	}
	checkNewIEFailures(t, ie.OnPathN6ConnectionInformation, func(i *ie.IE) (any, error) { return i.OnPathN6ConnectionInformation() }, 1)
	var f *ie.OnPathN6ConnectionInformationFields
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

func TestOnPathN6ConnectionInformationInvalidFields(t *testing.T) {
	valid := func() *ie.OnPathN6ConnectionInformationFields {
		return &ie.OnPathN6ConnectionInformationFields{Flags: 57, IPv4Address: netip.MustParseAddr("192.0.2.1"), IPv6Address: netip.MustParseAddr("2001:db8::1"), PortNumber: 1}
	}
	{
		f := valid()
		f.Flags |= 64
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags |= 64`)
		}
		assertNewIEInvalid(t, ie.NewOnPathN6ConnectionInformation(f.Flags, f.IPv4Address, f.IPv6Address, f.PortNumber), ie.OnPathN6ConnectionInformation)
	}
	{
		f := valid()
		f.Flags = f.Flags&^7 | 2
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags=f.Flags&^7|2`)
		}
		assertNewIEInvalid(t, ie.NewOnPathN6ConnectionInformation(f.Flags, f.IPv4Address, f.IPv6Address, f.PortNumber), ie.OnPathN6ConnectionInformation)
	}
	{
		f := valid()
		f.Flags &^= 8
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 8`)
		}
		assertNewIEInvalid(t, ie.NewOnPathN6ConnectionInformation(f.Flags, f.IPv4Address, f.IPv6Address, f.PortNumber), ie.OnPathN6ConnectionInformation)
	}
	{
		f := valid()
		f.IPv4Address = netip.Addr{}
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.IPv4Address=netip.Addr{}`)
		}
		assertNewIEInvalid(t, ie.NewOnPathN6ConnectionInformation(f.Flags, f.IPv4Address, f.IPv6Address, f.PortNumber), ie.OnPathN6ConnectionInformation)
	}
	{
		f := valid()
		f.IPv4Address = netip.MustParseAddr("2001:db8::1")
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.IPv4Address=netip.MustParseAddr("2001:db8::1")`)
		}
		assertNewIEInvalid(t, ie.NewOnPathN6ConnectionInformation(f.Flags, f.IPv4Address, f.IPv6Address, f.PortNumber), ie.OnPathN6ConnectionInformation)
	}
	{
		f := valid()
		f.Flags &^= 16
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 16`)
		}
		assertNewIEInvalid(t, ie.NewOnPathN6ConnectionInformation(f.Flags, f.IPv4Address, f.IPv6Address, f.PortNumber), ie.OnPathN6ConnectionInformation)
	}
	{
		f := valid()
		f.IPv6Address = netip.Addr{}
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.IPv6Address=netip.Addr{}`)
		}
		assertNewIEInvalid(t, ie.NewOnPathN6ConnectionInformation(f.Flags, f.IPv4Address, f.IPv6Address, f.PortNumber), ie.OnPathN6ConnectionInformation)
	}
	{
		f := valid()
		f.IPv6Address = netip.MustParseAddr("192.0.2.1")
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.IPv6Address=netip.MustParseAddr("192.0.2.1")`)
		}
		assertNewIEInvalid(t, ie.NewOnPathN6ConnectionInformation(f.Flags, f.IPv4Address, f.IPv6Address, f.PortNumber), ie.OnPathN6ConnectionInformation)
	}
	{
		f := valid()
		f.IPv6Address = netip.MustParseAddr("fe80::1%zone")
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.IPv6Address=netip.MustParseAddr("fe80::1%zone")`)
		}
		assertNewIEInvalid(t, ie.NewOnPathN6ConnectionInformation(f.Flags, f.IPv4Address, f.IPv6Address, f.PortNumber), ie.OnPathN6ConnectionInformation)
	}
	{
		f := valid()
		f.Flags &^= 32
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 32`)
		}
		assertNewIEInvalid(t, ie.NewOnPathN6ConnectionInformation(f.Flags, f.IPv4Address, f.IPv6Address, f.PortNumber), ie.OnPathN6ConnectionInformation)
	}
}
