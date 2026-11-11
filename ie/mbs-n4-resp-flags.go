package ie

// MBSN4ResponseFlags encodes MBSN4Resp-Flags; TS 29.244 V20.1.0 clause 8.2.211.
type MBSN4ResponseFlags uint8

const (
	// MBSN4RespFlagsNN19DT is the NN19DT value defined in clause 8.2.211.
	MBSN4RespFlagsNN19DT MBSN4ResponseFlags = 1
	// MBSN4RespFlagsJMTI is the JMTI value defined in clause 8.2.211.
	MBSN4RespFlagsJMTI MBSN4ResponseFlags = 2
	// MBSN4RespFlagsN19DTR is the N19DTR value defined in clause 8.2.211.
	MBSN4RespFlagsN19DTR MBSN4ResponseFlags = 4
)

// NewMBSN4RespFlags creates a MBSN4Resp-Flags IE.
// TS 29.244 V20.1.0 clause 8.2.211.
func NewMBSN4RespFlags(value MBSN4ResponseFlags) *IE {
	if value > 7 {
		return newInvalidIE(MBSN4RespFlags, ErrInvalidValue)
	}
	return New(MBSN4RespFlags, []byte{byte(value)})
}

// MBSN4RespFlags returns the MBSN4Resp-Flags value, ignoring spare bits.
// TS 29.244 V20.1.0 clauses 8.2.211 and 8.1.2.
func (i *IE) MBSN4RespFlags() (MBSN4ResponseFlags, error) {
	b, err := addedValuePayload(i, MBSN4RespFlags, 1, false)
	if err != nil {
		return 0, err
	}
	return MBSN4ResponseFlags(b[0] & 7), nil
}
