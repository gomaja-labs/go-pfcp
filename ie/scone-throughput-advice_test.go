package ie_test

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestSCONEThroughputAdvice(t *testing.T) {
	for flags := uint8(0); flags <= 3; flags++ {
		for boundary := 0; boundary < 2; boundary++ {
			f := &ie.SCONEThroughputAdviceFields{Flags: flags}
			want := []byte{flags}
			if flags&1 != 0 {
				f.UplinkThroughputAdvice = uint32(boundary) * 4294967295
				v := uint64(f.UplinkThroughputAdvice)
				want = append(want, byte(v>>24))
				want = append(want, byte(v>>16))
				want = append(want, byte(v>>8))
				want = append(want, byte(v>>0))
			}
			if flags&2 != 0 {
				f.DownlinkThroughputAdvice = uint32(boundary) * 4294967295
				v := uint64(f.DownlinkThroughputAdvice)
				want = append(want, byte(v>>24))
				want = append(want, byte(v>>16))
				want = append(want, byte(v>>8))
				want = append(want, byte(v>>0))
			}
			i := ie.NewSCONEThroughputAdvice(flags, f.UplinkThroughputAdvice, f.DownlinkThroughputAdvice)
			if !(true) {
				assertNewIEInvalid(t, i, ie.SCONEThroughputAdvice)
				continue
			}
			checkNewIEWire(t, i, want)
			got, err := i.SCONEThroughputAdvice()
			if err != nil || !reflect.DeepEqual(f, got) {
				t.Fatalf("flags %d: got %#v want %#v error %v", flags, got, f, err)
			}
			for n := 0; n < len(want); n++ {
				if _, err = ie.ParseSCONEThroughputAdviceFields(want[:n]); err == nil {
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
			got, err = ie.ParseSCONEThroughputAdviceFields(extended)
			if err != nil || !reflect.DeepEqual(f, got) {
				t.Fatal("spares/extensions", err)
			}
			*got = ie.SCONEThroughputAdviceFields{Flags: 255}
			if err = got.UnmarshalBinary(want); err != nil || !reflect.DeepEqual(f, got) {
				t.Fatal("reused receiver", err)
			}
		}
	}
	checkNewIEFailures(t, ie.SCONEThroughputAdvice, func(i *ie.IE) (any, error) { return i.SCONEThroughputAdvice() }, 1)
	var f *ie.SCONEThroughputAdviceFields
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

func TestSCONEThroughputAdviceInvalidFields(t *testing.T) {
	valid := func() *ie.SCONEThroughputAdviceFields {
		return &ie.SCONEThroughputAdviceFields{Flags: 3, UplinkThroughputAdvice: 1, DownlinkThroughputAdvice: 1}
	}
	{
		f := valid()
		f.Flags |= 4
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags |= 4`)
		}
		assertNewIEInvalid(t, ie.NewSCONEThroughputAdvice(f.Flags, f.UplinkThroughputAdvice, f.DownlinkThroughputAdvice), ie.SCONEThroughputAdvice)
	}
	{
		f := valid()
		f.Flags &^= 1
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 1`)
		}
		assertNewIEInvalid(t, ie.NewSCONEThroughputAdvice(f.Flags, f.UplinkThroughputAdvice, f.DownlinkThroughputAdvice), ie.SCONEThroughputAdvice)
	}
	{
		f := valid()
		f.Flags &^= 2
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 2`)
		}
		assertNewIEInvalid(t, ie.NewSCONEThroughputAdvice(f.Flags, f.UplinkThroughputAdvice, f.DownlinkThroughputAdvice), ie.SCONEThroughputAdvice)
	}
}
