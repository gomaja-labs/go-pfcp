// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewApplicationInstanceID creates a new ApplicationInstanceID IE.
// The value is an octet string carried unchanged; no character set is imposed or checked.
// TS 29.244 V20.1.0 clause 8.2.60: no text repertoire is specified for this identifier.
func NewApplicationInstanceID(id string) *IE {
	return newStringIE(ApplicationInstanceID, id)
}

// ApplicationInstanceID returns the unchanged octet string.
func (i *IE) ApplicationInstanceID() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}
	if i.Type != ApplicationInstanceID {
		return "", &InvalidTypeError{Type: i.Type}
	}
	b, err := i.payloadAtLeast(1)
	return string(b), err
}
