package ie

import (
	"bytes"
)

// NewHeaderDetectionReference creates a Header Detection Reference IE.
// TS 29.244 V20.1.0 clause 8.2.251.
func NewHeaderDetectionReference(v []byte) *IE {
	if len(v) > 65535 {
		return invalidRecentIE(HeaderDetectionReference, "Header Detection Reference")
	}
	return New(HeaderDetectionReference, bytes.Clone(v))
}

// HeaderDetectionReference returns the Header Detection Reference value.
// TS 29.244 V20.1.0 clause 8.2.251: the entire variable-length value is significant.
func (i *IE) HeaderDetectionReference() ([]byte, error) {
	b, err := recentIEPayload(i, HeaderDetectionReference, 1, false)
	if err != nil {
		return nil, err
	}
	return b, nil
}
