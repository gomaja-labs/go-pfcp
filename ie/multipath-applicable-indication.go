// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// TS 29.244 V20.1.0 clause 8.2.181.

// NewMultipathApplicableIndication creates a new MultipathApplicableIndication IE.
func NewMultipathApplicableIndication(flag uint8) *IE {
	return newUint8ValIE(MultipathApplicableIndication, flag)
}

// MultipathApplicableIndication returns MultipathApplicableIndication in uint8 if the type of IE matches.
func (i *IE) MultipathApplicableIndication() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case MultipathApplicableIndication:
		return i.ValueAsUint8()
	case CreatePDR:
		ies, err := i.CreatePDR()
		if err != nil {
			return 0, err
		}
		for _, x := range ies {
			if x.Type == MultipathApplicableIndication {
				return x.MultipathApplicableIndication()
			}
		}
		return 0, ErrIENotFound
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// HasMAI reports whether an IE has MAI bit.
func (i *IE) HasMAI() bool {
	if i == nil {
		return false
	}

	switch i.Type {
	case MultipathApplicableIndication:
		if len(i.Payload) < 1 {
			return false
		}

		return has1stBit(i.Payload[0])
	default:
		return false
	}
}
