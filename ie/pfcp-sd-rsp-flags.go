package ie

// PFCPSDResponseFlags encodes PFCPSDRsp-Flags; TS 29.244 V20.1.0 clause 8.2.215.
type PFCPSDResponseFlags uint8

const (
	// PFCPSDRspFlagsPURU is the PURU value defined in clause 8.2.215.
	PFCPSDRspFlagsPURU PFCPSDResponseFlags = 1
)

// NewPFCPSDRspFlags creates a PFCPSDRsp-Flags IE.
// TS 29.244 V20.1.0 clause 8.2.215.
func NewPFCPSDRspFlags(value PFCPSDResponseFlags) *IE {
	if value > 1 {
		return newInvalidIE(PFCPSDRspFlags, ErrInvalidValue)
	}
	return New(PFCPSDRspFlags, []byte{byte(value)})
}

// PFCPSDRspFlags returns the PFCPSDRsp-Flags value, ignoring spare bits.
// TS 29.244 V20.1.0 clauses 8.2.215 and 8.1.2.
func (i *IE) PFCPSDRspFlags() (PFCPSDResponseFlags, error) {
	b, err := addedValuePayload(i, PFCPSDRspFlags, 1, false)
	if err != nil {
		return 0, err
	}
	return PFCPSDResponseFlags(b[0] & 1), nil
}
