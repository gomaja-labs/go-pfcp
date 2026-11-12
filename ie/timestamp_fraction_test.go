package ie_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestTimestampFractionDiscarded(t *testing.T) {
	constructors := []func(time.Time) *ie.IE{ie.NewRecoveryTimeStamp, ie.NewStartTime, ie.NewEndTime, ie.NewTimeOfFirstPacket, ie.NewTimeOfLastPacket, ie.NewMonitoringTime, ie.NewActivationTime, ie.NewDeactivationTime, ie.NewEventTimeStamp}
	for _, ctor := range constructors {
		typ := ctor(time.Unix(0, 0)).Type
		t.Run(typ.String(), func(t *testing.T) {
			for _, sec := range []int64{0, 1700000000, 2085978495, 2085978496, 4294967294, 4294967295} {
				for _, ns := range []int64{1, 500000000, 999999999} {
					i := ctor(time.Unix(sec, ns))
					b, e := i.Marshal()
					if e != nil {
						t.Fatalf("fraction rejected: %v", e)
					}
					parsed, e := ie.Parse(b)
					if e != nil {
						t.Fatal(e)
					}
					values := reflect.ValueOf(parsed).MethodByName(typ.String()).Call(nil)
					if !values[1].IsNil() || !values[0].Interface().(time.Time).Equal(time.Unix(sec, 0)) {
						t.Fatal("fraction not discarded at same whole second")
					}
				}
			}
			for _, ts := range []time.Time{time.Unix(-1, 999999999), time.Unix(4294967296, 0)} {
				_, e := ctor(ts).Marshal()
				if !errors.Is(e, ie.ErrInvalidValue) {
					t.Fatalf("out-of-window time accepted: %v", ts)
				}
			}
		})
	}
	if _, err := ie.NewRecoveryTimeStamp(time.Now()).Marshal(); err != nil {
		t.Fatal(err)
	}
}
