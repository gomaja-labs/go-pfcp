// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewDeactivatePredefinedRules creates a new DeactivatePredefinedRules IE.
// The value is an octet string carried unchanged; no character set is imposed or checked.
// TS 29.244 V20.1.0 clause 8.2.73: no text repertoire is specified for this identifier.
func NewDeactivatePredefinedRules(name string) *IE {
	return newStringIE(DeactivatePredefinedRules, name)
}

// DeactivatePredefinedRules returns the unchanged octet string.
func (i *IE) DeactivatePredefinedRules() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}
	if i.Type != DeactivatePredefinedRules {
		return "", &InvalidTypeError{Type: i.Type}
	}
	b, err := i.payloadAtLeast(1)
	return string(b), err
}
