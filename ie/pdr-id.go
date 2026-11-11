// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewPDRID creates a new PDRID IE.
func NewPDRID(id uint16) *IE {
	return newUint16ValIE(PDRID, id)
}

// PDRID returns PDRID in uint16 if the type of IE matches.
func (i *IE) PDRID() (uint16, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case PDRID:
		return i.ValueAsUint16()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
