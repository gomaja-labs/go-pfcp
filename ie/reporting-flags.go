package ie

// ReportingFlagBits encodes Reporting Flags; TS 29.244 V20.1.0 clause 8.2.204.
type ReportingFlagBits uint8

const (
	// ReportingFlagsDUPL is the DUPL value defined in clause 8.2.204.
	ReportingFlagsDUPL ReportingFlagBits = 1
)

// NewReportingFlags creates a Reporting Flags IE.
// TS 29.244 V20.1.0 clause 8.2.204.
func NewReportingFlags(value ReportingFlagBits) *IE {
	if value > 1 {
		return newInvalidIE(ReportingFlags, ErrInvalidValue)
	}
	return New(ReportingFlags, []byte{byte(value)})
}

// ReportingFlags returns the Reporting Flags value, ignoring spare bits.
// TS 29.244 V20.1.0 clauses 8.2.204 and 8.1.2.
func (i *IE) ReportingFlags() (ReportingFlagBits, error) {
	b, err := addedValuePayload(i, ReportingFlags, 1, false)
	if err != nil {
		return 0, err
	}
	return ReportingFlagBits(b[0] & 1), nil
}
