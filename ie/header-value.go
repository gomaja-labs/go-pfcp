package ie

import (
	"bytes"
)

// NewHeaderValue creates a Header Value IE.
// TS 29.244 V20.1.0 clause 8.2.257.
func NewHeaderValue(v []byte) *IE {
	if len(v) > 65535 {
		return invalidRecentIE(HeaderValue, "Header Value")
	}
	return New(HeaderValue, bytes.Clone(v))
}

// HeaderValue returns the Header Value value.
// TS 29.244 V20.1.0 clause 8.2.257: the entire variable-length value is significant.
func (i *IE) HeaderValue() ([]byte, error) {
	b, err := recentIEPayload(i, HeaderValue, 1, false)
	if err != nil {
		return nil, err
	}
	return b, nil
}
