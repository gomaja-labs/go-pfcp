// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// TS 29.244 V20.1.0 clause 8.2.158.

import (
	"fmt"
	"io"
	"net"
)

// NewLinkSpecificMultipathIPAddress creates a new LinkSpecificMultipathIPAddress IE.
func NewLinkSpecificMultipathIPAddress(v4, v6, nv4, nv6 net.IP) *IE {
	fields := NewLinkSpecificMultipathIPAddressFields(v4, v6, nv4, nv6)

	b, err := fields.Marshal()
	if err != nil {
		return newInvalidIE(LinkSpecificMultipathIPAddress, fmt.Errorf("%w: %w", ErrInvalidValue, err))
	}

	return New(LinkSpecificMultipathIPAddress, b)
}

// LinkSpecificMultipathIPAddress returns LinkSpecificMultipathIPAddress in structured format if the type of IE matches.
func (i *IE) LinkSpecificMultipathIPAddress() (*LinkSpecificMultipathIPAddressFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case LinkSpecificMultipathIPAddress:
		fields, err := ParseLinkSpecificMultipathIPAddressFields(i.Payload)
		if err != nil {
			return nil, err
		}

		return fields, nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// LinkSpecificMultipathIPAddressFields represents a fields contained in LinkSpecificMultipathIPAddress IE.
type LinkSpecificMultipathIPAddressFields struct {
	// Flags contains the presence and option bits.
	Flags uint8
	// UELinkSpecificIPv4AddressFor3GPPAccess is the ue link specific IPv4 address for3 gpp access.
	UELinkSpecificIPv4AddressFor3GPPAccess net.IP
	// UELinkSpecificIPv6AddressFor3GPPAccess is the ue link specific IPv6 address for3 gpp access.
	UELinkSpecificIPv6AddressFor3GPPAccess net.IP
	// UELinkSpecificIPv4AddressForNon3GPPAccess is the ue link specific IPv4 address for non3 gpp access.
	UELinkSpecificIPv4AddressForNon3GPPAccess net.IP
	// UELinkSpecificIPv6AddressForNon3GPPAccess is the ue link specific IPv6 address for non3 gpp access.
	UELinkSpecificIPv6AddressForNon3GPPAccess net.IP
}

// NewLinkSpecificMultipathIPAddressFields creates a new NewLinkSpecificMultipathIPAddressFields.
func NewLinkSpecificMultipathIPAddressFields(v4, v6, nv4, nv6 net.IP) *LinkSpecificMultipathIPAddressFields {
	f := &LinkSpecificMultipathIPAddressFields{Flags: 0x00}

	if v4 != nil {
		f.SetIPv4Flag()
		f.UELinkSpecificIPv4AddressFor3GPPAccess = v4
	}
	if v6 != nil {
		f.SetIPv6Flag()
		f.UELinkSpecificIPv6AddressFor3GPPAccess = v6
	}

	if nv4 != nil {
		f.SetNV4Flag()
		f.UELinkSpecificIPv4AddressForNon3GPPAccess = nv4
	}
	if nv6 != nil {
		f.SetNV6Flag()
		f.UELinkSpecificIPv6AddressForNon3GPPAccess = nv6
	}

	return f
}

// HasNV6 reports whether NV6 flag is set.
func (f *LinkSpecificMultipathIPAddressFields) HasNV6() bool {
	if f == nil {
		return false
	}

	return has4thBit(f.Flags)
}

// SetNV6Flag sets NV6 flag in LinkSpecificMultipathIPAddress.
func (f *LinkSpecificMultipathIPAddressFields) SetNV6Flag() {
	if f == nil {
		return
	}

	f.Flags |= 0x08
}

// HasNV4 reports whether NV4 flag is set.
func (f *LinkSpecificMultipathIPAddressFields) HasNV4() bool {
	if f == nil {
		return false
	}

	return has3rdBit(f.Flags)
}

// SetNV4Flag sets NV4 flag in LinkSpecificMultipathIPAddress.
func (f *LinkSpecificMultipathIPAddressFields) SetNV4Flag() {
	if f == nil {
		return
	}

	f.Flags |= 0x04
}

// HasIPv6 reports whether IPv6 flag is set.
func (f *LinkSpecificMultipathIPAddressFields) HasIPv6() bool {
	if f == nil {
		return false
	}

	return has2ndBit(f.Flags)
}

// SetIPv6Flag sets IPv6 flag in LinkSpecificMultipathIPAddress.
func (f *LinkSpecificMultipathIPAddressFields) SetIPv6Flag() {
	if f == nil {
		return
	}

	f.Flags |= 0x02
}

// HasIPv4 reports whether IPv4 flag is set.
func (f *LinkSpecificMultipathIPAddressFields) HasIPv4() bool {
	if f == nil {
		return false
	}

	return has1stBit(f.Flags)
}

// SetIPv4Flag sets IPv4 flag in LinkSpecificMultipathIPAddress.
func (f *LinkSpecificMultipathIPAddressFields) SetIPv4Flag() {
	if f == nil {
		return
	}

	f.Flags |= 0x01
}

// ParseLinkSpecificMultipathIPAddressFields parses b into LinkSpecificMultipathIPAddressFields.
func ParseLinkSpecificMultipathIPAddressFields(b []byte) (*LinkSpecificMultipathIPAddressFields, error) {
	f := &LinkSpecificMultipathIPAddressFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary parses b into IE.
func (f *LinkSpecificMultipathIPAddressFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = LinkSpecificMultipathIPAddressFields{}

	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	f.Flags = b[0]
	offset := 1

	if f.HasIPv4() {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		f.UELinkSpecificIPv4AddressFor3GPPAccess = net.IP(b[offset : offset+4])
		offset += 4
	}

	if f.HasIPv6() {
		if l < offset+16 {
			return io.ErrUnexpectedEOF
		}
		f.UELinkSpecificIPv6AddressFor3GPPAccess = net.IP(b[offset : offset+16])
		offset += 16
	}

	if f.HasNV4() {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		f.UELinkSpecificIPv4AddressForNon3GPPAccess = net.IP(b[offset : offset+4])
		offset += 4
	}

	if f.HasNV6() {
		if l < offset+16 {
			return io.ErrUnexpectedEOF
		}
		f.UELinkSpecificIPv6AddressForNon3GPPAccess = net.IP(b[offset : offset+16])
	}

	return nil
}

// Marshal returns the serialized bytes of LinkSpecificMultipathIPAddressFields.
func (f *LinkSpecificMultipathIPAddressFields) Marshal() ([]byte, error) {
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
func (f *LinkSpecificMultipathIPAddressFields) MarshalTo(b []byte) error {
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

	if f.HasIPv4() {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		copy(b[offset:offset+4], f.UELinkSpecificIPv4AddressFor3GPPAccess.To4())
		offset += 4
	}

	if f.HasIPv6() {
		if l < offset+16 {
			return io.ErrUnexpectedEOF
		}
		copy(b[offset:offset+16], f.UELinkSpecificIPv6AddressFor3GPPAccess.To16())
		offset += 16
	}

	if f.HasNV4() {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		copy(b[offset:offset+4], f.UELinkSpecificIPv4AddressForNon3GPPAccess.To4())
		offset += 4
	}

	if f.HasNV6() {
		if l < offset+16 {
			return io.ErrUnexpectedEOF
		}
		copy(b[offset:offset+16], f.UELinkSpecificIPv6AddressForNon3GPPAccess.To16())
	}

	return nil
}

// MarshalLen returns field length in integer.
func (f *LinkSpecificMultipathIPAddressFields) MarshalLen() int {
	if f == nil {
		return 0
	}

	l := 1
	if f.HasIPv4() {
		l += 4
	}
	if f.HasIPv6() {
		l += 16
	}
	if f.HasNV4() {
		l += 4
	}
	if f.HasNV6() {
		l += 16
	}

	return l
}
