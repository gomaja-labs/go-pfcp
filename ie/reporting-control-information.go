package ie

// ReportingControlInformation values are defined in TS 29.244 V20.1.0 clause 8.2.272.
const (
	// ReportingControlInformationUELM selects UELM.
	ReportingControlInformationUELM uint8 = 1
)

// NewReportingControlInformation creates a Reporting Control Information IE.
// TS 29.244 V20.1.0 clause 8.2.272.
func NewReportingControlInformation(v uint8) *IE {
	if v > 1 {
		return invalidRecentIE(ReportingControlInformation, "Reporting Control Information")
	}
	return New(ReportingControlInformation, []byte{v})
}

// ReportingControlInformation returns the Reporting Control Information value.
// TS 29.244 V20.1.0 clause 8.2.272; clause 8.1.2 governs spare bits and extensions.
func (i *IE) ReportingControlInformation() (uint8, error) {
	b, err := recentIEPayload(i, ReportingControlInformation, 1, false)
	if err != nil {
		return 0, err
	}
	v := b[0] & 1

	if v > 1 {
		return 0, &DecodeError{Type: ReportingControlInformation, Err: ErrInvalidValue}
	}
	return v, nil
}
