package ie_test

import (
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestTrafficParameterThresholdCurrent(t *testing.T) {
	for _, v := range []ie.TrafficParameterThresholdFields{{}, {Flags: ie.TrafficParameterThresholdDL}, {Flags: ie.TrafficParameterThresholdDL, DownlinkJitter: 0xffffffff * time.Millisecond}} {
		b := []byte{byte(v.Flags)}
		if v.Flags != 0 {
			b = append(b, byte(v.DownlinkJitter/time.Millisecond>>24), byte(v.DownlinkJitter/time.Millisecond>>16), byte(v.DownlinkJitter/time.Millisecond>>8), byte(v.DownlinkJitter/time.Millisecond))
		}
		assertAddedFields(t, ie.NewTrafficParameterThreshold(v), &v, b, "TrafficParameterThreshold", false)
	}

	var nilFields *ie.TrafficParameterThresholdFields
	assertAddedNilFields(t, nilFields)
}

func TestTrafficParameterThresholdInvalidFields(t *testing.T) {
	for _, v := range []ie.TrafficParameterThresholdFields{{Flags: 2}, {DownlinkJitter: time.Millisecond}, {Flags: ie.TrafficParameterThresholdDL, DownlinkJitter: -1}, {Flags: ie.TrafficParameterThresholdDL, DownlinkJitter: time.Nanosecond}, {Flags: ie.TrafficParameterThresholdDL, DownlinkJitter: 0x100000000 * time.Millisecond}} {
		assertAddedInvalid(t, ie.NewTrafficParameterThreshold(v))
		if _, err := v.Marshal(); err == nil {
			t.Fatal("invalid fields accepted")
		}
	}
}
