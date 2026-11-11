// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewPrecedence creates a new Precedence IE.
func NewPrecedence(id uint32) *IE {
	return newUint32ValIE(Precedence, id)
}

// Precedence returns Precedence in uint32 if the type of IE matches.
func (i *IE) Precedence() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case Precedence:
		return i.ValueAsUint32()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
