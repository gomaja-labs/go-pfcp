package ie

// NewHeaderHandlingControlID creates a Header Handling Control ID IE.
// TS 29.244 V20.1.0 clause 8.2.259.
func NewHeaderHandlingControlID(v uint8) *IE {
	return New(HeaderHandlingControlID, []byte{v})
}

// HeaderHandlingControlID returns the Header Handling Control ID value.
// TS 29.244 V20.1.0 clause 8.2.259; clause 8.1.2 governs spare bits and extensions.
func (i *IE) HeaderHandlingControlID() (uint8, error) {
	b, err := recentIEPayload(i, HeaderHandlingControlID, 1, true)
	if err != nil {
		return 0, err
	}
	v := b[0] & 255

	return v, nil
}
