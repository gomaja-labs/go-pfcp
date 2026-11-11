package ie

// TransportLevelMarkingIndications values are defined in TS 29.244 V20.1.0 clause 8.2.280.
const (
	TransportLevelMarkingIndicationsTLMIP uint8 = 1
)

// NewTransportLevelMarkingIndications creates a Transport Level Marking Indications IE.
// TS 29.244 V20.1.0 clause 8.2.280.
func NewTransportLevelMarkingIndications(v uint8) *IE {
	if v > 1 {
		return invalidRecentIE(TransportLevelMarkingIndications, "Transport Level Marking Indications")
	}
	return New(TransportLevelMarkingIndications, []byte{v})
}

// TransportLevelMarkingIndications returns the Transport Level Marking Indications value.
// TS 29.244 V20.1.0 clause 8.2.280; clause 8.1.2 governs spare bits and extensions.
func (i *IE) TransportLevelMarkingIndications() (uint8, error) {
	b, err := recentIEPayload(i, TransportLevelMarkingIndications, 1, true)
	if err != nil {
		return 0, err
	}
	v := b[0] & 1

	if v > 1 {
		return 0, &DecodeError{Type: TransportLevelMarkingIndications, Err: ErrInvalidValue}
	}
	return v, nil
}
