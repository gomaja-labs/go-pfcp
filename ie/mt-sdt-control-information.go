package ie

// MTSDTControlInformation values are defined in TS 29.244 V20.1.0 clause 8.2.239.
const (
	MTSDTControlInformationRDSI uint8 = 1
)

// NewMTSDTControlInformation creates a MT-SDT Control Information IE.
// TS 29.244 V20.1.0 clause 8.2.239.
func NewMTSDTControlInformation(v uint8) *IE {
	if v > 1 {
		return invalidRecentIE(MTSDTControlInformation, "MT-SDT Control Information")
	}
	return New(MTSDTControlInformation, []byte{v})
}

// MTSDTControlInformation returns the MT-SDT Control Information value.
// TS 29.244 V20.1.0 clause 8.2.239; clause 8.1.2 governs spare bits and extensions.
func (i *IE) MTSDTControlInformation() (uint8, error) {
	b, err := recentIEPayload(i, MTSDTControlInformation, 1, false)
	if err != nil {
		return 0, err
	}
	v := b[0] & 1

	if v > 1 {
		return 0, &DecodeError{Type: MTSDTControlInformation, Err: ErrInvalidValue}
	}
	return v, nil
}
