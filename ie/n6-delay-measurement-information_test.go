package ie_test

import (
	"errors"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestN6DelayMeasurementInformation(t *testing.T) {
	for flags := uint8(0); flags <= 3; flags++ {
		for boundary := 0; boundary < 2; boundary++ {
			f := &ie.N6DelayMeasurementInformationFields{Flags: flags}
			want := []byte{flags}
			if true {
				f.AveragePacketDelay = time.Duration(boundary) * 4294967295 * time.Millisecond
				v := uint32(f.AveragePacketDelay / time.Millisecond)
				want = append(want, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
			}
			if flags&1 != 0 {
				f.MinimumPacketDelay = time.Duration(boundary) * 4294967295 * time.Millisecond
				v := uint32(f.MinimumPacketDelay / time.Millisecond)
				want = append(want, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
			}
			if flags&2 != 0 {
				f.MaximumPacketDelay = time.Duration(boundary) * 4294967295 * time.Millisecond
				v := uint32(f.MaximumPacketDelay / time.Millisecond)
				want = append(want, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
			}
			i := ie.NewN6DelayMeasurementInformation(flags, f.AveragePacketDelay, f.MinimumPacketDelay, f.MaximumPacketDelay)
			if !(true) {
				assertNewIEInvalid(t, i, ie.N6DelayMeasurementInformation)
				continue
			}
			checkNewIEWire(t, i, want)
			got, err := i.N6DelayMeasurementInformation()
			if err != nil || !reflect.DeepEqual(f, got) {
				t.Fatalf("flags %d: got %#v want %#v error %v", flags, got, f, err)
			}
			for n := 0; n < len(want); n++ {
				if _, err = ie.ParseN6DelayMeasurementInformationFields(want[:n]); err == nil {
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
			got, err = ie.ParseN6DelayMeasurementInformationFields(extended)
			if err != nil || !reflect.DeepEqual(f, got) {
				t.Fatal("spares/extensions", err)
			}
			*got = ie.N6DelayMeasurementInformationFields{Flags: 255}
			if err = got.UnmarshalBinary(want); err != nil || !reflect.DeepEqual(f, got) {
				t.Fatal("reused receiver", err)
			}
		}
	}
	checkNewIEFailures(t, ie.N6DelayMeasurementInformation, func(i *ie.IE) (any, error) { return i.N6DelayMeasurementInformation() }, 1)
	var f *ie.N6DelayMeasurementInformationFields
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

func TestN6DelayMeasurementInformationInvalidFields(t *testing.T) {
	valid := func() *ie.N6DelayMeasurementInformationFields {
		return &ie.N6DelayMeasurementInformationFields{Flags: 3, AveragePacketDelay: time.Millisecond, MinimumPacketDelay: time.Millisecond, MaximumPacketDelay: time.Millisecond}
	}
	{
		f := valid()
		f.Flags |= 4
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags |= 4`)
		}
		assertNewIEInvalid(t, ie.NewN6DelayMeasurementInformation(f.Flags, f.AveragePacketDelay, f.MinimumPacketDelay, f.MaximumPacketDelay), ie.N6DelayMeasurementInformation)
	}
	{
		f := valid()
		f.AveragePacketDelay = -time.Millisecond
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.AveragePacketDelay=-time.Millisecond`)
		}
		assertNewIEInvalid(t, ie.NewN6DelayMeasurementInformation(f.Flags, f.AveragePacketDelay, f.MinimumPacketDelay, f.MaximumPacketDelay), ie.N6DelayMeasurementInformation)
	}
	{
		f := valid()
		f.AveragePacketDelay = time.Nanosecond
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.AveragePacketDelay=time.Nanosecond`)
		}
		assertNewIEInvalid(t, ie.NewN6DelayMeasurementInformation(f.Flags, f.AveragePacketDelay, f.MinimumPacketDelay, f.MaximumPacketDelay), ie.N6DelayMeasurementInformation)
	}
	{
		f := valid()
		f.AveragePacketDelay = (1 << 32) * time.Millisecond
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.AveragePacketDelay=(1<<32)*time.Millisecond`)
		}
		assertNewIEInvalid(t, ie.NewN6DelayMeasurementInformation(f.Flags, f.AveragePacketDelay, f.MinimumPacketDelay, f.MaximumPacketDelay), ie.N6DelayMeasurementInformation)
	}
	{
		f := valid()
		f.Flags &^= 1
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 1`)
		}
		assertNewIEInvalid(t, ie.NewN6DelayMeasurementInformation(f.Flags, f.AveragePacketDelay, f.MinimumPacketDelay, f.MaximumPacketDelay), ie.N6DelayMeasurementInformation)
	}
	{
		f := valid()
		f.MinimumPacketDelay = -time.Millisecond
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.MinimumPacketDelay=-time.Millisecond`)
		}
		assertNewIEInvalid(t, ie.NewN6DelayMeasurementInformation(f.Flags, f.AveragePacketDelay, f.MinimumPacketDelay, f.MaximumPacketDelay), ie.N6DelayMeasurementInformation)
	}
	{
		f := valid()
		f.MinimumPacketDelay = time.Nanosecond
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.MinimumPacketDelay=time.Nanosecond`)
		}
		assertNewIEInvalid(t, ie.NewN6DelayMeasurementInformation(f.Flags, f.AveragePacketDelay, f.MinimumPacketDelay, f.MaximumPacketDelay), ie.N6DelayMeasurementInformation)
	}
	{
		f := valid()
		f.MinimumPacketDelay = (1 << 32) * time.Millisecond
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.MinimumPacketDelay=(1<<32)*time.Millisecond`)
		}
		assertNewIEInvalid(t, ie.NewN6DelayMeasurementInformation(f.Flags, f.AveragePacketDelay, f.MinimumPacketDelay, f.MaximumPacketDelay), ie.N6DelayMeasurementInformation)
	}
	{
		f := valid()
		f.Flags &^= 2
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 2`)
		}
		assertNewIEInvalid(t, ie.NewN6DelayMeasurementInformation(f.Flags, f.AveragePacketDelay, f.MinimumPacketDelay, f.MaximumPacketDelay), ie.N6DelayMeasurementInformation)
	}
	{
		f := valid()
		f.MaximumPacketDelay = -time.Millisecond
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.MaximumPacketDelay=-time.Millisecond`)
		}
		assertNewIEInvalid(t, ie.NewN6DelayMeasurementInformation(f.Flags, f.AveragePacketDelay, f.MinimumPacketDelay, f.MaximumPacketDelay), ie.N6DelayMeasurementInformation)
	}
	{
		f := valid()
		f.MaximumPacketDelay = time.Nanosecond
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.MaximumPacketDelay=time.Nanosecond`)
		}
		assertNewIEInvalid(t, ie.NewN6DelayMeasurementInformation(f.Flags, f.AveragePacketDelay, f.MinimumPacketDelay, f.MaximumPacketDelay), ie.N6DelayMeasurementInformation)
	}
	{
		f := valid()
		f.MaximumPacketDelay = (1 << 32) * time.Millisecond
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.MaximumPacketDelay=(1<<32)*time.Millisecond`)
		}
		assertNewIEInvalid(t, ie.NewN6DelayMeasurementInformation(f.Flags, f.AveragePacketDelay, f.MinimumPacketDelay, f.MaximumPacketDelay), ie.N6DelayMeasurementInformation)
	}
}
