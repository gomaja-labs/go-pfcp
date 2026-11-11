package ie

// NewNodeLevelRuleID creates a Node Level Rule ID IE.
// TS 29.244 V20.1.0 clause 8.2.284.
func NewNodeLevelRuleID(v uint8) *IE {
	if v == 0 {
		return invalidRecentIE(NodeLevelRuleID, "Node Level Rule ID")
	}
	return New(NodeLevelRuleID, []byte{v})
}

// NodeLevelRuleID returns the Node Level Rule ID value.
// TS 29.244 V20.1.0 clause 8.2.284; clause 8.1.2 governs spare bits and extensions.
func (i *IE) NodeLevelRuleID() (uint8, error) {
	b, err := recentIEPayload(i, NodeLevelRuleID, 1, true)
	if err != nil {
		return 0, err
	}
	v := b[0] & 255

	if v == 0 {
		return 0, &DecodeError{Type: NodeLevelRuleID, Err: ErrInvalidValue}
	}
	return v, nil
}
