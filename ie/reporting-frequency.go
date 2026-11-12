// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewReportingFrequency creates a new ReportingFrequency IE.
func NewReportingFrequency(perio, evett bool) *IE {
	// TS 29.244 V20.1.0 clause 8.2.168: the former SESRL bit is reserved.
	var flags uint8
	if perio {
		flags |= 2
	}
	if evett {
		flags |= 1
	}
	if flags == 0 {
		return newInvalidIE(ReportingFrequency, ErrInvalidValue)
	}
	return newUint8ValIE(ReportingFrequency, flags)
}

// ReportingFrequency returns ReportingFrequency in uint8 if the type of IE matches.
func (i *IE) ReportingFrequency() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case ReportingFrequency:
		// TS 29.244 V20.1.0 clause 8.2.168: ignore received spare/reserved bits.
		v, err := i.ValueAsUint8()
		return v & 3, err
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}

}

// HasEVETT reports whether an IE has EVETT bit.
func (i *IE) HasEVETT() bool {
	if i == nil {
		return false
	}

	v, err := i.ReportingFrequency()
	if err != nil {
		return false
	}

	return has1stBit(v)
}
