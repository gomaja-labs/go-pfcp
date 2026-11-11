package ie

// MPQUICControlFlags encodes MPQUIC Control Information; TS 29.244 V20.1.0 clause 8.2.225.
type MPQUICControlFlags uint8

const (
	// MPQUICControlInformationCUDP is the CUDP value defined in clause 8.2.225.
	MPQUICControlInformationCUDP MPQUICControlFlags = 1
	// MPQUICControlInformationCIP is the CIP value defined in clause 8.2.225.
	MPQUICControlInformationCIP MPQUICControlFlags = 2
	// MPQUICControlInformationCETH is the CETH value defined in clause 8.2.225.
	MPQUICControlInformationCETH MPQUICControlFlags = 4
)

// NewMPQUICControlInformation creates a MPQUIC Control Information IE.
// TS 29.244 V20.1.0 clause 8.2.225.
func NewMPQUICControlInformation(value MPQUICControlFlags) *IE {
	if value < 1 || value > 7 {
		return newInvalidIE(MPQUICControlInformation, ErrInvalidValue)
	}
	return New(MPQUICControlInformation, []byte{byte(value)})
}

// MPQUICControlInformation returns the MPQUIC Control Information value, ignoring spare bits.
// TS 29.244 V20.1.0 clauses 8.2.225 and 8.1.2.
func (i *IE) MPQUICControlInformation() (MPQUICControlFlags, error) {
	b, err := addedValuePayload(i, MPQUICControlInformation, 1, false)
	if err != nil {
		return 0, err
	}
	v := b[0] & 7
	if v == 0 {
		return 0, &DecodeError{Type: i.Type, Err: ErrInvalidValue}
	}
	return MPQUICControlFlags(v), nil
}
