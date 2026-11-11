package ie

// NewRTPHeaderExtensionID creates a RTP Header Extension ID IE.
// TS 29.244 V20.1.0 clause 8.2.235.
func NewRTPHeaderExtensionID(value uint8) *IE {
	if value < 1 {
		return newInvalidIE(RTPHeaderExtensionID, ErrInvalidValue)
	}
	return New(RTPHeaderExtensionID, []byte{value})
}

// RTPHeaderExtensionID returns the RTP Header Extension ID value, ignoring spare bits.
// TS 29.244 V20.1.0 clauses 8.2.235 and 8.1.2.
func (i *IE) RTPHeaderExtensionID() (uint8, error) {
	b, err := addedValuePayload(i, RTPHeaderExtensionID, 1, true)
	if err != nil {
		return 0, err
	}
	v := b[0] & 255
	if v == 0 {
		return 0, &DecodeError{Type: i.Type, Err: ErrInvalidValue}
	}
	return v, nil
}
