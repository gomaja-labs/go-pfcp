package ie

// NewHeaderValue creates a Header Value IE.
// The value is an octet string carried unchanged; no character set is imposed or checked.
// TS 29.244 V20.1.0 clause 8.2.257.
func NewHeaderValue(v string) *IE {
	if len(v) > 65535 {
		return invalidRecentIE(HeaderValue, "Header Value")
	}
	return New(HeaderValue, []byte(v))
}

// HeaderValue returns the Header Value value.
// TS 29.244 V20.1.0 clause 8.2.257: the entire variable-length value is significant.
func (i *IE) HeaderValue() (string, error) {
	b, err := recentIEPayload(i, HeaderValue, 1, false)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
