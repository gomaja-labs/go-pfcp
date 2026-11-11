package ie

// TrafficParameterMeasurementFlags encodes Traffic Parameter Measurement Indication; TS 29.244 V20.1.0 clause 8.2.223.
type TrafficParameterMeasurementFlags uint8

const (
	// TrafficParameterMeasurementIndicationULPMI is the ULPMI value defined in clause 8.2.223.
	TrafficParameterMeasurementIndicationULPMI TrafficParameterMeasurementFlags = 1
	// TrafficParameterMeasurementIndicationDLPMI is the DLPMI value defined in clause 8.2.223.
	TrafficParameterMeasurementIndicationDLPMI TrafficParameterMeasurementFlags = 2
	// TrafficParameterMeasurementIndicationN6JMI is the N6JMI value defined in clause 8.2.223.
	TrafficParameterMeasurementIndicationN6JMI TrafficParameterMeasurementFlags = 4
)

// NewTrafficParameterMeasurementIndication creates a Traffic Parameter Measurement Indication IE.
// TS 29.244 V20.1.0 clause 8.2.223.
func NewTrafficParameterMeasurementIndication(value TrafficParameterMeasurementFlags) *IE {
	if value > 7 {
		return newInvalidIE(TrafficParameterMeasurementIndication, ErrInvalidValue)
	}
	return New(TrafficParameterMeasurementIndication, []byte{byte(value)})
}

// TrafficParameterMeasurementIndication returns the Traffic Parameter Measurement Indication value, ignoring spare bits.
// TS 29.244 V20.1.0 clauses 8.2.223 and 8.1.2.
func (i *IE) TrafficParameterMeasurementIndication() (TrafficParameterMeasurementFlags, error) {
	b, err := addedValuePayload(i, TrafficParameterMeasurementIndication, 1, false)
	if err != nil {
		return 0, err
	}
	return TrafficParameterMeasurementFlags(b[0] & 7), nil
}
