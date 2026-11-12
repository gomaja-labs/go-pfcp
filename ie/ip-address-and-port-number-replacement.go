package ie

import (
	"net/netip"
)

// IPReplacementFlags selects replacement instructions in octet 5.
type IPReplacementFlags uint8

const (
	// IPReplacementDestinationV4 selects Destination V4.
	IPReplacementDestinationV4 IPReplacementFlags = 1
	// IPReplacementDestinationV6 selects Destination V6.
	IPReplacementDestinationV6 IPReplacementFlags = 2
	// IPReplacementDestinationPort selects Destination Port.
	IPReplacementDestinationPort IPReplacementFlags = 4
	// IPReplacementSourceV4 selects Source V4.
	IPReplacementSourceV4 IPReplacementFlags = 8
	// IPReplacementSourceV6 selects Source V6.
	IPReplacementSourceV6 IPReplacementFlags = 16
	// IPReplacementSourcePort selects Source Port.
	IPReplacementSourcePort IPReplacementFlags = 32
	// IPReplacementUseMappedN6 selects Use Mapped N6.
	IPReplacementUseMappedN6 IPReplacementFlags = 64
)

// IPAddressAndPortNumberReplacementFields is the structured value of IPAddressAndPortNumberReplacement.
// TS 29.244 V20.1.0 clause 8.2.200.
// Flags occupy octet 5 only, as the clause's bit definitions specify; the
// figure's "5 to 6" label is an editorial error. Values start at octet 6.
// The IE constructor requires an address unless UMN6RS is set, and forbids
// address and port fields with UMN6RS. The decoder follows the received flags.
type IPAddressAndPortNumberReplacementFields struct {
	// Flags contains the presence and option bits.
	Flags IPReplacementFlags
	// DestinationIPv4 is the destination ipv4.
	DestinationIPv4 netip.Addr
	// DestinationIPv6 is the destination ipv6.
	DestinationIPv6 netip.Addr
	// SourceIPv4 is the source ipv4.
	SourceIPv4 netip.Addr
	// SourceIPv6 is the source ipv6.
	SourceIPv6 netip.Addr
	// DestinationPort is the destination port.
	DestinationPort uint16
	// SourcePort is the source port.
	SourcePort uint16
}

// NewIPAddressAndPortNumberReplacement creates an IE, retaining invalid arguments as an encoding error.
// TS 29.244 V20.1.0 clause 8.2.200.
func NewIPAddressAndPortNumberReplacement(fields IPAddressAndPortNumberReplacementFields) *IE {
	if fields.Flags&IPReplacementUseMappedN6 != 0 {
		if fields.Flags&63 != 0 {
			return newInvalidIE(IPAddressAndPortNumberReplacement, ErrInvalidValue)
		}
	} else if fields.Flags&27 == 0 {
		return newInvalidIE(IPAddressAndPortNumberReplacement, ErrInvalidValue)
	}
	b, err := fields.Marshal()
	return addedFieldsIE(IPAddressAndPortNumberReplacement, b, err)
}

// IPAddressAndPortNumberReplacement returns the known fields, ignoring explicitly allowed extensions.
// TS 29.244 V20.1.0 clauses 8.2.200 and 8.1.2.
func (i *IE) IPAddressAndPortNumberReplacement() (*IPAddressAndPortNumberReplacementFields, error) {
	b, e := addedValuePayload(i, IPAddressAndPortNumberReplacement, 1, false)
	if e != nil {
		return nil, e
	}
	f, e := ParseIPAddressAndPortNumberReplacementFields(b)
	if e != nil {
		return nil, &DecodeError{Type: IPAddressAndPortNumberReplacement, Err: e}
	}
	return f, nil
}

// ParseIPAddressAndPortNumberReplacementFields decodes fields according to clause 8.2.200.
func ParseIPAddressAndPortNumberReplacementFields(b []byte) (*IPAddressAndPortNumberReplacementFields, error) {
	f := &IPAddressAndPortNumberReplacementFields{}
	if e := f.UnmarshalBinary(b); e != nil {
		return nil, e
	}
	return f, nil
}

// UnmarshalBinary resets the value and decodes the known fields.
func (f *IPAddressAndPortNumberReplacementFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = IPAddressAndPortNumberReplacementFields{}
	r := &addedReader{b: b}
	f.Flags = IPReplacementFlags(r.u8() & 127)
	if f.Flags&1 != 0 {
		f.DestinationIPv4 = r.ip(4)
	}
	if f.Flags&2 != 0 {
		f.DestinationIPv6 = r.ip(16)
	}
	if f.Flags&4 != 0 {
		f.DestinationPort = r.u16()
	}
	if f.Flags&8 != 0 {
		f.SourceIPv4 = r.ip(4)
	}
	if f.Flags&16 != 0 {
		f.SourceIPv6 = r.ip(16)
	}
	if f.Flags&32 != 0 {
		f.SourcePort = r.u16()
	}
	return r.err
}
func (f *IPAddressAndPortNumberReplacementFields) valid() bool {
	return f.Flags&^127 == 0 && ((f.Flags&1 != 0 && addedAddressValid(f.DestinationIPv4, 32)) || (f.Flags&1 == 0 && !f.DestinationIPv4.IsValid())) && ((f.Flags&2 != 0 && addedAddressValid(f.DestinationIPv6, 128)) || (f.Flags&2 == 0 && !f.DestinationIPv6.IsValid())) && ((f.Flags&8 != 0 && addedAddressValid(f.SourceIPv4, 32)) || (f.Flags&8 == 0 && !f.SourceIPv4.IsValid())) && ((f.Flags&16 != 0 && addedAddressValid(f.SourceIPv6, 128)) || (f.Flags&16 == 0 && !f.SourceIPv6.IsValid())) && (f.Flags&4 != 0 || f.DestinationPort == 0) && (f.Flags&32 != 0 || f.SourcePort == 0)
}

// Marshal encodes known fields with spare bits set to zero.
func (f *IPAddressAndPortNumberReplacementFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	if !f.valid() {
		return nil, ErrInvalidValue
	}
	b := []byte{byte(f.Flags)}
	if f.Flags&1 != 0 {
		b = append(b, f.DestinationIPv4.AsSlice()...)
	}
	if f.Flags&2 != 0 {
		b = append(b, f.DestinationIPv6.AsSlice()...)
	}
	if f.Flags&4 != 0 {
		b = addedPut16(b, f.DestinationPort)
	}
	if f.Flags&8 != 0 {
		b = append(b, f.SourceIPv4.AsSlice()...)
	}
	if f.Flags&16 != 0 {
		b = append(b, f.SourceIPv6.AsSlice()...)
	}
	if f.Flags&32 != 0 {
		b = addedPut16(b, f.SourcePort)
	}
	return b, nil
}

// MarshalLen returns the encoded length, or zero for nil or invalid fields.
func (f *IPAddressAndPortNumberReplacementFields) MarshalLen() int {
	b, _ := f.Marshal()
	return len(b)
}

// MarshalTo writes the fields, rejecting undersized buffers.
func (f *IPAddressAndPortNumberReplacementFields) MarshalTo(b []byte) error {
	wire, e := f.Marshal()
	return addedCopyTo(b, wire, e)
}
