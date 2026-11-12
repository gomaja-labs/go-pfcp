package ie

// HeaderHandlingAction values are defined in TS 29.244 V20.1.0 clause 8.2.255.
const (
	// HeaderHandlingActionDetect selects Detect.
	HeaderHandlingActionDetect uint8 = 0
	// HeaderHandlingActionRemove selects Remove.
	HeaderHandlingActionRemove uint8 = 1
	// HeaderHandlingActionReplace selects Replace.
	HeaderHandlingActionReplace uint8 = 2
	// HeaderHandlingActionInsert selects Insert.
	HeaderHandlingActionInsert uint8 = 3
)

// NewHeaderHandlingAction creates a Header Handling Action IE.
// TS 29.244 V20.1.0 clause 8.2.255.
func NewHeaderHandlingAction(v uint8) *IE {
	if v > 3 {
		return invalidRecentIE(HeaderHandlingAction, "Header Handling Action")
	}
	return New(HeaderHandlingAction, []byte{v})
}

// HeaderHandlingAction returns the Header Handling Action value.
// TS 29.244 V20.1.0 clause 8.2.255; clause 8.1.2 governs spare bits and extensions.
func (i *IE) HeaderHandlingAction() (uint8, error) {
	b, err := recentIEPayload(i, HeaderHandlingAction, 1, true)
	if err != nil {
		return 0, err
	}
	v := b[0] & 15

	if v > 3 {
		return 0, &DecodeError{Type: HeaderHandlingAction, Err: ErrInvalidValue}
	}
	return v, nil
}
