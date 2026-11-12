// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// TS 29.244 V20.1.0 clause 8.2.73: no text repertoire is specified for this identifier.
// NewDeactivatePredefinedRules creates a new DeactivatePredefinedRules IE.
func NewDeactivatePredefinedRules(name []byte) *IE {
	return New(DeactivatePredefinedRules, name)
}

// DeactivatePredefinedRules returns DeactivatePredefinedRules as opaque octets if the type of IE matches.
func (i *IE) DeactivatePredefinedRules() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case DeactivatePredefinedRules:
		return i.payloadAtLeast(1)
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}
