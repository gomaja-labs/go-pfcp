// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"
	"net"
)

// NewUEIPAddress creates a new UEIPAddress IE.
func NewUEIPAddress(flags uint8, v4, v6 string, v6d, v6pl uint8) *IE {
	fields := NewUEIPAddressFields(flags, v4, v6, v6d, v6pl)
	b, err := fields.Marshal()
	if err != nil {
		return nil
	}

	return New(UEIPAddress, b)
}

// UEIPAddress returns UEIPAddress in *UEIPAddressFields if the type of IE matches.
func (i *IE) UEIPAddress() (*UEIPAddressFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case UEIPAddress:
		fields, err := ParseUEIPAddressFields(i.Payload)
		if err != nil {
			return nil, err
		}

		return fields, nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// HasCHV4 reports whether an IE has CHV4 bit.
// This bit has been introduced in release 16.2
// with the removal of the bit CH at the same position.
func (i *IE) HasCHV4() bool {
	if i == nil {
		return false
	}

	if i.Type != UEIPAddress {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has5thBit(i.Payload[0])
}

// HasCHV6 reports whether an IE has CHV6 bit.
// This bit has been introduced in release 16.2
func (i *IE) HasCHV6() bool {
	if i == nil {
		return false
	}

	if i.Type != UEIPAddress {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has6thBit(i.Payload[0])
}

// HasIP6PL reports whether an IE has IP6PL bit.
// This bit has been introduced in release 16.2
func (i *IE) HasIP6PL() bool {
	// TS 29.244 V20.1.0 clause 8.2.25, octet 8 bit 6.
	if i != nil && i.Type == UPFunctionFeatures {
		return payloadFlag(i, UPFunctionFeatures, 3, 0x20)
	}

	if i == nil {
		return false
	}

	if i.Type != UEIPAddress {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has7thBit(i.Payload[0])
}

// HasIPv6D reports whether an IE has IPv6D bit.
func (i *IE) HasIPv6D() bool {
	if i == nil {
		return false
	}

	if i.Type != UEIPAddress {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has4thBit(i.Payload[0])
}

// HasSD reports whether an IE has SD bit.
func (i *IE) HasSD() bool {
	if i == nil {
		return false
	}

	if i.Type != UEIPAddress {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has3rdBit(i.Payload[0])
}

// UEIPAddressFields represents a fields contained in UEIPAddress IE.
type UEIPAddressFields struct {
	Flags                    uint8
	IPv4Address              net.IP
	IPv6Address              net.IP
	IPv6PrefixDelegationBits uint8
	IPv6PrefixLength         uint8
}

// NewUEIPAddressFields creates a new UEIPAddressFields.
func NewUEIPAddressFields(flags uint8, v4, v6 string, v6d, v6pl uint8) *UEIPAddressFields {
	f := &UEIPAddressFields{Flags: flags}

	if has2ndBit(flags) && !has5thBit(flags) {
		f.IPv4Address = net.ParseIP(v4).To4()
	}

	if has1stBit(flags) && !has6thBit(flags) {
		f.IPv6Address = net.ParseIP(v6).To16()
	}

	if has4thBit(flags) {
		f.IPv6PrefixDelegationBits = v6d
	}

	if has7thBit(flags) {
		f.IPv6PrefixLength = v6pl
	}

	return f
}

// ParseUEIPAddressFields parses b into UEIPAddressFields.
func ParseUEIPAddressFields(b []byte) (*UEIPAddressFields, error) {
	f := &UEIPAddressFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary parses b into IE.
func (f *UEIPAddressFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = UEIPAddressFields{}

	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	f.Flags = b[0]
	offset := 1

	if has2ndBit(f.Flags) && !has5thBit(f.Flags) {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		f.IPv4Address = net.IP(b[offset : offset+4]).To4()
		offset += 4
	}

	if has1stBit(f.Flags) && !has6thBit(f.Flags) {
		if l < offset+16 {
			return io.ErrUnexpectedEOF
		}
		f.IPv6Address = net.IP(b[offset : offset+16]).To16()
		offset += 16
	}

	if has4thBit(f.Flags) {
		if l < offset+1 {
			return io.ErrUnexpectedEOF
		}
		f.IPv6PrefixDelegationBits = b[offset]
	}

	if has7thBit(f.Flags) {
		if l < offset+1 {
			return io.ErrUnexpectedEOF
		}
		f.IPv6PrefixLength = b[offset]
	}

	return nil
}

// Marshal returns the serialized bytes of UEIPAddressFields.
func (f *UEIPAddressFields) Marshal() ([]byte, error) {
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
func (f *UEIPAddressFields) MarshalTo(b []byte) error {
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

	if has2ndBit(f.Flags) && !has5thBit(f.Flags) {
		if f.IPv4Address.To4() == nil {
			return ErrMalformed
		}
		copy(b[offset:offset+4], f.IPv4Address.To4())
		offset += 4
	}

	if has1stBit(f.Flags) && !has6thBit(f.Flags) {
		if f.IPv6Address.To16() == nil {
			return ErrMalformed
		}
		copy(b[offset:offset+16], f.IPv6Address.To16())
		offset += 16
	}

	if has4thBit(f.Flags) {
		b[offset] = f.IPv6PrefixDelegationBits
		offset++
	}

	if has7thBit(f.Flags) {
		b[offset] = f.IPv6PrefixLength
	}

	return nil
}

// MarshalLen returns field length in integer.
func (f *UEIPAddressFields) MarshalLen() int {
	if f == nil {
		return 0
	}

	l := 1
	if has2ndBit(f.Flags) && !has5thBit(f.Flags) {
		l += 4
	}
	if has1stBit(f.Flags) && !has6thBit(f.Flags) {
		l += 16
	}
	if has4thBit(f.Flags) {
		l++
	}
	if has7thBit(f.Flags) {
		l++
	}

	return l
}
