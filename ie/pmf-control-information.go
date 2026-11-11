// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewPMFControlInformation creates a new PMFControlInformation IE.
func NewPMFControlInformation(pmfi uint8) *IE {
	return newUint8ValIE(PMFControlInformation, pmfi&0x01)
}

// PMFControlInformation returns PMFControlInformation in uint8 if the type of IE matches.
func (i *IE) PMFControlInformation() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case PMFControlInformation:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// HasPMFI reports whether an IE has PMFI bit.
func (i *IE) HasPMFI() bool {
	if i == nil {
		return false
	}

	v, err := i.PMFControlInformation()
	if err != nil {
		return false
	}

	return has1stBit(v)
}
