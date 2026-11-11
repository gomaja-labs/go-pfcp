// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewPFCPSEReqFlags creates a new PFCPSEReqFlags IE.
func NewPFCPSEReqFlags(flag uint8) *IE {
	return newUint8ValIE(PFCPSEReqFlags, flag)
}

// PFCPSEReqFlags returns PFCPSEReqFlags in uint8 if the type of IE matches.
func (i *IE) PFCPSEReqFlags() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	if i.Type != PFCPSEReqFlags {
		return 0, &InvalidTypeError{Type: i.Type}
	}

	return i.ValueAsUint8()
}

// HasRESTI reports whether an IE has RESTI bit.
func (i *IE) HasRESTI() bool {
	if i == nil {
		return false
	}

	v, err := i.PFCPSEReqFlags()
	if err != nil {
		return false
	}

	return has1stBit(v)
}

// HasEPRI reports whether the EPRI flag is set.
// TS 29.244 V20.1.0 clause 8.2.136, octet 5 bit 4.
func (i *IE) HasEPRI() bool { return payloadFlag(i, PFCPSEReqFlags, 0, 0x08) }
