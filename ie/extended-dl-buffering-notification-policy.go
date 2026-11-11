package ie

// ExtendedDLBufferingNotificationPolicy values are defined in TS 29.244 V20.1.0 clause 8.2.238.
const (
	ExtendedDLBufferingNotificationPolicyEDBN uint8 = 1
)

// NewExtendedDLBufferingNotificationPolicy creates a Extended DL Buffering Notification Policy IE.
// TS 29.244 V20.1.0 clause 8.2.238.
func NewExtendedDLBufferingNotificationPolicy(v uint8) *IE {
	if v > 1 {
		return invalidRecentIE(ExtendedDLBufferingNotificationPolicy, "Extended DL Buffering Notification Policy")
	}
	return New(ExtendedDLBufferingNotificationPolicy, []byte{v})
}

// ExtendedDLBufferingNotificationPolicy returns the Extended DL Buffering Notification Policy value.
// TS 29.244 V20.1.0 clause 8.2.238; clause 8.1.2 governs spare bits and extensions.
func (i *IE) ExtendedDLBufferingNotificationPolicy() (uint8, error) {
	b, err := recentIEPayload(i, ExtendedDLBufferingNotificationPolicy, 1, false)
	if err != nil {
		return 0, err
	}
	v := b[0] & 1

	if v > 1 {
		return 0, &DecodeError{Type: ExtendedDLBufferingNotificationPolicy, Err: ErrInvalidValue}
	}
	return v, nil
}
