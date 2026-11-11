package ie_test

import (
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestThresholdsCurrent(t *testing.T) {
	for flags := 1; flags <= 3; flags++ {
		for _, rtt := range []time.Duration{0, time.Millisecond, 65535 * time.Millisecond} {
			for _, plr := range []uint8{0, 1, 100} {
				v := ie.ThresholdsFields{Flags: ie.ThresholdFlags(flags)}
				b := []byte{byte(flags)}
				if flags&1 != 0 {
					v.RTT = rtt
					b = append(b, byte(rtt/time.Millisecond>>8), byte(rtt/time.Millisecond))
				}
				if flags&2 != 0 {
					v.PacketLossRate = plr
					b = append(b, plr)
				}
				assertAddedFields(t, ie.NewThresholds(v), &v, b, "Thresholds", false)
			}
		}
	}
	for v := 101; v <= 255; v++ {
		f, err := ie.New(ie.Thresholds, []byte{2, byte(v)}).Thresholds()
		if err != nil || f.PacketLossRate != 0 {
			t.Fatal("out-of-range loss rate must mean zero")
		}
	}

	var nilFields *ie.ThresholdsFields
	assertAddedNilFields(t, nilFields)
}

func TestThresholdsInvalidFields(t *testing.T) {
	for _, v := range []ie.ThresholdsFields{{}, {Flags: 4}, {Flags: ie.ThresholdRTT, RTT: -1}, {Flags: ie.ThresholdRTT, RTT: time.Nanosecond}, {Flags: ie.ThresholdRTT, RTT: 65536 * time.Millisecond}, {Flags: ie.ThresholdPLR, PacketLossRate: 101}, {Flags: ie.ThresholdRTT, PacketLossRate: 1}} {
		assertAddedInvalid(t, ie.NewThresholds(v))
		if _, err := v.Marshal(); err == nil {
			t.Fatal("invalid fields accepted")
		}
	}
}

var _ = registerConstructorFixture(ie.NewThresholds, ie.ThresholdsFields{Flags: 1, RTT: time.Millisecond})
