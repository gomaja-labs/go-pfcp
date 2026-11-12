// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"
)

// PMFControlInformationFields contains the PMF flags and the complete per-flow list.
// TS 29.244 V20.1.0 clause 8.2.156.
type PMFControlInformationFields struct {
	// Flags contains the presence and option bits.
	Flags uint8
	// QFIs is the qf is.
	QFIs []uint8
}

// NewPMFControlInformation creates PMF control information with an optional QFI list.
func NewPMFControlInformation(flags uint8, qfis ...uint8) *IE {
	f := &PMFControlInformationFields{Flags: flags, QFIs: qfis}
	b, err := f.Marshal()
	if err != nil {
		return newInvalidIE(PMFControlInformation, ErrInvalidValue)
	}
	return New(PMFControlInformation, b)
}

// PMFControlInformation returns the flags and all QoS flows.
func (i *IE) PMFControlInformation() (*PMFControlInformationFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}
	if i.Type != PMFControlInformation {
		return nil, &InvalidTypeError{Type: i.Type}
	}
	return ParsePMFControlInformationFields(i.Payload)
}

// ParsePMFControlInformationFields decodes PMF control information.
// Decoding may borrow b; keep b unchanged while the result is in use.
func ParsePMFControlInformationFields(b []byte) (*PMFControlInformationFields, error) {
	f := &PMFControlInformationFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary decodes the count-delimited QFI list and ignores extension octets.
// Decoding may borrow b; keep b unchanged while the result is in use.
func (f *PMFControlInformationFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = PMFControlInformationFields{}
	if len(b) < 1 {
		return io.ErrUnexpectedEOF
	}
	f.Flags = b[0] & 7
	if f.Flags&4 != 0 {
		if len(b) < 2 {
			return io.ErrUnexpectedEOF
		}
		n := int(b[1])
		if len(b) < 2+n {
			return io.ErrUnexpectedEOF
		}
		// TS 29.244 V20.1.0 clause 8.2.156: accept a received zero count.
		f.QFIs = make([]uint8, n)
		for k := range f.QFIs {
			f.QFIs[k] = b[2+k] & 0x3f
		}
	}
	return nil
}

// Marshal serializes the PMF control information.
func (f *PMFControlInformationFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	b := make([]byte, f.MarshalLen())
	if err := f.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalLen returns the encoded size.
func (f *PMFControlInformationFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	if f.Flags&4 != 0 {
		return 2 + len(f.QFIs)
	}
	return 1
}

// MarshalTo writes the complete PMF information, rejecting unrepresentable flows.
func (f *PMFControlInformationFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if f.Flags&^uint8(7) != 0 {
		return ErrInvalidValue
	}
	if f.Flags&4 == 0 && len(f.QFIs) != 0 {
		return ErrInvalidValue
	}
	if f.Flags&4 != 0 && (len(f.QFIs) == 0 || len(f.QFIs) > 255) {
		return ErrInvalidLength
	}
	for _, q := range f.QFIs {
		if q > 63 {
			return ErrInvalidValue
		}
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}
	b[0] = f.Flags
	if f.Flags&4 != 0 {
		b[1] = uint8(len(f.QFIs))
		copy(b[2:], f.QFIs)
	}
	return nil
}

// HasPMFI reports whether PMF functionality is required.
func (i *IE) HasPMFI() bool { f, e := i.PMFControlInformation(); return e == nil && f.Flags&1 != 0 }

// HasDRTTI reports whether PMF RTT measurements are disallowed.
func (i *IE) HasDRTTI() bool { f, e := i.PMFControlInformation(); return e == nil && f.Flags&2 != 0 }

// HasPQPM reports whether per-flow performance measurement applies.
func (i *IE) HasPQPM() bool { f, e := i.PMFControlInformation(); return e == nil && f.Flags&4 != 0 }
