// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewQERID creates a new QERID IE.
func NewQERID(id uint32) *IE {
	return newUint32ValIE(QERID, id)
}

// QERID returns QERID in uint32 if the type of IE matches.
func (i *IE) QERID() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case QERID:
		return i.ValueAsUint32()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
