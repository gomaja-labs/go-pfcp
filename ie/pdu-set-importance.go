package ie

// PDU Set Importance bitmap values, TS 29.244 V20.1.0 clause 8.2.268.
const (
	// PDUSetImportancePSI0 selects PSI0.
	PDUSetImportancePSI0 uint16 = 1 << iota
	// PDUSetImportancePSI1 selects PSI1.
	PDUSetImportancePSI1
	// PDUSetImportancePSI2 selects PSI2.
	PDUSetImportancePSI2
	// PDUSetImportancePSI3 selects PSI3.
	PDUSetImportancePSI3
	// PDUSetImportancePSI4 selects PSI4.
	PDUSetImportancePSI4
	// PDUSetImportancePSI5 selects PSI5.
	PDUSetImportancePSI5
	// PDUSetImportancePSI6 selects PSI6.
	PDUSetImportancePSI6
	// PDUSetImportancePSI7 selects PSI7.
	PDUSetImportancePSI7
	// PDUSetImportancePSI8 selects PSI8.
	PDUSetImportancePSI8
	// PDUSetImportancePSI9 selects PSI9.
	PDUSetImportancePSI9
	// PDUSetImportancePSI10 selects PSI10.
	PDUSetImportancePSI10
	// PDUSetImportancePSI11 selects PSI11.
	PDUSetImportancePSI11
	// PDUSetImportancePSI12 selects PSI12.
	PDUSetImportancePSI12
	// PDUSetImportancePSI13 selects PSI13.
	PDUSetImportancePSI13
	// PDUSetImportancePSI14 selects PSI14.
	PDUSetImportancePSI14
	// PDUSetImportancePSI15 selects PSI15.
	PDUSetImportancePSI15
)

// NewPDUSetImportance creates a bitmap of PSI values. At least one is required.
// TS 29.244 V20.1.0 clause 8.2.268: PSI0 is bit 1 of the first octet,
// PSI8 bit 1 of the second octet.
func NewPDUSetImportance(bitmap uint16) *IE {
	if bitmap == 0 {
		return invalidRecentIE(PDUSetImportance, "empty PSI bitmap")
	}
	return New(PDUSetImportance, []byte{byte(bitmap), byte(bitmap >> 8)})
}

// PDUSetImportance returns the PSI bitmap, with bit n representing PSI n.
// TS 29.244 V20.1.0 clause 8.2.268.
func (i *IE) PDUSetImportance() (uint16, error) {
	b, err := recentIEPayload(i, PDUSetImportance, 2, true)
	if err != nil {
		return 0, err
	}
	v := uint16(b[0]) | uint16(b[1])<<8
	if v == 0 {
		return 0, &DecodeError{Type: PDUSetImportance, Err: ErrInvalidValue}
	}
	return v, nil
}
