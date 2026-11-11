package ie

import (
	"bytes"
)

// NewHeaderHandlingControlReference creates a Header Handling Control Reference IE.
// TS 29.244 V20.1.0 clause 8.2.254.
func NewHeaderHandlingControlReference(v []byte) *IE {
	if len(v) > 65535 {
		return invalidRecentIE(HeaderHandlingControlReference, "Header Handling Control Reference")
	}
	return New(HeaderHandlingControlReference, bytes.Clone(v))
}

// HeaderHandlingControlReference returns the Header Handling Control Reference value.
// TS 29.244 V20.1.0 clause 8.2.254: the entire variable-length value is significant.
func (i *IE) HeaderHandlingControlReference() ([]byte, error) {
	b, err := recentIEPayload(i, HeaderHandlingControlReference, 1, false)
	if err != nil {
		return nil, err
	}
	return b, nil
}
