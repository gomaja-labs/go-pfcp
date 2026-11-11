package ie

// MBSN4mbRequestFlags encodes MBSN4mbReq-Flags; TS 29.244 V20.1.0 clause 8.2.208.
type MBSN4mbRequestFlags uint8

const (
	// MBSN4mbReqFlagsPLLSSM is the PLLSSM value defined in clause 8.2.208.
	MBSN4mbReqFlagsPLLSSM MBSN4mbRequestFlags = 1
	// MBSN4mbReqFlagsJMBSSM is the JMBSSM value defined in clause 8.2.208.
	MBSN4mbReqFlagsJMBSSM MBSN4mbRequestFlags = 2
	// MBSN4mbReqFlagsMBSRESTI is the MBSRESTI value defined in clause 8.2.208.
	MBSN4mbReqFlagsMBSRESTI MBSN4mbRequestFlags = 4
)

// NewMBSN4mbReqFlags creates a MBSN4mbReq-Flags IE.
// TS 29.244 V20.1.0 clause 8.2.208.
func NewMBSN4mbReqFlags(value MBSN4mbRequestFlags) *IE {
	if value > 7 {
		return newInvalidIE(MBSN4mbReqFlags, ErrInvalidValue)
	}
	return New(MBSN4mbReqFlags, []byte{byte(value)})
}

// MBSN4mbReqFlags returns the MBSN4mbReq-Flags value, ignoring spare bits.
// TS 29.244 V20.1.0 clauses 8.2.208 and 8.1.2.
func (i *IE) MBSN4mbReqFlags() (MBSN4mbRequestFlags, error) {
	b, err := addedValuePayload(i, MBSN4mbReqFlags, 1, false)
	if err != nil {
		return 0, err
	}
	return MBSN4mbRequestFlags(b[0] & 7), nil
}
