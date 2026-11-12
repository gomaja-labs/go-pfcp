// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
)

// NewDownlinkDataServiceInformation creates a new DownlinkDataServiceInformation IE.
func NewDownlinkDataServiceInformation(hasPPI, hasQFI bool, ppi, qfi uint8) *IE {
	f := &DownlinkDataServiceInformationFields{PagingPolicyIndication: ppi, QFI: qfi}
	if hasPPI {
		f.Flags |= 1
	}
	if hasQFI {
		f.Flags |= 2
	}
	return NewDownlinkDataServiceInformationFromFields(f)
}

// DownlinkDataServiceInformation returns all conditional service fields.
func (i *IE) DownlinkDataServiceInformation() (*DownlinkDataServiceInformationFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}
	if i.Type != DownlinkDataServiceInformation {
		return nil, &InvalidTypeError{Type: i.Type}
	}
	return ParseDownlinkDataServiceInformationFields(i.Payload)
}

// HasPPI reports whether an IE has PPI bit.
func (i *IE) HasPPI() bool {
	if i == nil {
		return false
	}

	if i.Type != DownlinkDataServiceInformation {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has1stBit(i.Payload[0])
}

// HasQFI reports whether an IE has QFI bit.
func (i *IE) HasQFI() bool {
	if i == nil {
		return false
	}

	if i.Type != DownlinkDataServiceInformation {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has2ndBit(i.Payload[0])
}

// PPI returns PPI in uint8 if the type of IE matches.
func (i *IE) PPI() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	if i.Type != DownlinkDataServiceInformation {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) == 0 {
		return 0, io.ErrUnexpectedEOF
	}
	// TS 29.244 V20.1.0 clause 8.2.27: QFI is not a missing PPI.
	if !i.HasPPI() {
		return 0, ErrIENotFound
	}
	if len(i.Payload) < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	return i.Payload[1] & 0x07, nil
}

// DownlinkDataServiceInformationFields carries independently optional PPI, QFI and
// packet size. TS 29.244 V20.1.0 clause 8.2.27.
type DownlinkDataServiceInformationFields struct {
	Flags                  uint8
	PagingPolicyIndication uint8
	QFI                    uint8
	DLDataPacketsSize      uint16
}

// NewDownlinkDataServiceInformationFromFields encodes every service field.
func NewDownlinkDataServiceInformationFromFields(f *DownlinkDataServiceInformationFields) *IE {
	b, e := f.Marshal()
	if e != nil {
		return newInvalidIE(DownlinkDataServiceInformation, ErrInvalidValue)
	}
	return New(DownlinkDataServiceInformation, b)
}

// ParseDownlinkDataServiceInformationFields decodes the service fields.
func ParseDownlinkDataServiceInformationFields(b []byte) (*DownlinkDataServiceInformationFields, error) {
	f := &DownlinkDataServiceInformationFields{}
	if e := f.UnmarshalBinary(b); e != nil {
		return nil, e
	}
	return f, nil
}

// UnmarshalBinary decodes values and ignores trailing extension octets.
func (f *DownlinkDataServiceInformationFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = DownlinkDataServiceInformationFields{}
	if len(b) < 1 {
		return io.ErrUnexpectedEOF
	}
	f.Flags = b[0] & 7
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}
	o := 1
	if f.Flags&1 != 0 {
		f.PagingPolicyIndication = b[o] & 7
		o++
	}
	if f.Flags&2 != 0 {
		f.QFI = b[o] & 0x3f
		o++
	}
	if f.Flags&4 != 0 {
		f.DLDataPacketsSize = binary.BigEndian.Uint16(b[o : o+2])
	}
	return nil
}

// Marshal serializes all present service fields.
func (f *DownlinkDataServiceInformationFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	b := make([]byte, f.MarshalLen())
	if e := f.MarshalTo(b); e != nil {
		return nil, e
	}
	return b, nil
}

// MarshalLen returns the encoded length.
func (f *DownlinkDataServiceInformationFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	n := 1
	if f.Flags&1 != 0 {
		n++
	}
	if f.Flags&2 != 0 {
		n++
	}
	if f.Flags&4 != 0 {
		n += 2
	}
	return n
}

// MarshalTo writes every present field in specification order.
func (f *DownlinkDataServiceInformationFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if f.Flags&1 == 0 && f.PagingPolicyIndication != 0 || f.Flags&2 == 0 && f.QFI != 0 || f.Flags&4 == 0 && f.DLDataPacketsSize != 0 {
		return ErrInvalidValue
	}
	if f.Flags&^uint8(7) != 0 || f.PagingPolicyIndication > 7 || f.QFI > 63 {
		return ErrInvalidValue
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}
	b[0] = f.Flags
	o := 1
	if f.Flags&1 != 0 {
		b[o] = f.PagingPolicyIndication
		o++
	}
	if f.Flags&2 != 0 {
		b[o] = f.QFI
		o++
	}
	if f.Flags&4 != 0 {
		binary.BigEndian.PutUint16(b[o:o+2], f.DLDataPacketsSize)
	}
	return nil
}

// HasDLPSI reports whether the nested packet size is present.
func (i *IE) HasDLPSI() bool {
	return i != nil && i.Type == DownlinkDataServiceInformation && len(i.Payload) > 0 && i.Payload[0]&4 != 0
}
