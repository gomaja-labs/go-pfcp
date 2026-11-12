// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
	"net/netip"
	"unicode/utf8"
)

// NewRedirectInformation creates a new RedirectInformation IE.
func NewRedirectInformation(addrType uint8, addrs ...string) *IE {
	fields := NewRedirectInformationFields(addrType, addrs...)

	b, err := fields.Marshal()
	if err != nil {
		return newInvalidIE(RedirectInformation, ErrInvalidValue)
	}

	return New(RedirectInformation, b)
}

// RedirectInformation returns RedirectInformation in structured format if the type of IE matches.
func (i *IE) RedirectInformation() (*RedirectInformationFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case RedirectInformation:
		f, err := ParseRedirectInformationFields(i.Payload)
		if err != nil {
			return nil, err
		}
		return f, nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// RedirectAddressType definitions.
const (
	// RedirectAddrIPv4 selects IPv4.
	RedirectAddrIPv4 uint8 = 0
	// RedirectAddrIPv6 selects IPv6.
	RedirectAddrIPv6 uint8 = 1
	// RedirectAddrURL selects URL.
	RedirectAddrURL uint8 = 2
	// RedirectAddrSIPURI selects SIPURI.
	RedirectAddrSIPURI uint8 = 3
	// RedirectAddrIPv4AndIPv6 selects IPv4 And IPv6.
	RedirectAddrIPv4AndIPv6 uint8 = 4
	// RedirectAddrPort selects Port.
	RedirectAddrPort uint8 = 5
	// RedirectAddrIPv4AndPort selects IPv4 And Port.
	RedirectAddrIPv4AndPort uint8 = 6
	// RedirectAddrIPv6AndPort selects IPv6 And Port.
	RedirectAddrIPv6AndPort uint8 = 7
	// RedirectAddrIPv4AndIPv6AndPort selects IPv4 And IPv6 And Port.
	RedirectAddrIPv4AndIPv6AndPort uint8 = 8
)

// RedirectInformationFields represents addresses and a conditional redirect port.
// TS 29.244 V20.1.0 clause 8.2.20. Address text follows RFC 8506
// section 8.38, which replaces the RFC 4006 reference in the 3GPP text.
type RedirectInformationFields struct {
	// RedirectAddressType is the redirect address type.
	RedirectAddressType uint8
	// ServerAddrLength is the encoded server addr length in octets.
	ServerAddrLength uint16
	// RedirectServerAddress is the redirect server address.
	RedirectServerAddress string
	// OtherServerAddrLength is the encoded other server addr length in octets.
	OtherServerAddrLength uint16
	// OtherRedirectServerAddress is the other redirect server address.
	OtherRedirectServerAddress string
	// RedirectPort is the redirect port.
	RedirectPort uint16
}

// NewRedirectInformationFields creates an address-only redirect. Port redirects
// use NewRedirectInformationFromFields to supply the port explicitly.
func NewRedirectInformationFields(addrType uint8, addrs ...string) *RedirectInformationFields {
	f := &RedirectInformationFields{RedirectAddressType: addrType}
	if len(addrs) > 0 {
		f.RedirectServerAddress = addrs[0]
		f.ServerAddrLength = uint16(len(addrs[0]))
	}
	if len(addrs) > 1 {
		f.OtherRedirectServerAddress = addrs[1]
		f.OtherServerAddrLength = uint16(len(addrs[1]))
	}
	if len(addrs) > 2 {
		f.RedirectAddressType = 255
	}
	return f
}

// ParseRedirectInformationFields decodes all current address types.
// Decoding may borrow b; keep b unchanged while the result is in use.
func ParseRedirectInformationFields(b []byte) (*RedirectInformationFields, error) {
	f := &RedirectInformationFields{}
	if e := f.UnmarshalBinary(b); e != nil {
		return nil, e
	}
	return f, nil
}

// UnmarshalBinary decodes addresses and port and ignores extension octets.
// Decoding may borrow b; keep b unchanged while the result is in use.
func (f *RedirectInformationFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = RedirectInformationFields{}
	if len(b) < 1 {
		return io.ErrUnexpectedEOF
	}
	f.RedirectAddressType = b[0] & 15
	if f.RedirectAddressType > 8 {
		return ErrInvalidValue
	}
	o := 1
	if f.RedirectAddressType != RedirectAddrPort {
		v, e := takeLength16(b, &o)
		if e != nil {
			return e
		}
		f.ServerAddrLength = uint16(len(v))
		f.RedirectServerAddress = string(v)
	}
	if f.hasOtherAddress() {
		v, e := takeLength16(b, &o)
		if e != nil {
			return e
		}
		f.OtherServerAddrLength = uint16(len(v))
		f.OtherRedirectServerAddress = string(v)
	}
	if f.RedirectAddressType >= RedirectAddrPort {
		if len(b) < o+2 {
			return io.ErrUnexpectedEOF
		}
		f.RedirectPort = binary.BigEndian.Uint16(b[o : o+2])
	}
	return f.validateAddresses(false)
}
func (f *RedirectInformationFields) hasOtherAddress() bool {
	return f.RedirectAddressType == RedirectAddrIPv4AndIPv6 || f.RedirectAddressType == RedirectAddrIPv4AndIPv6AndPort
}

// Marshal serializes redirect information.
func (f *RedirectInformationFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	b := make([]byte, f.MarshalLen())
	if e := f.MarshalTo(b); e != nil {
		return nil, e
	}
	return b, nil
}

// MarshalLen returns the serialized length.
func (f *RedirectInformationFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	n := 1
	if f.RedirectAddressType != RedirectAddrPort {
		n += 2 + len(f.RedirectServerAddress)
	}
	if f.hasOtherAddress() {
		n += 2 + len(f.OtherRedirectServerAddress)
	}
	if f.RedirectAddressType >= RedirectAddrPort {
		n += 2
	}
	return n
}

// MarshalTo writes the complete redirect information.
func (f *RedirectInformationFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if f.RedirectAddressType > 8 {
		return ErrInvalidValue
	}
	if len(f.RedirectServerAddress) != int(f.ServerAddrLength) || len(f.OtherRedirectServerAddress) != int(f.OtherServerAddrLength) || f.MarshalLen() > 65535 {
		return ErrInvalidLength
	}
	if err := f.validateAddresses(true); err != nil {
		return err
	}
	if f.RedirectAddressType == RedirectAddrPort && f.RedirectServerAddress != "" || !f.hasOtherAddress() && f.OtherRedirectServerAddress != "" || f.RedirectAddressType < RedirectAddrPort && f.RedirectPort != 0 {
		return ErrInvalidValue
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}
	b[0] = f.RedirectAddressType
	o := 1
	if f.RedirectAddressType != RedirectAddrPort {
		binary.BigEndian.PutUint16(b[o:o+2], f.ServerAddrLength)
		o += 2
		copy(b[o:], f.RedirectServerAddress)
		o += len(f.RedirectServerAddress)
	}
	if f.hasOtherAddress() {
		binary.BigEndian.PutUint16(b[o:o+2], f.OtherServerAddrLength)
		o += 2
		copy(b[o:], f.OtherRedirectServerAddress)
		o += len(f.OtherRedirectServerAddress)
	}
	if f.RedirectAddressType >= RedirectAddrPort {
		binary.BigEndian.PutUint16(b[o:o+2], f.RedirectPort)
	}
	return nil
}

// NewRedirectInformationFromFields creates a redirect with explicit port/address fields.
func NewRedirectInformationFromFields(fields *RedirectInformationFields) *IE {
	b, e := fields.Marshal()
	if e != nil {
		return newInvalidIE(RedirectInformation, ErrInvalidValue)
	}
	return New(RedirectInformation, b)
}

// TS 29.244 V20.1.0 clause 8.2.20 and RFC 8506 section 8.38:
// addresses are nonempty UTF-8 and use the specified IP family. Send IPv6
// in RFC 5952 canonical form; receive other valid spellings without changing them.
func (f *RedirectInformationFields) validateAddresses(sending bool) error {
	if f.RedirectAddressType > 8 {
		return ErrInvalidValue
	}
	if f.RedirectAddressType == RedirectAddrPort {
		return nil
	}
	values := []struct {
		value  string
		family uint8
	}{{f.RedirectServerAddress, 0}}
	switch f.RedirectAddressType {
	case RedirectAddrIPv4, RedirectAddrIPv4AndPort:
		values[0].family = 4
	case RedirectAddrIPv6, RedirectAddrIPv6AndPort:
		values[0].family = 6
	case RedirectAddrIPv4AndIPv6, RedirectAddrIPv4AndIPv6AndPort:
		values[0].family = 4
		values = append(values, struct {
			value  string
			family uint8
		}{f.OtherRedirectServerAddress, 6})
	}
	for _, v := range values {
		if v.value == "" || !utf8.ValidString(v.value) {
			return ErrInvalidValue
		}
		if v.family == 0 {
			continue
		}
		addr, err := netip.ParseAddr(v.value)
		if err != nil || addr.Zone() != "" || v.family == 4 && !addr.Is4() || v.family == 6 && !addr.Is6() {
			return ErrInvalidValue
		}
		if sending && v.family == 6 && addr.String() != v.value {
			return ErrInvalidValue
		}
	}
	return nil
}
