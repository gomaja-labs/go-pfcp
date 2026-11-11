package ie_test

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestReportingThresholds(t *testing.T) {
	for flags := uint8(0); flags <= 15; flags++ {
		for boundary := 0; boundary < 2; boundary++ {
			f := &ie.ReportingThresholdsFields{Flags: flags}
			want := []byte{flags}
			if flags&1 != 0 {
				f.DownlinkCongestion = uint16(boundary) * 10000
				v := uint64(f.DownlinkCongestion)
				want = append(want, byte(v>>8))
				want = append(want, byte(v>>0))
			}
			if flags&2 != 0 {
				f.UplinkCongestion = uint16(boundary) * 10000
				v := uint64(f.UplinkCongestion)
				want = append(want, byte(v>>8))
				want = append(want, byte(v>>0))
			}
			if flags&4 != 0 {
				f.DownlinkDataRate = uint64(boundary) * 1099511627775
				v := f.DownlinkDataRate
				want = append(want, byte(v>>32))
				want = append(want, byte(v>>24))
				want = append(want, byte(v>>16))
				want = append(want, byte(v>>8))
				want = append(want, byte(v>>0))
			}
			if flags&8 != 0 {
				f.UplinkDataRate = uint64(boundary) * 1099511627775
				v := f.UplinkDataRate
				want = append(want, byte(v>>32))
				want = append(want, byte(v>>24))
				want = append(want, byte(v>>16))
				want = append(want, byte(v>>8))
				want = append(want, byte(v>>0))
			}
			i := ie.NewReportingThresholds(flags, f.DownlinkCongestion, f.UplinkCongestion, f.DownlinkDataRate, f.UplinkDataRate)
			if flags == 0 {
				assertNewIEInvalid(t, i, ie.ReportingThresholds)
				continue
			}
			checkNewIEWire(t, i, want)
			got, err := i.ReportingThresholds()
			if err != nil || !reflect.DeepEqual(f, got) {
				t.Fatalf("flags %d: got %#v want %#v error %v", flags, got, f, err)
			}
			for n := 0; n < len(want); n++ {
				if _, err = ie.ParseReportingThresholdsFields(want[:n]); err == nil {
					t.Fatalf("accepted truncation %d flags %d", n, flags)
				}
			}
			for n := 0; n < f.MarshalLen(); n++ {
				if err = f.MarshalTo(make([]byte, n)); !errors.Is(err, io.ErrUnexpectedEOF) {
					t.Fatal("short Fields buffer", err)
				}
			}
			extended := append(append([]byte{}, want...), 0xa5, 0x5a)
			extended[0] |= 240
			got, err = ie.ParseReportingThresholdsFields(extended)
			if err != nil || !reflect.DeepEqual(f, got) {
				t.Fatal("spares/extensions", err)
			}
			*got = ie.ReportingThresholdsFields{Flags: 255}
			if err = got.UnmarshalBinary(want); err != nil || !reflect.DeepEqual(f, got) {
				t.Fatal("reused receiver", err)
			}
		}
	}
	checkNewIEFailures(t, ie.ReportingThresholds, func(i *ie.IE) (any, error) { return i.ReportingThresholds() }, 1)
	var f *ie.ReportingThresholdsFields
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

func TestReportingThresholdsInvalidFields(t *testing.T) {
	valid := func() *ie.ReportingThresholdsFields {
		return &ie.ReportingThresholdsFields{Flags: 15, DownlinkCongestion: 1, UplinkCongestion: 1, DownlinkDataRate: 1, UplinkDataRate: 1}
	}
	{
		f := valid()
		f.Flags |= 16
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags |= 16`)
		}
		assertNewIEInvalid(t, ie.NewReportingThresholds(f.Flags, f.DownlinkCongestion, f.UplinkCongestion, f.DownlinkDataRate, f.UplinkDataRate), ie.ReportingThresholds)
	}
	{
		f := valid()
		f.DownlinkCongestion = 10001
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.DownlinkCongestion=10001`)
		}
		assertNewIEInvalid(t, ie.NewReportingThresholds(f.Flags, f.DownlinkCongestion, f.UplinkCongestion, f.DownlinkDataRate, f.UplinkDataRate), ie.ReportingThresholds)
	}
	{
		f := valid()
		f.UplinkCongestion = 10001
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.UplinkCongestion=10001`)
		}
		assertNewIEInvalid(t, ie.NewReportingThresholds(f.Flags, f.DownlinkCongestion, f.UplinkCongestion, f.DownlinkDataRate, f.UplinkDataRate), ie.ReportingThresholds)
	}
	{
		f := valid()
		f.DownlinkDataRate = 1 << 40
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.DownlinkDataRate=1<<40`)
		}
		assertNewIEInvalid(t, ie.NewReportingThresholds(f.Flags, f.DownlinkCongestion, f.UplinkCongestion, f.DownlinkDataRate, f.UplinkDataRate), ie.ReportingThresholds)
	}
	{
		f := valid()
		f.UplinkDataRate = 1 << 40
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.UplinkDataRate=1<<40`)
		}
		assertNewIEInvalid(t, ie.NewReportingThresholds(f.Flags, f.DownlinkCongestion, f.UplinkCongestion, f.DownlinkDataRate, f.UplinkDataRate), ie.ReportingThresholds)
	}
	{
		f := valid()
		f.Flags &^= 1
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 1`)
		}
		assertNewIEInvalid(t, ie.NewReportingThresholds(f.Flags, f.DownlinkCongestion, f.UplinkCongestion, f.DownlinkDataRate, f.UplinkDataRate), ie.ReportingThresholds)
	}
	{
		f := valid()
		f.Flags &^= 2
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 2`)
		}
		assertNewIEInvalid(t, ie.NewReportingThresholds(f.Flags, f.DownlinkCongestion, f.UplinkCongestion, f.DownlinkDataRate, f.UplinkDataRate), ie.ReportingThresholds)
	}
	{
		f := valid()
		f.Flags &^= 4
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 4`)
		}
		assertNewIEInvalid(t, ie.NewReportingThresholds(f.Flags, f.DownlinkCongestion, f.UplinkCongestion, f.DownlinkDataRate, f.UplinkDataRate), ie.ReportingThresholds)
	}
	{
		f := valid()
		f.Flags &^= 8
		if _, err := f.Marshal(); err == nil {
			t.Fatal(`accepted invalid Fields: f.Flags &^= 8`)
		}
		assertNewIEInvalid(t, ie.NewReportingThresholds(f.Flags, f.DownlinkCongestion, f.UplinkCongestion, f.DownlinkDataRate, f.UplinkDataRate), ie.ReportingThresholds)
	}
}
