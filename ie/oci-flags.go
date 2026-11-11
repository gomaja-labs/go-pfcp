// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewOCIFlags creates a new OCIFlags IE.
func NewOCIFlags(flags uint8) *IE {
	return newUint8ValIE(OCIFlags, flags)
}

// OCIFlags returns OCIFlags in uint8 if the type of IE matches.
func (i *IE) OCIFlags() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case OCIFlags:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// HasAOCI reports whether the AOCI flag is set.
// TS 29.244 V20.1.0 clause 8.2.76, octet 5 bit 1.
func (i *IE) HasAOCI() bool { return payloadFlag(i, OCIFlags, 0, 0x01) }
