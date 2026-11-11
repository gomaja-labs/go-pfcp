// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewEthernetPDUSessionInformation creates a new EthernetPDUSessionInformation IE.
func NewEthernetPDUSessionInformation(info uint8) *IE {
	return newUint8ValIE(EthernetPDUSessionInformation, info)
}

// EthernetPDUSessionInformation returns EthernetPDUSessionInformation in uint8 if the type of IE matches.
func (i *IE) EthernetPDUSessionInformation() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case EthernetPDUSessionInformation:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// HasETHI reports whether an IE has ETHI bit.
func (i *IE) HasETHI() bool {
	if i == nil {
		return false
	}

	v, err := i.EthernetPDUSessionInformation()
	if err != nil {
		return false
	}

	return has1stBit(v)
}
