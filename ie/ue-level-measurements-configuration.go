package ie

import (
	"io"
	"time"
)

// Job and measurement values, TS 29.244 V20.1.0 tables 8.2.245-1 and 8.2.245-2.
const (
	// UELevelMeasurementsConfigurationMeasurementsOnly selects Measurements Only.
	UELevelMeasurementsConfigurationMeasurementsOnly uint8 = 1
	// UELevelMeasurementsConfigurationTraceAndMeasurements selects Trace And Measurements.
	UELevelMeasurementsConfigurationTraceAndMeasurements uint8 = 2
	// UELevelMeasurementsConfigurationAverageDLPSAUPFToUE selects Average DLPSAUPF To UE.
	UELevelMeasurementsConfigurationAverageDLPSAUPFToUE uint8 = 1
	// UELevelMeasurementsConfigurationAverageULPSAUPFToUEExcludingD1 selects Average ULPSAUPF To UE Excluding D1.
	UELevelMeasurementsConfigurationAverageULPSAUPFToUEExcludingD1 uint8 = 2
	// UELevelMeasurementsConfigurationAverageULPSAUPFToUEIncludingD1 selects Average ULPSAUPF To UE Including D1.
	UELevelMeasurementsConfigurationAverageULPSAUPFToUEIncludingD1 uint8 = 3
	// UELevelMeasurementsConfigurationAverageULPSAUPFToNGRAN selects Average ULPSAUPF To NGRAN.
	UELevelMeasurementsConfigurationAverageULPSAUPFToNGRAN uint8 = 4
	// UELevelMeasurementsConfigurationAverageDLPSAUPFToNGRAN selects Average DLPSAUPF To NGRAN.
	UELevelMeasurementsConfigurationAverageDLPSAUPFToNGRAN uint8 = 5
	// UELevelMeasurementsConfigurationAverageRTTPSAUPFToNGRAN selects Average RTTPSAUPF To NGRAN.
	UELevelMeasurementsConfigurationAverageRTTPSAUPFToNGRAN uint8 = 6
	// UELevelMeasurementsConfigurationAverageRTTPSAUPFToNGRANForSNSSAI selects Average RTTPSAUPF To NGRAN For SNSSAI.
	UELevelMeasurementsConfigurationAverageRTTPSAUPFToNGRANForSNSSAI uint8 = 7
	// UELevelMeasurementsConfigurationAverageRTTPSAUPFToNGRANForQFI selects Average RTTPSAUPF To NGRAN For QFI.
	UELevelMeasurementsConfigurationAverageRTTPSAUPFToNGRANForQFI uint8 = 8
	// UELevelMeasurementsConfigurationAverageRTTPSAUPFToNGRANForSNSSAIAndQFI selects Average RTTPSAUPF To NGRAN For SNSSAI And QFI.
	UELevelMeasurementsConfigurationAverageRTTPSAUPFToNGRANForSNSSAIAndQFI uint8 = 9
)

// NewUELevelMeasurementsConfiguration creates a UE Level Measurements Configuration IE.
// TS 29.244 V20.1.0 clause 8.2.245.
func NewUELevelMeasurementsConfiguration(jobType uint8, measurementTypes []uint8, timePeriod time.Duration) *IE {
	return recentFieldsIE(UELevelMeasurementsConfiguration, NewUELevelMeasurementsConfigurationFields(jobType, append([]uint8(nil), measurementTypes...), timePeriod))
}

// UELevelMeasurementsConfiguration returns fields defined in TS 29.244 V20.1.0 clause 8.2.245.
func (i *IE) UELevelMeasurementsConfiguration() (*UELevelMeasurementsConfigurationFields, error) {
	b, err := recentIEPayload(i, UELevelMeasurementsConfiguration, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseUELevelMeasurementsConfigurationFields(b)
	if err != nil {
		return nil, &DecodeError{Type: UELevelMeasurementsConfiguration, Err: err}
	}
	return f, nil
}

// UELevelMeasurementsConfigurationFields models TS 29.244 V20.1.0 clause 8.2.245.
type UELevelMeasurementsConfigurationFields struct {
	// JobType selects measurements only or trace and measurements.
	JobType uint8
	// MeasurementTypes contains the requested measurement bit mask.
	MeasurementTypes []uint8
	// TimePeriod is the measurement time period.
	TimePeriod time.Duration
}

// NewUELevelMeasurementsConfigurationFields retains all supplied arguments; Marshal validates them.
func NewUELevelMeasurementsConfigurationFields(jobType uint8, measurementTypes []uint8, timePeriod time.Duration) *UELevelMeasurementsConfigurationFields {
	return &UELevelMeasurementsConfigurationFields{jobType, append([]uint8(nil), measurementTypes...), timePeriod}
}

// ParseUELevelMeasurementsConfigurationFields parses the known fields and ignores trailing extension octets.
func ParseUELevelMeasurementsConfigurationFields(b []byte) (*UELevelMeasurementsConfigurationFields, error) {
	f := &UELevelMeasurementsConfigurationFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}
func (f *UELevelMeasurementsConfigurationFields) validate() error {
	if f == nil {
		return ErrIENotFound
	}
	if f.JobType < 1 || f.JobType > 2 || len(f.MeasurementTypes) > 255 || f.TimePeriod < 0 || f.TimePeriod%time.Second != 0 || f.TimePeriod/time.Second > 0xffffff {
		return ErrInvalidValue
	}
	for _, v := range f.MeasurementTypes {
		if v < 1 || v > 9 {
			return ErrInvalidValue
		}
	}
	return nil
}

// MarshalLen returns the known field size in octets.
func (f *UELevelMeasurementsConfigurationFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 5 + len(f.MeasurementTypes)
}

// Marshal serializes the fields.
func (f *UELevelMeasurementsConfigurationFields) Marshal() ([]byte, error) {
	if err := f.validate(); err != nil {
		return nil, err
	}
	b := make([]byte, f.MarshalLen())
	if err := f.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes the fields into b, rejecting short buffers and invalid values.
func (f *UELevelMeasurementsConfigurationFields) MarshalTo(b []byte) error {
	if err := f.validate(); err != nil {
		return err
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}
	b[0] = f.JobType
	b[1] = byte(len(f.MeasurementTypes))
	copy(b[2:], f.MeasurementTypes)
	putRecentUint(b[2+len(f.MeasurementTypes):f.MarshalLen()], uint64(f.TimePeriod/time.Second))
	return nil
}

// UnmarshalBinary replaces the receiver with the decoded fields.
func (f *UELevelMeasurementsConfigurationFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = UELevelMeasurementsConfigurationFields{}
	if len(b) < 2 {
		return io.ErrUnexpectedEOF
	}
	n := int(b[1])
	if len(b) < 5+n {
		return io.ErrUnexpectedEOF
	}
	f.JobType = b[0]
	f.MeasurementTypes = append([]uint8(nil), b[2:2+n]...)
	f.TimePeriod = time.Duration(recentUint(b[2+n:5+n])) * time.Second
	return f.validate()
}
