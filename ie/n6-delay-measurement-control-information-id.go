package ie

import "encoding/binary"

// NewN6DelayMeasurementControlInformationID creates a N6 Delay Measurement Control Information ID IE.
// TS 29.244 V20.1.0 clause 8.2.264.
func NewN6DelayMeasurementControlInformationID(v uint16) *IE {
	return New(N6DelayMeasurementControlInformationID, []byte{byte(v >> 8), byte(v)})
}

// N6DelayMeasurementControlInformationID returns the N6 Delay Measurement Control Information ID value.
// TS 29.244 V20.1.0 clause 8.2.264; clause 8.1.2 governs spare bits and extensions.
func (i *IE) N6DelayMeasurementControlInformationID() (uint16, error) {
	b, err := recentIEPayload(i, N6DelayMeasurementControlInformationID, 2, true)
	if err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint16(b) & 65535

	return v, nil
}
