// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewGTPUPathInterfaceType creates a new GTPUPathInterfaceType IE.
func NewGTPUPathInterfaceType(n3, n9 int) *IE {
	return newUint8ValIE(GTPUPathInterfaceType, uint8((n3<<1)|n9))
}

// GTPUPathInterfaceType returns GTPUPathInterfaceType in uint8 if the type of IE matches.
func (i *IE) GTPUPathInterfaceType() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case GTPUPathInterfaceType:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}

}

// HasN3 reports whether an IE has N3 bit.
func (i *IE) HasN3() bool {
	if i == nil {
		return false
	}

	v, err := i.GTPUPathInterfaceType()
	if err != nil {
		return false
	}

	return has2ndBit(v)
}

// HasN9 reports whether an IE has N9 bit.
func (i *IE) HasN9() bool {
	if i == nil {
		return false
	}

	v, err := i.GTPUPathInterfaceType()
	if err != nil {
		return false
	}

	return has1stBit(v)
}
