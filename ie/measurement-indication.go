package ie

// MeasurementIndicationFlags encodes Measurement Indication; TS 29.244 V20.1.0 clause 8.2.231.
type MeasurementIndicationFlags uint8

const (
	// MeasurementIndicationDQFI is the DQFI value defined in clause 8.2.231.
	MeasurementIndicationDQFI MeasurementIndicationFlags = 1
)

// NewMeasurementIndication creates a Measurement Indication IE.
// TS 29.244 V20.1.0 clause 8.2.231.
func NewMeasurementIndication(value MeasurementIndicationFlags) *IE {
	if value > 1 {
		return newInvalidIE(MeasurementIndication, ErrInvalidValue)
	}
	return New(MeasurementIndication, []byte{byte(value)})
}

// MeasurementIndication returns the Measurement Indication value, ignoring spare bits.
// TS 29.244 V20.1.0 clauses 8.2.231 and 8.1.2.
func (i *IE) MeasurementIndication() (MeasurementIndicationFlags, error) {
	b, err := addedValuePayload(i, MeasurementIndication, 1, false)
	if err != nil {
		return 0, err
	}
	return MeasurementIndicationFlags(b[0] & 1), nil
}
