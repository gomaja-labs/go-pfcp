package ie

// N6DelayMeasurementFailureInformation values are defined in TS 29.244 V20.1.0 clause 8.2.263.
const (
	// N6DelayMeasurementFailureInformationUnspecified selects Unspecified.
	N6DelayMeasurementFailureInformationUnspecified uint8 = 0
	// N6DelayMeasurementFailureInformationConnectionEstablishmentFailure selects Connection Establishment Failure.
	N6DelayMeasurementFailureInformationConnectionEstablishmentFailure uint8 = 1
	// N6DelayMeasurementFailureInformationConnectionLost selects Connection Lost.
	N6DelayMeasurementFailureInformationConnectionLost uint8 = 2
	// N6DelayMeasurementFailureInformationAuthenticationFailure selects Authentication Failure.
	N6DelayMeasurementFailureInformationAuthenticationFailure uint8 = 3
)

// NewN6DelayMeasurementFailureInformation creates a N6 Delay Measurement Failure Information IE.
// TS 29.244 V20.1.0 clause 8.2.263.
func NewN6DelayMeasurementFailureInformation(v uint8) *IE {
	if v > 3 {
		return invalidRecentIE(N6DelayMeasurementFailureInformation, "N6 Delay Measurement Failure Information")
	}
	return New(N6DelayMeasurementFailureInformation, []byte{v})
}

// N6DelayMeasurementFailureInformation returns the N6 Delay Measurement Failure Information value.
// TS 29.244 V20.1.0 clause 8.2.263; clause 8.1.2 governs spare bits and extensions.
func (i *IE) N6DelayMeasurementFailureInformation() (uint8, error) {
	b, err := recentIEPayload(i, N6DelayMeasurementFailureInformation, 1, false)
	if err != nil {
		return 0, err
	}
	v := b[0] & 255
	if v > 3 {
		v = 0
	}
	return v, nil
}
