package ie

// TransportModeValue encodes Transport Mode; TS 29.244 V20.1.0 clause 8.2.227.
type TransportModeValue uint8

const (
	// TransportModeDatagram1 is the Datagram1 value defined in clause 8.2.227.
	TransportModeDatagram1 TransportModeValue = 0
	// TransportModeDatagram2 is the Datagram2 value defined in clause 8.2.227.
	TransportModeDatagram2 TransportModeValue = 1
	// TransportModeStream is the Stream value defined in clause 8.2.227.
	TransportModeStream TransportModeValue = 2
)

// NewTransportMode creates a Transport Mode IE.
// TS 29.244 V20.1.0 clause 8.2.227.
func NewTransportMode(value TransportModeValue) *IE {
	if value > 2 {
		return newInvalidIE(TransportMode, ErrInvalidValue)
	}
	return New(TransportMode, []byte{byte(value)})
}

// TransportMode returns the Transport Mode value, ignoring spare bits.
// TS 29.244 V20.1.0 clauses 8.2.227 and 8.1.2.
func (i *IE) TransportMode() (TransportModeValue, error) {
	b, err := addedValuePayload(i, TransportMode, 1, false)
	if err != nil {
		return 0, err
	}
	return TransportModeValue(b[0] & 15), nil
}
