// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewActivatePredefinedRules creates a new ActivatePredefinedRules IE.
func NewActivatePredefinedRules(name string) *IE {
	return newStringIE(ActivatePredefinedRules, name)
}

// ActivatePredefinedRules returns ActivatePredefinedRules in string if the type of IE matches.
func (i *IE) ActivatePredefinedRules() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}

	switch i.Type {
	case ActivatePredefinedRules:
		return i.ValueAsString()
	default:
		return "", &InvalidTypeError{Type: i.Type}
	}
}
