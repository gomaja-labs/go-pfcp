package ie

import (
	"bytes"
)

// NewHeaderDetectionSupportInformation creates a Header Detection Support Information IE.
// TS 29.244 V20.1.0 clause 8.2.252.
func NewHeaderDetectionSupportInformation(v []byte) *IE {
	if len(v) > 65535 {
		return invalidRecentIE(HeaderDetectionSupportInformation, "Header Detection Support Information")
	}
	return New(HeaderDetectionSupportInformation, bytes.Clone(v))
}

// HeaderDetectionSupportInformation returns the Header Detection Support Information value.
// TS 29.244 V20.1.0 clause 8.2.252: the entire variable-length value is significant.
func (i *IE) HeaderDetectionSupportInformation() ([]byte, error) {
	b, err := recentIEPayload(i, HeaderDetectionSupportInformation, 1, false)
	if err != nil {
		return nil, err
	}
	return b, nil
}
