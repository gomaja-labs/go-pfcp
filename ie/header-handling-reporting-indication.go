package ie

// HeaderHandlingReportingIndication values are defined in TS 29.244 V20.1.0 clause 8.2.266.
const (
	HeaderHandlingReportingIndicationOTRI uint8 = 1
)

// NewHeaderHandlingReportingIndication creates a Header Handling Reporting Indication IE.
// TS 29.244 V20.1.0 clause 8.2.266.
func NewHeaderHandlingReportingIndication(v uint8) *IE {
	if v > 1 {
		return invalidRecentIE(HeaderHandlingReportingIndication, "Header Handling Reporting Indication")
	}
	return New(HeaderHandlingReportingIndication, []byte{v})
}

// HeaderHandlingReportingIndication returns the Header Handling Reporting Indication value.
// TS 29.244 V20.1.0 clause 8.2.266; clause 8.1.2 governs spare bits and extensions.
func (i *IE) HeaderHandlingReportingIndication() (uint8, error) {
	b, err := recentIEPayload(i, HeaderHandlingReportingIndication, 1, false)
	if err != nil {
		return 0, err
	}
	v := b[0] & 1

	if v > 1 {
		return 0, &DecodeError{Type: HeaderHandlingReportingIndication, Err: ErrInvalidValue}
	}
	return v, nil
}
