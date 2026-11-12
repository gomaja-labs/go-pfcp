// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// TS 29.244 V20.1.0 clause 8.2.161: no text repertoire is specified for this identifier.
// NewDataNetworkAccessIdentifier creates a new DataNetworkAccessIdentifier IE.
func NewDataNetworkAccessIdentifier(id []byte) *IE {
	return New(DataNetworkAccessIdentifier, id)
}

// DataNetworkAccessIdentifier returns DataNetworkAccessIdentifier as opaque octets if the type of IE matches.
func (i *IE) DataNetworkAccessIdentifier() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case DataNetworkAccessIdentifier:
		return i.payloadAtLeast(1)
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}
