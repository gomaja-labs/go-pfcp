package ie

// NewHeaderDetectionReference creates a Header Detection Reference IE.
// The value is an octet string carried unchanged; no character set is imposed or checked.
// TS 29.244 V20.1.0 clause 8.2.251.
func NewHeaderDetectionReference(v string) *IE {
	if len(v) > 65535 {
		return invalidRecentIE(HeaderDetectionReference, "Header Detection Reference")
	}
	return New(HeaderDetectionReference, []byte(v))
}

// HeaderDetectionReference returns the Header Detection Reference value.
// TS 29.244 V20.1.0 clause 8.2.251: the entire variable-length value is significant.
func (i *IE) HeaderDetectionReference() (string, error) {
	b, err := recentIEPayload(i, HeaderDetectionReference, 1, false)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
