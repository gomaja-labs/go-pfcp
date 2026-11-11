package ie

// NewReportingEndpointID creates a Reporting Endpoint ID IE.
// TS 29.244 V20.1.0 clause 8.2.253.
func NewReportingEndpointID(v uint8) *IE {
	return New(ReportingEndpointID, []byte{v})
}

// ReportingEndpointID returns the Reporting Endpoint ID value.
// TS 29.244 V20.1.0 clause 8.2.253; clause 8.1.2 governs spare bits and extensions.
func (i *IE) ReportingEndpointID() (uint8, error) {
	b, err := recentIEPayload(i, ReportingEndpointID, 1, true)
	if err != nil {
		return 0, err
	}
	v := b[0] & 255

	return v, nil
}
