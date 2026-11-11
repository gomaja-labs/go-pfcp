// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewEthernetFilterProperties creates a new EthernetFilterProperties IE.
func NewEthernetFilterProperties(props uint8) *IE {
	return newUint8ValIE(EthernetFilterProperties, props)
}

// EthernetFilterProperties returns EthernetFilterProperties in uint8 if the type of IE matches.
func (i *IE) EthernetFilterProperties() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case EthernetFilterProperties:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// HasBIDE reports whether an IE has BIDE bit.
func (i *IE) HasBIDE() bool {
	if i == nil {
		return false
	}

	v, err := i.EthernetFilterProperties()
	if err != nil {
		return false
	}

	return has1stBit(v)
}
