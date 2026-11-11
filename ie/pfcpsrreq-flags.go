// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewPFCPSRReqFlags creates a new PFCPSRReqFlags IE.
func NewPFCPSRReqFlags(flag uint8) *IE {
	return newUint8ValIE(PFCPSRReqFlags, flag)
}

// PFCPSRReqFlags returns PFCPSRReqFlags in uint8 if the type of IE matches.
func (i *IE) PFCPSRReqFlags() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	if i.Type != PFCPSRReqFlags {
		return 0, &InvalidTypeError{Type: i.Type}
	}

	return i.ValueAsUint8()
}

// HasPSDBU reports whether an IE has PSDBU bit.
func (i *IE) HasPSDBU() bool {
	if i == nil {
		return false
	}

	v, err := i.PFCPSRReqFlags()
	if err != nil {
		return false
	}

	return has1stBit(v)
}

// HasPPPCP reports whether the PPPCP flag is set.
// TS 29.244 V20.1.0 clause 8.2.119, octet 5 bit 2.
func (i *IE) HasPPPCP() bool { return payloadFlag(i, PFCPSRReqFlags, 0, 0x02) }
