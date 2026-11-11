package ie

import (
	"bytes"
	"encoding/base64"
)

// NewMetadata creates a Metadata IE from an already base64-encoded octet string.
// An empty value creates a null IE. TS 29.244 V20.1.0 clause 8.2.219;
// TS 29.571 clause 5.5.2 (Bytes).
func NewMetadata(value []byte) *IE {
	if len(value) > 65535 || !addedMetadataValid(value) {
		return newInvalidIE(Metadata, ErrInvalidValue)
	}
	return New(Metadata, bytes.Clone(value))
}

// Metadata returns the base64-encoded octet string, preserving its exact encoding.
// TS 29.244 V20.1.0 clause 8.2.219.
func (i *IE) Metadata() ([]byte, error) {
	b, err := addedValuePayload(i, Metadata, 1, false)
	if err != nil {
		return nil, err
	}
	if !addedMetadataValid(b) {
		return nil, &DecodeError{Type: Metadata, Err: ErrInvalidValue}
	}
	return bytes.Clone(b), nil
}
func addedMetadataValid(b []byte) bool {
	if bytes.ContainsAny(b, "\r\n") {
		return false
	}
	_, err := base64.StdEncoding.Strict().DecodeString(string(b))
	return err == nil
}
