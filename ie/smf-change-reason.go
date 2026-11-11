package ie

// SMFChangeReason values are defined in TS 29.244 V20.1.0 clause 8.2.267.
const (
	SMFChangeReasonMPSNR uint8 = 1
	SMFChangeReasonMPSSF uint8 = 2
	SMFChangeReasonMPSAI uint8 = 4
)

// NewSMFChangeReason creates a SMF Change Reason IE.
// TS 29.244 V20.1.0 clause 8.2.267.
func NewSMFChangeReason(v uint8) *IE {
	if v > 7 {
		return invalidRecentIE(SMFChangeReason, "SMF Change Reason")
	}
	return New(SMFChangeReason, []byte{v})
}

// SMFChangeReason returns the SMF Change Reason value.
// TS 29.244 V20.1.0 clause 8.2.267; clause 8.1.2 governs spare bits and extensions.
func (i *IE) SMFChangeReason() (uint8, error) {
	b, err := recentIEPayload(i, SMFChangeReason, 1, false)
	if err != nil {
		return 0, err
	}
	v := b[0] & 7

	if v > 7 {
		return 0, &DecodeError{Type: SMFChangeReason, Err: ErrInvalidValue}
	}
	return v, nil
}
