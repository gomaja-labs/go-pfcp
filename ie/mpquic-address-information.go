package ie

import (
	"fmt"
	"net/netip"
)

// MPQUICProxyType identifies the proxy service.
// TS 24.193 V19.6.0 Table 6.1.4.2-1; TS 29.244 V20.1.0 clause 8.2.226.
type MPQUICProxyType uint8

const (
	// MPQUICProxyUDP selects UDP.
	MPQUICProxyUDP MPQUICProxyType = 2
	// MPQUICProxyIP selects IP.
	MPQUICProxyIP MPQUICProxyType = 3
	// MPQUICProxyUDPAndIP selects UDP And IP.
	MPQUICProxyUDPAndIP MPQUICProxyType = 4
	// MPQUICProxyEthernet selects Ethernet.
	MPQUICProxyEthernet MPQUICProxyType = 5
)

// MPQUICAddressFlags identifies the address fields in clause 8.2.226.
type MPQUICAddressFlags uint8

const (
	// MPQUICAddressV4 marks the presence of an IPv4 address.
	MPQUICAddressV4 MPQUICAddressFlags = 1
	// MPQUICAddressV6 marks the presence of an IPv6 address.
	MPQUICAddressV6 MPQUICAddressFlags = 2
)

// MPQUICAddressInformationFields is the structured value of MPQUIC Address Information.
// TS 29.244 V20.1.0 clause 8.2.226.
// Proxy type and port are mandatory; address fields follow V4 and V6 flags.
type MPQUICAddressInformationFields struct {
	// Flags contains the presence and option bits.
	Flags MPQUICAddressFlags
	// ProxyType is the proxy type.
	ProxyType MPQUICProxyType
	// Port is the port.
	Port uint16
	// IPv4Address is the IPv4 address, when present.
	IPv4Address netip.Addr
	// IPv6Address is the IPv6 address, when present.
	IPv6Address netip.Addr
}

// NewMPQUICAddressInformation creates a MPQUIC Address Information IE from its fields.
// TS 29.244 V20.1.0 clause 8.2.226.
func NewMPQUICAddressInformation(fields MPQUICAddressInformationFields) *IE {
	b, err := fields.Marshal()
	return addedFieldsIE(MPQUICAddressInformation, b, err)
}

// MPQUICAddressInformation returns the structured MPQUIC Address Information value.
// TS 29.244 V20.1.0 clauses 8.2.226 and 8.1.2.
func (i *IE) MPQUICAddressInformation() (*MPQUICAddressInformationFields, error) {
	b, err := addedValuePayload(i, MPQUICAddressInformation, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseMPQUICAddressInformationFields(b)
	if err != nil {
		return nil, &DecodeError{Type: MPQUICAddressInformation, Err: err}
	}
	return f, nil
}

// ParseMPQUICAddressInformationFields decodes the known fields; unknown extension octets are ignored.
// TS 29.244 V20.1.0 clauses 8.2.226 and 8.1.2.
func ParseMPQUICAddressInformationFields(b []byte) (*MPQUICAddressInformationFields, error) {
	f := &MPQUICAddressInformationFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary decodes MPQUIC Address Information and resets previously decoded fields.
func (f *MPQUICAddressInformationFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = MPQUICAddressInformationFields{}
	r := &addedReader{b: b}
	f.Flags = MPQUICAddressFlags(r.u8() & 3)
	f.ProxyType = MPQUICProxyType(r.u8())
	f.Port = r.u16()
	if f.Flags&MPQUICAddressV4 != 0 {
		f.IPv4Address = r.ip(4)
	}
	if f.Flags&MPQUICAddressV6 != 0 {
		f.IPv6Address = r.ip(16)
	}
	if r.err != nil {
		return fmt.Errorf("%w: %w", ErrMalformed, r.err)
	}

	if !f.valid() {
		return fmt.Errorf("%w: %w", ErrMalformed, ErrInvalidValue)
	}
	return nil
}
func (f *MPQUICAddressInformationFields) valid() bool {
	return f.Flags&^3 == 0 && f.ProxyType >= 2 && f.ProxyType <= 5 && ((f.Flags&MPQUICAddressV4 != 0 && addedAddressValid(f.IPv4Address, 32)) || (f.Flags&MPQUICAddressV4 == 0 && !f.IPv4Address.IsValid())) && ((f.Flags&MPQUICAddressV6 != 0 && addedAddressValid(f.IPv6Address, 128)) || (f.Flags&MPQUICAddressV6 == 0 && !f.IPv6Address.IsValid()))
}

// Marshal encodes the known fields with spare bits set to zero.
func (f *MPQUICAddressInformationFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	if !f.valid() {
		return nil, ErrInvalidValue
	}
	b := addedPut16([]byte{byte(f.Flags), byte(f.ProxyType)}, f.Port)
	if f.Flags&MPQUICAddressV4 != 0 {
		b = append(b, f.IPv4Address.AsSlice()...)
	}
	if f.Flags&MPQUICAddressV6 != 0 {
		b = append(b, f.IPv6Address.AsSlice()...)
	}
	if len(b) > 65535 {
		return nil, ErrInvalidValue
	}
	return b, nil
}

// MarshalLen returns the encoded value length, or zero for nil or invalid fields.
func (f *MPQUICAddressInformationFields) MarshalLen() int { b, _ := f.Marshal(); return len(b) }

// MarshalTo writes the encoded fields to b, rejecting an undersized buffer.
func (f *MPQUICAddressInformationFields) MarshalTo(b []byte) error {
	wire, err := f.Marshal()
	return addedCopyTo(b, wire, err)
}
