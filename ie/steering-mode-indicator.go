package ie

// SteeringModeFlags encodes Steering Mode Indicator; TS 29.244 V20.1.0 clause 8.2.197.
type SteeringModeFlags uint8

const (
	// SteeringModeIndicatorALBI is the ALBI value defined in clause 8.2.197.
	SteeringModeIndicatorALBI SteeringModeFlags = 1
	// SteeringModeIndicatorUEAI is the UEAI value defined in clause 8.2.197.
	SteeringModeIndicatorUEAI SteeringModeFlags = 2
)

// NewSteeringModeIndicator creates a Steering Mode Indicator IE.
// TS 29.244 V20.1.0 clause 8.2.197.
func NewSteeringModeIndicator(value SteeringModeFlags) *IE {
	if value > 2 {
		return newInvalidIE(SteeringModeIndicator, ErrInvalidValue)
	}
	return New(SteeringModeIndicator, []byte{byte(value)})
}

// SteeringModeIndicator returns the Steering Mode Indicator value, ignoring spare bits.
// TS 29.244 V20.1.0 clauses 8.2.197 and 8.1.2.
func (i *IE) SteeringModeIndicator() (SteeringModeFlags, error) {
	b, err := addedValuePayload(i, SteeringModeIndicator, 1, false)
	if err != nil {
		return 0, err
	}
	v := b[0] & 3
	if v > 2 {
		return 0, &DecodeError{Type: i.Type, Err: ErrInvalidValue}
	}
	return SteeringModeFlags(v), nil
}
