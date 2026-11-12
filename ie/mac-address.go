// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"fmt"
	"io"
	"net"
)

// NewMACAddress creates a new MACAddress IE.
func NewMACAddress(src, dst, upperSrc, upperDst net.HardwareAddr) *IE {
	fields := NewMACAddressFields(src, dst, upperSrc, upperDst)

	b, err := fields.Marshal()
	if err != nil {
		return newInvalidIE(MACAddress, fmt.Errorf("%w: %w", ErrInvalidValue, err))
	}

	return New(MACAddress, b)
}

// MACAddress returns MACAddress in structured format if the type of IE matches.
func (i *IE) MACAddress() (*MACAddressFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case MACAddress:
		fields, err := ParseMACAddressFields(i.Payload)
		if err != nil {
			return nil, err
		}

		return fields, nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// MACAddressFields represents a fields contained in MACAddress IE.
type MACAddressFields struct {
	// Flags contains the presence and option bits.
	Flags uint8
	// SourceMACAddress is the source mac address.
	SourceMACAddress net.HardwareAddr
	// DestinationMACAddress is the destination mac address.
	DestinationMACAddress net.HardwareAddr
	// UpperSourceMACAddress is the upper source mac address.
	UpperSourceMACAddress net.HardwareAddr
	// UpperDestinationMACAddress is the upper destination mac address.
	UpperDestinationMACAddress net.HardwareAddr
}

// NewMACAddressFields creates a new NewMACAddressFields.
func NewMACAddressFields(src, dst, upperSrc, upperDst net.HardwareAddr) *MACAddressFields {
	f := &MACAddressFields{}

	if src != nil {
		f.SetSOURFlag()
		f.SourceMACAddress = src
	}

	if dst != nil {
		f.SetDESTFlag()
		f.DestinationMACAddress = dst
	}

	if upperSrc != nil {
		f.SetUSOUFlag()
		f.UpperSourceMACAddress = upperSrc
	}

	if upperDst != nil {
		f.SetUDESFlag()
		f.UpperDestinationMACAddress = upperDst
	}

	return f
}

// HasUDES reports whether UDES flag is set.
func (f *MACAddressFields) HasUDES() bool {
	if f == nil {
		return false
	}

	return has4thBit(f.Flags)
}

// SetUDESFlag sets UDES flag in MACAddress.
func (f *MACAddressFields) SetUDESFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x08
}

// HasUSOU reports whether USOU flag is set.
func (f *MACAddressFields) HasUSOU() bool {
	if f == nil {
		return false
	}

	return has3rdBit(f.Flags)
}

// SetUSOUFlag sets USOU flag in MACAddress.
func (f *MACAddressFields) SetUSOUFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x04
}

// HasDEST reports whether DEST flag is set.
func (f *MACAddressFields) HasDEST() bool {
	if f == nil {
		return false
	}

	return has2ndBit(f.Flags)
}

// SetDESTFlag sets DEST flag in MACAddress.
func (f *MACAddressFields) SetDESTFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x02
}

// HasSOUR reports whether SOUR flag is set.
func (f *MACAddressFields) HasSOUR() bool {
	if f == nil {
		return false
	}

	return has1stBit(f.Flags)
}

// SetSOURFlag sets SOUR flag in MACAddress.
func (f *MACAddressFields) SetSOURFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x01
}

// ParseMACAddressFields parses b into MACAddressFields.
// Decoding may borrow b; keep b unchanged while the result is in use.
func ParseMACAddressFields(b []byte) (*MACAddressFields, error) {
	f := &MACAddressFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary parses b into IE.
// Decoding may borrow b; keep b unchanged while the result is in use.
func (f *MACAddressFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = MACAddressFields{}

	l := len(b)
	if l < 2 {
		return io.ErrUnexpectedEOF
	}

	f.Flags = b[0]
	offset := 1

	if f.HasSOUR() {
		if l < offset+6 {
			return io.ErrUnexpectedEOF
		}
		f.SourceMACAddress = make(net.HardwareAddr, 6)
		copy(f.SourceMACAddress, b[offset:offset+6])
		offset += 6
	}

	if f.HasDEST() {
		if l < offset+6 {
			return io.ErrUnexpectedEOF
		}
		f.DestinationMACAddress = make(net.HardwareAddr, 6)
		copy(f.DestinationMACAddress, b[offset:offset+6])
		offset += 6
	}

	if f.HasUSOU() {
		if l < offset+6 {
			return io.ErrUnexpectedEOF
		}
		f.UpperSourceMACAddress = make(net.HardwareAddr, 6)
		copy(f.UpperSourceMACAddress, b[offset:offset+6])
		offset += 6
	}

	if f.HasUDES() {
		if l < offset+6 {
			return io.ErrUnexpectedEOF
		}
		f.UpperDestinationMACAddress = make(net.HardwareAddr, 6)
		copy(f.UpperDestinationMACAddress, b[offset:offset+6])
	}

	return nil
}

// Marshal returns the serialized bytes of MACAddressFields.
func (f *MACAddressFields) Marshal() ([]byte, error) {
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
func (f *MACAddressFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	// TS 29.244 V20.1.0 clause 8.2.93: each present MAC address occupies six octets.
	for _, addr := range []net.HardwareAddr{f.SourceMACAddress, f.DestinationMACAddress, f.UpperSourceMACAddress, f.UpperDestinationMACAddress} {
		if addr != nil && len(addr) != 6 {
			return ErrInvalidValue
		}
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

	if f.HasSOUR() {
		if l < offset+6 {
			return io.ErrUnexpectedEOF
		}
		copy(b[offset:offset+6], f.SourceMACAddress[:6])
		offset += 6
	}

	if f.HasDEST() {
		if l < offset+6 {
			return io.ErrUnexpectedEOF
		}
		copy(b[offset:offset+6], f.DestinationMACAddress[:6])
		offset += 6
	}

	if f.HasUSOU() {
		if l < offset+6 {
			return io.ErrUnexpectedEOF
		}
		copy(b[offset:offset+6], f.UpperSourceMACAddress[:6])
		offset += 6
	}

	if f.HasUDES() {
		if l < offset+6 {
			return io.ErrUnexpectedEOF
		}
		copy(b[offset:offset+6], f.UpperDestinationMACAddress[:6])
	}

	return nil
}

// MarshalLen returns field length in integer.
func (f *MACAddressFields) MarshalLen() int {
	if f == nil {
		return 0
	}

	l := 1
	if f.SourceMACAddress != nil {
		l += 6
	}
	if f.DestinationMACAddress != nil {
		l += 6
	}
	if f.UpperSourceMACAddress != nil {
		l += 6
	}
	if f.UpperDestinationMACAddress != nil {
		l += 6
	}
	return l
}
