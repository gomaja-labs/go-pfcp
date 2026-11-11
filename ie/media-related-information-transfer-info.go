package ie

// MediaRelatedInformationTransferInfo values are defined in TS 29.244 V20.1.0 clause 8.2.271.
const (
	MediaRelatedInformationTransferInfoUnspecified uint8 = 0
	MediaRelatedInformationTransferInfoUDPOption   uint8 = 1
)

// NewMediaRelatedInformationTransferInfo creates a Media Related Information Transfer Info IE.
// TS 29.244 V20.1.0 clause 8.2.271.
func NewMediaRelatedInformationTransferInfo(v uint8) *IE {
	if v > 1 {
		return invalidRecentIE(MediaRelatedInformationTransferInfo, "Media Related Information Transfer Info")
	}
	return New(MediaRelatedInformationTransferInfo, []byte{v})
}

// MediaRelatedInformationTransferInfo returns the Media Related Information Transfer Info value.
// TS 29.244 V20.1.0 clause 8.2.271; clause 8.1.2 governs spare bits and extensions.
func (i *IE) MediaRelatedInformationTransferInfo() (uint8, error) {
	b, err := recentIEPayload(i, MediaRelatedInformationTransferInfo, 1, false)
	if err != nil {
		return 0, err
	}
	v := b[0] & 15

	if v > 1 {
		return 0, &DecodeError{Type: MediaRelatedInformationTransferInfo, Err: ErrInvalidValue}
	}
	return v, nil
}
