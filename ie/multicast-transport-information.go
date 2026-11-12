package ie

import (
	"net/netip"
)

// MulticastTransportInformationFields is the structured value of Multicast Transport Information.
// TS 29.244 V20.1.0 clause 8.2.207.
type MulticastTransportInformationFields struct {
	// CommonTEID is the common teid.
	CommonTEID uint32
	// MulticastAddress is the multicast address.
	MulticastAddress netip.Addr
	// SourceAddress is the source address.
	SourceAddress netip.Addr
}

// NewMulticastTransportInformation creates a Multicast Transport Information IE from its fields.
// TS 29.244 V20.1.0 clause 8.2.207.
func NewMulticastTransportInformation(fields MulticastTransportInformationFields) *IE {
	b, err := fields.Marshal()
	return addedFieldsIE(MulticastTransportInformation, b, err)
}

// MulticastTransportInformation returns the structured Multicast Transport Information value.
// TS 29.244 V20.1.0 clauses 8.2.207 and 8.1.2.
func (i *IE) MulticastTransportInformation() (*MulticastTransportInformationFields, error) {
	b, err := addedValuePayload(i, MulticastTransportInformation, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseMulticastTransportInformationFields(b)
	if err != nil {
		return nil, &DecodeError{Type: MulticastTransportInformation, Err: err}
	}
	return f, nil
}

// ParseMulticastTransportInformationFields decodes the known fields; unknown extension octets are ignored.
// TS 29.244 V20.1.0 clauses 8.2.207 and 8.1.2.
func ParseMulticastTransportInformationFields(b []byte) (*MulticastTransportInformationFields, error) {
	f := &MulticastTransportInformationFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary decodes Multicast Transport Information and resets previously decoded fields.
func (f *MulticastTransportInformationFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = MulticastTransportInformationFields{}
	r := &addedReader{b: b}
	_ = r.u8()
	f.CommonTEID = r.u32()
	f.MulticastAddress, f.SourceAddress = addedReadSSM(r)
	if r.err != nil {
		return r.err
	}

	if !f.valid() {
		return ErrInvalidValue
	}
	return nil
}
func (f *MulticastTransportInformationFields) valid() bool {
	return addedSSMValid(f.MulticastAddress, f.SourceAddress)
}

// Marshal encodes the known fields with spare bits set to zero.
func (f *MulticastTransportInformationFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	if !f.valid() {
		return nil, ErrInvalidValue
	}
	b := addedPut32([]byte{0}, f.CommonTEID)
	b = addedAppendSSM(b, f.MulticastAddress, f.SourceAddress)
	if len(b) > 65535 {
		return nil, ErrInvalidValue
	}
	return b, nil
}

// MarshalLen returns the encoded value length, or zero for nil or invalid fields.
func (f *MulticastTransportInformationFields) MarshalLen() int { b, _ := f.Marshal(); return len(b) }

// MarshalTo writes the encoded fields to b, rejecting an undersized buffer.
func (f *MulticastTransportInformationFields) MarshalTo(b []byte) error {
	wire, err := f.Marshal()
	return addedCopyTo(b, wire, err)
}
func addedSSMValid(dst, src netip.Addr) bool {
	return dst.IsValid() && src.IsValid() && dst.Zone() == "" && src.Zone() == "" && dst.IsMulticast() && !src.IsMulticast()
}
func addedAppendSSM(b []byte, dst, src netip.Addr) []byte {
	for _, a := range []netip.Addr{dst, src} {
		tag := byte(4)
		if a.Is6() {
			tag = 0x50
		}
		b = append(b, tag)
		b = append(b, a.AsSlice()...)
	}
	return b
}
func addedReadSSM(r *addedReader) (netip.Addr, netip.Addr) {
	read := func() netip.Addr {
		tag := r.u8()
		switch tag {
		case 4:
			return r.ip(4)
		case 0x50:
			return r.ip(16)
		default:
			if r.err == nil {
				r.err = ErrInvalidValue
			}
			return netip.Addr{}
		}
	}
	dst := read()
	return dst, read()
}
