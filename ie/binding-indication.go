package ie

// NewBindingIndication creates a Binding Indication IE.
// The value is an octet string carried unchanged; no character set is imposed or checked.
// TS 29.244 V20.1.0 clause 8.2.281.
// TS 29.500 clause 5.2.3.2.6 defines the opaque header value.
func NewBindingIndication(v string) *IE {
	if len(v) > 65535 {
		return invalidRecentIE(BindingIndication, "Binding Indication")
	}
	return New(BindingIndication, []byte(v))
}

// BindingIndication returns the Binding Indication value.
// TS 29.244 V20.1.0 clause 8.2.281: the entire variable-length value is significant.
func (i *IE) BindingIndication() (string, error) {
	b, err := recentIEPayload(i, BindingIndication, 1, false)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
