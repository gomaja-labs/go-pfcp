// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewMTEDTControlInformation creates a new MTEDTControlInformation IE.
func NewMTEDTControlInformation(rdsi uint8) *IE {
	return newUint8ValIE(MTEDTControlInformation, rdsi&0x01)
}

// MTEDTControlInformation returns MTEDTControlInformation in uint8 if the type of IE matches.
func (i *IE) MTEDTControlInformation() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case MTEDTControlInformation:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// HasRDSI reports whether an IE has RDSI bit.
func (i *IE) HasRDSI() bool {
	if i == nil {
		return false
	}

	v, err := i.MTEDTControlInformation()
	if err != nil {
		return false
	}

	return has1stBit(v)
}
