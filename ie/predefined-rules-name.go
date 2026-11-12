package ie

// NewPredefinedRulesName creates a Predefined Rules Name IE. An empty value creates a null IE.
// The value is an octet string carried unchanged; no character set is imposed or checked.
// TS 29.244 V20.1.0 clauses 8.2.205 and 8.1.2.
func NewPredefinedRulesName(value string) *IE {
	if len(value) > 65535 {
		return newInvalidIE(PredefinedRulesName, ErrInvalidValue)
	}
	return New(PredefinedRulesName, []byte(value))
}

// PredefinedRulesName returns the Predefined Rules Name value.
// TS 29.244 V20.1.0 clause 8.2.205.
func (i *IE) PredefinedRulesName() (string, error) {
	b, err := addedValuePayload(i, PredefinedRulesName, 1, false)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
