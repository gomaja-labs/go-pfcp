// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewPFCPSMReqFlags creates a new PFCPSMReqFlags IE.
func NewPFCPSMReqFlags(flag uint8) *IE {
	return newUint8ValIE(PFCPSMReqFlags, flag)
}

// PFCPSMReqFlags returns PFCPSMReqFlags in uint8 if the type of IE matches.
func (i *IE) PFCPSMReqFlags() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case PFCPSMReqFlags:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// HasDROBU reports whether an IE has DROBU bit.
func (i *IE) HasDROBU() bool {
	if i == nil {
		return false
	}

	if len(i.Payload) < 1 {
		return false
	}

	switch i.Type {
	case PFCPSMReqFlags:
		v, err := i.PFCPSMReqFlags()
		if err != nil {
			return false
		}

		return has1stBit(v)
	case PFCPSRRspFlags:
		v, err := i.PFCPSRRspFlags()
		if err != nil {
			return false
		}

		return has1stBit(v)
	default:
		return false
	}
}

// HasSNDEM reports whether an IE has SNDEM bit.
func (i *IE) HasSNDEM() bool {
	if i == nil {
		return false
	}

	if len(i.Payload) < 1 {
		return false
	}

	switch i.Type {
	case PFCPSMReqFlags:
		v, err := i.PFCPSMReqFlags()
		if err != nil {
			return false
		}

		return has2ndBit(v)
	default:
		return false
	}
}

// HasQAURR reports whether an IE has QAURR bit.
func (i *IE) HasQAURR() bool {
	if i == nil {
		return false
	}

	if len(i.Payload) < 1 {
		return false
	}

	switch i.Type {
	case PFCPSMReqFlags:
		v, err := i.PFCPSMReqFlags()
		if err != nil {
			return false
		}

		return has3rdBit(v)
	default:
		return false
	}
}

// HasSUMPC reports whether the SUMPC flag is set.
// TS 29.244 V20.1.0 clause 8.2.31, octet 5 bit 4.
func (i *IE) HasSUMPC() bool {
	if payloadFlag(i, PFCPSEReqFlags, 0, 0x02) {
		return true
	}
	return payloadFlag(i, PFCPSMReqFlags, 0, 0x08)
}

// HasRUMUC reports whether the RUMUC flag is set.
// TS 29.244 V20.1.0 clause 8.2.31, octet 5 bit 5.
func (i *IE) HasRUMUC() bool { return payloadFlag(i, PFCPSMReqFlags, 0, 0x10) }

// HasDETEID reports whether the DETEID flag is set.
// TS 29.244 V20.1.0 clause 8.2.31, octet 5 bit 6.
func (i *IE) HasDETEID() bool { return payloadFlag(i, PFCPSMReqFlags, 0, 0x20) }

// HasHRSBOM reports whether the HRSBOM flag is set.
// TS 29.244 V20.1.0 clause 8.2.31, octet 5 bit 7.
func (i *IE) HasHRSBOM() bool {
	// TS 29.244 V20.1.0 clause 8.2.136 also assigns HRSBOM to octet 5 bit 3.
	if payloadFlag(i, PFCPSEReqFlags, 0, 0x04) {
		return true
	}
	return payloadFlag(i, PFCPSMReqFlags, 0, 0x40)
}
