package ie_test

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

type durationCase struct {
	name string
	new  func(time.Duration) *ie.IE
	get  func(*ie.IE) (time.Duration, error)
	unit time.Duration
	max  int64
}

func durationCases() []durationCase {
	return []durationCase{
		{"8.2.28/DownlinkDataNotificationDelay", ie.NewDownlinkDataNotificationDelay, (*ie.IE).DownlinkDataNotificationDelay, 50 * time.Millisecond, 255},
		{"8.2.83/UserPlaneInactivityTimer", ie.NewUserPlaneInactivityTimer, (*ie.IE).UserPlaneInactivityTimer, time.Second, math.MaxUint32},
		{"8.2.42/MeasurementPeriod", ie.NewMeasurementPeriod, (*ie.IE).MeasurementPeriod, time.Second, math.MaxUint32},
		{"8.2.14/TimeThreshold", ie.NewTimeThreshold, (*ie.IE).TimeThreshold, time.Second, math.MaxUint32},
		{"8.2.51/TimeQuota", ie.NewTimeQuota, (*ie.IE).TimeQuota, time.Second, math.MaxUint32},
		{"8.2.48/QuotaHoldingTime", ie.NewQuotaHoldingTime, (*ie.IE).QuotaHoldingTime, time.Second, math.MaxUint32},
		{"8.2.132/QuotaValidityTime", ie.NewQuotaValidityTime, (*ie.IE).QuotaValidityTime, time.Second, math.MaxUint32},
		{"8.2.17/SubsequentTimeThreshold", ie.NewSubsequentTimeThreshold, (*ie.IE).SubsequentTimeThreshold, time.Second, math.MaxUint32},
		{"8.2.87/SubsequentTimeQuota", ie.NewSubsequentTimeQuota, (*ie.IE).SubsequentTimeQuota, time.Second, math.MaxUint32},
		{"8.2.45/DurationMeasurement", ie.NewDurationMeasurement, (*ie.IE).DurationMeasurement, time.Second, math.MaxUint32},
		{"8.2.170/MinimumWaitTime", ie.NewMinimumWaitTime, (*ie.IE).MinimumWaitTime, time.Second, math.MaxUint32},
		{"8.2.105/EthernetInactivityTimer", ie.NewEthernetInactivityTimer, (*ie.IE).EthernetInactivityTimer, time.Second, math.MaxUint32},
		{"8.2.162/AveragePacketDelay", ie.NewAveragePacketDelay, (*ie.IE).AveragePacketDelay, time.Millisecond, math.MaxUint32},
		{"8.2.163/MinimumPacketDelay", ie.NewMinimumPacketDelay, (*ie.IE).MinimumPacketDelay, time.Millisecond, math.MaxUint32},
		{"8.2.164/MaximumPacketDelay", ie.NewMaximumPacketDelay, (*ie.IE).MaximumPacketDelay, time.Millisecond, math.MaxUint32},
		{"8.2.184/ValidityTimer", ie.NewValidityTimer, (*ie.IE).ValidityTimer, time.Second, math.MaxUint16},
		{"8.2.221/DLPeriodicity", ie.NewDLPeriodicity, (*ie.IE).DLPeriodicity, time.Millisecond, math.MaxUint32},
		{"8.2.222/ULPeriodicity", ie.NewULPeriodicity, (*ie.IE).ULPeriodicity, time.Millisecond, math.MaxUint32},
		{"8.2.81/TimeQuotaMechanism", func(d time.Duration) *ie.IE { return ie.NewTimeQuotaMechanism(ie.BTITCTP, d) }, func(i *ie.IE) (time.Duration, error) {
			b, err := i.TimeQuotaMechanism()
			if err != nil {
				return 0, err
			}
			return time.Duration(binary.BigEndian.Uint32(b[1:])) * time.Second, nil
		}, time.Second, math.MaxUint32},
	}
}

func TestDurationConstructorBoundaries(t *testing.T) {
	for _, tc := range durationCases() {
		t.Run(tc.name, func(t *testing.T) {
			for _, d := range []time.Duration{0, tc.unit, 31 * tc.unit, 32 * tc.unit, time.Duration(tc.max) * tc.unit, time.Duration(tc.max+1) * tc.unit, -time.Nanosecond, time.Nanosecond, tc.unit - time.Nanosecond, tc.unit + time.Nanosecond, 1500 * time.Millisecond, 12750 * time.Millisecond} {
				t.Run(d.String(), func(t *testing.T) {
					if d < 0 || d%tc.unit != 0 || d/tc.unit > time.Duration(tc.max) {
						assertInvalidConstructor(t, tc.new(d))
						return
					}
					b, err := tc.new(d).Marshal()
					if err != nil {
						t.Fatal(err)
					}
					parsed, err := ie.Parse(b)
					if err != nil {
						t.Fatal(err)
					}
					got, err := tc.get(parsed)
					if err != nil || got != d {
						t.Fatalf("got %v, %v; want %v", got, err, d)
					}
				})
			}
		})
	}
	for _, btit := range []uint8{2, 3, 4, 255} {
		assertInvalidConstructor(t, ie.NewTimeQuotaMechanism(btit, time.Second))
	}
	for _, btit := range []uint8{0, 1} {
		b, err := ie.NewTimeQuotaMechanism(btit, time.Second).Marshal()
		if err != nil || b[4] != btit {
			t.Fatalf("BTIT %d: %v", btit, err)
		}
	}
}

func timerCases() []struct {
	name string
	typ  ie.IEType
	new  func(time.Duration) *ie.IE
	get  func(*ie.IE) (time.Duration, error)
} {
	return []struct {
		name string
		typ  ie.IEType
		new  func(time.Duration) *ie.IE
		get  func(*ie.IE) (time.Duration, error)
	}{
		{"8.2.29/DLBufferingDuration", ie.DLBufferingDuration, ie.NewDLBufferingDuration, (*ie.IE).DLBufferingDuration},
		{"8.2.35/Timer", ie.Timer, ie.NewTimer, (*ie.IE).Timer},
		{"8.2.78/GracefulReleasePeriod", ie.GracefulReleasePeriod, ie.NewGracefulReleasePeriod, (*ie.IE).GracefulReleasePeriod},
	}
}

func TestTimerExactCoding(t *testing.T) {
	for _, tc := range timerCases() {
		t.Run(tc.name, func(t *testing.T) {
			for _, test := range []struct {
				d     time.Duration
				octet byte
			}{
				{0, 0}, {time.Duration(math.MaxInt64), 0xe0}, {2 * time.Second, 1}, {60 * time.Second, 0x21}, {62 * time.Second, 0x1f},
				{31 * time.Minute, 0x3f}, {310 * time.Minute, 0x5f}, {31 * time.Hour, 0x7f}, {310 * time.Hour, 0x9f},
			} {
				b, err := tc.new(test.d).Marshal()
				if err != nil {
					t.Fatal(err)
				}
				if b[4] != test.octet {
					t.Errorf("duration %v encoded incorrectly", test.d)
				}
				got, err := tc.get(tc.new(test.d))
				if err != nil || got != test.d {
					t.Fatalf("got %v, %v", got, err)
				}
			}
			for _, d := range []time.Duration{-time.Nanosecond, time.Nanosecond, time.Second, 64 * time.Second, 90 * time.Second, 32 * time.Minute, 320 * time.Minute, 32 * time.Hour, 320 * time.Hour, 400 * time.Hour, 12750 * time.Millisecond} {
				t.Run("invalid/"+d.String(), func(t *testing.T) { assertInvalidConstructor(t, tc.new(d)) })
			}
			// Every received octet, including unassigned minute units and noncanonical zeros.
			units := [8]time.Duration{2 * time.Second, time.Minute, 10 * time.Minute, time.Hour, 10 * time.Hour, time.Minute, time.Minute, 0}
			for n := 0; n < 256; n++ {
				want := time.Duration(n&31) * units[n>>5]
				if n>>5 == 7 {
					want = time.Duration(math.MaxInt64)
				}
				got, err := tc.get(ie.New(tc.typ, []byte{byte(n)}))
				if err != nil || got != want {
					t.Fatalf("octet %d: got %v, %v; want %v", n, got, err, want)
				}
				encoded := tc.new(got)
				if _, err := encoded.Marshal(); err != nil {
					t.Fatalf("octet %d: %v", n, err)
				}
				again, err := tc.get(encoded)
				if err != nil || again != got {
					t.Fatalf("octet %d: changed duration", n)
				}
			}
		})
	}
}

func TestSignedTimeOffsetDurations(t *testing.T) {
	// TS 29.244 V20.1.0 clauses 8.2.147 and 8.2.149: signed64 nanoseconds.
	for _, tc := range []durationCase{
		{name: "threshold", new: ie.NewTimeOffsetThreshold, get: (*ie.IE).TimeOffsetThreshold},
		{name: "measurement", new: ie.NewTimeOffsetMeasurement, get: (*ie.IE).TimeOffsetMeasurement},
	} {
		for _, d := range []time.Duration{math.MinInt64, -1, 0, 1, math.MaxInt64} {
			t.Run(fmt.Sprint(tc.name, d), func(t *testing.T) {
				b, err := tc.new(d).Marshal()
				if err != nil {
					t.Fatal(err)
				}
				parsed, err := ie.Parse(b)
				if err != nil {
					t.Fatal(err)
				}
				got, err := tc.get(parsed)
				if got != d || err != nil {
					t.Fatalf("got %v, %v", got, err)
				}
			})
		}
	}
}

func TestStructuredDurationBoundaries(t *testing.T) {
	for _, d := range []time.Duration{0, time.Second, time.Duration(0xffffff) * time.Second, time.Duration(0x1000000) * time.Second, -1, 1} {
		f := ie.NewUELevelMeasurementsConfigurationFields(1, []uint8{1}, d)
		if d < 0 || d%time.Second != 0 || d/time.Second > 0xffffff {
			assertInvalidConstructor(t, ie.NewUELevelMeasurementsConfiguration(1, []uint8{1}, d))
			if _, err := f.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
				t.Errorf("UE fields %v: %v", d, err)
			}
			if err := f.MarshalTo(make([]byte, 16)); !errors.Is(err, ie.ErrInvalidValue) {
				t.Errorf("UE fields MarshalTo %v: %v", d, err)
			}
		} else {
			b, err := f.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			got, err := ie.ParseUELevelMeasurementsConfigurationFields(b)
			if err != nil || got.TimePeriod != d {
				t.Fatalf("UE fields %v: %v", d, err)
			}
		}
	}
	for slot := 0; slot < 3; slot++ {
		for _, d := range []time.Duration{0, time.Millisecond, time.Duration(math.MaxUint32) * time.Millisecond, time.Duration(math.MaxUint32+1) * time.Millisecond, -1, 1} {
			values := [3]time.Duration{}
			values[slot] = d
			f := ie.NewN6DelayMeasurementInformationFields(3, values[0], values[1], values[2])
			if d < 0 || d%time.Millisecond != 0 || d/time.Millisecond > math.MaxUint32 {
				assertInvalidConstructor(t, ie.NewN6DelayMeasurementInformation(3, values[0], values[1], values[2]))
				if _, err := f.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
					t.Errorf("N6 slot %d %v: %v", slot, d, err)
				}
				if err := f.MarshalTo(make([]byte, 16)); !errors.Is(err, ie.ErrInvalidValue) {
					t.Errorf("N6 fields MarshalTo %v: %v", d, err)
				}
			} else {
				b, err := f.Marshal()
				if err != nil {
					t.Fatal(err)
				}
				got, err := ie.ParseN6DelayMeasurementInformationFields(b)
				if err != nil || *got != *f {
					t.Fatalf("N6 fields %v: %v", d, err)
				}
			}
		}
	}
}
