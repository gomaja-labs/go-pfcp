package ie_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestUELevelMeasurementsConfiguration(t *testing.T) {
	for _, job := range []uint8{1, 2} {
		for _, count := range []int{0, 1, 9, 255} {
			for _, sec := range []uint32{0, 1, 0xffffff} {
				types := make([]uint8, count)
				for n := range types {
					types[n] = uint8(n%9) + 1
				}
				f := ie.NewUELevelMeasurementsConfigurationFields(job, types, time.Duration(sec)*time.Second)
				want := append([]byte{job, byte(count)}, types...)
				want = append(want, byte(sec>>16), byte(sec>>8), byte(sec))
				i := ie.NewUELevelMeasurementsConfiguration(job, types, time.Duration(sec)*time.Second)
				checkNewIEWire(t, i, want)
				checkNewFields(t, f, func(b []byte) (newFields, error) { return ie.ParseUELevelMeasurementsConfigurationFields(b) })
				got, err := i.UELevelMeasurementsConfiguration()
				if err != nil || !reflect.DeepEqual(f, got) {
					t.Fatal("fields", err)
				}
			}
		}
	}
	for _, f := range []*ie.UELevelMeasurementsConfigurationFields{
		{JobType: 0}, {JobType: 3}, {JobType: 1, MeasurementTypes: []uint8{0}}, {JobType: 1, MeasurementTypes: []uint8{10}},
		{JobType: 1, MeasurementTypes: make([]uint8, 256)}, {JobType: 1, TimePeriod: -1}, {JobType: 1, TimePeriod: time.Nanosecond}, {JobType: 1, TimePeriod: 0x1000000 * time.Second},
	} {
		assertNewIEInvalid(t, ie.NewUELevelMeasurementsConfiguration(f.JobType, f.MeasurementTypes, f.TimePeriod), ie.UELevelMeasurementsConfiguration)
	}
	checkNewFieldsNil(t, (*ie.UELevelMeasurementsConfigurationFields)(nil))
	checkNewIEFailures(t, ie.UELevelMeasurementsConfiguration, func(i *ie.IE) (any, error) { return i.UELevelMeasurementsConfiguration() }, 5)
}
