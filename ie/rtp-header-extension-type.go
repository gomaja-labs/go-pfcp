package ie

// RTPHeaderExtensionTypeValue encodes RTP Header Extension Type; TS 29.244 V20.1.0 clause 8.2.234.
type RTPHeaderExtensionTypeValue uint8

const (
	// RTPHeaderExtensionTypePDUSetMarking is the PDUSetMarking value defined in clause 8.2.234.
	RTPHeaderExtensionTypePDUSetMarking RTPHeaderExtensionTypeValue = 1
	// RTPHeaderExtensionTypeDynamicTrafficCharacteristics is the DynamicTrafficCharacteristics value defined in clause 8.2.234.
	RTPHeaderExtensionTypeDynamicTrafficCharacteristics RTPHeaderExtensionTypeValue = 2
	// RTPHeaderExtensionTypeExpeditedTransfer is the ExpeditedTransfer value defined in clause 8.2.234.
	RTPHeaderExtensionTypeExpeditedTransfer RTPHeaderExtensionTypeValue = 3
)

// NewRTPHeaderExtensionType creates a RTP Header Extension Type IE.
// TS 29.244 V20.1.0 clause 8.2.234.
func NewRTPHeaderExtensionType(value RTPHeaderExtensionTypeValue) *IE {
	if value < 1 || value > 3 {
		return newInvalidIE(RTPHeaderExtensionType, ErrInvalidValue)
	}
	return New(RTPHeaderExtensionType, []byte{byte(value)})
}

// RTPHeaderExtensionType returns the RTP Header Extension Type value, ignoring spare bits.
// TS 29.244 V20.1.0 clauses 8.2.234 and 8.1.2.
func (i *IE) RTPHeaderExtensionType() (RTPHeaderExtensionTypeValue, error) {
	b, err := addedValuePayload(i, RTPHeaderExtensionType, 1, true)
	if err != nil {
		return 0, err
	}
	return RTPHeaderExtensionTypeValue(b[0] & 255), nil
}
