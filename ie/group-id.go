package ie

// NewGroupID creates a Group ID IE. An empty value creates a null IE.
// The value is an octet string carried unchanged; no character set is imposed or checked.
// TS 29.244 V20.1.0 clauses 8.2.198 and 8.1.2.
func NewGroupID(value string) *IE {
	if len(value) > 65535 {
		return newInvalidIE(GroupID, ErrInvalidValue)
	}
	return New(GroupID, []byte(value))
}

// GroupID returns the Group ID value.
// TS 29.244 V20.1.0 clause 8.2.198.
func (i *IE) GroupID() (string, error) {
	b, err := addedValuePayload(i, GroupID, 1, false)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
