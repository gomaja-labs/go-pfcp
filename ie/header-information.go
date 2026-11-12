package ie

// NewHeaderInformation creates a Header Information IE.
// The value is an octet string carried unchanged; no character set is imposed or checked.
// TS 29.244 V20.1.0 clause 8.2.256.
func NewHeaderInformation(v string) *IE {
	if len(v) > 65535 {
		return invalidRecentIE(HeaderInformation, "Header Information")
	}
	return New(HeaderInformation, []byte(v))
}

// HeaderInformation returns the Header Information value.
// TS 29.244 V20.1.0 clause 8.2.256: the entire variable-length value is significant.
func (i *IE) HeaderInformation() (string, error) {
	b, err := recentIEPayload(i, HeaderInformation, 1, false)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
