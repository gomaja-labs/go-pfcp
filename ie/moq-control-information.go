package ie

// MoQControlInformation values are defined in TS 29.244 V20.1.0 clause 8.2.269.
const (
	// MoQControlInformationMQREI selects MQREI.
	MoQControlInformationMQREI uint8 = 1
	// MoQControlInformationMQRAI selects MQRAI.
	MoQControlInformationMQRAI uint8 = 2
)

// NewMoQControlInformation creates a MoQ Control Information IE.
// TS 29.244 V20.1.0 clause 8.2.269.
func NewMoQControlInformation(v uint8) *IE {
	if v > 3 {
		return invalidRecentIE(MoQControlInformation, "MoQ Control Information")
	}
	return New(MoQControlInformation, []byte{v})
}

// MoQControlInformation returns the MoQ Control Information value.
// TS 29.244 V20.1.0 clause 8.2.269; clause 8.1.2 governs spare bits and extensions.
func (i *IE) MoQControlInformation() (uint8, error) {
	b, err := recentIEPayload(i, MoQControlInformation, 1, false)
	if err != nil {
		return 0, err
	}
	v := b[0] & 3

	if v > 3 {
		return 0, &DecodeError{Type: MoQControlInformation, Err: ErrInvalidValue}
	}
	return v, nil
}
