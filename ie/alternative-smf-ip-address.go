// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"fmt"
	"io"
	"net"
)

// NewAlternativeSMFIPAddress creates a new AlternativeSMFIPAddress IE.
func NewAlternativeSMFIPAddress(v4, v6 net.IP) *IE {
	fields := NewAlternativeSMFIPAddressFields(v4, v6)

	b, err := fields.Marshal()
	if err != nil {
		return newInvalidIE(AlternativeSMFIPAddress, fmt.Errorf("%w: %w", ErrInvalidValue, err))
	}

	return New(AlternativeSMFIPAddress, b)
}

// AlternativeSMFIPAddress returns AlternativeSMFIPAddress in structured format if the type of IE matches.
func (i *IE) AlternativeSMFIPAddress() (*AlternativeSMFIPAddressFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	if i.Type != AlternativeSMFIPAddress {
		return nil, &InvalidTypeError{Type: i.Type}
	}

	fields, err := ParseAlternativeSMFIPAddressFields(i.Payload)
	if err != nil {
		return nil, err
	}

	return fields, nil
}

// AlternativeSMFIPAddressFields represents a fields contained in AlternativeSMFIPAddress IE.
type AlternativeSMFIPAddressFields struct {
	Flags       uint8
	IPv4Address net.IP
	IPv6Address net.IP
}

// NewAlternativeSMFIPAddressFields creates a new NewAlternativeSMFIPAddressFields.
func NewAlternativeSMFIPAddressFields(v4, v6 net.IP) *AlternativeSMFIPAddressFields {
	f := &AlternativeSMFIPAddressFields{}

	if v4 != nil {
		f.IPv4Address = v4
		f.SetIPv4Flag()
	}
	if v6 != nil {
		f.IPv6Address = v6
		f.SetIPv6Flag()
	}

	return f
}

// HasIPv4 reports whether IPv4 flag is set.
func (f *AlternativeSMFIPAddressFields) HasIPv4() bool {
	if f == nil {
		return false
	}

	return has2ndBit(f.Flags)
}

// SetIPv4Flag sets IPv4 flag in AlternativeSMFIPAddress.
func (f *AlternativeSMFIPAddressFields) SetIPv4Flag() {
	if f == nil {
		return
	}

	f.Flags |= 0x02
}

// HasIPv6 reports whether IPv6 flag is set.
func (f *AlternativeSMFIPAddressFields) HasIPv6() bool {
	if f == nil {
		return false
	}

	return has1stBit(f.Flags)
}

// SetIPv6Flag sets IPv6 flag in AlternativeSMFIPAddress.
func (f *AlternativeSMFIPAddressFields) SetIPv6Flag() {
	if f == nil {
		return
	}

	f.Flags |= 0x01
}

// ParseAlternativeSMFIPAddressFields parses b into AlternativeSMFIPAddressFields.
func ParseAlternativeSMFIPAddressFields(b []byte) (*AlternativeSMFIPAddressFields, error) {
	f := &AlternativeSMFIPAddressFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary parses b into IE.
func (f *AlternativeSMFIPAddressFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = AlternativeSMFIPAddressFields{}

	l := len(b)
	if l < 2 {
		return io.ErrUnexpectedEOF
	}

	f.Flags = b[0]
	offset := 1

	if f.HasIPv4() {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		f.IPv4Address = net.IP(b[offset : offset+4])
		offset += 4
	}

	if f.HasIPv6() {
		if l < offset+16 {
			return io.ErrUnexpectedEOF
		}
		f.IPv6Address = net.IP(b[offset : offset+16])
	}

	return nil
}

// Marshal returns the serialized bytes of AlternativeSMFIPAddressFields.
func (f *AlternativeSMFIPAddressFields) Marshal() ([]byte, error) {
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
func (f *AlternativeSMFIPAddressFields) MarshalTo(b []byte) error {
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

	if f.IPv4Address != nil {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		copy(b[offset:offset+4], f.IPv4Address.To4())
		offset += 4
	}
	if f.IPv6Address != nil {
		if l < offset+16 {
			return io.ErrUnexpectedEOF
		}
		copy(b[offset:offset+16], f.IPv6Address.To16())
	}

	return nil
}

// MarshalLen returns field length in integer.
func (f *AlternativeSMFIPAddressFields) MarshalLen() int {
	if f == nil {
		return 0
	}

	l := 1
	if f.IPv4Address != nil {
		l += 4
	}
	if f.IPv6Address != nil {
		l += 16
	}

	return l
}

// HasPPE reports whether this is the preferred alternative PFCP entity.
// TS 29.244 V20.1.0 clause 8.2.129, octet 5 bit 3.
func (i *IE) HasPPE() bool { return payloadFlag(i, AlternativeSMFIPAddress, 0, 0x04) }

// HasPPE reports whether the preferred PFCP entity flag is set.
// TS 29.244 V20.1.0 clause 8.2.129, octet 5 bit 3.
func (f *AlternativeSMFIPAddressFields) HasPPE() bool { return f != nil && f.Flags&0x04 != 0 }

// SetPPEFlag marks the alternative SMF address as the preferred PFCP entity.
// TS 29.244 V20.1.0 clause 8.2.129, octet 5 bit 3.
func (f *AlternativeSMFIPAddressFields) SetPPEFlag() {
	if f != nil {
		f.Flags |= 0x04
	}
}
