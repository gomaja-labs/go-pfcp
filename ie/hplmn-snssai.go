package ie

// HPLMNSNSSAIFields is the structured value of HPLMN S-NSSAI.
// TS 29.244 V20.1.0 clause 8.2.232.
type HPLMNSNSSAIFields struct {
	// SST is the slice/service type.
	SST uint8
	// SD contains the slice differentiator octets.
	SD uint32
}

// NewHPLMNSNSSAI creates a HPLMN S-NSSAI IE from its fields.
// TS 29.244 V20.1.0 clause 8.2.232.
func NewHPLMNSNSSAI(fields HPLMNSNSSAIFields) *IE {
	b, err := fields.Marshal()
	return addedFieldsIE(HPLMNSNSSAI, b, err)
}

// HPLMNSNSSAI returns the structured HPLMN S-NSSAI value.
// TS 29.244 V20.1.0 clauses 8.2.232 and 8.1.2.
func (i *IE) HPLMNSNSSAI() (*HPLMNSNSSAIFields, error) {
	b, err := addedValuePayload(i, HPLMNSNSSAI, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseHPLMNSNSSAIFields(b)
	if err != nil {
		return nil, &DecodeError{Type: HPLMNSNSSAI, Err: err}
	}
	return f, nil
}

// ParseHPLMNSNSSAIFields decodes exactly four value octets.
// TS 29.244 V20.1.0 clauses 8.2.232 and 8.1.2.
// Decoding may borrow b; keep b unchanged while the result is in use.
func ParseHPLMNSNSSAIFields(b []byte) (*HPLMNSNSSAIFields, error) {
	f := &HPLMNSNSSAIFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary decodes HPLMN S-NSSAI and resets previously decoded fields.
// Decoding may borrow b; keep b unchanged while the result is in use.
func (f *HPLMNSNSSAIFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = HPLMNSNSSAIFields{}
	r := &addedReader{b: b}
	f.SST = r.u8()
	v := r.take(3)
	if r.err == nil {
		f.SD = uint32(v[0])<<16 | uint32(v[1])<<8 | uint32(v[2])
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
func (f *HPLMNSNSSAIFields) valid() bool { return f.SD <= 0xffffff }

// Marshal encodes the known fields with spare bits set to zero.
func (f *HPLMNSNSSAIFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	if !f.valid() {
		return nil, ErrInvalidValue
	}
	b := []byte{f.SST, byte(f.SD >> 16), byte(f.SD >> 8), byte(f.SD)}
	if len(b) > 65535 {
		return nil, ErrInvalidValue
	}
	return b, nil
}

// MarshalLen returns the encoded value length, or zero for nil or invalid fields.
func (f *HPLMNSNSSAIFields) MarshalLen() int { b, _ := f.Marshal(); return len(b) }

// MarshalTo writes the encoded fields to b, rejecting an undersized buffer.
func (f *HPLMNSNSSAIFields) MarshalTo(b []byte) error {
	wire, err := f.Marshal()
	return addedCopyTo(b, wire, err)
}
