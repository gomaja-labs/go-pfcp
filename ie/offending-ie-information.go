package ie

import (
	"encoding/binary"
	"io"
	"net/netip"
	"time"
)

// OffendingIEInformationFields is the structured value of Offending IE Information.
// TS 29.244 V20.1.0 clause 8.2.185.
type OffendingIEInformationFields struct {
	// Type identifies the offending IE type.
	Type IEType
	// Value contains the offending IE value octets.
	Value []byte
}

// NewOffendingIEInformation creates a Offending IE Information IE from its fields.
// TS 29.244 V20.1.0 clause 8.2.185.
func NewOffendingIEInformation(fields OffendingIEInformationFields) *IE {
	b, err := fields.Marshal()
	return addedFieldsIE(OffendingIEInformation, b, err)
}

// OffendingIEInformation returns the structured Offending IE Information value.
// TS 29.244 V20.1.0 clauses 8.2.185 and 8.1.2.
func (i *IE) OffendingIEInformation() (*OffendingIEInformationFields, error) {
	b, err := addedValuePayload(i, OffendingIEInformation, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseOffendingIEInformationFields(b)
	if err != nil {
		return nil, &DecodeError{Type: OffendingIEInformation, Err: err}
	}
	return f, nil
}

// ParseOffendingIEInformationFields decodes the type and the remaining offending value.
// TS 29.244 V20.1.0 clauses 8.2.185 and 8.1.2.
func ParseOffendingIEInformationFields(b []byte) (*OffendingIEInformationFields, error) {
	f := &OffendingIEInformationFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary decodes Offending IE Information and resets previously decoded fields.
func (f *OffendingIEInformationFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = OffendingIEInformationFields{}
	r := &addedReader{b: b}
	f.Type = IEType(r.u16())
	if r.err == nil {
		f.Value = append([]byte(nil), r.take(len(b)-r.pos)...)
	}
	if r.err != nil {
		return r.err
	}

	if !f.valid() {
		return ErrInvalidValue
	}
	return nil
}
func (f *OffendingIEInformationFields) valid() bool { return len(f.Value) <= 65533 }

// Marshal encodes the known fields with spare bits set to zero.
func (f *OffendingIEInformationFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	if !f.valid() {
		return nil, ErrInvalidValue
	}
	b := addedPut16(nil, uint16(f.Type))
	b = append(b, f.Value...)
	if len(b) > 65535 {
		return nil, ErrInvalidValue
	}
	return b, nil
}

// MarshalLen returns the encoded value length, or zero for nil or invalid fields.
func (f *OffendingIEInformationFields) MarshalLen() int { b, _ := f.Marshal(); return len(b) }

// MarshalTo writes the encoded fields to b, rejecting an undersized buffer.
func (f *OffendingIEInformationFields) MarshalTo(b []byte) error {
	wire, err := f.Marshal()
	return addedCopyTo(b, wire, err)
}

// addedValuePayload applies the null, fixed-length and extension rules of
// TS 29.244 V20.1.0 clause 8.1.2 without changing the retained wire payload.
func addedValuePayload(i *IE, typ IEType, minimum int, fixed bool) ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}
	if i.Type != typ {
		return nil, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < minimum {
		return nil, &DecodeError{Type: typ, Err: io.ErrUnexpectedEOF}
	}
	if fixed && len(i.Payload) != minimum {
		return nil, &DecodeError{Type: typ, Err: ErrInvalidLength}
	}
	return i.Payload, nil
}

func addedDurationFits(v, unit time.Duration, max uint64) bool {
	return v >= 0 && v%unit == 0 && uint64(v/unit) <= max
}

func addedFieldsIE(typ IEType, b []byte, err error) *IE {
	if err != nil {
		return newInvalidIE(typ, ErrInvalidValue)
	}
	return New(typ, b)
}

// addedReader bounds each conditional field before slicing it.
type addedReader struct {
	b   []byte
	pos int
	err error
}

func (r *addedReader) take(n int) []byte {
	if r.err != nil {
		return nil
	}
	if n < 0 || n > len(r.b)-r.pos {
		r.err = io.ErrUnexpectedEOF
		return nil
	}
	b := r.b[r.pos : r.pos+n : r.pos+n]
	r.pos += n
	return b
}
func (r *addedReader) u8() uint8 {
	b := r.take(1)
	if b == nil {
		return 0
	}
	return b[0]
}
func (r *addedReader) u16() uint16 {
	b := r.take(2)
	if b == nil {
		return 0
	}
	return binary.BigEndian.Uint16(b)
}
func (r *addedReader) u32() uint32 {
	b := r.take(4)
	if b == nil {
		return 0
	}
	return binary.BigEndian.Uint32(b)
}
func (r *addedReader) ip(n int) netip.Addr {
	b := r.take(n)
	if b == nil {
		return netip.Addr{}
	}
	a, _ := netip.AddrFromSlice(b)
	return a
}
func addedPut16(b []byte, v uint16) []byte { return append(b, byte(v>>8), byte(v)) }
func addedPut32(b []byte, v uint32) []byte {
	return append(b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}
func addedCopyTo(dst, src []byte, err error) error {
	if err != nil {
		return err
	}
	if len(dst) < len(src) {
		return io.ErrUnexpectedEOF
	}
	copy(dst, src)
	return nil
}
func addedAddressValid(a netip.Addr, bits int) bool {
	return a.IsValid() && a.Zone() == "" && a.BitLen() == bits
}
