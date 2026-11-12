package ie

// DSCPToPPIMappingInformationFields is the structured value of DSCPToPPIMappingInformation.
// TS 29.244 V20.1.0 clause 8.2.214.
// PPI occupies three bits in octet 5 and each following octet contains one
// six-bit DSCP. The list contains 1 to 63 values. Extension octets cannot be
// distinguished from DSCP values because no count or terminator is encoded.
type DSCPToPPIMappingInformationFields struct {
	// PPI is the paging policy indicator.
	PPI uint8
	// DSCP contains the differentiated services code points.
	DSCP []uint8
}

// NewDSCPToPPIMappingInformation creates an IE, retaining invalid arguments as an encoding error.
// TS 29.244 V20.1.0 clause 8.2.214.
func NewDSCPToPPIMappingInformation(fields DSCPToPPIMappingInformationFields) *IE {

	b, err := fields.Marshal()
	return addedFieldsIE(DSCPToPPIMappingInformation, b, err)
}

// DSCPToPPIMappingInformation returns the PPI and all DSCP values.
// The payload has no distinguishable extension boundary.
// TS 29.244 V20.1.0 clauses 8.2.214 and 8.1.2.
func (i *IE) DSCPToPPIMappingInformation() (*DSCPToPPIMappingInformationFields, error) {
	b, e := addedValuePayload(i, DSCPToPPIMappingInformation, 0, false)
	if e != nil {
		return nil, e
	}
	f, e := ParseDSCPToPPIMappingInformationFields(b)
	if e != nil {
		return nil, &DecodeError{Type: DSCPToPPIMappingInformation, Err: e}
	}
	return f, nil
}

// ParseDSCPToPPIMappingInformationFields decodes fields according to clause 8.2.214.
// Decoding may borrow b; keep b unchanged while the result is in use.
func ParseDSCPToPPIMappingInformationFields(b []byte) (*DSCPToPPIMappingInformationFields, error) {
	f := &DSCPToPPIMappingInformationFields{}
	if e := f.UnmarshalBinary(b); e != nil {
		return nil, e
	}
	return f, nil
}

// UnmarshalBinary resets the value and decodes the known fields.
// Decoding may borrow b; keep b unchanged while the result is in use.
func (f *DSCPToPPIMappingInformationFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = DSCPToPPIMappingInformationFields{}
	if len(b) < 2 || len(b) > 64 {
		return ErrMalformed
	}
	f.PPI = b[0] & 7
	f.DSCP = make([]uint8, len(b)-1)
	for j, v := range b[1:] {
		f.DSCP[j] = v & 63
	}
	return nil
}
func (f *DSCPToPPIMappingInformationFields) valid() bool {
	return f.PPI <= 7 && len(f.DSCP) >= 1 && len(f.DSCP) <= 63 && f.validDSCP()
}

// Marshal encodes known fields with spare bits set to zero.
func (f *DSCPToPPIMappingInformationFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	if !f.valid() {
		return nil, ErrInvalidValue
	}
	return append([]byte{f.PPI}, f.DSCP...), nil
}

// MarshalLen returns the encoded length, or zero for nil or invalid fields.
func (f *DSCPToPPIMappingInformationFields) MarshalLen() int { b, _ := f.Marshal(); return len(b) }

// MarshalTo writes the fields, rejecting undersized buffers.
func (f *DSCPToPPIMappingInformationFields) MarshalTo(b []byte) error {
	wire, e := f.Marshal()
	return addedCopyTo(b, wire, e)
}

func (f *DSCPToPPIMappingInformationFields) validDSCP() bool {
	for _, v := range f.DSCP {
		if v > 63 {
			return false
		}
	}
	return true
}
