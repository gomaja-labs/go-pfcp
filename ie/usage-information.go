// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewUsageInformation creates a new UsageInformation IE.
func NewUsageInformation(bef, aft, uae, ube int) *IE {
	return newUint8ValIE(UsageInformation, uint8((ube<<3)|(uae<<2)|(aft<<1)|(bef)))
}

// UsageInformation returns UsageInformation in uint8 if the type of IE matches.
func (i *IE) UsageInformation() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case UsageInformation:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// HasUBE reports whether and IE has UBE bit.
func (i *IE) HasUBE() bool {
	if i == nil {
		return false
	}

	v, err := i.UsageInformation()
	if err != nil {
		return false
	}

	return has4thBit(v)
}

// HasUAE reports whether and IE has UAE bit.
func (i *IE) HasUAE() bool {
	if i == nil {
		return false
	}

	v, err := i.UsageInformation()
	if err != nil {
		return false
	}

	return has3rdBit(v)
}

// HasAFT reports whether and IE has AFT bit.
func (i *IE) HasAFT() bool {
	if i == nil {
		return false
	}

	v, err := i.UsageInformation()
	if err != nil {
		return false
	}

	return has2ndBit(v)
}

// HasBEF reports whether and IE has BEF bit.
func (i *IE) HasBEF() bool {
	if i == nil {
		return false
	}

	v, err := i.UsageInformation()
	if err != nil {
		return false
	}

	return has1stBit(v)
}
