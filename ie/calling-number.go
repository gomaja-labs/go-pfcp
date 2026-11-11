package ie

import "unicode/utf8"

// NewCallingNumber creates a Calling Number IE. An empty value creates a null IE.
// TS 29.244 V20.1.0 clauses 8.2.190 and 8.1.2.
func NewCallingNumber(value string) *IE {
	if len(value) > 65535 || !utf8.ValidString(value) {
		return newInvalidIE(CallingNumber, ErrInvalidValue)
	}
	return New(CallingNumber, append([]byte(nil), []byte(value)...))
}

// CallingNumber returns the Calling Number value.
// TS 29.244 V20.1.0 clause 8.2.190.
func (i *IE) CallingNumber() (string, error) {
	b, err := addedValuePayload(i, CallingNumber, 1, false)
	if err != nil {
		return "", err
	}
	if !utf8.Valid(b) {
		return "", &DecodeError{Type: i.Type, Err: ErrInvalidValue}
	}
	return string(b), nil
}
