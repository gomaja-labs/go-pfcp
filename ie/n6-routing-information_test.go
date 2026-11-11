package ie_test

import (
	"errors"
	"io"
	"net/netip"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestN6RoutingInformation(t *testing.T) {
	for flags := uint8(0); flags <= 63; flags++ {
		for boundary := 0; boundary < 2; boundary++ {
			f := &ie.N6RoutingInformationFields{Flags: flags}
			want := []byte{flags}
			if flags&1 != 0 {
				f.SourceIPv4Address = netip.MustParseAddr("192.0.2.129")
				b := f.SourceIPv4Address.AsSlice()
				want = append(want, b...)
			}
			if flags&2 != 0 {
				f.SourceIPv6Address = netip.MustParseAddr("2001:db8::1234")
				b := f.SourceIPv6Address.AsSlice()
				want = append(want, b...)
			}
			if flags&4 != 0 {
				f.SourcePort = uint16(boundary) * 65535
				v := uint64(f.SourcePort)
				want = append(want, byte(v>>8))
				want = append(want, byte(v>>0))
			}
			if flags&8 != 0 {
				f.DestinationIPv4Address = netip.MustParseAddr("192.0.2.129")
				b := f.DestinationIPv4Address.AsSlice()
				want = append(want, b...)
			}
			if flags&16 != 0 {
				f.DestinationIPv6Address = netip.MustParseAddr("2001:db8::1234")
				b := f.DestinationIPv6Address.AsSlice()
				want = append(want, b...)
			}
			if flags&32 != 0 {
				f.DestinationPort = uint16(boundary) * 65535
				v := uint64(f.DestinationPort)
				want = append(want, byte(v>>8))
				want = append(want, byte(v>>0))
			}
			i := ie.NewN6RoutingInformation(flags, f.SourceIPv4Address, f.SourceIPv6Address, f.SourcePort, f.DestinationIPv4Address, f.DestinationIPv6Address, f.DestinationPort)
			if !(true) {
				assertNewIEInvalid(t, i, ie.N6RoutingInformation)
				continue
			}
			checkNewIEWire(t, i, want)
			got, err := i.N6RoutingInformation()
			if err != nil || !reflect.DeepEqual(f, got) {
				t.Fatalf("flags %d: got %#v want %#v error %v", flags, got, f, err)
			}
			for n := 0; n < len(want); n++ {
				if _, err = ie.ParseN6RoutingInformationFields(want[:n]); err == nil {
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
			got, err = ie.ParseN6RoutingInformationFields(extended)
			if err != nil || !reflect.DeepEqual(f, got) {
				t.Fatal("spares/extensions", err)
			}
			*got = ie.N6RoutingInformationFields{Flags: 255}
			if err = got.UnmarshalBinary(want); err != nil || !reflect.DeepEqual(f, got) {
				t.Fatal("reused receiver", err)
			}
		}
	}
	checkNewIEFailures(t, ie.N6RoutingInformation, func(i *ie.IE) (any, error) { return i.N6RoutingInformation() }, 1)
	var f *ie.N6RoutingInformationFields
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

func TestN6RoutingInformationInvalidFields(t *testing.T) {
	valid := func() *ie.N6RoutingInformationFields {
		return &ie.N6RoutingInformationFields{Flags: 63, SourceIPv4Address: netip.MustParseAddr("192.0.2.1"), SourceIPv6Address: netip.MustParseAddr("2001:db8::1"), SourcePort: 1, DestinationIPv4Address: netip.MustParseAddr("192.0.2.1"), DestinationIPv6Address: netip.MustParseAddr("2001:db8::1"), DestinationPort: 1}
	}
	{
		f := valid()
		f.Flags |= 64
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags |= 64`)
		}
		assertNewIEInvalid(t, ie.NewN6RoutingInformation(f.Flags, f.SourceIPv4Address, f.SourceIPv6Address, f.SourcePort, f.DestinationIPv4Address, f.DestinationIPv6Address, f.DestinationPort), ie.N6RoutingInformation)
	}
	{
		f := valid()
		f.Flags &^= 1
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 1`)
		}
		assertNewIEInvalid(t, ie.NewN6RoutingInformation(f.Flags, f.SourceIPv4Address, f.SourceIPv6Address, f.SourcePort, f.DestinationIPv4Address, f.DestinationIPv6Address, f.DestinationPort), ie.N6RoutingInformation)
	}
	{
		f := valid()
		f.SourceIPv4Address = netip.Addr{}
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.SourceIPv4Address=netip.Addr{}`)
		}
		assertNewIEInvalid(t, ie.NewN6RoutingInformation(f.Flags, f.SourceIPv4Address, f.SourceIPv6Address, f.SourcePort, f.DestinationIPv4Address, f.DestinationIPv6Address, f.DestinationPort), ie.N6RoutingInformation)
	}
	{
		f := valid()
		f.SourceIPv4Address = netip.MustParseAddr("2001:db8::1")
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.SourceIPv4Address=netip.MustParseAddr("2001:db8::1")`)
		}
		assertNewIEInvalid(t, ie.NewN6RoutingInformation(f.Flags, f.SourceIPv4Address, f.SourceIPv6Address, f.SourcePort, f.DestinationIPv4Address, f.DestinationIPv6Address, f.DestinationPort), ie.N6RoutingInformation)
	}
	{
		f := valid()
		f.Flags &^= 2
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 2`)
		}
		assertNewIEInvalid(t, ie.NewN6RoutingInformation(f.Flags, f.SourceIPv4Address, f.SourceIPv6Address, f.SourcePort, f.DestinationIPv4Address, f.DestinationIPv6Address, f.DestinationPort), ie.N6RoutingInformation)
	}
	{
		f := valid()
		f.SourceIPv6Address = netip.Addr{}
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.SourceIPv6Address=netip.Addr{}`)
		}
		assertNewIEInvalid(t, ie.NewN6RoutingInformation(f.Flags, f.SourceIPv4Address, f.SourceIPv6Address, f.SourcePort, f.DestinationIPv4Address, f.DestinationIPv6Address, f.DestinationPort), ie.N6RoutingInformation)
	}
	{
		f := valid()
		f.SourceIPv6Address = netip.MustParseAddr("192.0.2.1")
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.SourceIPv6Address=netip.MustParseAddr("192.0.2.1")`)
		}
		assertNewIEInvalid(t, ie.NewN6RoutingInformation(f.Flags, f.SourceIPv4Address, f.SourceIPv6Address, f.SourcePort, f.DestinationIPv4Address, f.DestinationIPv6Address, f.DestinationPort), ie.N6RoutingInformation)
	}
	{
		f := valid()
		f.SourceIPv6Address = netip.MustParseAddr("fe80::1%zone")
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.SourceIPv6Address=netip.MustParseAddr("fe80::1%zone")`)
		}
		assertNewIEInvalid(t, ie.NewN6RoutingInformation(f.Flags, f.SourceIPv4Address, f.SourceIPv6Address, f.SourcePort, f.DestinationIPv4Address, f.DestinationIPv6Address, f.DestinationPort), ie.N6RoutingInformation)
	}
	{
		f := valid()
		f.Flags &^= 4
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 4`)
		}
		assertNewIEInvalid(t, ie.NewN6RoutingInformation(f.Flags, f.SourceIPv4Address, f.SourceIPv6Address, f.SourcePort, f.DestinationIPv4Address, f.DestinationIPv6Address, f.DestinationPort), ie.N6RoutingInformation)
	}
	{
		f := valid()
		f.Flags &^= 8
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 8`)
		}
		assertNewIEInvalid(t, ie.NewN6RoutingInformation(f.Flags, f.SourceIPv4Address, f.SourceIPv6Address, f.SourcePort, f.DestinationIPv4Address, f.DestinationIPv6Address, f.DestinationPort), ie.N6RoutingInformation)
	}
	{
		f := valid()
		f.DestinationIPv4Address = netip.Addr{}
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.DestinationIPv4Address=netip.Addr{}`)
		}
		assertNewIEInvalid(t, ie.NewN6RoutingInformation(f.Flags, f.SourceIPv4Address, f.SourceIPv6Address, f.SourcePort, f.DestinationIPv4Address, f.DestinationIPv6Address, f.DestinationPort), ie.N6RoutingInformation)
	}
	{
		f := valid()
		f.DestinationIPv4Address = netip.MustParseAddr("2001:db8::1")
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.DestinationIPv4Address=netip.MustParseAddr("2001:db8::1")`)
		}
		assertNewIEInvalid(t, ie.NewN6RoutingInformation(f.Flags, f.SourceIPv4Address, f.SourceIPv6Address, f.SourcePort, f.DestinationIPv4Address, f.DestinationIPv6Address, f.DestinationPort), ie.N6RoutingInformation)
	}
	{
		f := valid()
		f.Flags &^= 16
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 16`)
		}
		assertNewIEInvalid(t, ie.NewN6RoutingInformation(f.Flags, f.SourceIPv4Address, f.SourceIPv6Address, f.SourcePort, f.DestinationIPv4Address, f.DestinationIPv6Address, f.DestinationPort), ie.N6RoutingInformation)
	}
	{
		f := valid()
		f.DestinationIPv6Address = netip.Addr{}
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.DestinationIPv6Address=netip.Addr{}`)
		}
		assertNewIEInvalid(t, ie.NewN6RoutingInformation(f.Flags, f.SourceIPv4Address, f.SourceIPv6Address, f.SourcePort, f.DestinationIPv4Address, f.DestinationIPv6Address, f.DestinationPort), ie.N6RoutingInformation)
	}
	{
		f := valid()
		f.DestinationIPv6Address = netip.MustParseAddr("192.0.2.1")
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.DestinationIPv6Address=netip.MustParseAddr("192.0.2.1")`)
		}
		assertNewIEInvalid(t, ie.NewN6RoutingInformation(f.Flags, f.SourceIPv4Address, f.SourceIPv6Address, f.SourcePort, f.DestinationIPv4Address, f.DestinationIPv6Address, f.DestinationPort), ie.N6RoutingInformation)
	}
	{
		f := valid()
		f.DestinationIPv6Address = netip.MustParseAddr("fe80::1%zone")
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.DestinationIPv6Address=netip.MustParseAddr("fe80::1%zone")`)
		}
		assertNewIEInvalid(t, ie.NewN6RoutingInformation(f.Flags, f.SourceIPv4Address, f.SourceIPv6Address, f.SourcePort, f.DestinationIPv4Address, f.DestinationIPv6Address, f.DestinationPort), ie.N6RoutingInformation)
	}
	{
		f := valid()
		f.Flags &^= 32
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 32`)
		}
		assertNewIEInvalid(t, ie.NewN6RoutingInformation(f.Flags, f.SourceIPv4Address, f.SourceIPv6Address, f.SourcePort, f.DestinationIPv4Address, f.DestinationIPv6Address, f.DestinationPort), ie.N6RoutingInformation)
	}
}
