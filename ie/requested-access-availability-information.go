// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewRequestedAccessAvailabilityInformation creates a new RequestedAccessAvailabilityInformation IE.
func NewRequestedAccessAvailabilityInformation(rrca uint8) *IE {
	return newUint8ValIE(RequestedAccessAvailabilityInformation, rrca&0x01)
}

// RequestedAccessAvailabilityInformation returns RequestedAccessAvailabilityInformation in uint8 if the type of IE matches.
func (i *IE) RequestedAccessAvailabilityInformation() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case RequestedAccessAvailabilityInformation:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// HasRRCA reports whether an IE has RRCA bit.
func (i *IE) HasRRCA() bool {
	if i == nil {
		return false
	}

	v, err := i.RequestedAccessAvailabilityInformation()
	if err != nil {
		return false
	}

	return has1stBit(v)
}
