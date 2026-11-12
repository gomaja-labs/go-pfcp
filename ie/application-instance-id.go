// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// TS 29.244 V20.1.0 clause 8.2.60: no text repertoire is specified for this identifier.
// NewApplicationInstanceID creates a new ApplicationInstanceID IE.
func NewApplicationInstanceID(id []byte) *IE {
	return New(ApplicationInstanceID, id)
}

// ApplicationInstanceID returns ApplicationInstanceID as opaque octets if the type of IE matches.
func (i *IE) ApplicationInstanceID() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case ApplicationInstanceID:
		return i.payloadAtLeast(1)
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}
