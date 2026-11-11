package ie

// SecurityModeSTAMP values are defined in TS 29.244 V20.1.0 clause 8.2.273.
const (
	SecurityModeSTAMPINTP uint8 = 1
	SecurityModeSTAMPCONP uint8 = 2
)

// NewSecurityModeSTAMP creates a Security Mode (STAMP) IE.
// TS 29.244 V20.1.0 clause 8.2.273.
// RFC 8762 sections 4.4 and 4.5 define authenticated and secure transport modes.
func NewSecurityModeSTAMP(v uint8) *IE {
	if v > 3 {
		return invalidRecentIE(SecurityModeSTAMP, "Security Mode (STAMP)")
	}
	return New(SecurityModeSTAMP, []byte{v})
}

// SecurityModeSTAMP returns the Security Mode (STAMP) value.
// TS 29.244 V20.1.0 clause 8.2.273; clause 8.1.2 governs spare bits and extensions.
func (i *IE) SecurityModeSTAMP() (uint8, error) {
	b, err := recentIEPayload(i, SecurityModeSTAMP, 1, true)
	if err != nil {
		return 0, err
	}
	v := b[0] & 3

	if v > 3 {
		return 0, &DecodeError{Type: SecurityModeSTAMP, Err: ErrInvalidValue}
	}
	return v, nil
}
