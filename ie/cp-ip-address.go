package ie

import (
	"net/netip"
)

// CPIPAddressFlags identifies the address fields in clause 8.2.199.
type CPIPAddressFlags uint8

const (
	// CPIPAddressV6 marks the presence of an IPv6 address.
	CPIPAddressV6 CPIPAddressFlags = 1
	// CPIPAddressV4 marks the presence of an IPv4 address.
	CPIPAddressV4 CPIPAddressFlags = 2
)

// CPIPAddressFields is the structured value of CP IP Address.
// TS 29.244 V20.1.0 clause 8.2.199.
type CPIPAddressFields struct {
	// Flags contains the presence and option bits.
	Flags CPIPAddressFlags
	// IPv4Address is the IPv4 address, when present.
	IPv4Address netip.Addr
	// IPv6Address is the IPv6 address, when present.
	IPv6Address netip.Addr
}

// NewCPIPAddress creates a CP IP Address IE from its fields.
// TS 29.244 V20.1.0 clause 8.2.199.
func NewCPIPAddress(fields CPIPAddressFields) *IE {
	b, err := fields.Marshal()
	return addedFieldsIE(CPIPAddress, b, err)
}

// CPIPAddress returns the structured CP IP Address value.
// TS 29.244 V20.1.0 clauses 8.2.199 and 8.1.2.
func (i *IE) CPIPAddress() (*CPIPAddressFields, error) {
	b, err := addedValuePayload(i, CPIPAddress, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseCPIPAddressFields(b)
	if err != nil {
		return nil, &DecodeError{Type: CPIPAddress, Err: err}
	}
	return f, nil
}

// ParseCPIPAddressFields decodes the fields and rejects trailing octets.
// TS 29.244 V20.1.0 clauses 8.2.199 and 8.1.2.
// Decoding may borrow b; keep b unchanged while the result is in use.
func ParseCPIPAddressFields(b []byte) (*CPIPAddressFields, error) {
	f := &CPIPAddressFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary decodes CP IP Address and resets previously decoded fields.
// Decoding may borrow b; keep b unchanged while the result is in use.
func (f *CPIPAddressFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = CPIPAddressFields{}
	r := &addedReader{b: b}
	f.Flags = CPIPAddressFlags(r.u8() & 3)
	if f.Flags&CPIPAddressV4 != 0 {
		f.IPv4Address = r.ip(4)
	}
	if f.Flags&CPIPAddressV6 != 0 {
		f.IPv6Address = r.ip(16)
	}
	if r.err != nil {
		return r.err
	}

	if r.pos != len(b) {
		return ErrInvalidLength
	}
	if !f.valid() {
		return ErrInvalidValue
	}
	return nil
}
func (f *CPIPAddressFields) valid() bool {
	return f.Flags >= 1 && f.Flags <= 3 && ((f.Flags&CPIPAddressV4 != 0 && addedAddressValid(f.IPv4Address, 32)) || (f.Flags&CPIPAddressV4 == 0 && !f.IPv4Address.IsValid())) && ((f.Flags&CPIPAddressV6 != 0 && addedAddressValid(f.IPv6Address, 128)) || (f.Flags&CPIPAddressV6 == 0 && !f.IPv6Address.IsValid()))
}

// Marshal encodes the known fields with spare bits set to zero.
func (f *CPIPAddressFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	if !f.valid() {
		return nil, ErrInvalidValue
	}
	b := []byte{byte(f.Flags)}
	if f.Flags&CPIPAddressV4 != 0 {
		b = append(b, f.IPv4Address.AsSlice()...)
	}
	if f.Flags&CPIPAddressV6 != 0 {
		b = append(b, f.IPv6Address.AsSlice()...)
	}
	if len(b) > 65535 {
		return nil, ErrInvalidValue
	}
	return b, nil
}

// MarshalLen returns the encoded value length, or zero for nil or invalid fields.
func (f *CPIPAddressFields) MarshalLen() int { b, _ := f.Marshal(); return len(b) }

// MarshalTo writes the encoded fields to b, rejecting an undersized buffer.
func (f *CPIPAddressFields) MarshalTo(b []byte) error {
	wire, err := f.Marshal()
	return addedCopyTo(b, wire, err)
}
