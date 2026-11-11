package ie

import "unicode/utf8"

// NewTunnelPassword creates a Tunnel Password IE. An empty value creates a null IE.
// TS 29.244 V20.1.0 clauses 8.2.212 and 8.1.2.
func NewTunnelPassword(value string) *IE {
	if len(value) > 65535 || !utf8.ValidString(value) {
		return newInvalidIE(TunnelPassword, ErrInvalidValue)
	}
	return New(TunnelPassword, append([]byte(nil), []byte(value)...))
}

// TunnelPassword returns the Tunnel Password value.
// TS 29.244 V20.1.0 clause 8.2.212.
func (i *IE) TunnelPassword() (string, error) {
	b, err := addedValuePayload(i, TunnelPassword, 1, false)
	if err != nil {
		return "", err
	}
	if !utf8.Valid(b) {
		return "", &DecodeError{Type: i.Type, Err: ErrInvalidValue}
	}
	return string(b), nil
}
