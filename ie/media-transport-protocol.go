package ie

// MediaTransportProtocolValue encodes Media Transport Protocol; TS 29.244 V20.1.0 clause 8.2.233.
type MediaTransportProtocolValue uint8

const (
	// MediaTransportProtocolUnspecified is the Unspecified value defined in clause 8.2.233.
	MediaTransportProtocolUnspecified MediaTransportProtocolValue = 0
	// MediaTransportProtocolRTP is the RTP value defined in clause 8.2.233.
	MediaTransportProtocolRTP MediaTransportProtocolValue = 1
	// MediaTransportProtocolSRTP is the SRTP value defined in clause 8.2.233.
	MediaTransportProtocolSRTP MediaTransportProtocolValue = 2
	// MediaTransportProtocolMOQT is the MOQT value defined in clause 8.2.233.
	MediaTransportProtocolMOQT MediaTransportProtocolValue = 3
)

// NewMediaTransportProtocol creates a Media Transport Protocol IE.
// TS 29.244 V20.1.0 clause 8.2.233.
func NewMediaTransportProtocol(value MediaTransportProtocolValue) *IE {
	if value > 3 {
		return newInvalidIE(MediaTransportProtocol, ErrInvalidValue)
	}
	return New(MediaTransportProtocol, []byte{byte(value)})
}

// MediaTransportProtocol returns the Media Transport Protocol value, ignoring spare bits.
// TS 29.244 V20.1.0 clauses 8.2.233 and 8.1.2.
func (i *IE) MediaTransportProtocol() (MediaTransportProtocolValue, error) {
	b, err := addedValuePayload(i, MediaTransportProtocol, 1, false)
	if err != nil {
		return 0, err
	}
	return MediaTransportProtocolValue(b[0] & 15), nil
}
