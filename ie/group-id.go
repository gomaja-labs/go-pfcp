package ie

// NewGroupID creates a Group ID IE. An empty value creates a null IE.
// TS 29.244 V20.1.0 clauses 8.2.198 and 8.1.2.
func NewGroupID(value []byte) *IE {
	if len(value) > 65535 {
		return newInvalidIE(GroupID, ErrInvalidValue)
	}
	return New(GroupID, append([]byte(nil), value...))
}

// GroupID returns the Group ID value.
// TS 29.244 V20.1.0 clause 8.2.198.
func (i *IE) GroupID() ([]byte, error) {
	b, err := addedValuePayload(i, GroupID, 1, false)
	if err != nil {
		return nil, err
	}

	return append([]byte(nil), b...), nil
}
