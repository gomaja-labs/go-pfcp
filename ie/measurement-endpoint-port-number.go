package ie

import "encoding/binary"

// NewMeasurementEndpointPortNumber creates a Measurement Endpoint Port Number IE.
// TS 29.244 V20.1.0 clause 8.2.265.
func NewMeasurementEndpointPortNumber(v uint16) *IE {
	return New(MeasurementEndpointPortNumber, []byte{byte(v >> 8), byte(v)})
}

// MeasurementEndpointPortNumber returns the Measurement Endpoint Port Number value.
// TS 29.244 V20.1.0 clause 8.2.265; clause 8.1.2 governs spare bits and extensions.
func (i *IE) MeasurementEndpointPortNumber() (uint16, error) {
	b, err := recentIEPayload(i, MeasurementEndpointPortNumber, 2, true)
	if err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint16(b) & 65535

	return v, nil
}
