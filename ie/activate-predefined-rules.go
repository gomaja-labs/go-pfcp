// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewActivatePredefinedRules creates a new ActivatePredefinedRules IE.
// The value is an octet string carried unchanged; no character set is imposed or checked.
// TS 29.244 V20.1.0 clause 8.2.72: no text repertoire is specified for this identifier.
func NewActivatePredefinedRules(name string) *IE {
	return newStringIE(ActivatePredefinedRules, name)
}

// ActivatePredefinedRules returns the unchanged octet string.
func (i *IE) ActivatePredefinedRules() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}
	if i.Type != ActivatePredefinedRules {
		return "", &InvalidTypeError{Type: i.Type}
	}
	b, err := i.payloadAtLeast(1)
	return string(b), err
}
