package ie

// Security mode and feature bits for TS 29.244 V20.1.0 clause 8.2.275.
// RFC 4656 section 3.1 defines the base modes. TWAMP extensions use RFC 5618
// section 3, RFC 5938 section 3, RFC 6038 section 4.1, RFC 7717 section 7,
// RFC 7750 section 2 and RFC 8186 section 2. OWAMP's IKE bit is different.
const (
	// SecurityModeOWAMPOrTWAMPStop selects Stop.
	SecurityModeOWAMPOrTWAMPStop uint32 = 0
	// SecurityModeOWAMPOrTWAMPUnauthenticated selects Unauthenticated.
	SecurityModeOWAMPOrTWAMPUnauthenticated uint32 = 1
	// SecurityModeOWAMPOrTWAMPAuthenticated selects Authenticated.
	SecurityModeOWAMPOrTWAMPAuthenticated uint32 = 2
	// SecurityModeOWAMPOrTWAMPEncrypted selects Encrypted.
	SecurityModeOWAMPOrTWAMPEncrypted uint32 = 4
	// SecurityModeOWAMPOrTWAMPMixed selects Mixed.
	SecurityModeOWAMPOrTWAMPMixed uint32 = 8
	// SecurityModeOWAMPOrTWAMPOWAMPIKEv2Derived selects OWAMP IKEv2 Derived.
	SecurityModeOWAMPOrTWAMPOWAMPIKEv2Derived uint32 = 16
	// SecurityModeOWAMPOrTWAMPTWAMPIndividualSessionControl selects TWAMP Individual Session Control.
	SecurityModeOWAMPOrTWAMPTWAMPIndividualSessionControl uint32 = 16
	// SecurityModeOWAMPOrTWAMPTWAMPReflectOctets selects TWAMP Reflect Octets.
	SecurityModeOWAMPOrTWAMPTWAMPReflectOctets uint32 = 32
	// SecurityModeOWAMPOrTWAMPTWAMPSymmetricalSize selects TWAMP Symmetrical Size.
	SecurityModeOWAMPOrTWAMPTWAMPSymmetricalSize uint32 = 64
	// SecurityModeOWAMPOrTWAMPTWAMPIKEv2Derived selects TWAMP IKEv2 Derived.
	SecurityModeOWAMPOrTWAMPTWAMPIKEv2Derived uint32 = 128
	// SecurityModeOWAMPOrTWAMPTWAMPDSCPECNMonitoring selects TWAMP DSCP/ECN Monitoring.
	SecurityModeOWAMPOrTWAMPTWAMPDSCPECNMonitoring uint32 = 256
	// SecurityModeOWAMPOrTWAMPTWAMPPTPv2Timestamp selects TWAMP PTPv2 Timestamp.
	SecurityModeOWAMPOrTWAMPTWAMPPTPv2Timestamp uint32 = 512
)

// NewSecurityModeOWAMPOrTWAMP creates the requested Mode field.
// TS 29.244 V20.1.0 clause 8.2.275, RFC 4656 section 3.1 and its mode extensions.
// One base security mode (or zero to stop) is selected. Feature bits must be
// interpreted with the measurement protocol; bit 4 differs between OWAMP and TWAMP.
func NewSecurityModeOWAMPOrTWAMP(v uint32) *IE {
	if !validDelaySecurityMode(v) {
		return invalidRecentIE(SecurityModeOWAMPOrTWAMP, "security mode")
	}
	return New(SecurityModeOWAMPOrTWAMP, []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)})
}

// SecurityModeOWAMPOrTWAMP returns the selected mode and assigned feature bits.
// TS 29.244 V20.1.0 clause 8.2.275; RFC 4656 section 3.1 ignores unassigned bits.
func (i *IE) SecurityModeOWAMPOrTWAMP() (uint32, error) {
	b, err := recentIEPayload(i, SecurityModeOWAMPOrTWAMP, 4, true)
	if err != nil {
		return 0, err
	}
	v := uint32(recentUint(b)) & 0x3ff
	if !validDelaySecurityMode(v) {
		return 0, &DecodeError{Type: SecurityModeOWAMPOrTWAMP, Err: ErrInvalidValue}
	}
	return v, nil
}
func validDelaySecurityMode(v uint32) bool {
	mode := v & 15
	return v&^uint32(0x3ff) == 0 && (v == 0 || mode == 1 || mode == 2 || mode == 4 || mode == 8)
}
