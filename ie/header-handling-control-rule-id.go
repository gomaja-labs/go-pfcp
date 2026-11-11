package ie

// NewHeaderHandlingControlRuleID creates a Header Handling Control Rule ID IE.
// TS 29.244 V20.1.0 clause 8.2.260.
func NewHeaderHandlingControlRuleID(v uint8) *IE {
	return New(HeaderHandlingControlRuleID, []byte{v})
}

// HeaderHandlingControlRuleID returns the Header Handling Control Rule ID value.
// TS 29.244 V20.1.0 clause 8.2.260; clause 8.1.2 governs spare bits and extensions.
func (i *IE) HeaderHandlingControlRuleID() (uint8, error) {
	b, err := recentIEPayload(i, HeaderHandlingControlRuleID, 1, true)
	if err != nil {
		return 0, err
	}
	v := b[0] & 255

	return v, nil
}
