// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"fmt"
	"io"
	"net"
)

// NewIPMulticastAddress creates a new IPMulticastAddress IE.
func NewIPMulticastAddress(sv4, sv6, ev4, ev6 net.IP) *IE {
	fields := NewIPMulticastAddressFields(sv4, sv6, ev4, ev6)

	b, err := fields.Marshal()
	if err != nil {
		return newInvalidIE(IPMulticastAddress, fmt.Errorf("%w: %w", ErrInvalidValue, err))
	}

	return New(IPMulticastAddress, b)
}

// IPMulticastAddress returns IPMulticastAddress in structured format if the type of IE matches.
func (i *IE) IPMulticastAddress() (*IPMulticastAddressFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case IPMulticastAddress:
		fields, err := ParseIPMulticastAddressFields(i.Payload)
		if err != nil {
			return nil, err
		}

		return fields, nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// IPMulticastAddressFields represents a fields contained in IPMulticastAddress IE.
type IPMulticastAddressFields struct {
	// Flags contains the presence and option bits.
	Flags uint8
	// StartIPv4Address is the start IPv4 address.
	StartIPv4Address net.IP
	// StartIPv6Address is the start IPv6 address.
	StartIPv6Address net.IP
	// EndIPv4Address is the end IPv4 address.
	EndIPv4Address net.IP
	// EndIPv6Address is the end IPv6 address.
	EndIPv6Address net.IP
}

// NewIPMulticastAddressFields creates a new NewIPMulticastAddressFields.
func NewIPMulticastAddressFields(sv4, sv6, ev4, ev6 net.IP) *IPMulticastAddressFields {
	f := &IPMulticastAddressFields{Flags: 0x00}
	if sv4 == nil && sv6 == nil {
		f.SetAnyFlag()
		return f
	}

	if sv4 != nil {
		f.SetIPv4Flag()
		f.StartIPv4Address = sv4
	}
	if sv6 != nil {
		f.SetIPv6Flag()
		f.StartIPv6Address = sv6
	}

	if ev4 != nil {
		f.SetRangeFlag()
		f.EndIPv4Address = ev4
	}

	if ev6 != nil {
		f.SetRangeFlag()
		f.EndIPv6Address = ev6
	}

	return f
}

// HasAny reports whether Any flag is set.
func (f *IPMulticastAddressFields) HasAny() bool {
	if f == nil {
		return false
	}

	return has4thBit(f.Flags)
}

// SetAnyFlag sets Any flag in IPMulticastAddress.
func (f *IPMulticastAddressFields) SetAnyFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x08
}

// HasRange reports whether Range flag is set.
func (f *IPMulticastAddressFields) HasRange() bool {
	if f == nil {
		return false
	}

	return has3rdBit(f.Flags)
}

// SetRangeFlag sets Range flag in IPMulticastAddress.
func (f *IPMulticastAddressFields) SetRangeFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x04
}

// HasIPv4 reports whether IPv4 flag is set.
func (f *IPMulticastAddressFields) HasIPv4() bool {
	if f == nil {
		return false
	}

	return has2ndBit(f.Flags)
}

// SetIPv4Flag sets IPv4 flag in IPMulticastAddress.
func (f *IPMulticastAddressFields) SetIPv4Flag() {
	if f == nil {
		return
	}

	f.Flags |= 0x02
}

// HasIPv6 reports whether IPv6 flag is set.
func (f *IPMulticastAddressFields) HasIPv6() bool {
	if f == nil {
		return false
	}

	return has1stBit(f.Flags)
}

// SetIPv6Flag sets IPv6 flag in IPMulticastAddress.
func (f *IPMulticastAddressFields) SetIPv6Flag() {
	if f == nil {
		return
	}

	f.Flags |= 0x01
}

// ParseIPMulticastAddressFields parses b into IPMulticastAddressFields.
// Decoding may borrow b; keep b unchanged while the result is in use.
func ParseIPMulticastAddressFields(b []byte) (*IPMulticastAddressFields, error) {
	f := &IPMulticastAddressFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary parses b into IE.
// Decoding may borrow b; keep b unchanged while the result is in use.
func (f *IPMulticastAddressFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = IPMulticastAddressFields{}

	// TS 29.244 V20.1.0 clause 8.2.137: Any has no address fields.
	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	f.Flags = b[0]
	offset := 1

	if f.HasAny() {
		return nil
	}

	if f.HasIPv4() {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		f.StartIPv4Address = net.IP(b[offset : offset+4])
		offset += 4
	}

	if f.HasIPv6() {
		if l < offset+16 {
			return io.ErrUnexpectedEOF
		}
		f.StartIPv6Address = net.IP(b[offset : offset+16])
		offset += 16
	}

	if f.HasRange() {

		if f.HasIPv4() {
			if l < offset+4 {
				return io.ErrUnexpectedEOF
			}
			f.EndIPv4Address = net.IP(b[offset : offset+4])
			offset += 4
		}

		if f.HasIPv6() {
			if l < offset+16 {
				return io.ErrUnexpectedEOF
			}
			f.EndIPv6Address = net.IP(b[offset : offset+16])
		}
	}

	return nil
}

// Marshal returns the serialized bytes of IPMulticastAddressFields.
func (f *IPMulticastAddressFields) Marshal() ([]byte, error) {
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
func (f *IPMulticastAddressFields) MarshalTo(b []byte) error {
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

	if f.HasAny() {
		return nil
	}

	if f.HasIPv4() {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		copy(b[offset:offset+4], f.StartIPv4Address.To4())
		offset += 4
	}

	if f.HasIPv6() {
		if l < offset+16 {
			return io.ErrUnexpectedEOF
		}
		copy(b[offset:offset+16], f.StartIPv6Address.To16())
		offset += 16
	}

	if f.HasRange() {
		if f.HasIPv4() {
			if l < offset+4 {
				return io.ErrUnexpectedEOF
			}
			copy(b[offset:offset+4], f.EndIPv4Address.To4())
			offset += 4
		}

		if f.HasIPv6() {
			if l < offset+16 {
				return io.ErrUnexpectedEOF
			}
			copy(b[offset:offset+16], f.EndIPv6Address.To16())
		}
	}

	return nil
}

// MarshalLen returns field length in integer.
func (f *IPMulticastAddressFields) MarshalLen() int {
	if f == nil {
		return 0
	}

	l := 1
	if f.StartIPv4Address != nil {
		l += 4
	}
	if f.StartIPv6Address != nil {
		l += 16
	}
	if f.EndIPv4Address != nil {
		l += 4
	}
	if f.EndIPv6Address != nil {
		l += 16
	}

	return l
}
