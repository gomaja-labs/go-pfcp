package ie

// ConfiguredTimeDomainFlags encodes Configured Time Domain; TS 29.244 V20.1.0 clause 8.2.218.
type ConfiguredTimeDomainFlags uint8

const (
	// ConfiguredTimeDomainCTDI is the CTDI value defined in clause 8.2.218.
	ConfiguredTimeDomainCTDI ConfiguredTimeDomainFlags = 1
)

// NewConfiguredTimeDomain creates a Configured Time Domain IE.
// TS 29.244 V20.1.0 clause 8.2.218.
func NewConfiguredTimeDomain(value ConfiguredTimeDomainFlags) *IE {
	if value > 1 {
		return newInvalidIE(ConfiguredTimeDomain, ErrInvalidValue)
	}
	return New(ConfiguredTimeDomain, []byte{byte(value)})
}

// ConfiguredTimeDomain returns the Configured Time Domain value, ignoring spare bits.
// TS 29.244 V20.1.0 clauses 8.2.218 and 8.1.2.
func (i *IE) ConfiguredTimeDomain() (ConfiguredTimeDomainFlags, error) {
	b, err := addedValuePayload(i, ConfiguredTimeDomain, 1, false)
	if err != nil {
		return 0, err
	}
	return ConfiguredTimeDomainFlags(b[0] & 1), nil
}
