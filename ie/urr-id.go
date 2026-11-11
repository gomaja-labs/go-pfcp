// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewURRID creates a new URRID IE.
func NewURRID(id uint32) *IE {
	return newUint32ValIE(URRID, id)
}

// URRID returns URRID in uint32 if the type of IE matches.
func (i *IE) URRID() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case URRID:
		return i.ValueAsUint32()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// IsAllocatedByCPFunction reports whether URRID is allocated by CP Function.
func (i *IE) IsAllocatedByCPFunction() bool {
	if i == nil {
		return false
	}

	if i.Type != URRID {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return (i.Payload[0]>>7)&0x01 == 1
}

// IsAllocatedByUPFunction reports whether URRID is allocated by UP Function.
func (i *IE) IsAllocatedByUPFunction() bool {
	if i == nil {
		return false
	}

	if i.Type != URRID {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return (i.Payload[0]>>7)&0x01 != 1
}
