package ie

// SessionReflectorModeSTAMP values are defined in TS 29.244 V20.1.0 clause 8.2.279.
const (
	// SessionReflectorModeSTAMPStateless selects Stateless.
	SessionReflectorModeSTAMPStateless uint8 = 0
	// SessionReflectorModeSTAMPStateful selects Stateful.
	SessionReflectorModeSTAMPStateful uint8 = 1
)

// NewSessionReflectorModeSTAMP creates a Session Reflector Mode (STAMP) IE.
// TS 29.244 V20.1.0 clause 8.2.279.
func NewSessionReflectorModeSTAMP(v uint8) *IE {
	if v > 1 {
		return invalidRecentIE(SessionReflectorModeSTAMP, "Session Reflector Mode (STAMP)")
	}
	return New(SessionReflectorModeSTAMP, []byte{v})
}

// SessionReflectorModeSTAMP returns the Session Reflector Mode (STAMP) value.
// TS 29.244 V20.1.0 clause 8.2.279; clause 8.1.2 governs spare bits and extensions.
func (i *IE) SessionReflectorModeSTAMP() (uint8, error) {
	b, err := recentIEPayload(i, SessionReflectorModeSTAMP, 1, true)
	if err != nil {
		return 0, err
	}
	v := b[0] & 1

	if v > 1 {
		return 0, &DecodeError{Type: SessionReflectorModeSTAMP, Err: ErrInvalidValue}
	}
	return v, nil
}
