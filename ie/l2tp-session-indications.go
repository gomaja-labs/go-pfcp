package ie

// L2TPSessionIndicationFlags encodes L2TP Session Indications; TS 29.244 V20.1.0 clause 8.2.192.
type L2TPSessionIndicationFlags uint8

const (
	// L2TPSessionIndicationsREUIA is the REUIA value defined in clause 8.2.192.
	L2TPSessionIndicationsREUIA L2TPSessionIndicationFlags = 1
	// L2TPSessionIndicationsREDSA is the REDSA value defined in clause 8.2.192.
	L2TPSessionIndicationsREDSA L2TPSessionIndicationFlags = 2
	// L2TPSessionIndicationsRENSA is the RENSA value defined in clause 8.2.192.
	L2TPSessionIndicationsRENSA L2TPSessionIndicationFlags = 4
)

// NewL2TPSessionIndications creates a L2TP Session Indications IE.
// TS 29.244 V20.1.0 clause 8.2.192.
func NewL2TPSessionIndications(value L2TPSessionIndicationFlags) *IE {
	if value > 7 {
		return newInvalidIE(L2TPSessionIndications, ErrInvalidValue)
	}
	return New(L2TPSessionIndications, []byte{byte(value)})
}

// L2TPSessionIndications returns the L2TP Session Indications value, ignoring spare bits.
// TS 29.244 V20.1.0 clauses 8.2.192 and 8.1.2.
func (i *IE) L2TPSessionIndications() (L2TPSessionIndicationFlags, error) {
	b, err := addedValuePayload(i, L2TPSessionIndications, 1, false)
	if err != nil {
		return 0, err
	}
	return L2TPSessionIndicationFlags(b[0] & 7), nil
}
