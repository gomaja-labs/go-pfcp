// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewEthertype creates a new Ethertype IE.
func NewEthertype(typ uint16) *IE {
	return newUint16ValIE(Ethertype, typ)
}

// Ethertype returns Ethertype in uint16 if the type of IE matches.
func (i *IE) Ethertype() (uint16, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case Ethertype:
		return i.ValueAsUint16()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}

}
