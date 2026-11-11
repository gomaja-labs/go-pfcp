// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewApplicationID creates a new ApplicationID IE.
func NewApplicationID(instance string) *IE {
	return newStringIE(ApplicationID, instance)
}

// ApplicationID returns ApplicationID in string if the type of IE matches.
func (i *IE) ApplicationID() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}

	switch i.Type {
	case ApplicationID:
		return i.ValueAsString()
	default:
		return "", &InvalidTypeError{Type: i.Type}
	}
}
