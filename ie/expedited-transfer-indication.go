package ie

// ExpeditedTransferIndication values are defined in TS 29.244 V20.1.0 clause 8.2.278.
const (
	// ExpeditedTransferIndicationTrue selects True.
	ExpeditedTransferIndicationTrue uint8 = 1
	// ExpeditedTransferIndicationFalse selects False.
	ExpeditedTransferIndicationFalse uint8 = 2
)

// NewExpeditedTransferIndication creates a Expedited Transfer Indication IE.
// TS 29.244 V20.1.0 clause 8.2.278.
func NewExpeditedTransferIndication(v uint8) *IE {
	if v != 1 && v != 2 {
		return invalidRecentIE(ExpeditedTransferIndication, "Expedited Transfer Indication")
	}
	return New(ExpeditedTransferIndication, []byte{v})
}

// ExpeditedTransferIndication returns the Expedited Transfer Indication value.
// TS 29.244 V20.1.0 clause 8.2.278; clause 8.1.2 governs spare bits and extensions.
func (i *IE) ExpeditedTransferIndication() (uint8, error) {
	b, err := recentIEPayload(i, ExpeditedTransferIndication, 1, true)
	if err != nil {
		return 0, err
	}
	v := b[0] & 3

	if v != 1 && v != 2 {
		return 0, &DecodeError{Type: ExpeditedTransferIndication, Err: ErrInvalidValue}
	}
	return v, nil
}
