package ie

import (
	"fmt"
	"net/netip"
)

// MBSSessionIdentifierFlags identifies the fields in clause 8.2.206.
type MBSSessionIdentifierFlags uint8

const (
	// MBSSessionTMGI selects TMGI.
	MBSSessionTMGI MBSSessionIdentifierFlags = 1
	// MBSSessionSSMI selects SSMI.
	MBSSessionSSMI MBSSessionIdentifierFlags = 2
	// MBSSessionNIDI selects NIDI.
	MBSSessionNIDI MBSSessionIdentifierFlags = 4
)

// MBSSessionIdentifierFields is the structured value of MBS Session Identifier.
// TS 29.244 V20.1.0 clause 8.2.206.
// All value fields follow TMGI, SSMI and NIDI flags; no field is inferred.
type MBSSessionIdentifierFields struct {
	// Flags contains the presence and option bits.
	Flags MBSSessionIdentifierFlags
	// TMGI is the tmgi.
	TMGI [6]byte
	// MulticastAddress is the multicast address.
	MulticastAddress netip.Addr
	// SourceAddress is the source address.
	SourceAddress netip.Addr
	// NID is the nid.
	NID uint64
}

// NewMBSSessionIdentifier creates a MBS Session Identifier IE from its fields.
// TS 29.244 V20.1.0 clause 8.2.206.
func NewMBSSessionIdentifier(fields MBSSessionIdentifierFields) *IE {
	b, err := fields.Marshal()
	return addedFieldsIE(MBSSessionIdentifier, b, err)
}

// MBSSessionIdentifier returns the structured MBS Session Identifier value.
// TS 29.244 V20.1.0 clauses 8.2.206 and 8.1.2.
func (i *IE) MBSSessionIdentifier() (*MBSSessionIdentifierFields, error) {
	b, err := addedValuePayload(i, MBSSessionIdentifier, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseMBSSessionIdentifierFields(b)
	if err != nil {
		return nil, &DecodeError{Type: MBSSessionIdentifier, Err: err}
	}
	return f, nil
}

// ParseMBSSessionIdentifierFields decodes the known fields; unknown extension octets are ignored.
// TS 29.244 V20.1.0 clauses 8.2.206 and 8.1.2.
// Decoding may borrow b; keep b unchanged while the result is in use.
func ParseMBSSessionIdentifierFields(b []byte) (*MBSSessionIdentifierFields, error) {
	f := &MBSSessionIdentifierFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary decodes MBS Session Identifier and resets previously decoded fields.
// Decoding may borrow b; keep b unchanged while the result is in use.
func (f *MBSSessionIdentifierFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = MBSSessionIdentifierFields{}
	r := &addedReader{b: b}
	f.Flags = MBSSessionIdentifierFlags(r.u8() & 7)
	if f.Flags&MBSSessionTMGI != 0 {
		copy(f.TMGI[:], r.take(6))
	}
	if f.Flags&MBSSessionSSMI != 0 {
		f.MulticastAddress, f.SourceAddress = addedReadSSM(r)
	}
	if f.Flags&MBSSessionNIDI != 0 {
		v := r.take(6)
		if r.err == nil {
			for _, o := range v {
				f.NID = f.NID<<8 | uint64(o)
			}
			f.NID >>= 4
		}
	}
	if r.err != nil {
		return fmt.Errorf("%w: %w", ErrMalformed, r.err)
	}

	if !f.valid() {
		return fmt.Errorf("%w: %w", ErrMalformed, ErrInvalidValue)
	}
	return nil
}
func (f *MBSSessionIdentifierFields) valid() bool {
	return f.Flags&^7 == 0 && f.NID <= 0xfffffffffff && (f.Flags&MBSSessionNIDI != 0 || f.NID == 0) && (f.Flags&MBSSessionTMGI != 0 || f.TMGI == [6]byte{}) && ((f.Flags&MBSSessionSSMI != 0 && addedSSMValid(f.MulticastAddress, f.SourceAddress)) || (f.Flags&MBSSessionSSMI == 0 && !f.MulticastAddress.IsValid() && !f.SourceAddress.IsValid()))
}

// Marshal encodes the known fields with spare bits set to zero.
func (f *MBSSessionIdentifierFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	if !f.valid() {
		return nil, ErrInvalidValue
	}
	b := []byte{byte(f.Flags)}
	if f.Flags&MBSSessionTMGI != 0 {
		b = append(b, f.TMGI[:]...)
	}
	if f.Flags&MBSSessionSSMI != 0 {
		b = addedAppendSSM(b, f.MulticastAddress, f.SourceAddress)
	}
	if f.Flags&MBSSessionNIDI != 0 {
		v := f.NID << 4
		b = append(b, byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
	if len(b) > 65535 {
		return nil, ErrInvalidValue
	}
	return b, nil
}

// MarshalLen returns the encoded value length, or zero for nil or invalid fields.
func (f *MBSSessionIdentifierFields) MarshalLen() int { b, _ := f.Marshal(); return len(b) }

// MarshalTo writes the encoded fields to b, rejecting an undersized buffer.
func (f *MBSSessionIdentifierFields) MarshalTo(b []byte) error {
	wire, err := f.Marshal()
	return addedCopyTo(b, wire, err)
}
