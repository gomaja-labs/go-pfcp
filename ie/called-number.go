package ie

import "unicode/utf8"

// NewCalledNumber creates a Called Number IE. An empty value creates a null IE.
// TS 29.244 V20.1.0 clauses 8.2.191 and 8.1.2.
func NewCalledNumber(value string) *IE {
	if len(value) > 65535 || !utf8.ValidString(value) {
		return newInvalidIE(CalledNumber, ErrInvalidValue)
	}
	return New(CalledNumber, append([]byte(nil), []byte(value)...))
}

// CalledNumber returns the Called Number value.
// TS 29.244 V20.1.0 clause 8.2.191.
func (i *IE) CalledNumber() (string, error) {
	b, err := addedValuePayload(i, CalledNumber, 1, false)
	if err != nil {
		return "", err
	}
	if !utf8.Valid(b) {
		return "", &DecodeError{Type: i.Type, Err: ErrInvalidValue}
	}
	return string(b), nil
}
