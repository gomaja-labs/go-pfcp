package ie

import (
	"io"
	"net/netip"
)

// MoQRelayIPAddress flags and values, TS 29.244 V20.1.0 clause 8.2.270.
const (
	// MoQRelayIPAddressV4 marks the presence of an IPv4 address.
	MoQRelayIPAddressV4 uint8 = 2
	// MoQRelayIPAddressV6 marks the presence of an IPv6 address.
	MoQRelayIPAddressV6 uint8 = 1
)

// NewMoQRelayIPAddress creates a MoQ Relay IP Address IE.
// TS 29.244 V20.1.0 clause 8.2.270.
func NewMoQRelayIPAddress(flags uint8, ipv4Address netip.Addr, ipv6Address netip.Addr) *IE {
	return recentFieldsIE(MoQRelayIPAddress, NewMoQRelayIPAddressFields(flags, ipv4Address, ipv6Address))
}

// MoQRelayIPAddress returns the MoQ Relay IP Address fields.
// TS 29.244 V20.1.0 clause 8.2.270.
func (i *IE) MoQRelayIPAddress() (*MoQRelayIPAddressFields, error) {
	b, err := recentIEPayload(i, MoQRelayIPAddress, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseMoQRelayIPAddressFields(b)
	if err != nil {
		return nil, &DecodeError{Type: MoQRelayIPAddress, Err: err}
	}
	return f, nil
}

// MoQRelayIPAddressFields models TS 29.244 V20.1.0 clause 8.2.270.
// Flags select optional fields; nonzero values for absent fields are rejected.
type MoQRelayIPAddressFields struct {
	// Flags contains the presence and option bits.
	Flags uint8
	// IPv4Address is the IPv4 address, when present.
	IPv4Address netip.Addr
	// IPv6Address is the IPv6 address, when present.
	IPv6Address netip.Addr
}

// NewMoQRelayIPAddressFields stores all arguments; Marshal validates their representation.
func NewMoQRelayIPAddressFields(flags uint8, ipv4Address netip.Addr, ipv6Address netip.Addr) *MoQRelayIPAddressFields {
	return &MoQRelayIPAddressFields{flags, ipv4Address, ipv6Address}
}

// ParseMoQRelayIPAddressFields decodes fields, ignoring spare bits and trailing extension octets.
// Decoding may borrow b; keep b unchanged while the result is in use.
func ParseMoQRelayIPAddressFields(b []byte) (*MoQRelayIPAddressFields, error) {
	f := &MoQRelayIPAddressFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}
func (f *MoQRelayIPAddressFields) validate() error {
	if f == nil {
		return ErrIENotFound
	}
	if f.Flags & ^uint8(3) != 0 {
		return ErrInvalidValue
	}
	if f.Flags&2 == 0 && f.IPv4Address.IsValid() {
		return ErrInvalidValue
	}
	if f.Flags&2 != 0 && (!f.IPv4Address.Is4() || f.IPv4Address.Zone() != "") {
		return ErrInvalidValue
	}
	if f.Flags&1 == 0 && f.IPv6Address.IsValid() {
		return ErrInvalidValue
	}
	if f.Flags&1 != 0 && (!f.IPv6Address.Is6() || f.IPv6Address.Zone() != "") {
		return ErrInvalidValue
	}
	return nil
}

// MarshalLen returns the known field size in octets.
func (f *MoQRelayIPAddressFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	n := 1
	if f.Flags&2 != 0 {
		n += 4
	}
	if f.Flags&1 != 0 {
		n += 16
	}
	return n
}

// Marshal serializes fields as defined in clause 8.2.270.
func (f *MoQRelayIPAddressFields) Marshal() ([]byte, error) {
	if err := f.validate(); err != nil {
		return nil, err
	}
	b := make([]byte, f.MarshalLen())
	if err := f.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo writes fields to b, rejecting invalid values and short buffers.
func (f *MoQRelayIPAddressFields) MarshalTo(b []byte) error {
	if err := f.validate(); err != nil {
		return err
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}
	b[0] = f.Flags
	offset := 1
	if f.Flags&2 != 0 {
		copy(b[offset:offset+4], f.IPv4Address.AsSlice())
		offset += 4
	}
	if f.Flags&1 != 0 {
		copy(b[offset:offset+16], f.IPv6Address.AsSlice())
	}
	return nil
}

// UnmarshalBinary replaces the receiver with the decoded known fields.
// Decoding may borrow b; keep b unchanged while the result is in use.
func (f *MoQRelayIPAddressFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = MoQRelayIPAddressFields{}
	if len(b) < 1 {
		return io.ErrUnexpectedEOF
	}
	f.Flags = b[0] & 3
	offset := 1
	if f.Flags&2 != 0 {
		if len(b)-offset < 4 {
			return io.ErrUnexpectedEOF
		}
		f.IPv4Address, _ = netip.AddrFromSlice(b[offset : offset+4])
		offset += 4
	}
	if f.Flags&1 != 0 {
		if len(b)-offset < 16 {
			return io.ErrUnexpectedEOF
		}
		f.IPv6Address, _ = netip.AddrFromSlice(b[offset : offset+16])
	}
	return f.validate()
}
