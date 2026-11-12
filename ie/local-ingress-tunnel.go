package ie

import (
	"net/netip"
)

// LocalIngressTunnelFlags identifies the tunnel allocation and family in clause 8.2.209.
type LocalIngressTunnelFlags uint8

const (
	// LocalIngressV4 marks the presence of an IPv4 address.
	LocalIngressV4 LocalIngressTunnelFlags = 1
	// LocalIngressV6 marks the presence of an IPv6 address.
	LocalIngressV6 LocalIngressTunnelFlags = 2
	// LocalIngressCH selects CH.
	LocalIngressCH LocalIngressTunnelFlags = 4
)

// LocalIngressTunnelFields is the structured value of Local Ingress Tunnel.
// TS 29.244 V20.1.0 clause 8.2.209.
type LocalIngressTunnelFields struct {
	// Flags contains the presence and option bits.
	Flags LocalIngressTunnelFlags
	// Port is the port.
	Port uint16
	// Address is the address.
	Address netip.Addr
}

// NewLocalIngressTunnel creates a Local Ingress Tunnel IE from its fields.
// TS 29.244 V20.1.0 clause 8.2.209.
func NewLocalIngressTunnel(fields LocalIngressTunnelFields) *IE {
	b, err := fields.Marshal()
	return addedFieldsIE(LocalIngressTunnel, b, err)
}

// LocalIngressTunnel returns the structured Local Ingress Tunnel value.
// TS 29.244 V20.1.0 clauses 8.2.209 and 8.1.2.
func (i *IE) LocalIngressTunnel() (*LocalIngressTunnelFields, error) {
	b, err := addedValuePayload(i, LocalIngressTunnel, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseLocalIngressTunnelFields(b)
	if err != nil {
		return nil, &DecodeError{Type: LocalIngressTunnel, Err: err}
	}
	return f, nil
}

// ParseLocalIngressTunnelFields decodes the known fields; unknown extension octets are ignored.
// TS 29.244 V20.1.0 clauses 8.2.209 and 8.1.2.
// Decoding may borrow b; keep b unchanged while the result is in use.
func ParseLocalIngressTunnelFields(b []byte) (*LocalIngressTunnelFields, error) {
	f := &LocalIngressTunnelFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary decodes Local Ingress Tunnel and resets previously decoded fields.
// Decoding may borrow b; keep b unchanged while the result is in use.
func (f *LocalIngressTunnelFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = LocalIngressTunnelFields{}
	r := &addedReader{b: b}
	f.Flags = LocalIngressTunnelFlags(r.u8() & 7)
	if f.Flags&LocalIngressCH == 0 {
		f.Port = r.u16()
		if f.Flags&LocalIngressV4 != 0 {
			f.Address = r.ip(4)
		} else if f.Flags&LocalIngressV6 != 0 {
			f.Address = r.ip(16)
		}
	}
	if r.err != nil {
		return r.err
	}

	if !f.valid() {
		return ErrInvalidValue
	}
	return nil
}
func (f *LocalIngressTunnelFields) valid() bool {
	return f.Flags&^7 == 0 && (f.Flags&3 == 1 || f.Flags&3 == 2) && ((f.Flags&LocalIngressCH != 0 && f.Port == 0 && !f.Address.IsValid()) || (f.Flags&LocalIngressCH == 0 && ((f.Flags&LocalIngressV4 != 0 && addedAddressValid(f.Address, 32)) || (f.Flags&LocalIngressV6 != 0 && addedAddressValid(f.Address, 128)))))
}

// Marshal encodes the known fields with spare bits set to zero.
func (f *LocalIngressTunnelFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	if !f.valid() {
		return nil, ErrInvalidValue
	}
	b := []byte{byte(f.Flags)}
	if f.Flags&LocalIngressCH == 0 {
		b = addedPut16(b, f.Port)
		b = append(b, f.Address.AsSlice()...)
	}
	if len(b) > 65535 {
		return nil, ErrInvalidValue
	}
	return b, nil
}

// MarshalLen returns the encoded value length, or zero for nil or invalid fields.
func (f *LocalIngressTunnelFields) MarshalLen() int { b, _ := f.Marshal(); return len(b) }

// MarshalTo writes the encoded fields to b, rejecting an undersized buffer.
func (f *LocalIngressTunnelFields) MarshalTo(b []byte) error {
	wire, err := f.Marshal()
	return addedCopyTo(b, wire, err)
}
