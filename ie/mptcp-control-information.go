// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewMPTCPControlInformation creates a new MPTCPControlInformation IE.
func NewMPTCPControlInformation(tci uint8) *IE {
	return newUint8ValIE(MPTCPControlInformation, tci&0x01)
}

// MPTCPControlInformation returns MPTCPControlInformation in uint8 if the type of IE matches.
func (i *IE) MPTCPControlInformation() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case MPTCPControlInformation:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// HasTCI reports whether an IE has TCI bit.
func (i *IE) HasTCI() bool {
	if i == nil {
		return false
	}

	v, err := i.MPTCPControlInformation()
	if err != nil {
		return false
	}

	return has1stBit(v)
}
