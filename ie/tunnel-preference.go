package ie

// NewTunnelPreference creates a Tunnel Preference IE.
// TS 29.244 V20.1.0 clause 8.2.189.
func NewTunnelPreference(value uint32) *IE {
	if value > 16777215 {
		return newInvalidIE(TunnelPreference, ErrInvalidValue)
	}
	return New(TunnelPreference, []byte{byte(value >> 16), byte(value >> 8), byte(value)})
}

// TunnelPreference returns the Tunnel Preference value.
// TS 29.244 V20.1.0 clauses 8.2.189 and 8.1.2.
func (i *IE) TunnelPreference() (uint32, error) {
	b, err := addedValuePayload(i, TunnelPreference, 3, true)
	if err != nil {
		return 0, err
	}
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2]), nil
}
