package ie

// RemainingDataReportingIndication values are defined in TS 29.244 V20.1.0 clause 8.2.277.
const (
	// RemainingDataReportingIndicationDiscard selects Discard.
	RemainingDataReportingIndicationDiscard uint8 = 0
	// RemainingDataReportingIndicationSend selects Send.
	RemainingDataReportingIndicationSend uint8 = 1
)

// NewRemainingDataReportingIndication creates a Remaining Data Reporting Indication IE.
// TS 29.244 V20.1.0 clause 8.2.277.
func NewRemainingDataReportingIndication(v uint8) *IE {
	if v > 1 {
		return invalidRecentIE(RemainingDataReportingIndication, "Remaining Data Reporting Indication")
	}
	return New(RemainingDataReportingIndication, []byte{v})
}

// RemainingDataReportingIndication returns the Remaining Data Reporting Indication value.
// TS 29.244 V20.1.0 clause 8.2.277; clause 8.1.2 governs spare bits and extensions.
func (i *IE) RemainingDataReportingIndication() (uint8, error) {
	b, err := recentIEPayload(i, RemainingDataReportingIndication, 1, false)
	if err != nil {
		return 0, err
	}
	v := b[0] & 15

	if v > 1 {
		return 0, &DecodeError{Type: RemainingDataReportingIndication, Err: ErrInvalidValue}
	}
	return v, nil
}
