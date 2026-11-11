package ie

// HeaderHandlingCondition values are defined in TS 29.244 V20.1.0 clause 8.2.258.
const (
	HeaderHandlingConditionEveryMatch     uint8 = 0
	HeaderHandlingConditionFirstMatchOnly uint8 = 1
)

// NewHeaderHandlingCondition creates a Header Handling Condition IE.
// TS 29.244 V20.1.0 clause 8.2.258.
func NewHeaderHandlingCondition(v uint8) *IE {
	if v > 1 {
		return invalidRecentIE(HeaderHandlingCondition, "Header Handling Condition")
	}
	return New(HeaderHandlingCondition, []byte{v})
}

// HeaderHandlingCondition returns the Header Handling Condition value.
// TS 29.244 V20.1.0 clause 8.2.258; clause 8.1.2 governs spare bits and extensions.
func (i *IE) HeaderHandlingCondition() (uint8, error) {
	b, err := recentIEPayload(i, HeaderHandlingCondition, 1, true)
	if err != nil {
		return 0, err
	}
	v := b[0] & 15

	if v > 1 {
		return 0, &DecodeError{Type: HeaderHandlingCondition, Err: ErrInvalidValue}
	}
	return v, nil
}
