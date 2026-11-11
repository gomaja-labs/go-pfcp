package ie

import (
	"io"
	"net/netip"
)

// MeasurementEndpointAddress flags and values, TS 29.244 V20.1.0 clause 8.2.248.
const (
	MeasurementEndpointAddressV4 uint8 = 1
	MeasurementEndpointAddressV6 uint8 = 2
)

// NewMeasurementEndpointAddress creates a Measurement Endpoint Address IE.
// TS 29.244 V20.1.0 clause 8.2.248.
func NewMeasurementEndpointAddress(flags uint8, ipv4Address netip.Addr, ipv6Address netip.Addr) *IE {
	return recentFieldsIE(MeasurementEndpointAddress, NewMeasurementEndpointAddressFields(flags, ipv4Address, ipv6Address))
}

// MeasurementEndpointAddress returns the Measurement Endpoint Address fields.
// TS 29.244 V20.1.0 clause 8.2.248.
func (i *IE) MeasurementEndpointAddress() (*MeasurementEndpointAddressFields, error) {
	b, err := recentIEPayload(i, MeasurementEndpointAddress, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseMeasurementEndpointAddressFields(b)
	if err != nil {
		return nil, &DecodeError{Type: MeasurementEndpointAddress, Err: err}
	}
	return f, nil
}

// MeasurementEndpointAddressFields models TS 29.244 V20.1.0 clause 8.2.248.
// Flags select optional fields; nonzero values for absent fields are rejected.
type MeasurementEndpointAddressFields struct {
	Flags       uint8
	IPv4Address netip.Addr
	IPv6Address netip.Addr
}

// NewMeasurementEndpointAddressFields stores all arguments; Marshal validates their representation.
func NewMeasurementEndpointAddressFields(flags uint8, ipv4Address netip.Addr, ipv6Address netip.Addr) *MeasurementEndpointAddressFields {
	return &MeasurementEndpointAddressFields{flags, ipv4Address, ipv6Address}
}

// ParseMeasurementEndpointAddressFields decodes fields, ignoring spare bits and trailing extension octets.
func ParseMeasurementEndpointAddressFields(b []byte) (*MeasurementEndpointAddressFields, error) {
	f := &MeasurementEndpointAddressFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}
func (f *MeasurementEndpointAddressFields) validate() error {
	if f == nil {
		return ErrIENotFound
	}
	if f.Flags & ^uint8(3) != 0 || f.Flags == 0 {
		return ErrInvalidValue
	}
	if f.Flags&1 == 0 && f.IPv4Address.IsValid() {
		return ErrInvalidValue
	}
	if f.Flags&1 != 0 && (!f.IPv4Address.Is4() || f.IPv4Address.Zone() != "") {
		return ErrInvalidValue
	}
	if f.Flags&2 == 0 && f.IPv6Address.IsValid() {
		return ErrInvalidValue
	}
	if f.Flags&2 != 0 && (!f.IPv6Address.Is6() || f.IPv6Address.Zone() != "") {
		return ErrInvalidValue
	}
	return nil
}

// MarshalLen returns the known field size in octets.
func (f *MeasurementEndpointAddressFields) MarshalLen() int {
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
	return n
}

// Marshal serializes fields as defined in clause 8.2.248.
func (f *MeasurementEndpointAddressFields) Marshal() ([]byte, error) {
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
func (f *MeasurementEndpointAddressFields) MarshalTo(b []byte) error {
	if err := f.validate(); err != nil {
		return err
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}
	b[0] = f.Flags
	offset := 1
	if f.Flags&1 != 0 {
		copy(b[offset:offset+4], f.IPv4Address.AsSlice())
		offset += 4
	}
	if f.Flags&2 != 0 {
		copy(b[offset:offset+16], f.IPv6Address.AsSlice())
	}
	return nil
}

// UnmarshalBinary replaces the receiver with the decoded known fields.
func (f *MeasurementEndpointAddressFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = MeasurementEndpointAddressFields{}
	if len(b) < 1 {
		return io.ErrUnexpectedEOF
	}
	f.Flags = b[0] & 3
	offset := 1
	if f.Flags&1 != 0 {
		if len(b)-offset < 4 {
			return io.ErrUnexpectedEOF
		}
		f.IPv4Address, _ = netip.AddrFromSlice(b[offset : offset+4])
		offset += 4
	}
	if f.Flags&2 != 0 {
		if len(b)-offset < 16 {
			return io.ErrUnexpectedEOF
		}
		f.IPv6Address, _ = netip.AddrFromSlice(b[offset : offset+16])
	}
	return f.validate()
}
