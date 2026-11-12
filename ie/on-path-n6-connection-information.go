package ie

import (
	"io"
	"net/netip"
)

// OnPathN6ConnectionInformation flags and values, TS 29.244 V20.1.0 clause 8.2.261.
const (
	// OnPathN6ConnectionInformationV4 marks the presence of an IPv4 address.
	OnPathN6ConnectionInformationV4 uint8 = 8
	// OnPathN6ConnectionInformationV6 marks the presence of an IPv6 address.
	OnPathN6ConnectionInformationV6 uint8 = 16
	// OnPathN6ConnectionInformationPN selects PN.
	OnPathN6ConnectionInformationPN uint8 = 32
	// OnPathN6ConnectionInformationUnspecified selects Unspecified.
	OnPathN6ConnectionInformationUnspecified uint8 = 0
	// OnPathN6ConnectionInformationConnectUDP selects Connect UDP.
	OnPathN6ConnectionInformationConnectUDP uint8 = 1
)

// NewOnPathN6ConnectionInformation creates a On-path N6 Connection Information IE.
// TS 29.244 V20.1.0 clause 8.2.261.
func NewOnPathN6ConnectionInformation(flags uint8, ipv4Address netip.Addr, ipv6Address netip.Addr, portNumber uint16) *IE {
	return recentFieldsIE(OnPathN6ConnectionInformation, NewOnPathN6ConnectionInformationFields(flags, ipv4Address, ipv6Address, portNumber))
}

// OnPathN6ConnectionInformation returns the On-path N6 Connection Information fields.
// TS 29.244 V20.1.0 clause 8.2.261.
func (i *IE) OnPathN6ConnectionInformation() (*OnPathN6ConnectionInformationFields, error) {
	b, err := recentIEPayload(i, OnPathN6ConnectionInformation, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseOnPathN6ConnectionInformationFields(b)
	if err != nil {
		return nil, &DecodeError{Type: OnPathN6ConnectionInformation, Err: err}
	}
	return f, nil
}

// OnPathN6ConnectionInformationFields models TS 29.244 V20.1.0 clause 8.2.261.
// Flags select optional fields; nonzero values for absent fields are rejected.
type OnPathN6ConnectionInformationFields struct {
	// Flags contains the presence and option bits.
	Flags uint8
	// IPv4Address is the IPv4 address, when present.
	IPv4Address netip.Addr
	// IPv6Address is the IPv6 address, when present.
	IPv6Address netip.Addr
	// PortNumber is the port number.
	PortNumber uint16
}

// NewOnPathN6ConnectionInformationFields stores all arguments; Marshal validates their representation.
func NewOnPathN6ConnectionInformationFields(flags uint8, ipv4Address netip.Addr, ipv6Address netip.Addr, portNumber uint16) *OnPathN6ConnectionInformationFields {
	return &OnPathN6ConnectionInformationFields{flags, ipv4Address, ipv6Address, portNumber}
}

// ParseOnPathN6ConnectionInformationFields decodes fields, ignoring spare bits and trailing extension octets.
func ParseOnPathN6ConnectionInformationFields(b []byte) (*OnPathN6ConnectionInformationFields, error) {
	f := &OnPathN6ConnectionInformationFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}
func (f *OnPathN6ConnectionInformationFields) validate() error {
	if f == nil {
		return ErrIENotFound
	}
	if f.Flags & ^uint8(63) != 0 || f.Flags&7 > 1 {
		return ErrInvalidValue
	}
	if f.Flags&8 == 0 && f.IPv4Address.IsValid() {
		return ErrInvalidValue
	}
	if f.Flags&8 != 0 && (!f.IPv4Address.Is4() || f.IPv4Address.Zone() != "") {
		return ErrInvalidValue
	}
	if f.Flags&16 == 0 && f.IPv6Address.IsValid() {
		return ErrInvalidValue
	}
	if f.Flags&16 != 0 && (!f.IPv6Address.Is6() || f.IPv6Address.Zone() != "") {
		return ErrInvalidValue
	}
	if f.Flags&32 == 0 && f.PortNumber != 0 {
		return ErrInvalidValue
	}
	return nil
}

// MarshalLen returns the known field size in octets.
func (f *OnPathN6ConnectionInformationFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	n := 1
	if f.Flags&8 != 0 {
		n += 4
	}
	if f.Flags&16 != 0 {
		n += 16
	}
	if f.Flags&32 != 0 {
		n += 2
	}
	return n
}

// Marshal serializes fields as defined in clause 8.2.261.
func (f *OnPathN6ConnectionInformationFields) Marshal() ([]byte, error) {
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
func (f *OnPathN6ConnectionInformationFields) MarshalTo(b []byte) error {
	if err := f.validate(); err != nil {
		return err
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}
	b[0] = f.Flags
	offset := 1
	if f.Flags&8 != 0 {
		copy(b[offset:offset+4], f.IPv4Address.AsSlice())
		offset += 4
	}
	if f.Flags&16 != 0 {
		copy(b[offset:offset+16], f.IPv6Address.AsSlice())
		offset += 16
	}
	if f.Flags&32 != 0 {
		putRecentUint(b[offset:offset+2], uint64(f.PortNumber))
	}
	return nil
}

// UnmarshalBinary replaces the receiver with the decoded known fields.
func (f *OnPathN6ConnectionInformationFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = OnPathN6ConnectionInformationFields{}
	if len(b) < 1 {
		return io.ErrUnexpectedEOF
	}
	f.Flags = b[0] & 63
	offset := 1
	if f.Flags&8 != 0 {
		if len(b)-offset < 4 {
			return io.ErrUnexpectedEOF
		}
		f.IPv4Address, _ = netip.AddrFromSlice(b[offset : offset+4])
		offset += 4
	}
	if f.Flags&16 != 0 {
		if len(b)-offset < 16 {
			return io.ErrUnexpectedEOF
		}
		f.IPv6Address, _ = netip.AddrFromSlice(b[offset : offset+16])
		offset += 16
	}
	if f.Flags&32 != 0 {
		if len(b)-offset < 2 {
			return io.ErrUnexpectedEOF
		}
		f.PortNumber = uint16(recentUint(b[offset : offset+2]))
	}
	return f.validate()
}
