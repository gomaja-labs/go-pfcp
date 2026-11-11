package ie

import (
	"io"
	"net/netip"
)

// MappedN6IPAddress flags and values, TS 29.244 V20.1.0 clause 8.2.242.
const (
	MappedN6IPAddressCHV4 uint8 = 1
	MappedN6IPAddressV4   uint8 = 2
)

// NewMappedN6IPAddress creates a Mapped N6 IP Address IE.
// TS 29.244 V20.1.0 clause 8.2.242.
func NewMappedN6IPAddress(flags uint8, ipv4Address netip.Addr) *IE {
	return recentFieldsIE(MappedN6IPAddress, NewMappedN6IPAddressFields(flags, ipv4Address))
}

// MappedN6IPAddress returns the Mapped N6 IP Address fields.
// TS 29.244 V20.1.0 clause 8.2.242.
func (i *IE) MappedN6IPAddress() (*MappedN6IPAddressFields, error) {
	b, err := recentIEPayload(i, MappedN6IPAddress, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseMappedN6IPAddressFields(b)
	if err != nil {
		return nil, &DecodeError{Type: MappedN6IPAddress, Err: err}
	}
	return f, nil
}

// MappedN6IPAddressFields models TS 29.244 V20.1.0 clause 8.2.242.
// Flags select optional fields; nonzero values for absent fields are rejected.
type MappedN6IPAddressFields struct {
	Flags       uint8
	IPv4Address netip.Addr
}

// NewMappedN6IPAddressFields stores all arguments; Marshal validates their representation.
func NewMappedN6IPAddressFields(flags uint8, ipv4Address netip.Addr) *MappedN6IPAddressFields {
	return &MappedN6IPAddressFields{flags, ipv4Address}
}

// ParseMappedN6IPAddressFields decodes fields, ignoring spare bits and trailing extension octets.
func ParseMappedN6IPAddressFields(b []byte) (*MappedN6IPAddressFields, error) {
	f := &MappedN6IPAddressFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}
func (f *MappedN6IPAddressFields) validate() error {
	if f == nil {
		return ErrIENotFound
	}
	if f.Flags & ^uint8(3) != 0 || f.Flags == 3 {
		return ErrInvalidValue
	}
	if f.Flags&2 == 0 && f.IPv4Address.IsValid() {
		return ErrInvalidValue
	}
	if f.Flags&2 != 0 && (!f.IPv4Address.Is4() || f.IPv4Address.Zone() != "") {
		return ErrInvalidValue
	}
	return nil
}

// MarshalLen returns the known field size in octets.
func (f *MappedN6IPAddressFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	n := 1
	if f.Flags&2 != 0 {
		n += 4
	}
	return n
}

// Marshal serializes fields as defined in clause 8.2.242.
func (f *MappedN6IPAddressFields) Marshal() ([]byte, error) {
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
func (f *MappedN6IPAddressFields) MarshalTo(b []byte) error {
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
	}
	return nil
}

// UnmarshalBinary replaces the receiver with the decoded known fields.
func (f *MappedN6IPAddressFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = MappedN6IPAddressFields{}
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
	}
	return f.validate()
}
