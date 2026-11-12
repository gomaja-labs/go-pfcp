package ie

// NewNotificationCorrelationID creates a Notification Correlation ID IE. An empty value creates a null IE.
// The value is an octet string carried unchanged; no character set is imposed or checked.
// TS 29.244 V20.1.0 clauses 8.2.203 and 8.1.2.
func NewNotificationCorrelationID(value string) *IE {
	if len(value) > 65535 {
		return newInvalidIE(NotificationCorrelationID, ErrInvalidValue)
	}
	return New(NotificationCorrelationID, []byte(value))
}

// NotificationCorrelationID returns the Notification Correlation ID value.
// TS 29.244 V20.1.0 clause 8.2.203.
func (i *IE) NotificationCorrelationID() (string, error) {
	b, err := addedValuePayload(i, NotificationCorrelationID, 1, false)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
