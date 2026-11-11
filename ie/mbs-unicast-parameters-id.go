package ie

// NewMBSUnicastParametersID creates a MBS Unicast Parameters ID IE.
// TS 29.244 V20.1.0 clause 8.2.210.
func NewMBSUnicastParametersID(value uint16) *IE {
	return New(MBSUnicastParametersID, []byte{byte(value >> 8), byte(value)})
}

// MBSUnicastParametersID returns the MBS Unicast Parameters ID value.
// TS 29.244 V20.1.0 clauses 8.2.210 and 8.1.2.
func (i *IE) MBSUnicastParametersID() (uint16, error) {
	b, err := addedValuePayload(i, MBSUnicastParametersID, 2, false)
	if err != nil {
		return 0, err
	}
	return uint16(b[0])<<8 | uint16(b[1]), nil
}
