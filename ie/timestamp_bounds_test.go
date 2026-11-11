package ie_test

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestTimestampEraWindow(t *testing.T) {
	constructors := []func(time.Time) *ie.IE{ie.NewRecoveryTimeStamp, ie.NewStartTime, ie.NewEndTime, ie.NewTimeOfFirstPacket, ie.NewTimeOfLastPacket, ie.NewMonitoringTime, ie.NewActivationTime, ie.NewDeactivationTime, ie.NewEventTimeStamp}
	for _, ctor := range constructors {
		for _, unix := range []int64{-1, 1 << 32} {
			if ctor(time.Unix(unix, 0)) != nil {
				t.Error("out-of-window timestamp wrapped")
			}
		}
		for _, unix := range []int64{0, 2085978495, 2085978496, 4294967295} {
			i := ctor(time.Unix(unix, 999))
			if i == nil {
				t.Fatal("valid timestamp rejected")
			}
			v := binary.BigEndian.Uint32(i.Payload)
			if v != uint32(unix+2208988800) {
				t.Error("wrong NTP seconds")
			}
			got := reflectTimestamp(t, i)
			if got.Unix() != unix {
				t.Error("timestamp era lost", unix, got.Unix())
			}
		}
	}
	for _, v := range []uint32{0, 2208988799, 2208988800, 4294967295} {
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, v)
		got, e := ie.New(ie.RecoveryTimeStamp, b).RecoveryTimeStamp()
		if e != nil {
			t.Fatal(e)
		}
		want := int64(v) - 2208988800
		if v < 2208988800 {
			want += 1 << 32
		}
		if got.Unix() != want {
			t.Error("era rule", v, got.Unix(), want)
		}
	}
}

func reflectTimestamp(t *testing.T, i *ie.IE) time.Time {
	t.Helper()
	// These types share the same seconds encoding; RecoveryTimeStamp exposes it.
	got, e := ie.New(ie.RecoveryTimeStamp, i.Payload).RecoveryTimeStamp()
	if e != nil {
		t.Fatal(e)
	}
	return got
}
