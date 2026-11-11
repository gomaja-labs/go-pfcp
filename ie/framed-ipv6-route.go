// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewFramedIPv6Route creates a new FramedIPv6Route IE.
func NewFramedIPv6Route(name string) *IE {
	return newStringIE(FramedIPv6Route, name)
}

// FramedIPv6Route returns FramedIPv6Route in string if the type of IE matches.
func (i *IE) FramedIPv6Route() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}

	switch i.Type {
	case FramedIPv6Route:
		return i.ValueAsString()
	default:
		return "", &InvalidTypeError{Type: i.Type}
	}
}
