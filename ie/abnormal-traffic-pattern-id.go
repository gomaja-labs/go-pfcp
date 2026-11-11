package ie

// AbnormalTrafficPatternID values are defined in TS 29.244 V20.1.0 clause 8.2.285.
const (
	AbnormalTrafficPatternIDMalformedPackets    uint8 = 0
	AbnormalTrafficPatternIDDuplicatePackets    uint8 = 1
	AbnormalTrafficPatternIDFragmentedPackets   uint8 = 2
	AbnormalTrafficPatternIDDDOS                uint8 = 3
	AbnormalTrafficPatternIDOperatorSpecificMin uint8 = 32
	AbnormalTrafficPatternIDOperatorSpecificMax uint8 = 63
)

// NewAbnormalTrafficPatternID creates a Abnormal Traffic Pattern ID IE.
// TS 29.244 V20.1.0 clause 8.2.285.
func NewAbnormalTrafficPatternID(v uint8) *IE {
	if v > 3 && (v < 32 || v > 63) {
		return invalidRecentIE(AbnormalTrafficPatternID, "Abnormal Traffic Pattern ID")
	}
	return New(AbnormalTrafficPatternID, []byte{v})
}

// AbnormalTrafficPatternID returns the Abnormal Traffic Pattern ID value.
// TS 29.244 V20.1.0 clause 8.2.285; clause 8.1.2 governs spare bits and extensions.
func (i *IE) AbnormalTrafficPatternID() (uint8, error) {
	b, err := recentIEPayload(i, AbnormalTrafficPatternID, 1, true)
	if err != nil {
		return 0, err
	}
	v := b[0] & 63

	if v > 3 && (v < 32 || v > 63) {
		return 0, &DecodeError{Type: AbnormalTrafficPatternID, Err: ErrInvalidValue}
	}
	return v, nil
}
