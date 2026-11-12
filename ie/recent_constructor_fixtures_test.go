package ie_test

import (
	"net/netip"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// Explicit valid inputs for the constructors in TS 29.244 V20.1.0 clauses 8.2.236–286.
var _ = registerConstructorFixture(ie.NewRTPPayloadType, uint8(96))
var _ = registerConstructorFixture(ie.NewRTPPayloadFormat, uint8(1))
var _ = registerConstructorFixture(ie.NewExtendedDLBufferingNotificationPolicy, uint8(1))
var _ = registerConstructorFixture(ie.NewMTSDTControlInformation, uint8(1))
var _ = registerConstructorFixture(ie.NewN6DelayMeasurementProtocols, uint8(1))
var _ = registerConstructorFixture(ie.NewPacketInspectionFunctionality, uint8(1))
var _ = registerConstructorFixture(ie.NewReportingEndpointID, uint8(1))
var _ = registerConstructorFixture(ie.NewHeaderHandlingAction, uint8(1))
var _ = registerConstructorFixture(ie.NewHeaderHandlingCondition, uint8(1))
var _ = registerConstructorFixture(ie.NewHeaderHandlingControlID, uint8(1))
var _ = registerConstructorFixture(ie.NewHeaderHandlingControlRuleID, uint8(1))
var _ = registerConstructorFixture(ie.NewMeasurementReportingType, uint8(5))
var _ = registerConstructorFixture(ie.NewN6DelayMeasurementFailureInformation, uint8(1))
var _ = registerConstructorFixture(ie.NewN6DelayMeasurementControlInformationID, uint16(1))
var _ = registerConstructorFixture(ie.NewMeasurementEndpointPortNumber, uint16(1))
var _ = registerConstructorFixture(ie.NewHeaderHandlingReportingIndication, uint8(1))
var _ = registerConstructorFixture(ie.NewSMFChangeReason, uint8(1))
var _ = registerConstructorFixture(ie.NewMoQControlInformation, uint8(1))
var _ = registerConstructorFixture(ie.NewMediaRelatedInformationTransferInfo, uint8(1))
var _ = registerConstructorFixture(ie.NewReportingControlInformation, uint8(1))
var _ = registerConstructorFixture(ie.NewSecurityModeSTAMP, uint8(1))
var _ = registerConstructorFixture(ie.NewRemainingDataReportingIndication, uint8(1))
var _ = registerConstructorFixture(ie.NewExpeditedTransferIndication, uint8(1))
var _ = registerConstructorFixture(ie.NewSessionReflectorModeSTAMP, uint8(1))
var _ = registerConstructorFixture(ie.NewTransportLevelMarkingIndications, uint8(1))
var _ = registerConstructorFixture(ie.NewPDUSetImportanceForN6UnmarkedPDUs, uint8(1))
var _ = registerConstructorFixture(ie.NewNodeLevelRuleID, uint8(1))
var _ = registerConstructorFixture(ie.NewAbnormalTrafficPatternID, uint8(1))
var _ = registerConstructorFixture(ie.NewOperatorConfigurableUPFCapability, "capability")
var _ = registerConstructorFixture(ie.NewHeaderDetectionReference, []byte("value"))
var _ = registerConstructorFixture(ie.NewHeaderDetectionSupportInformation, []byte("value"))
var _ = registerConstructorFixture(ie.NewHeaderHandlingControlReference, []byte("value"))
var _ = registerConstructorFixture(ie.NewHeaderInformation, []byte("value"))
var _ = registerConstructorFixture(ie.NewHeaderValue, []byte("value"))
var _ = registerConstructorFixture(ie.NewHMACKeySTAMP, []byte("value"))
var _ = registerConstructorFixture(ie.NewBindingIndication, []byte("value"))
var _ = registerConstructorFixture(ie.NewReportingThresholds, uint8(15), uint16(7500), uint16(9000), uint64(0x0102030405), uint64(0x060708090a))
var _ = registerConstructorFixture(ie.NewReportingThresholdsFields, uint8(15), uint16(7500), uint16(9000), uint64(0x0102030405), uint64(0x060708090a))
var _ = registerConstructorFixture(ie.NewRTPHeaderExtensionAdditionalInformation, uint8(7), uint8(2), true, true)
var _ = registerConstructorFixture(ie.NewRTPHeaderExtensionAdditionalInformationFields, uint8(7), uint8(2), true, true)
var _ = registerConstructorFixture(ie.NewMappedN6IPAddress, uint8(2), netip.MustParseAddr("192.0.2.129"))
var _ = registerConstructorFixture(ie.NewMappedN6IPAddressFields, uint8(2), netip.MustParseAddr("192.0.2.129"))
var _ = registerConstructorFixture(ie.NewN6RoutingInformation, uint8(63), netip.MustParseAddr("192.0.2.1"), netip.MustParseAddr("2001:db8::1"), uint16(12345), netip.MustParseAddr("192.0.2.2"), netip.MustParseAddr("2001:db8::2"), uint16(54321))
var _ = registerConstructorFixture(ie.NewN6RoutingInformationFields, uint8(63), netip.MustParseAddr("192.0.2.1"), netip.MustParseAddr("2001:db8::1"), uint16(12345), netip.MustParseAddr("192.0.2.2"), netip.MustParseAddr("2001:db8::2"), uint16(54321))
var _ = registerConstructorFixture(ie.NewURI, []byte("https://example.test/trace?x=1"))
var _ = registerConstructorFixture(ie.NewUELevelMeasurementsConfiguration, uint8(2), []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}, 0x123456*time.Second)
var _ = registerConstructorFixture(ie.NewUELevelMeasurementsConfigurationFields, uint8(2), []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}, 0x123456*time.Second)
var _ = registerConstructorFixture(ie.NewN6DelayMeasurementInformation, uint8(3), 12*time.Millisecond, 4*time.Millisecond, 30*time.Millisecond)
var _ = registerConstructorFixture(ie.NewN6DelayMeasurementInformationFields, uint8(3), 12*time.Millisecond, 4*time.Millisecond, 30*time.Millisecond)
var _ = registerConstructorFixture(ie.NewMeasurementEndpointAddress, uint8(3), netip.MustParseAddr("192.0.2.3"), netip.MustParseAddr("2001:db8::3"))
var _ = registerConstructorFixture(ie.NewMeasurementEndpointAddressFields, uint8(3), netip.MustParseAddr("192.0.2.3"), netip.MustParseAddr("2001:db8::3"))
var _ = registerConstructorFixture(ie.NewOnPathN6ConnectionInformation, uint8(57), netip.MustParseAddr("192.0.2.4"), netip.MustParseAddr("2001:db8::4"), uint16(443))
var _ = registerConstructorFixture(ie.NewOnPathN6ConnectionInformationFields, uint8(57), netip.MustParseAddr("192.0.2.4"), netip.MustParseAddr("2001:db8::4"), uint16(443))
var _ = registerConstructorFixture(ie.NewPDUSetImportance, uint16(0xa501))
var _ = registerConstructorFixture(ie.NewMoQRelayIPAddress, uint8(3), netip.MustParseAddr("192.0.2.5"), netip.MustParseAddr("2001:db8::5"))
var _ = registerConstructorFixture(ie.NewMoQRelayIPAddressFields, uint8(3), netip.MustParseAddr("192.0.2.5"), netip.MustParseAddr("2001:db8::5"))
var _ = registerConstructorFixture(ie.NewSecurityModeOWAMPOrTWAMP, uint32(2))
var _ = registerConstructorFixture(ie.NewKeyIDAndSharedSecretOWAMPOrTWAMP, "test-key", "synthetic-secret")
var _ = registerConstructorFixture(ie.NewKeyIDAndSharedSecretOWAMPOrTWAMPFields, "test-key", "synthetic-secret")
var _ = registerConstructorFixture(ie.NewSCONEThroughputAdvice, uint8(3), uint32(1234567), uint32(7654321))
var _ = registerConstructorFixture(ie.NewSCONEThroughputAdviceFields, uint8(3), uint32(1234567), uint32(7654321))
var _ = registerConstructorFixture(ie.NewMitigationActions, ie.MitigationActionsIDDN)

// optionalFixture supplies a value only when its presence bit is set.
func optionalFixture[T any](flags, mask uint8, value T) T {
	if flags&mask != 0 {
		return value
	}
	var zero T
	return zero
}

var _ = registerConstructorFlagFixture(ie.NewReportingThresholdsFields, func(flags uint8) []any {
	if flags > 15 || flags == 0 {
		return nil
	}
	return []any{flags, optionalFixture(flags, 1, uint16(7500)), optionalFixture(flags, 2, uint16(9000)), optionalFixture(flags, 4, uint64(0x0102030405)), optionalFixture(flags, 8, uint64(0x060708090a))}
})

var _ = registerConstructorFlagFixture(ie.NewRTPHeaderExtensionAdditionalInformationFields, func(flags uint8) []any {
	if flags > 7 {
		return nil
	}
	return []any{flags, optionalFixture(flags, 1, uint8(2)), flags&2 != 0, flags&4 != 0}
})

var _ = registerConstructorFlagFixture(ie.NewMappedN6IPAddressFields, func(flags uint8) []any {
	if flags > 3 || flags == 3 {
		return nil
	}
	return []any{flags, optionalFixture(flags, 2, netip.MustParseAddr("192.0.2.1"))}
})

var _ = registerConstructorFlagFixture(ie.NewN6RoutingInformationFields, func(flags uint8) []any {
	if flags > 63 {
		return nil
	}
	return []any{flags, optionalFixture(flags, 1, netip.MustParseAddr("192.0.2.1")), optionalFixture(flags, 2, netip.MustParseAddr("2001:db8::1")), optionalFixture(flags, 4, uint16(12345)), optionalFixture(flags, 8, netip.MustParseAddr("192.0.2.2")), optionalFixture(flags, 16, netip.MustParseAddr("2001:db8::2")), optionalFixture(flags, 32, uint16(54321))}
})

var _ = registerConstructorFlagFixture(ie.NewN6DelayMeasurementInformationFields, func(flags uint8) []any {
	if flags > 3 {
		return nil
	}
	return []any{flags, 12 * time.Millisecond, optionalFixture(flags, 1, 4*time.Millisecond), optionalFixture(flags, 2, 30*time.Millisecond)}
})

var _ = registerConstructorFlagFixture(ie.NewMeasurementEndpointAddressFields, func(flags uint8) []any {
	if flags > 3 || flags == 0 {
		return nil
	}
	return []any{flags, optionalFixture(flags, 1, netip.MustParseAddr("192.0.2.1")), optionalFixture(flags, 2, netip.MustParseAddr("2001:db8::1"))}
})

var _ = registerConstructorFlagFixture(ie.NewOnPathN6ConnectionInformationFields, func(flags uint8) []any {
	if flags > 63 || flags&7 > 1 {
		return nil
	}
	return []any{flags, optionalFixture(flags, 8, netip.MustParseAddr("192.0.2.1")), optionalFixture(flags, 16, netip.MustParseAddr("2001:db8::1")), optionalFixture(flags, 32, uint16(443))}
})

var _ = registerConstructorFlagFixture(ie.NewMoQRelayIPAddressFields, func(flags uint8) []any {
	if flags > 3 {
		return nil
	}
	return []any{flags, optionalFixture(flags, 2, netip.MustParseAddr("192.0.2.1")), optionalFixture(flags, 1, netip.MustParseAddr("2001:db8::1"))}
})

var _ = registerConstructorFlagFixture(ie.NewSCONEThroughputAdviceFields, func(flags uint8) []any {
	if flags > 3 {
		return nil
	}
	return []any{flags, optionalFixture(flags, 1, uint32(1234567)), optionalFixture(flags, 2, uint32(7654321))}
})
