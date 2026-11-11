package ie

import (
	"unicode/utf8"
)

// NewOperatorConfigurableUPFCapability creates a Operator Configurable UPF Capability IE.
// TS 29.244 V20.1.0 clause 8.2.249.
// The value is the string OpConfiguredCapability in TS 29.571 clause 5.2.2.
func NewOperatorConfigurableUPFCapability(v string) *IE {
	if len(v) > 65535 || !utf8.ValidString(v) {
		return invalidRecentIE(OperatorConfigurableUPFCapability, "Operator Configurable UPF Capability")
	}
	return New(OperatorConfigurableUPFCapability, []byte(v))
}

// OperatorConfigurableUPFCapability returns the Operator Configurable UPF Capability value.
// TS 29.244 V20.1.0 clause 8.2.249: the entire variable-length value is significant.
func (i *IE) OperatorConfigurableUPFCapability() (string, error) {
	b, err := recentIEPayload(i, OperatorConfigurableUPFCapability, 1, false)
	if err != nil {
		return "", err
	}
	if !utf8.Valid(b) {
		return "", &DecodeError{Type: OperatorConfigurableUPFCapability, Err: ErrInvalidValue}
	}
	return string(b), nil
}
