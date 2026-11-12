package ie

import (
	"io"
	"net/netip"
)

// N6RoutingInformation flags and values, TS 29.244 V20.1.0 clause 8.2.243.
const (
	// N6RoutingInformationSIPV4 selects SIPV4.
	N6RoutingInformationSIPV4 uint8 = 1
	// N6RoutingInformationSIPV6 selects SIPV6.
	N6RoutingInformationSIPV6 uint8 = 2
	// N6RoutingInformationSPO selects SPO.
	N6RoutingInformationSPO uint8 = 4
	// N6RoutingInformationDIPV4 selects DIPV4.
	N6RoutingInformationDIPV4 uint8 = 8
	// N6RoutingInformationDIPV6 selects DIPV6.
	N6RoutingInformationDIPV6 uint8 = 16
	// N6RoutingInformationDPO selects DPO.
	N6RoutingInformationDPO uint8 = 32
)

// NewN6RoutingInformation creates a N6 Routing Information IE.
// TS 29.244 V20.1.0 clause 8.2.243.
func NewN6RoutingInformation(flags uint8, sourceIPv4Address netip.Addr, sourceIPv6Address netip.Addr, sourcePort uint16, destinationIPv4Address netip.Addr, destinationIPv6Address netip.Addr, destinationPort uint16) *IE {
	return recentFieldsIE(N6RoutingInformation, NewN6RoutingInformationFields(flags, sourceIPv4Address, sourceIPv6Address, sourcePort, destinationIPv4Address, destinationIPv6Address, destinationPort))
}

// N6RoutingInformation returns the N6 Routing Information fields.
// TS 29.244 V20.1.0 clause 8.2.243.
func (i *IE) N6RoutingInformation() (*N6RoutingInformationFields, error) {
	b, err := recentIEPayload(i, N6RoutingInformation, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseN6RoutingInformationFields(b)
	if err != nil {
		return nil, &DecodeError{Type: N6RoutingInformation, Err: err}
	}
	return f, nil
}

// N6RoutingInformationFields models TS 29.244 V20.1.0 clause 8.2.243.
// Flags select optional fields; nonzero values for absent fields are rejected.
type N6RoutingInformationFields struct {
	// Flags contains the presence and option bits.
	Flags uint8
	// SourceIPv4Address is the source IPv4 address.
	SourceIPv4Address netip.Addr
	// SourceIPv6Address is the source IPv6 address.
	SourceIPv6Address netip.Addr
	// SourcePort is the source port.
	SourcePort uint16
	// DestinationIPv4Address is the destination IPv4 address.
	DestinationIPv4Address netip.Addr
	// DestinationIPv6Address is the destination IPv6 address.
	DestinationIPv6Address netip.Addr
	// DestinationPort is the destination port.
	DestinationPort uint16
}

// NewN6RoutingInformationFields stores all arguments; Marshal validates their representation.
func NewN6RoutingInformationFields(flags uint8, sourceIPv4Address netip.Addr, sourceIPv6Address netip.Addr, sourcePort uint16, destinationIPv4Address netip.Addr, destinationIPv6Address netip.Addr, destinationPort uint16) *N6RoutingInformationFields {
	return &N6RoutingInformationFields{flags, sourceIPv4Address, sourceIPv6Address, sourcePort, destinationIPv4Address, destinationIPv6Address, destinationPort}
}

// ParseN6RoutingInformationFields decodes fields, ignoring spare bits and trailing extension octets.
func ParseN6RoutingInformationFields(b []byte) (*N6RoutingInformationFields, error) {
	f := &N6RoutingInformationFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}
func (f *N6RoutingInformationFields) validate() error {
	if f == nil {
		return ErrIENotFound
	}
	if f.Flags & ^uint8(63) != 0 {
		return ErrInvalidValue
	}
	if f.Flags&1 == 0 && f.SourceIPv4Address.IsValid() {
		return ErrInvalidValue
	}
	if f.Flags&1 != 0 && (!f.SourceIPv4Address.Is4() || f.SourceIPv4Address.Zone() != "") {
		return ErrInvalidValue
	}
	if f.Flags&2 == 0 && f.SourceIPv6Address.IsValid() {
		return ErrInvalidValue
	}
	if f.Flags&2 != 0 && (!f.SourceIPv6Address.Is6() || f.SourceIPv6Address.Zone() != "") {
		return ErrInvalidValue
	}
	if f.Flags&4 == 0 && f.SourcePort != 0 {
		return ErrInvalidValue
	}
	if f.Flags&8 == 0 && f.DestinationIPv4Address.IsValid() {
		return ErrInvalidValue
	}
	if f.Flags&8 != 0 && (!f.DestinationIPv4Address.Is4() || f.DestinationIPv4Address.Zone() != "") {
		return ErrInvalidValue
	}
	if f.Flags&16 == 0 && f.DestinationIPv6Address.IsValid() {
		return ErrInvalidValue
	}
	if f.Flags&16 != 0 && (!f.DestinationIPv6Address.Is6() || f.DestinationIPv6Address.Zone() != "") {
		return ErrInvalidValue
	}
	if f.Flags&32 == 0 && f.DestinationPort != 0 {
		return ErrInvalidValue
	}
	return nil
}

// MarshalLen returns the known field size in octets.
func (f *N6RoutingInformationFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	n := 1
	if f.Flags&1 != 0 {
		n += 4
	}
	if f.Flags&2 != 0 {
		n += 16
	}
	if f.Flags&4 != 0 {
		n += 2
	}
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

// Marshal serializes fields as defined in clause 8.2.243.
func (f *N6RoutingInformationFields) Marshal() ([]byte, error) {
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
func (f *N6RoutingInformationFields) MarshalTo(b []byte) error {
	if err := f.validate(); err != nil {
		return err
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}
	b[0] = f.Flags
	offset := 1
	if f.Flags&1 != 0 {
		copy(b[offset:offset+4], f.SourceIPv4Address.AsSlice())
		offset += 4
	}
	if f.Flags&2 != 0 {
		copy(b[offset:offset+16], f.SourceIPv6Address.AsSlice())
		offset += 16
	}
	if f.Flags&4 != 0 {
		putRecentUint(b[offset:offset+2], uint64(f.SourcePort))
		offset += 2
	}
	if f.Flags&8 != 0 {
		copy(b[offset:offset+4], f.DestinationIPv4Address.AsSlice())
		offset += 4
	}
	if f.Flags&16 != 0 {
		copy(b[offset:offset+16], f.DestinationIPv6Address.AsSlice())
		offset += 16
	}
	if f.Flags&32 != 0 {
		putRecentUint(b[offset:offset+2], uint64(f.DestinationPort))
	}
	return nil
}

// UnmarshalBinary replaces the receiver with the decoded known fields.
func (f *N6RoutingInformationFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = N6RoutingInformationFields{}
	if len(b) < 1 {
		return io.ErrUnexpectedEOF
	}
	f.Flags = b[0] & 63
	offset := 1
	if f.Flags&1 != 0 {
		if len(b)-offset < 4 {
			return io.ErrUnexpectedEOF
		}
		f.SourceIPv4Address, _ = netip.AddrFromSlice(b[offset : offset+4])
		offset += 4
	}
	if f.Flags&2 != 0 {
		if len(b)-offset < 16 {
			return io.ErrUnexpectedEOF
		}
		f.SourceIPv6Address, _ = netip.AddrFromSlice(b[offset : offset+16])
		offset += 16
	}
	if f.Flags&4 != 0 {
		if len(b)-offset < 2 {
			return io.ErrUnexpectedEOF
		}
		f.SourcePort = uint16(recentUint(b[offset : offset+2]))
		offset += 2
	}
	if f.Flags&8 != 0 {
		if len(b)-offset < 4 {
			return io.ErrUnexpectedEOF
		}
		f.DestinationIPv4Address, _ = netip.AddrFromSlice(b[offset : offset+4])
		offset += 4
	}
	if f.Flags&16 != 0 {
		if len(b)-offset < 16 {
			return io.ErrUnexpectedEOF
		}
		f.DestinationIPv6Address, _ = netip.AddrFromSlice(b[offset : offset+16])
		offset += 16
	}
	if f.Flags&32 != 0 {
		if len(b)-offset < 2 {
			return io.ErrUnexpectedEOF
		}
		f.DestinationPort = uint16(recentUint(b[offset : offset+2]))
	}
	return f.validate()
}
