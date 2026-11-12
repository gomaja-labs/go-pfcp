package ie_test

import (
	"net"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2 defines the owned leaf value layouts.
// Each expected structure states the selected fields independently of Parse.
func helperExtendedSemanticValues() []helperSemanticCase {
	v4 := net.IP{192, 0, 2, 1}
	v6 := net.ParseIP("2001:db8::1")
	mac := net.HardwareAddr{2, 0, 0, 0, 0, 1}
	return []helperSemanticCase{
		// TS 29.244 V20.1.0 clause 8.2.167.
		{name: "RequestedQoSMonitoring", input: ie.NewRequestedQoSMonitoringFlags(0x3ff), want: uint16(0x3ff), minimum: 2},
		// TS 29.244 V20.1.0 clause 8.2.43.
		{name: "FQCSIDNodeType", input: ie.NewFQCSID("192.0.2.1", 5, 0x1234), want: uint8(5), minimum: 8},
		// TS 29.244 V20.1.0 clause 8.2.27.
		{name: "DownlinkDataServiceInformation", input: ie.NewDownlinkDataServiceInformationFromFields(&ie.DownlinkDataServiceInformationFields{Flags: 7, PagingPolicyIndication: 5, QFI: 9, DLDataPacketsSize: 123}), want: &ie.DownlinkDataServiceInformationFields{Flags: 7, PagingPolicyIndication: 5, QFI: 9, DLDataPacketsSize: 123}, minimum: 5},
		// TS 29.244 V20.1.0 clause 8.2.7.
		{name: "GateStatusULDL", input: ie.NewGateStatus(1, 0), want: []any{uint8(1), uint8(0)}, minimum: 1},
		{name: "ValueAsUint8", input: ie.New(0xffff, []byte{0x12}), want: uint8(0x12), minimum: 1, generic: true},
		{name: "ValueAsUint16", input: ie.New(0xffff, []byte{0x12, 0x34}), want: uint16(0x1234), minimum: 2, generic: true},
		{name: "ValueAsUint32", input: ie.New(0xffff, []byte{0x12, 0x34, 0x56, 0x78}), want: uint32(0x12345678), minimum: 4, generic: true},
		{name: "ValueAsUint64", input: ie.New(0xffff, []byte{1, 2, 3, 4, 5, 6, 7, 8}), want: uint64(0x0102030405060708), minimum: 8, generic: true},
		{name: "ValueAsString", input: ie.New(0xffff, []byte("abc")), want: "abc", minimum: 1, generic: true},
		{name: "ValueAsFQDN", input: ie.New(0xffff, []byte{3, 'a', 'b', 'c'}), want: "abc", minimum: 4, generic: true},
		// TS 29.244 V20.1.0 clause 8.2.129.
		{name: "AlternativeSMFIPAddress", input: ie.NewAlternativeSMFIPAddress(v4, nil), want: &ie.AlternativeSMFIPAddressFields{Flags: 2, IPv4Address: v4}, minimum: 5},
		// TS 29.244 V20.1.0 clause 8.2.135.
		{name: "CPPFCPEntityIPAddress", input: ie.NewCPPFCPEntityIPAddress(v4, nil), want: &ie.CPPFCPEntityIPAddressFields{Flags: 2, IPv4Address: v4}, minimum: 5},
		// TS 29.244 V20.1.0 clause 8.2.43.
		{name: "CSIDs", input: ie.NewFQCSID("192.0.2.1", 5, 0x1234, 0x5678), want: []uint16{0x1234, 0x5678}, minimum: 10},
		// TS 29.244 V20.1.0 clause 8.2.94.
		{name: "CTAG", input: ie.NewCTAG(7, 3, 1, 0xabc), want: &ie.CTAGFields{Flags: 7, PCP: 3, DEIFlag: 1, CVID: 0xabc}, minimum: 3},
		// TS 29.244 V20.1.0 clause 8.2.150.
		{name: "CumulativeRateRatioMeasurement", input: ie.NewCumulativeRateRatioMeasurement(-123), want: int32(-123), minimum: 4},
		// TS 29.244 V20.1.0 clause 8.2.66.
		{name: "DLFlowLevelMarking", input: ie.NewDLFlowLevelMarking(3, []byte{0x12, 0x34}, []byte{0x56, 0x78}), want: &ie.DLFlowLevelMarkingFields{Flags: 3, ToSTrafficClass: []byte{0x12, 0x34}, ServiceClassIdentifier: []byte{0x56, 0x78}}, minimum: 5},
		// TS 29.244 V20.1.0 clause 8.2.27.
		{name: "DownlinkDataServiceInformation", input: ie.NewDownlinkDataServiceInformation(true, true, 5, 9), want: &ie.DownlinkDataServiceInformationFields{Flags: 3, PagingPolicyIndication: 5, QFI: 9}, minimum: 3},
		// TS 29.244 V20.1.0 clause 8.2.49.
		{name: "DroppedDLTrafficThreshold", input: ie.NewDroppedDLTrafficThreshold(true, false, 123, 0), want: &ie.DroppedDLTrafficThresholdFields{Flags: 1, DownlinkPackets: 123}, minimum: 9},
		// TS 29.244 V20.1.0 clause 8.2.74.
		{name: "FARID", input: ie.NewFARID(0x12345678), want: 0x12345678, minimum: 4},
		// TS 29.244 V20.1.0 clause 8.2.43.
		{name: "FQCSID", input: ie.NewFQCSID("192.0.2.1", 5, 0x1234), want: []byte{1, 192, 0, 2, 1, 0x12, 0x34, 5}, minimum: 8},
		// TS 29.244 V20.1.0 clause 8.2.37.
		{name: "FSEID", input: ie.NewFSEID(0x12345678, v4, nil), want: &ie.FSEIDFields{Flags: 2, SEID: 0x12345678, IPv4Address: v4}, minimum: 13},
		// TS 29.244 V20.1.0 clause 8.2.3.
		{name: "FTEID", input: ie.NewFTEID(1, 0x12345678, v4, nil, 0), want: &ie.FTEIDFields{Flags: 1, TEID: 0x12345678, IPv4Address: v4}, minimum: 9},
		// TS 29.244 V20.1.0 clause 8.2.80.
		{name: "FailedRuleID", input: ie.NewFailedRuleID(5, 0x1234), want: 0x1234, minimum: 3},
		// TS 29.244 V20.1.0 clause 8.2.143.
		{name: "FiveGSUserPlaneNodeID", input: ie.NewFiveGSUserPlaneNodeID(true, 0x0102030405060708), want: &ie.FiveGSUserPlaneNodeIDFields{BID: true, Value: 0x0102030405060708}, minimum: 9},
		// TS 29.244 V20.1.0 clause 8.2.61.
		{name: "FlowInformation", input: ie.NewFlowInformation(2, "abc"), want: []byte{2, 0, 3, 'a', 'b', 'c'}, minimum: 6},
		// TS 29.244 V20.1.0 clause 8.2.23.
		{name: "ForwardingPolicy", input: ie.NewForwardingPolicy("abc"), want: []byte{3, 'a', 'b', 'c'}, minimum: 4},
		// TS 29.244 V20.1.0 clause 8.2.110.
		{name: "FramedRouting", input: ie.NewFramedRouting(3), want: 3, minimum: 4},
		// TS 29.244 V20.1.0 clause 8.2.9.
		{name: "GBR", input: ie.NewGBR(0x0102030405, 0x060708090a), want: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, minimum: 10},
		// TS 29.244 V20.1.0 clause 8.2.9.
		{name: "GBRUL", input: ie.NewGBR(0x0102030405, 0x060708090a), want: uint64(0x0102030405), minimum: 10},
		// TS 29.244 V20.1.0 clause 8.2.9.
		{name: "GBRDL", input: ie.NewGBR(0x0102030405, 0x060708090a), want: uint64(0x060708090a), minimum: 10},
		// TS 29.244 V20.1.0 clause 8.2.67.
		{name: "HeaderEnrichment", input: ie.NewHeaderEnrichment(0, "abc", "def"), want: &ie.HeaderEnrichmentFields{HeaderType: 0, NameLength: 3, HeaderFieldName: "abc", ValueLength: 3, HeaderFieldValue: "def"}, minimum: 9},
		// TS 29.244 V20.1.0 clause 8.2.137.
		{name: "IPMulticastAddress", input: ie.NewIPMulticastAddress(v4, nil, nil, nil), want: &ie.IPMulticastAddressFields{Flags: 2, StartIPv4Address: v4}, minimum: 5},
		// TS 29.244 V20.1.0 clause 8.2.18.
		{name: "InactivityDetectionTime", input: ie.NewInactivityDetectionTime(123), want: 123, minimum: 4},
		// TS 29.244 V20.1.0 clause 8.2.158.
		{name: "LinkSpecificMultipathIPAddress", input: ie.NewLinkSpecificMultipathIPAddress(v4, nil, nil, nil), want: &ie.LinkSpecificMultipathIPAddressFields{Flags: 1, UELinkSpecificIPv4AddressFor3GPPAccess: v4}, minimum: 5},
		// TS 29.244 V20.1.0 clause 8.2.55.
		{name: "LinkedURRID", input: ie.NewLinkedURRID(0x12345678), want: 0x12345678, minimum: 4},
		// TS 29.244 V20.1.0 clause 8.2.93.
		{name: "MACAddress", input: ie.NewMACAddress(mac, nil, nil, nil), want: &ie.MACAddressFields{Flags: 1, SourceMACAddress: mac}, minimum: 7},
		// TS 29.244 V20.1.0 clause 8.2.103.
		{name: "MACAddressesDetected", input: ie.NewMACAddressesDetected(nil, nil, mac), want: &ie.MACAddressesDetectedFields{NumberOfMACAddresses: 1, MACAddresses: []net.HardwareAddr{mac}}, minimum: 9},
		// TS 29.244 V20.1.0 clause 8.2.104.
		{name: "MACAddressesRemoved", input: ie.NewMACAddressesRemoved(nil, nil, mac), want: &ie.MACAddressesRemovedFields{NumberOfMACAddresses: 1, MACAddresses: []net.HardwareAddr{mac}}, minimum: 9},
		// TS 29.244 V20.1.0 clause 8.2.8.
		{name: "MBR", input: ie.NewMBR(0x0102030405, 0x060708090a), want: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, minimum: 10},
		// TS 29.244 V20.1.0 clause 8.2.8.
		{name: "MBRUL", input: ie.NewMBR(0x0102030405, 0x060708090a), want: uint64(0x0102030405), minimum: 10},
		// TS 29.244 V20.1.0 clause 8.2.8.
		{name: "MBRDL", input: ie.NewMBR(0x0102030405, 0x060708090a), want: uint64(0x060708090a), minimum: 10},
		// TS 29.244 V20.1.0 clause 8.2.157.
		{name: "MPTCPAddressInformation", input: ie.NewMPTCPAddressInformation(1, 8080, v4, nil), want: &ie.MPTCPAddressInformationFields{Flags: 1, MPTCPProxyType: 1, MPTCPProxyPort: 8080, MPTCPIPv4Address: v4}, minimum: 8},
		// TS 29.244 V20.1.0 clause 8.2.40.
		{name: "MeasurementMethod", input: ie.NewMeasurementMethod(1, 0, 1), want: 5, minimum: 1},
		// TS 29.244 V20.1.0 clause 8.2.84.
		{name: "Multiplier", input: ie.NewMultiplier(1, 2), want: []byte{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 2}, minimum: 12},
		// TS 29.244 V20.1.0 clause 8.2.175.
		{name: "NFInstanceID", input: ie.NewNFInstanceID([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}), want: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, minimum: 16},
		// TS 29.244 V20.1.0 clause 8.2.142.
		{name: "NWTTPortNumber", input: ie.NewNWTTPortNumber(123), want: 123, minimum: 4},
		// TS 29.244 V20.1.0 clause 8.2.43.
		{name: "NodeAddress", input: ie.NewFQCSID("192.0.2.1", 5, 0x1234), want: []byte{192, 0, 2, 1}, minimum: 5},
		// TS 29.244 V20.1.0 clause 8.2.183.
		{name: "NumberOfUEIPAddresses", input: ie.NewNumberOfUEIPAddresses(3, 123, 456), want: &ie.NumberOfUEIPAddressesFields{Flags: 3, NumberOfUEIPv4Addresses: 123, NumberOfUEIPv6Addresses: 456}, minimum: 9},
		// TS 29.244 V20.1.0 clause 8.2.22.
		{name: "OffendingIE", input: ie.NewOffendingIE(ie.Cause), want: ie.Cause, minimum: 2},
		// TS 29.244 V20.1.0 clause 8.2.56.
		{name: "OuterHeaderCreation", input: ie.NewOuterHeaderCreation(0x0100, 0x12345678, "192.0.2.1", "", 0, 0, 0), want: &ie.OuterHeaderCreationFields{OuterHeaderCreationDescription: 0x0100, TEID: 0x12345678, IPv4Address: v4}, minimum: 10},
		// TS 29.244 V20.1.0 clause 8.2.64.
		{name: "OuterHeaderRemoval", input: ie.NewOuterHeaderRemoval(1, 2), want: []byte{1, 2}, minimum: 1},
		// TS 29.244 V20.1.0 clause 8.2.31.
		{name: "PFCPSMReqFlags", input: ie.NewPFCPSMReqFlags(5), want: 5, minimum: 1},
		// TS 29.244 V20.1.0 clause 8.2.119.
		{name: "PFCPSRReqFlags", input: ie.NewPFCPSRReqFlags(1), want: 1, minimum: 1},
		// TS 29.244 V20.1.0 clause 8.2.39.
		{name: "PFDContents", input: ie.NewPFDContents("abc", "", "", nil, "", nil, nil, nil), want: &ie.PFDContentsFields{Flags: 1, FDLength: 3, FlowDescription: "abc"}, minimum: 7},
		// TS 29.244 V20.1.0 clause 8.2.159.
		{name: "PMFAddressInformation", input: ie.NewPMFAddressInformation(v4, nil, 0, 0, nil, nil), want: &ie.PMFAddressInformationFields{Flags: 1, PMFIPv4Address: v4}, minimum: 5},
		// TS 29.244 V20.1.0 clause 8.2.169.
		{name: "PacketDelayThresholds", input: ie.NewPacketDelayThresholds(7, 123, 456, 789), want: &ie.PacketDelayThresholdsFields{Flags: 7, DownlinkPacketDelayThresholds: 123, UplinkPacketDelayThresholds: 456, RoundTripPacketDelayThresholds: 789}, minimum: 13},
		// TS 29.244 V20.1.0 clause 8.2.63.
		{name: "PacketRate", input: ie.NewPacketRate(3, 1, 123, 2, 456), want: &ie.PacketRateFields{Flags: 3, UplinkTimeUnit: 1, UplinkPacketRate: 123, DownlinkTimeUnit: 2, DownlinkPacketRate: 456}, minimum: 7},
		// TS 29.244 V20.1.0 clause 8.2.139.
		{name: "PacketRateStatus", input: ie.NewPacketRateStatus(3, 123, 0, 456, 0, time.Unix(1700000000, 0).UTC()), want: &ie.PacketRateStatusFields{Flags: 3, NumberOfRemainingUplinkPacketsAllowed: 123, NumberOfRemainingDownlinkPacketsAllowed: 456, RateControlStatusValidityTime: time.Unix(1700000000, 0).UTC()}, minimum: 13},
		// TS 29.244 V20.1.0 clause 8.2.11.
		{name: "Precedence", input: ie.NewPrecedence(0x12345678), want: 0x12345678, minimum: 4},
		// TS 29.244 V20.1.0 clause 8.2.10.
		{name: "QERCorrelationID", input: ie.NewQERCorrelationID(0x12345678), want: 0x12345678, minimum: 4},
		// TS 29.244 V20.1.0 clause 8.2.75.
		{name: "QERID", input: ie.NewQERID(0x12345678), want: 0x12345678, minimum: 4},
		// TS 29.244 V20.1.0 clause 8.2.171.
		{name: "QoSMonitoringMeasurement", input: ie.NewQoSMonitoringMeasurement(7, 123, 456, 789), want: &ie.QoSMonitoringMeasurementFields{Flags: 7, DownlinkPacketDelay: 123, UplinkPacketDelay: 456, RoundTripPacketDelay: 789}, minimum: 13},
		// TS 29.244 V20.1.0 clause 8.2.90.
		{name: "QueryURRReference", input: ie.NewQueryURRReference(0x12345678), want: 0x12345678, minimum: 4},
		{name: "RedirectInformation", input: ie.NewRedirectInformation(0, "192.0.2.1"), want: &ie.RedirectInformationFields{RedirectAddressType: 0, ServerAddrLength: 9, RedirectServerAddress: "192.0.2.1"}, minimum: 12},
		// TS 29.244 V20.1.0 clause 8.2.70.
		{name: "RemoteGTPUPeer", input: ie.NewRemoteGTPUPeer(2, "192.0.2.1", "", 0, "", time.Time{}), want: &ie.RemoteGTPUPeerFields{Flags: 2, IPv4Address: v4}, minimum: 5},
		// TS 29.244 V20.1.0 clause 8.2.19.
		{name: "ReportingTriggers", input: ie.NewReportingTriggers(1, 2), want: []byte{1, 2}, minimum: 2},
		// TS 29.244 V20.1.0 clause 8.2.5.
		{name: "SDFFilter", input: ie.NewSDFFilter("abc", nil, nil, nil, 0), want: &ie.SDFFilterFields{Flags: 1, FDLength: 3, FlowDescription: "abc"}, minimum: 7},
		// TS 29.244 V20.1.0 clause 8.2.176.
		{name: "SNSSAI", input: ie.NewSNSSAI(1, 0x234567), want: []byte{1, 0x23, 0x45, 0x67}, minimum: 4},
		// TS 29.244 V20.1.0 clause 8.2.95.
		{name: "STAG", input: ie.NewSTAG(7, 3, 1, 0xabc), want: &ie.STAGFields{Flags: 7, PCP: 3, DEIFlag: 1, CVID: 0xabc}, minimum: 3},
		// TS 29.244 V20.1.0 clause 8.2.33.
		{name: "SequenceNumber", input: ie.NewSequenceNumber(0x12345678), want: 0x12345678, minimum: 4},
		// TS 29.244 V20.1.0 clause 8.2.138.
		{name: "SourceIPAddress", input: ie.NewSourceIPAddress(v4, nil, 0), want: &ie.SourceIPAddressFields{Flags: 2, IPv4Address: v4}, minimum: 5},
		// TS 29.244 V20.1.0 clause 8.2.106.
		{name: "SubsequentEventQuota", input: ie.NewSubsequentEventQuota(123), want: 123, minimum: 4},
		// TS 29.244 V20.1.0 clause 8.2.107.
		{name: "SubsequentEventThreshold", input: ie.NewSubsequentEventThreshold(123), want: 123, minimum: 4},
		// TS 29.244 V20.1.0 clause 8.2.86.
		{name: "SubsequentVolumeQuota", input: ie.NewSubsequentVolumeQuota(7, 123, 456, 789), want: &ie.SubsequentVolumeQuotaFields{Flags: 7, TotalVolume: 123, UplinkVolume: 456, DownlinkVolume: 789}, minimum: 25},
		// TS 29.244 V20.1.0 clause 8.2.16.
		{name: "SubsequentVolumeThreshold", input: ie.NewSubsequentVolumeThreshold(7, 123, 456, 789), want: &ie.SubsequentVolumeThresholdFields{Flags: 7, TotalVolume: 123, UplinkVolume: 456, DownlinkVolume: 789}, minimum: 25},
		// TS 29.244 V20.1.0 clause 8.2.81.
		{name: "TimeQuotaMechanism", input: ie.NewTimeQuotaMechanism(1, 10*time.Second), want: []byte{1, 0, 0, 0, 10}, minimum: 5},
		// TS 29.244 V20.1.0 clause 8.2.108.
		{name: "TraceInformation", input: ie.NewTraceInformation("001", "01", "abc", []byte{1}, 2, []byte{3}, v4), want: &ie.TraceInformationFields{MCC: "001", MNC: "01", TraceID: "abc", TriggeringEventsLength: 1, TriggeringEvents: []byte{1}, SessionTraceDepth: 2, ListOfInterfacesLength: 1, ListOfInterfaces: []byte{3}, IPAddressOfTraceCollectionEntityLength: 4, IPAddressOfTraceCollectionEntity: v4}, minimum: 16},
		// TS 29.244 V20.1.0 clause 8.2.62.
		{name: "UEIPAddress", input: ie.NewUEIPAddress(1, "", "2001:db8::1", 0, 0), want: &ie.UEIPAddressFields{Flags: 1, IPv6Address: v6}, minimum: 17},
		// TS 29.244 V20.1.0 clause 8.2.128.
		{name: "UEIPAddressPoolIdentity", input: ie.NewUEIPAddressPoolIdentity("abc"), want: "abc", minimum: 5},
		// TS 29.244 V20.1.0 clause 8.2.25.
		{name: "UPFunctionFeatures", input: ie.NewUPFunctionFeatures(1, 2), want: []byte{1, 2}, minimum: 2},
		// TS 29.244 V20.1.0 clause 8.2.54.
		{name: "URRID", input: ie.NewURRID(0x12345678), want: 0x12345678, minimum: 4},
		// TS 29.244 V20.1.0 clause 8.2.71.
		{name: "URSEQN", input: ie.NewURSEQN(0x12345678), want: 0x12345678, minimum: 4},
		// TS 29.244 V20.1.0 clause 8.2.41.
		{name: "UsageReportTrigger", input: ie.NewUsageReportTrigger(1, 2), want: []byte{1, 2}, minimum: 2},
		// TS 29.244 V20.1.0 clause 8.2.101.
		{name: "UserID", input: ie.NewUserID(8, "", "", "", "abc"), want: &ie.UserIDFields{Flags: 8, NAILength: 3, NAI: "abc"}, minimum: 5},
		// TS 29.244 V20.1.0 clause 8.2.44.
		{name: "VolumeMeasurement", input: ie.NewVolumeMeasurement(63, 123, 456, 789, 12, 34, 56), want: &ie.VolumeMeasurementFields{Flags: 63, TotalVolume: 123, UplinkVolume: 456, DownlinkVolume: 789, TotalNumberOfPackets: 12, UplinkNumberOfPackets: 34, DownlinkNumberOfPackets: 56}, minimum: 49},
		// TS 29.244 V20.1.0 clause 8.2.50.
		{name: "VolumeQuota", input: ie.NewVolumeQuota(7, 123, 456, 789), want: &ie.VolumeQuotaFields{Flags: 7, TotalVolume: 123, UplinkVolume: 456, DownlinkVolume: 789}, minimum: 25},
		// TS 29.244 V20.1.0 clause 8.2.13.
		{name: "VolumeThreshold", input: ie.NewVolumeThreshold(7, 123, 456, 789), want: &ie.VolumeThresholdFields{Flags: 7, TotalVolume: 123, UplinkVolume: 456, DownlinkVolume: 789}, minimum: 25},
	}
}
