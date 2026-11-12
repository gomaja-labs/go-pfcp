package ie

// NewTransportLevelMarking creates the two-octet ToS/Traffic Class marking.
// TS 29.244 V20.1.0 clause 8.2.12; the first octet is the value, the second its mask.
func NewTransportLevelMarking(tos []byte) *IE {
	if len(tos) != 0 && len(tos) != 2 {
		return newInvalidIE(TransportLevelMarking, ErrInvalidValue)
	}
	return New(TransportLevelMarking, tos)
}

// TransportLevelMarking returns the two marking octets, ignoring extensions.
// TS 29.244 V20.1.0 clauses 8.2.12 and 8.1.2.
func (i *IE) TransportLevelMarking() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}
	if i.Type != TransportLevelMarking {
		return nil, &InvalidTypeError{Type: i.Type}
	}
	b, err := i.payloadAtLeast(2)
	if err != nil {
		return nil, err
	}
	return b[:2], nil
}
