package ie

// N6DelayMeasurementProtocols values are defined in TS 29.244 V20.1.0 clause 8.2.246.
const (
	// N6DelayMeasurementProtocolsTWAMP selects TWAMP.
	N6DelayMeasurementProtocolsTWAMP uint8 = 1
	// N6DelayMeasurementProtocolsOWAMP selects OWAMP.
	N6DelayMeasurementProtocolsOWAMP uint8 = 2
	// N6DelayMeasurementProtocolsSTAMP selects STAMP.
	N6DelayMeasurementProtocolsSTAMP uint8 = 4
	// N6DelayMeasurementProtocolsOther selects Other.
	N6DelayMeasurementProtocolsOther uint8 = 8
)

// NewN6DelayMeasurementProtocols creates a N6 Delay Measurement Protocols IE.
// TS 29.244 V20.1.0 clause 8.2.246.
func NewN6DelayMeasurementProtocols(v uint8) *IE {
	if v > 15 {
		return invalidRecentIE(N6DelayMeasurementProtocols, "N6 Delay Measurement Protocols")
	}
	return New(N6DelayMeasurementProtocols, []byte{v})
}

// N6DelayMeasurementProtocols returns the N6 Delay Measurement Protocols value.
// TS 29.244 V20.1.0 clause 8.2.246; clause 8.1.2 governs spare bits and extensions.
func (i *IE) N6DelayMeasurementProtocols() (uint8, error) {
	b, err := recentIEPayload(i, N6DelayMeasurementProtocols, 1, false)
	if err != nil {
		return 0, err
	}
	v := b[0] & 15

	if v > 15 {
		return 0, &DecodeError{Type: N6DelayMeasurementProtocols, Err: ErrInvalidValue}
	}
	return v, nil
}
