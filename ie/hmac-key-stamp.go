package ie

import (
	"bytes"
)

// NewHMACKeySTAMP creates a HMAC Key (STAMP) IE.
// TS 29.244 V20.1.0 clause 8.2.274.
// RFC 8762 section 4.4 defines the HMAC key usage.
func NewHMACKeySTAMP(v []byte) *IE {
	if len(v) > 65535 {
		return invalidRecentIE(HMACKeySTAMP, "HMAC Key (STAMP)")
	}
	return New(HMACKeySTAMP, bytes.Clone(v))
}

// HMACKeySTAMP returns the HMAC Key (STAMP) value.
// TS 29.244 V20.1.0 clause 8.2.274: the entire variable-length value is significant.
func (i *IE) HMACKeySTAMP() ([]byte, error) {
	b, err := recentIEPayload(i, HMACKeySTAMP, 1, false)
	if err != nil {
		return nil, err
	}
	return b, nil
}
