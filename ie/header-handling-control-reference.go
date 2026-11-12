package ie

// NewHeaderHandlingControlReference creates a Header Handling Control Reference IE.
// The value is an octet string carried unchanged; no character set is imposed or checked.
// TS 29.244 V20.1.0 clause 8.2.254.
func NewHeaderHandlingControlReference(v string) *IE {
	if len(v) > 65535 {
		return invalidRecentIE(HeaderHandlingControlReference, "Header Handling Control Reference")
	}
	return New(HeaderHandlingControlReference, []byte(v))
}

// HeaderHandlingControlReference returns the Header Handling Control Reference value.
// TS 29.244 V20.1.0 clause 8.2.254: the entire variable-length value is significant.
func (i *IE) HeaderHandlingControlReference() (string, error) {
	b, err := recentIEPayload(i, HeaderHandlingControlReference, 1, false)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
