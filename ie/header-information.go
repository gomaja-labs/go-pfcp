package ie

import (
	"bytes"
)

// NewHeaderInformation creates a Header Information IE.
// TS 29.244 V20.1.0 clause 8.2.256.
func NewHeaderInformation(v []byte) *IE {
	if len(v) > 65535 {
		return invalidRecentIE(HeaderInformation, "Header Information")
	}
	return New(HeaderInformation, bytes.Clone(v))
}

// HeaderInformation returns the Header Information value.
// TS 29.244 V20.1.0 clause 8.2.256: the entire variable-length value is significant.
func (i *IE) HeaderInformation() ([]byte, error) {
	b, err := recentIEPayload(i, HeaderInformation, 1, false)
	if err != nil {
		return nil, err
	}
	return b, nil
}
