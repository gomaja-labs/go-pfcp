package ie

import (
	"bytes"
)

// NewBindingIndication creates a Binding Indication IE.
// TS 29.244 V20.1.0 clause 8.2.281.
// TS 29.500 clause 5.2.3.2.6 defines the opaque header value.
func NewBindingIndication(v []byte) *IE {
	if len(v) > 65535 {
		return invalidRecentIE(BindingIndication, "Binding Indication")
	}
	return New(BindingIndication, bytes.Clone(v))
}

// BindingIndication returns the Binding Indication value.
// TS 29.244 V20.1.0 clause 8.2.281: the entire variable-length value is significant.
func (i *IE) BindingIndication() ([]byte, error) {
	b, err := recentIEPayload(i, BindingIndication, 1, false)
	if err != nil {
		return nil, err
	}
	return b, nil
}
