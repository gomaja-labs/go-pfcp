package ie

// NewTLContainer creates a TL-Container IE. An empty value creates a null IE.
// TS 29.244 V20.1.0 clauses 8.2.230 and 8.1.2.
func NewTLContainer(value []byte) *IE {
	if len(value) > 65535 {
		return newInvalidIE(TLContainer, ErrInvalidValue)
	}
	return New(TLContainer, append([]byte(nil), value...))
}

// TLContainer returns the TL-Container value.
// TS 29.244 V20.1.0 clause 8.2.230.
func (i *IE) TLContainer() ([]byte, error) {
	b, err := addedValuePayload(i, TLContainer, 1, false)
	if err != nil {
		return nil, err
	}

	return append([]byte(nil), b...), nil
}
