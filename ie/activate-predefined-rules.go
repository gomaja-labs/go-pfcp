// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewActivatePredefinedRules creates a new ActivatePredefinedRules IE.
// TS 29.244 V20.1.0 clause 8.2.72: no text repertoire is specified for this identifier.
func NewActivatePredefinedRules(name []byte) *IE {
	return New(ActivatePredefinedRules, name)
}

// ActivatePredefinedRules returns ActivatePredefinedRules as opaque octets if the type of IE matches.
func (i *IE) ActivatePredefinedRules() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case ActivatePredefinedRules:
		return i.payloadAtLeast(1)
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}
