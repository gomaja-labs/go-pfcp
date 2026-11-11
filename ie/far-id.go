// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewFARID creates a new FARID IE.
func NewFARID(id uint32) *IE {
	return newUint32ValIE(FARID, id)
}

// FARID returns FARID in uint32 if the type of IE matches.
func (i *IE) FARID() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case FARID:
		return i.ValueAsUint32()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
