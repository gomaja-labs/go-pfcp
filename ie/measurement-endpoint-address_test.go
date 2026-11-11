package ie_test

import (
	"errors"
	"io"
	"net/netip"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestMeasurementEndpointAddress(t *testing.T) {
	for flags := uint8(0); flags <= 3; flags++ {
		for boundary := 0; boundary < 2; boundary++ {
			f := &ie.MeasurementEndpointAddressFields{Flags: flags}
			want := []byte{flags}
			if flags&1 != 0 {
				f.IPv4Address = netip.MustParseAddr("192.0.2.129")
				b := f.IPv4Address.AsSlice()
				want = append(want, b...)
			}
			if flags&2 != 0 {
				f.IPv6Address = netip.MustParseAddr("2001:db8::1234")
				b := f.IPv6Address.AsSlice()
				want = append(want, b...)
			}
			i := ie.NewMeasurementEndpointAddress(flags, f.IPv4Address, f.IPv6Address)
			if flags == 0 {
				assertNewIEInvalid(t, i, ie.MeasurementEndpointAddress)
				continue
			}
			checkNewIEWire(t, i, want)
			got, err := i.MeasurementEndpointAddress()
			if err != nil || !reflect.DeepEqual(f, got) {
				t.Fatalf("flags %d: got %#v want %#v error %v", flags, got, f, err)
			}
			for n := 0; n < len(want); n++ {
				if _, err = ie.ParseMeasurementEndpointAddressFields(want[:n]); err == nil {
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
			got, err = ie.ParseMeasurementEndpointAddressFields(extended)
			if err != nil || !reflect.DeepEqual(f, got) {
				t.Fatal("spares/extensions", err)
			}
			*got = ie.MeasurementEndpointAddressFields{Flags: 255}
			if err = got.UnmarshalBinary(want); err != nil || !reflect.DeepEqual(f, got) {
				t.Fatal("reused receiver", err)
			}
		}
	}
	checkNewIEFailures(t, ie.MeasurementEndpointAddress, func(i *ie.IE) (any, error) { return i.MeasurementEndpointAddress() }, 1)
	var f *ie.MeasurementEndpointAddressFields
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

func TestMeasurementEndpointAddressInvalidFields(t *testing.T) {
	valid := func() *ie.MeasurementEndpointAddressFields {
		return &ie.MeasurementEndpointAddressFields{Flags: 3, IPv4Address: netip.MustParseAddr("192.0.2.1"), IPv6Address: netip.MustParseAddr("2001:db8::1")}
	}
	{
		f := valid()
		f.Flags |= 4
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags |= 4`)
		}
		assertNewIEInvalid(t, ie.NewMeasurementEndpointAddress(f.Flags, f.IPv4Address, f.IPv6Address), ie.MeasurementEndpointAddress)
	}
	{
		f := valid()
		f.Flags &^= 1
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 1`)
		}
		assertNewIEInvalid(t, ie.NewMeasurementEndpointAddress(f.Flags, f.IPv4Address, f.IPv6Address), ie.MeasurementEndpointAddress)
	}
	{
		f := valid()
		f.IPv4Address = netip.Addr{}
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.IPv4Address=netip.Addr{}`)
		}
		assertNewIEInvalid(t, ie.NewMeasurementEndpointAddress(f.Flags, f.IPv4Address, f.IPv6Address), ie.MeasurementEndpointAddress)
	}
	{
		f := valid()
		f.IPv4Address = netip.MustParseAddr("2001:db8::1")
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.IPv4Address=netip.MustParseAddr("2001:db8::1")`)
		}
		assertNewIEInvalid(t, ie.NewMeasurementEndpointAddress(f.Flags, f.IPv4Address, f.IPv6Address), ie.MeasurementEndpointAddress)
	}
	{
		f := valid()
		f.Flags &^= 2
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 2`)
		}
		assertNewIEInvalid(t, ie.NewMeasurementEndpointAddress(f.Flags, f.IPv4Address, f.IPv6Address), ie.MeasurementEndpointAddress)
	}
	{
		f := valid()
		f.IPv6Address = netip.Addr{}
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.IPv6Address=netip.Addr{}`)
		}
		assertNewIEInvalid(t, ie.NewMeasurementEndpointAddress(f.Flags, f.IPv4Address, f.IPv6Address), ie.MeasurementEndpointAddress)
	}
	{
		f := valid()
		f.IPv6Address = netip.MustParseAddr("192.0.2.1")
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.IPv6Address=netip.MustParseAddr("192.0.2.1")`)
		}
		assertNewIEInvalid(t, ie.NewMeasurementEndpointAddress(f.Flags, f.IPv4Address, f.IPv6Address), ie.MeasurementEndpointAddress)
	}
	{
		f := valid()
		f.IPv6Address = netip.MustParseAddr("fe80::1%zone")
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.IPv6Address=netip.MustParseAddr("fe80::1%zone")`)
		}
		assertNewIEInvalid(t, ie.NewMeasurementEndpointAddress(f.Flags, f.IPv4Address, f.IPv6Address), ie.MeasurementEndpointAddress)
	}
}
