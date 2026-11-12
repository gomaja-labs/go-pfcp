package ie

// VendorSpecificNodeReportTypeFields is the structured value of Vendor-Specific Node Report Type.
// TS 29.244 V20.1.0 clause 8.2.217.
type VendorSpecificNodeReportTypeFields struct {
	// EnterpriseID is the vendor enterprise identifier.
	EnterpriseID uint16
	// ReportTypes is the report types.
	ReportTypes []byte
}

// NewVendorSpecificNodeReportType creates a Vendor-Specific Node Report Type IE from its fields.
// TS 29.244 V20.1.0 clause 8.2.217.
func NewVendorSpecificNodeReportType(fields VendorSpecificNodeReportTypeFields) *IE {
	b, err := fields.Marshal()
	return addedFieldsIE(VendorSpecificNodeReportType, b, err)
}

// VendorSpecificNodeReportType returns the structured Vendor-Specific Node Report Type value.
// TS 29.244 V20.1.0 clauses 8.2.217 and 8.1.2.
func (i *IE) VendorSpecificNodeReportType() (*VendorSpecificNodeReportTypeFields, error) {
	b, err := addedValuePayload(i, VendorSpecificNodeReportType, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseVendorSpecificNodeReportTypeFields(b)
	if err != nil {
		return nil, &DecodeError{Type: VendorSpecificNodeReportType, Err: err}
	}
	return f, nil
}

// ParseVendorSpecificNodeReportTypeFields decodes the enterprise ID and the full vendor bitmask.
// TS 29.244 V20.1.0 clauses 8.2.217 and 8.1.2.
func ParseVendorSpecificNodeReportTypeFields(b []byte) (*VendorSpecificNodeReportTypeFields, error) {
	f := &VendorSpecificNodeReportTypeFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary decodes Vendor-Specific Node Report Type and resets previously decoded fields.
func (f *VendorSpecificNodeReportTypeFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = VendorSpecificNodeReportTypeFields{}
	r := &addedReader{b: b}
	f.EnterpriseID = r.u16()
	if r.err == nil {
		f.ReportTypes = append([]byte(nil), r.take(len(b)-r.pos)...)
	}
	if r.err != nil {
		return r.err
	}

	if !f.valid() {
		return ErrInvalidValue
	}
	return nil
}
func (f *VendorSpecificNodeReportTypeFields) valid() bool {
	return len(f.ReportTypes) >= 1 && len(f.ReportTypes) <= 65533 && addedAnyBit(f.ReportTypes)
}

// Marshal encodes the known fields with spare bits set to zero.
func (f *VendorSpecificNodeReportTypeFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	if !f.valid() {
		return nil, ErrInvalidValue
	}
	b := addedPut16(nil, f.EnterpriseID)
	b = append(b, f.ReportTypes...)
	if len(b) > 65535 {
		return nil, ErrInvalidValue
	}
	return b, nil
}

// MarshalLen returns the encoded value length, or zero for nil or invalid fields.
func (f *VendorSpecificNodeReportTypeFields) MarshalLen() int { b, _ := f.Marshal(); return len(b) }

// MarshalTo writes the encoded fields to b, rejecting an undersized buffer.
func (f *VendorSpecificNodeReportTypeFields) MarshalTo(b []byte) error {
	wire, err := f.Marshal()
	return addedCopyTo(b, wire, err)
}
func addedAnyBit(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return true
		}
	}
	return false
}
