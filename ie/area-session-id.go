package ie

// NewAreaSessionID creates a Area Session ID IE.
// TS 29.244 V20.1.0 clause 8.2.213.
func NewAreaSessionID(value uint16) *IE {
	return New(AreaSessionID, []byte{byte(value >> 8), byte(value)})
}

// AreaSessionID returns the Area Session ID value.
// TS 29.244 V20.1.0 clauses 8.2.213 and 8.1.2.
func (i *IE) AreaSessionID() (uint16, error) {
	b, err := addedValuePayload(i, AreaSessionID, 2, true)
	if err != nil {
		return 0, err
	}
	return uint16(b[0])<<8 | uint16(b[1]), nil
}
