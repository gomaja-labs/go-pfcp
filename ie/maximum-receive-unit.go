package ie

// NewMaximumReceiveUnit creates a Maximum Receive Unit IE.
// TS 29.244 V20.1.0 clause 8.2.195.
func NewMaximumReceiveUnit(value uint16) *IE {
	return New(MaximumReceiveUnit, []byte{byte(value >> 8), byte(value)})
}

// MaximumReceiveUnit returns the Maximum Receive Unit value.
// TS 29.244 V20.1.0 clauses 8.2.195 and 8.1.2.
func (i *IE) MaximumReceiveUnit() (uint16, error) {
	b, err := addedValuePayload(i, MaximumReceiveUnit, 2, true)
	if err != nil {
		return 0, err
	}
	return uint16(b[0])<<8 | uint16(b[1]), nil
}
