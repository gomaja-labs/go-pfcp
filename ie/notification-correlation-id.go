package ie

// NewNotificationCorrelationID creates a Notification Correlation ID IE. An empty value creates a null IE.
// TS 29.244 V20.1.0 clauses 8.2.203 and 8.1.2.
func NewNotificationCorrelationID(value []byte) *IE {
	if len(value) > 65535 {
		return newInvalidIE(NotificationCorrelationID, ErrInvalidValue)
	}
	return New(NotificationCorrelationID, append([]byte(nil), value...))
}

// NotificationCorrelationID returns the Notification Correlation ID value.
// TS 29.244 V20.1.0 clause 8.2.203.
func (i *IE) NotificationCorrelationID() ([]byte, error) {
	b, err := addedValuePayload(i, NotificationCorrelationID, 1, false)
	if err != nil {
		return nil, err
	}

	return append([]byte(nil), b...), nil
}
