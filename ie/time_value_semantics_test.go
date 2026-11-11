package ie_test

import (
	"encoding/binary"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestSignedCumulativeRateRatio(t *testing.T) {
	ctor := reflect.ValueOf(ie.NewCumulativeRateRatioMeasurement)
	for _, want := range []int32{-1 << 31, -1, 0, 1, 1<<31 - 1} {
		i := ctor.Call([]reflect.Value{reflect.ValueOf(want).Convert(ctor.Type().In(0))})[0].Interface().(*ie.IE)
		if binary.BigEndian.Uint32(i.Payload) != uint32(want) {
			t.Fatal("signed constructor changed wire bits")
		}
		for _, candidate := range []*ie.IE{i, ie.New(ie.CumulativeRateRatioMeasurement, append(append([]byte{}, i.Payload...), 0xaa))} {
			got := reflect.ValueOf(candidate).MethodByName("CumulativeRateRatioMeasurement").Call(nil)
			if got[0].Kind() != reflect.Int32 {
				t.Fatalf("measurement is %s; want int32", got[0].Kind())
			}
			if !got[1].IsNil() || got[0].Int() != int64(want) {
				t.Fatalf("measurement=%v error=%v; want %d", got[0], got[1], want)
			}
		}
	}
	for _, candidate := range []*ie.IE{nil, ie.New(0xffff, nil), ie.New(ie.CumulativeRateRatioMeasurement, nil), ie.New(ie.CumulativeRateRatioMeasurement, []byte{1, 2, 3})} {
		got := reflect.ValueOf(candidate).MethodByName("CumulativeRateRatioMeasurement").Call(nil)
		if got[1].IsNil() {
			t.Fatal("invalid measurement accepted")
		}
	}
}

func TestTimestampConstructorPolicy(t *testing.T) {
	ctors := []func(time.Time) *ie.IE{ie.NewMonitoringTime, ie.NewTimeOfFirstPacket, ie.NewTimeOfLastPacket, ie.NewStartTime, ie.NewEndTime, ie.NewRecoveryTimeStamp, ie.NewEventTimeStamp, ie.NewActivationTime, ie.NewDeactivationTime}
	dates := []struct {
		year                      int
		month                     time.Month
		day, hour, minute, second int
		valid                     bool
	}{
		{1900, 1, 1, 0, 0, 0, false}, {1969, 12, 31, 23, 59, 59, false}, {1970, 1, 1, 0, 0, 0, true},
		{2036, 2, 7, 6, 28, 15, true}, {2036, 2, 7, 6, 28, 16, true}, {2105, 1, 1, 0, 0, 0, true},
	}
	for _, ctor := range ctors {
		typ := ctor(time.Unix(0, 0)).Type
		t.Run(typ.String(), func(t *testing.T) {
			for _, d := range dates {
				want := time.Date(d.year, d.month, d.day, d.hour, d.minute, d.second, 0, time.UTC)
				i := ctor(want)
				if !d.valid {
					requireInvalidValue(t, want.String(), typ, i)
					continue
				}
				b, e := i.Marshal()
				if e != nil {
					t.Fatal(e)
				}
				decoded, e := ie.Parse(b)
				if e != nil {
					t.Fatal(e)
				}
				got := reflect.ValueOf(decoded).MethodByName(typ.String()).Call(nil)
				if !got[1].IsNil() || !got[0].Interface().(time.Time).Equal(want) {
					t.Fatalf("timestamp changed: want %v", want)
				}
			}
			requireInvalidValue(t, "fractional seconds", typ, ctor(time.Unix(1700000000, 1)))
			for _, i := range []*ie.IE{nil, ie.New(0xffff, nil), ie.New(typ, nil), ie.New(typ, []byte{1, 2, 3})} {
				got := reflect.ValueOf(i).MethodByName(typ.String()).Call(nil)
				if got[1].IsNil() {
					t.Fatal("invalid timestamp accepted")
				}
			}
			_, e := ctor(time.Unix(1700000000, 1)).Marshal()
			if !errors.Is(e, ie.ErrInvalidValue) {
				t.Fatal("fractional timestamp did not retain constructor error")
			}
		})
	}
}

func TestPacketRateStatusTimestamp(t *testing.T) {
	want := time.Date(2026, 1, 1, 0, 0, 0, 500000000, time.UTC)
	i := ie.NewPacketRateStatus(1, 1, 0, 0, 0, want)
	if i == nil {
		t.Fatal("nil constructor")
	}
	if len(i.Payload) != 11 || binary.BigEndian.Uint64(i.Payload[3:]) != uint64(uint32(want.Unix()+2208988800))<<32|0x80000000 {
		t.Fatal("timestamp is not NTP 32.32")
	}
	f, e := i.PacketRateStatus()
	if e != nil || !f.RateControlStatusValidityTime.Equal(want) {
		t.Fatalf("timestamp lost: %v", e)
	}
	for _, fraction := range []uint32{1, 2, 3, 0x80000000, 0xfffffffe, 0xffffffff} {
		p := append([]byte{}, i.Payload...)
		binary.BigEndian.PutUint32(p[7:], fraction)
		f, e := ie.ParsePacketRateStatusFields(p)
		if e != nil {
			t.Fatal(e)
		}
		b, e := f.Marshal()
		if e != nil || !reflect.DeepEqual(b, p) {
			t.Fatalf("fraction %d did not round trip, error=%v", fraction, e)
		}
		f.RateControlStatusValidityTime = want
		b, e = f.Marshal()
		if e != nil || binary.BigEndian.Uint32(b[7:]) != 0x80000000 {
			t.Fatalf("edited time did not replace retained fraction: %v", e)
		}
	}
	for _, ns := range []int{0, 1, 999999999, 123456789} {
		want := time.Unix(1700000000, int64(ns))
		f, e := ie.NewPacketRateStatus(1, 1, 0, 0, 0, want).PacketRateStatus()
		if e != nil || !f.RateControlStatusValidityTime.Equal(want) {
			t.Fatalf("nanoseconds %d lost, error=%v", ns, e)
		}
	}
	for _, p := range [][]byte{{0}, {4}, {0x80}} {
		f, e := ie.ParsePacketRateStatusFields(p)
		if e != nil || f.MarshalLen() != 1 {
			t.Fatalf("no-limit parse: %v", e)
		}
	}
	i = ie.NewPacketRateStatus(0, 0, 0, 0, 0, time.Time{})
	if i == nil {
		t.Fatal("nil no-limit IE")
	}
	if b, e := i.Marshal(); e != nil || len(b) != 5 {
		t.Fatalf("no-limit encoding error=%v", e)
	}
	for _, args := range []struct {
		flags            byte
		ul, aul, dl, adl uint16
		ts               time.Time
	}{
		{0, 1, 0, 0, 0, time.Time{}}, {0, 0, 1, 0, 0, time.Time{}}, {0, 0, 0, 1, 0, time.Time{}}, {0, 0, 0, 0, 1, time.Time{}},
		{1, 1, 0, 1, 0, want}, {1, 1, 0, 0, 1, want}, {2, 0, 1, 1, 0, want}, {2, 0, 0, 1, 1, want},
		{1, 1, 0, 0, 0, time.Unix(1<<32, 0)}, {0, 0, 0, 0, 0, want}, {1, 1, 1, 0, 0, want}, {1, 1, 0, 0, 0, time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)},
		{8, 0, 0, 0, 0, time.Time{}},
	} {
		requireInvalidValue(t, "packet-rate argument", ie.PacketRateStatus, ie.NewPacketRateStatus(args.flags, args.ul, args.aul, args.dl, args.adl, args.ts))
	}
}

func TestPacketRateStatusWireFractionBoundary(t *testing.T) {
	// Rounding the final NTP fraction can cross a second or the chosen era window.
	// The unchanged Fields value must still preserve all received fraction bits.
	for _, seconds := range []uint32{0, 2208988800, 2208988799, 0xffffffff} {
		p := make([]byte, 11)
		p[0] = 1
		binary.BigEndian.PutUint32(p[3:], seconds)
		binary.BigEndian.PutUint32(p[7:], 0xffffffff)
		f, e := ie.ParsePacketRateStatusFields(p)
		if e != nil {
			t.Fatal(e)
		}
		b, e := f.Marshal()
		if e != nil || !reflect.DeepEqual(b, p) {
			t.Fatalf("seconds %d lost exact wire fraction: %v", seconds, e)
		}
	}
}
