package ie

// MeasurementReportingType values are defined in TS 29.244 V20.1.0 clause 8.2.262.
const (
	MeasurementReportingTypePER uint8 = 1
	MeasurementReportingTypeEVT uint8 = 2
	MeasurementReportingTypeIRE uint8 = 4
)

// NewMeasurementReportingType creates a Measurement Reporting Type IE.
// TS 29.244 V20.1.0 clause 8.2.262.
func NewMeasurementReportingType(v uint8) *IE {
	if v == 0 || v > 7 || v&3 == 3 {
		return invalidRecentIE(MeasurementReportingType, "Measurement Reporting Type")
	}
	return New(MeasurementReportingType, []byte{v})
}

// MeasurementReportingType returns the Measurement Reporting Type value.
// TS 29.244 V20.1.0 clause 8.2.262; clause 8.1.2 governs spare bits and extensions.
func (i *IE) MeasurementReportingType() (uint8, error) {
	b, err := recentIEPayload(i, MeasurementReportingType, 1, false)
	if err != nil {
		return 0, err
	}
	v := b[0] & 7

	if v == 0 || v > 7 || v&3 == 3 {
		return 0, &DecodeError{Type: MeasurementReportingType, Err: ErrInvalidValue}
	}
	return v, nil
}
