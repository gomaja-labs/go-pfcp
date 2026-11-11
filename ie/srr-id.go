// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewSRRID creates a new SRRID IE.
func NewSRRID(id uint8) *IE {
	return newUint8ValIE(SRRID, id)
}

// SRRID returns SRRID in uint8 if the type of IE matches.
func (i *IE) SRRID() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case SRRID:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
