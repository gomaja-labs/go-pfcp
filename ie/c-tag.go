// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"fmt"
	"io"
)

// NewCTAG creates a new CTAG IE.
func NewCTAG(flags, pcp, deiFlag uint8, cvid uint16) *IE {
	fields := NewCTAGFields(flags, pcp, deiFlag, cvid)

	b, err := fields.Marshal()
	if err != nil {
		return newInvalidIE(CTAG, fmt.Errorf("%w: %w", ErrInvalidValue, err))
	}

	return New(CTAG, b)
}

// CTAG returns CTAG in structured format if the type of IE matches.
func (i *IE) CTAG() (*CTAGFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case CTAG:
		fields, err := ParseCTAGFields(i.Payload)
		if err != nil {
			return nil, err
		}

		return fields, nil

	case MACAddressesDetected:
		x, err := i.MACAddressesDetected()
		if err != nil {
			return nil, err
		}
		return ParseCTAGFields(x.CTAG)
	case MACAddressesRemoved:
		x, err := i.MACAddressesRemoved()
		if err != nil {
			return nil, err
		}
		return ParseCTAGFields(x.CTAG)
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// CTAGFields represents a fields contained in CTAG IE.
type CTAGFields struct {
	// Flags contains the presence and option bits.
	Flags uint8
	// PCP is the three-bit priority code point.
	PCP uint8
	// DEIFlag is the one-bit drop eligible indicator.
	DEIFlag uint8
	// CVID is the twelve-bit VLAN identifier.
	CVID uint16
}

// NewCTAGFields creates a new NewCTAGFields.
func NewCTAGFields(flags, pcp, deiFlag uint8, cvid uint16) *CTAGFields {
	return &CTAGFields{
		Flags:   flags,
		PCP:     pcp,
		DEIFlag: deiFlag,
		CVID:    cvid,
	}
}

// HasVID reports whether VID flag is set.
func (f *CTAGFields) HasVID() bool {
	if f == nil {
		return false
	}

	return has3rdBit(f.Flags)
}

// SetVIDFlag sets VID flag in CTAG.
func (f *CTAGFields) SetVIDFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x04
}

// HasDEI reports whether DEI flag is set.
func (f *CTAGFields) HasDEI() bool {
	if f == nil {
		return false
	}

	return has2ndBit(f.Flags)
}

// SetDEIFlag sets DEI flag in CTAG.
func (f *CTAGFields) SetDEIFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x02
}

// HasPCP reports whether PCP flag is set.
func (f *CTAGFields) HasPCP() bool {
	if f == nil {
		return false
	}

	return has1stBit(f.Flags)
}

// SetPCPFlag sets PCP flag in CTAG.
func (f *CTAGFields) SetPCPFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x01
}

// ParseCTAGFields parses b into CTAGFields.
func ParseCTAGFields(b []byte) (*CTAGFields, error) {
	f := &CTAGFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary parses b into IE.
func (f *CTAGFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = CTAGFields{}

	l := len(b)
	if l < 3 {
		return io.ErrUnexpectedEOF
	}

	f.Flags = b[0]
	offset := 1

	if f.HasPCP() {
		f.PCP = b[offset] & 0x07
	}

	if f.HasDEI() {
		f.DEIFlag = (b[offset] >> 3) & 0x01
	}

	if f.HasVID() {
		// TS 29.244 V20.1.0 clause 8.2.94: VID straddles the flag nibble.
		f.CVID = uint16(b[offset]&0xf0)<<4 | uint16(b[offset+1])
	}

	return nil
}

// Marshal returns the serialized bytes of CTAGFields.
func (f *CTAGFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}

	b := make([]byte, f.MarshalLen())
	if err := f.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (f *CTAGFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if f.MarshalLen() > 65535 {
		return ErrInvalidLength
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}

	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	b[0] = f.Flags
	offset := 1

	binary.BigEndian.PutUint16(b[offset:offset+2], ((f.CVID<<4)&0xf000|f.CVID&0xff)|(uint16(f.DEIFlag<<3|f.PCP)<<8))
	return nil
}

// MarshalLen returns field length in integer.
func (f *CTAGFields) MarshalLen() int {
	if f == nil {
		return 0
	}

	return 3
}
