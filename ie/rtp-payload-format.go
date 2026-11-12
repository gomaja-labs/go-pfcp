package ie

// RTPPayloadFormat values are defined in TS 29.244 V20.1.0 clause 8.2.237.
const (
	// RTPPayloadFormatH264 selects H264.
	RTPPayloadFormatH264 uint8 = 1
	// RTPPayloadFormatH265 selects H265.
	RTPPayloadFormatH265 uint8 = 2
)

// NewRTPPayloadFormat creates a RTP Payload Format IE.
// TS 29.244 V20.1.0 clause 8.2.237.
func NewRTPPayloadFormat(v uint8) *IE {
	if v != 1 && v != 2 {
		return invalidRecentIE(RTPPayloadFormat, "RTP Payload Format")
	}
	return New(RTPPayloadFormat, []byte{v})
}

// RTPPayloadFormat returns the RTP Payload Format value.
// TS 29.244 V20.1.0 clause 8.2.237; clause 8.1.2 governs spare bits and extensions.
func (i *IE) RTPPayloadFormat() (uint8, error) {
	b, err := recentIEPayload(i, RTPPayloadFormat, 1, true)
	if err != nil {
		return 0, err
	}
	v := b[0] & 255

	if v != 1 && v != 2 {
		return 0, &DecodeError{Type: RTPPayloadFormat, Err: ErrInvalidValue}
	}
	return v, nil
}
