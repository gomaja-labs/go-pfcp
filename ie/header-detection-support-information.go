package ie

// NewHeaderDetectionSupportInformation creates a Header Detection Support Information IE.
// The value is an octet string carried unchanged; no character set is imposed or checked.
// TS 29.244 V20.1.0 clause 8.2.252.
func NewHeaderDetectionSupportInformation(v string) *IE {
	if len(v) > 65535 {
		return invalidRecentIE(HeaderDetectionSupportInformation, "Header Detection Support Information")
	}
	return New(HeaderDetectionSupportInformation, []byte(v))
}

// HeaderDetectionSupportInformation returns the Header Detection Support Information value.
// TS 29.244 V20.1.0 clause 8.2.252: the entire variable-length value is significant.
func (i *IE) HeaderDetectionSupportInformation() (string, error) {
	b, err := recentIEPayload(i, HeaderDetectionSupportInformation, 1, false)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
