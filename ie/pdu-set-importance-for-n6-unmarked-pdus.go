package ie

// NewPDUSetImportanceForN6UnmarkedPDUs creates a PDU Set Importance for N6-unmarked PDUs IE.
// TS 29.244 V20.1.0 clause 8.2.282.
func NewPDUSetImportanceForN6UnmarkedPDUs(v uint8) *IE {
	if v == 0 || v > 15 {
		return invalidRecentIE(PDUSetImportanceForN6UnmarkedPDUs, "PDU Set Importance for N6-unmarked PDUs")
	}
	return New(PDUSetImportanceForN6UnmarkedPDUs, []byte{v})
}

// PDUSetImportanceForN6UnmarkedPDUs returns the PDU Set Importance for N6-unmarked PDUs value.
// TS 29.244 V20.1.0 clause 8.2.282; clause 8.1.2 governs spare bits and extensions.
func (i *IE) PDUSetImportanceForN6UnmarkedPDUs() (uint8, error) {
	b, err := recentIEPayload(i, PDUSetImportanceForN6UnmarkedPDUs, 1, true)
	if err != nil {
		return 0, err
	}
	v := b[0] & 15

	if v == 0 || v > 15 {
		return 0, &DecodeError{Type: PDUSetImportanceForN6UnmarkedPDUs, Err: ErrInvalidValue}
	}
	return v, nil
}
