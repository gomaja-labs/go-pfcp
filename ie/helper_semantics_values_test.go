package ie_test

import (
	"net/netip"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// These value expectations retain the explicit decoded values from the existing
// scalar, string, duration, timestamp and octet-array vectors. They are inputs
// to the direct contract matrix, never values computed by the decoder.
func helperSemanticValues() []helperSemanticCase {
	values := []helperSemanticCase{
		// TS 29.244 V20.1.0 clause 8.2.117.
		{name: "APNDNN", input: ie.NewAPNDNN("some.apn.example"), want: "some.apn.example"},
		// TS 29.244 V20.1.0 clause 8.2.155.
		{name: "ATSSSLLControlInformation", input: ie.NewATSSSLLControlInformation(1), want: 1},
		// TS 29.244 V20.1.0 clause 8.2.160.
		{name: "ATSSSLLInformation", input: ie.NewATSSSLLInformation(1), want: 1},
		// TS 29.244 V20.1.0 clause 8.2.153.
		{name: "AccessAvailabilityInformation", input: ie.NewAccessAvailabilityInformation(3, 3), want: 0x0f},
		// TS 29.244 V20.1.0 clause 8.2.153.
		{name: "AccessType", input: ie.NewAccessAvailabilityInformation(3, 3), want: 3},
		// TS 29.244 V20.1.0 clause 8.2.72.
		{name: "ActivatePredefinedRules", input: ie.NewActivatePredefinedRules("go-pfcp"), want: "go-pfcp"},
		// TS 29.244 V20.1.0 clause 8.2.121.
		{name: "ActivationTime", input: ie.NewActivationTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)), want: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)},
		// TS 29.244 V20.1.0 clause 8.2.91.
		{name: "AdditionalUsageReportsInformation", input: ie.NewAdditionalUsageReportsInformation(0x00ff), want: 0x00ff},
		// TS 29.244 V20.1.0 clause 8.2.85.
		{name: "AggregatedURRID", input: ie.NewAggregatedURRID(0xffffffff), want: uint32(0xffffffff)},
		// TS 29.244 V20.1.0 clause 8.2.6.
		{name: "ApplicationID", input: ie.NewApplicationID("https://example.com/apps/go-pfcp/"), want: "https://example.com/apps/go-pfcp/"},
		// TS 29.244 V20.1.0 clause 8.2.60.
		{name: "ApplicationInstanceID", input: ie.NewApplicationInstanceID("go-pfcp"), want: "go-pfcp"},
		// TS 29.244 V20.1.0 clause 8.2.26.
		{name: "ApplyAction", input: ie.NewApplyAction(0x04), want: []byte{0x04}},
		// TS 29.244 V20.1.0 clause 8.2.153.
		{name: "AvailabilityStatus", input: ie.NewAccessAvailabilityInformation(3, 3), want: 3},
		// TS 29.244 V20.1.0 clause 8.2.162.
		{name: "AveragePacketDelay", input: ie.NewAveragePacketDelay(10 * time.Second), want: 10 * time.Second},
		// TS 29.244 V20.1.0 clause 8.2.115.
		{name: "AveragingWindow", input: ie.NewAveragingWindow(0xffffffff), want: uint32(0xffffffff)},
		// TS 29.244 V20.1.0 clause 8.2.57.
		{name: "BARID", input: ie.NewBARID(0xff), want: 0xff},
		// TS 29.244 V20.1.0 clause 8.2.58.
		{name: "CPFunctionFeatures", input: ie.NewCPFunctionFeatures(0x3f), want: []byte{0x3f}},
		// TS 29.244 V20.1.0 clause 8.2.1.
		{name: "Cause", input: ie.NewCause(ie.CauseRequestAccepted), want: ie.CauseRequestAccepted},
		// TS 29.244 V20.1.0 clause 8.2.140.
		{name: "CreateBridgeRouterInfo", input: ie.NewCreateBridgeRouterInfo(1), want: 1},
		// TS 29.244 V20.1.0 clause 8.2.148.
		{name: "CumulativeRateRatioThreshold", input: ie.NewCumulativeRateRatioThreshold(0xffffffff), want: uint32(0xffffffff)},
		// TS 29.244 V20.1.0 clause 8.2.29.
		{name: "DLBufferingDuration", input: ie.NewDLBufferingDuration(20 * time.Hour), want: 20 * time.Hour},
		// TS 29.244 V20.1.0 clause 8.2.30.
		{name: "DLBufferingSuggestedPacketCount", input: ie.NewDLBufferingSuggestedPacketCount(0xffff), want: 0xffff},
		// TS 29.244 V20.1.0 clause 8.2.173.
		{name: "DLDataPacketsSize", input: ie.NewDLDataPacketsSize(0xffff), want: 0xffff},
		// TS 29.244 V20.1.0 clause 8.2.161.
		{name: "DataNetworkAccessIdentifier", input: ie.NewDataNetworkAccessIdentifier("go-pfcp"), want: "go-pfcp"},
		// TS 29.244 V20.1.0 clause 8.2.179.
		{name: "DataStatus", input: ie.NewDataStatus(0x03), want: 0x03},
		// TS 29.244 V20.1.0 clause 8.2.73.
		{name: "DeactivatePredefinedRules", input: ie.NewDeactivatePredefinedRules("go-pfcp"), want: "go-pfcp"},
		// TS 29.244 V20.1.0 clause 8.2.122.
		{name: "DeactivationTime", input: ie.NewDeactivationTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)), want: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)},
		// TS 29.244 V20.1.0 clause 8.2.24.
		{name: "DestinationInterface", input: ie.NewDestinationInterface(ie.DstInterfaceAccess), want: ie.DstInterfaceAccess},
		// TS 29.244 V20.1.0 clause 8.2.28.
		{name: "DownlinkDataNotificationDelay", input: ie.NewDownlinkDataNotificationDelay(100 * time.Millisecond), want: 100 * time.Millisecond},
		// TS 29.244 V20.1.0 clause 8.2.45.
		{name: "DurationMeasurement", input: ie.NewDurationMeasurement(10 * time.Second), want: 10 * time.Second},
		// TS 29.244 V20.1.0 clause 8.2.53.
		{name: "EndTime", input: ie.NewEndTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)), want: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)},
		// TS 29.244 V20.1.0 clause 8.2.98.
		{name: "EthernetFilterID", input: ie.NewEthernetFilterID(0xffffffff), want: uint32(0xffffffff)},
		// TS 29.244 V20.1.0 clause 8.2.99.
		{name: "EthernetFilterProperties", input: ie.NewEthernetFilterProperties(0x01), want: 0x01},
		// TS 29.244 V20.1.0 clause 8.2.105.
		{name: "EthernetInactivityTimer", input: ie.NewEthernetInactivityTimer(10 * time.Second), want: 10 * time.Second},
		// TS 29.244 V20.1.0 clause 8.2.102.
		{name: "EthernetPDUSessionInformation", input: ie.NewEthernetPDUSessionInformation(0x01), want: 0x01},
		// TS 29.244 V20.1.0 clause 8.2.96.
		{name: "Ethertype", input: ie.NewEthertype(0xffff), want: 0xffff},
		// TS 29.244 V20.1.0 clause 8.2.112.
		{name: "EventQuota", input: ie.NewEventQuota(0xffffffff), want: uint32(0xffffffff)},
		// TS 29.244 V20.1.0 clause 8.2.113.
		{name: "EventThreshold", input: ie.NewEventThreshold(0xffffffff), want: uint32(0xffffffff)},
		// TS 29.244 V20.1.0 clause 8.2.114.
		{name: "EventTimeStamp", input: ie.NewEventTimeStamp(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)), want: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)},
		// TS 29.244 V20.1.0 clause 8.2.84.
		{name: "Exponent", input: ie.NewMultiplier(0xffffffffffffffff, 0x11223344), want: 0x11223344},
		// TS 29.244 V20.1.0 clause 8.2.61.
		{name: "FlowDescription", input: ie.NewFlowInformation(ie.FlowDirectionDownlink, "go-pfcp"), want: "go-pfcp"},
		// TS 29.244 V20.1.0 clause 8.2.61.
		{name: "FlowDirection", input: ie.NewFlowInformation(ie.FlowDirectionDownlink, "go-pfcp"), want: ie.FlowDirectionDownlink},
		// TS 29.244 V20.1.0 clause 8.2.23.
		{name: "ForwardingPolicyIdentifier", input: ie.NewForwardingPolicy("go-pfcp"), want: "go-pfcp"},
		// TS 29.244 V20.1.0 clause 8.2.111.
		{name: "FramedIPv6Route", input: ie.NewFramedIPv6Route("go-pfcp"), want: "go-pfcp"},
		// TS 29.244 V20.1.0 clause 8.2.109.
		{name: "FramedRoute", input: ie.NewFramedRoute("go-pfcp"), want: "go-pfcp"},
		// TS 29.244 V20.1.0 clause 8.2.64.
		{name: "GTPUExtensionHeaderDeletion", input: ie.NewOuterHeaderRemoval(0x01, 0x02), want: 0x02},
		// TS 29.244 V20.1.0 clause 8.2.166.
		{name: "GTPUPathInterfaceType", input: ie.NewGTPUPathInterfaceType(1, 1), want: 3},
		// TS 29.244 V20.1.0 clause 8.2.7.
		{name: "GateStatus", input: ie.NewGateStatus(ie.GateStatusOpen, ie.GateStatusOpen), want: 0},
		// TS 29.244 V20.1.0 clause 8.2.7.
		{name: "GateStatusDL", input: ie.NewGateStatus(ie.GateStatusOpen, ie.GateStatusOpen), want: ie.GateStatusOpen},
		// TS 29.244 V20.1.0 clause 8.2.7.
		{name: "GateStatusUL", input: ie.NewGateStatus(ie.GateStatusOpen, ie.GateStatusOpen), want: ie.GateStatusOpen},
		// TS 29.244 V20.1.0 clause 8.2.78.
		{name: "GracefulReleasePeriod", input: ie.NewGracefulReleasePeriod(20 * time.Hour), want: 20 * time.Hour},
		// TS 29.244 V20.1.0 clause 8.2.177.
		{name: "IPVersion", input: ie.NewIPVersion(true, false), want: 0x01},
		// TS 29.244 V20.1.0 clause 8.2.123.
		{name: "MARID", input: ie.NewMARID(0x1111), want: 0x1111},
		// TS 29.244 V20.1.0 clause 8.2.154.
		{name: "MPTCPControlInformation", input: ie.NewMPTCPControlInformation(1), want: 1},
		// TS 29.244 V20.1.0 clause 8.2.172.
		{name: "MTEDTControlInformation", input: ie.NewMTEDTControlInformation(1), want: 1},
		// TS 29.244 V20.1.0 clause 8.2.164.
		{name: "MaximumPacketDelay", input: ie.NewMaximumPacketDelay(10 * time.Second), want: 10 * time.Second},
		// TS 29.244 V20.1.0 clause 8.2.68.
		{name: "MeasurementInformation", input: ie.NewMeasurementInformation(0x1f), want: 0x1f},
		// TS 29.244 V20.1.0 clause 8.2.42.
		{name: "MeasurementPeriod", input: ie.NewMeasurementPeriod(10 * time.Second), want: 10 * time.Second},
		// TS 29.244 V20.1.0 clause 8.2.34.
		{name: "Metric", input: ie.NewMetric(0x01), want: 0x01},
		// TS 29.244 V20.1.0 clause 8.2.163.
		{name: "MinimumPacketDelay", input: ie.NewMinimumPacketDelay(10 * time.Second), want: 10 * time.Second},
		// TS 29.244 V20.1.0 clause 8.2.170.
		{name: "MinimumWaitTime", input: ie.NewMinimumWaitTime(10 * time.Second), want: 10 * time.Second},
		// TS 29.244 V20.1.0 clause 8.2.15.
		{name: "MonitoringTime", input: ie.NewMonitoringTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)), want: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)},
		// TS 29.244 V20.1.0 clause 8.2.181.
		{name: "MultipathApplicableIndication", input: ie.NewMultipathApplicableIndication(0x01), want: 0x01},
		// TS 29.244 V20.1.0 clause 8.2.4.
		{name: "NetworkInstance", input: ie.NewNetworkInstance("go-pfcp"), want: "go-pfcp"},
		{name: "NetworkInstanceFQDN", input: ie.NewNetworkInstanceFQDN("go-pfcp.com"), want: "go-pfcp.com"},
		// TS 29.244 V20.1.0 clause 8.2.4.
		{name: "NetworkInstanceHeuristic", input: ie.NewNetworkInstance("go-pfcp"), want: "go-pfcp"},
		// TS 29.244 V20.1.0 clause 8.2.38.
		{name: "NodeID", input: ie.NewNodeIDAddress(netip.MustParseAddr("192.0.2.1")), want: ie.NodeIDValue{Type: ie.NodeIDIPv4Address, Address: netip.MustParseAddr("192.0.2.1")}},
		// TS 29.244 V20.1.0 clause 8.2.43.
		{name: "NodeIDType", input: ie.NewFQCSID("127.0.0.1", 5, 1), want: ie.NodeIDIPv4Address},
		// TS 29.244 V20.1.0 clause 8.2.69.
		{name: "NodeReportType", input: ie.NewNodeReportType(0x01), want: 0x01},
		// TS 29.244 V20.1.0 clause 8.2.133.
		{name: "NumberOfReports", input: ie.NewNumberOfReports(0x1111), want: 0x1111},
		// TS 29.244 V20.1.0 clause 8.2.76.
		{name: "OCIFlags", input: ie.NewOCIFlags(0x01), want: 0x01},
		// TS 29.244 V20.1.0 clause 8.2.64.
		{name: "OuterHeaderRemovalDescription", input: ie.NewOuterHeaderRemoval(0x01, 0x02), want: 0x01},
		// TS 29.244 V20.1.0 clause 8.2.79.
		{name: "PDNType", input: ie.NewPDNType(ie.PDNTypeIPv4), want: ie.PDNTypeIPv4},
		// TS 29.244 V20.1.0 clause 8.2.36.
		{name: "PDRID", input: ie.NewPDRID(0xffff), want: 0xffff},
		// TS 29.244 V20.1.0 clause 8.2.178.
		{name: "PFCPASReqFlags", input: ie.NewPFCPASReqFlags(0x01), want: 0x01},
		// TS 29.244 V20.1.0 clause 8.2.134.
		{name: "PFCPASRspFlags", input: ie.NewPFCPASRspFlags(0x01), want: 0x01},
		// TS 29.244 V20.1.0 clause 8.2.120.
		{name: "PFCPAUReqFlags", input: ie.NewPFCPAUReqFlags(0x01), want: 0x01},
		// TS 29.244 V20.1.0 clause 8.2.77.
		{name: "PFCPAssociationReleaseRequest", input: ie.NewPFCPAssociationReleaseRequest(1, 1), want: 3},
		// TS 29.244 V20.1.0 clause 8.2.136.
		{name: "PFCPSEReqFlags", input: ie.NewPFCPSEReqFlags(0x01), want: 0x01},
		// TS 29.244 V20.1.0 clause 8.2.32.
		{name: "PFCPSRRspFlags", input: ie.NewPFCPSRRspFlags(0x01), want: 0x01},
		// TS 29.244 V20.1.0 clause 8.2.156.
		{name: "PMFControlInformation", input: ie.NewPMFControlInformation(7, 1, 9), want: &ie.PMFControlInformationFields{Flags: 7, QFIs: []uint8{1, 9}}},
		// TS 29.244 V20.1.0 clause 8.2.27.
		{name: "PPI", input: ie.NewDownlinkDataServiceInformation(true, true, 1, 1), want: 1},
		// TS 29.244 V20.1.0 clause 8.2.130.
		{name: "PacketReplicationAndDetectionCarryOnInformation", input: ie.NewPacketReplicationAndDetectionCarryOnInformation(0x0f), want: 0x0f},
		// TS 29.244 V20.1.0 clause 8.2.116.
		{name: "PagingPolicyIndicator", input: ie.NewPagingPolicyIndicator(1), want: 1},
		// TS 29.244 V20.1.0 clause 8.2.144.
		{name: "PortManagementInformationContainer", input: ie.NewPortManagementInformationContainer([]byte("go-pfcp")), want: []byte("go-pfcp")},
		// TS 29.244 V20.1.0 clause 8.2.141.
		{name: "PortNumber", input: ie.NewPortNumber(0xffffffff), want: uint32(0xffffffff)},
		// TS 29.244 V20.1.0 clause 8.2.127.
		{name: "Priority", input: ie.NewPriority(ie.PriorityActive), want: ie.PriorityActive},
		// TS 29.244 V20.1.0 clause 8.2.97.
		{name: "Proxying", input: ie.NewProxying(1, 1), want: 0x03},
		// TS 29.244 V20.1.0 clause 8.2.174.
		{name: "QERControlIndications", input: ie.NewQERControlIndications(0, 0, 1), want: 0x01},
		// TS 29.244 V20.1.0 clause 8.2.27.
		{name: "QFI", input: ie.NewDownlinkDataServiceInformation(true, true, 1, 1), want: 1},
		// TS 29.244 V20.1.0 clause 8.2.165.
		{name: "QoSReportTrigger", input: ie.NewQoSReportTrigger(1, 1, 1), want: 0x07},
		// TS 29.244 V20.1.0 clause 8.2.48.
		{name: "QuotaHoldingTime", input: ie.NewQuotaHoldingTime(10 * time.Second), want: 10 * time.Second},
		// TS 29.244 V20.1.0 clause 8.2.132.
		{name: "QuotaValidityTime", input: ie.NewQuotaValidityTime(10 * time.Second), want: 10 * time.Second},
		// TS 29.244 V20.1.0 clause 8.2.180.
		{name: "RDSConfigurationInformation", input: ie.NewRDSConfigurationInformation(0x01), want: 0x01},
		// TS 29.244 V20.1.0 clause 8.2.88.
		{name: "RQI", input: ie.NewRQI(0x01), want: 0x01},
		// TS 29.244 V20.1.0 clause 8.2.65.
		{name: "RecoveryTimeStamp", input: ie.NewRecoveryTimeStamp(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)), want: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)},
		// TS 29.244 V20.1.0 clause 8.2.21.
		{name: "ReportType", input: ie.NewReportType(1, 1, 1, 1), want: 0x0f},
		// TS 29.244 V20.1.0 clause 8.2.168.
		{name: "ReportingFrequency", input: ie.NewReportingFrequency(0, 1, 1), want: 0x03},
		// TS 29.244 V20.1.0 clause 8.2.152.
		{name: "RequestedAccessAvailabilityInformation", input: ie.NewRequestedAccessAvailabilityInformation(1), want: 1},
		// TS 29.244 V20.1.0 clause 8.2.145.
		{name: "RequestedClockDriftInformation", input: ie.NewRequestedClockDriftInformation(1, 1), want: 0x03},
		// TS 29.244 V20.1.0 clause 8.2.167.
		{name: "RequestedQoSMonitoring", input: ie.NewRequestedQoSMonitoring(1, 1, 1), want: 0x07},
		// TS 29.244 V20.1.0 clause 8.2.80.
		{name: "RuleIDType", input: ie.NewFailedRuleID(ie.RuleIDTypePDR, 0xffff), want: ie.RuleIDTypePDR},
		// TS 29.244 V20.1.0 clause 8.2.176.
		{name: "SD", input: ie.NewSNSSAI(0x11, 0x223344), want: 0x223344},
		// TS 29.244 V20.1.0 clause 8.2.131.
		{name: "SMFSetID", input: ie.NewSMFSetID("go-pfcp"), want: "go-pfcp"},
		// TS 29.244 V20.1.0 clause 8.2.151.
		{name: "SRRID", input: ie.NewSRRID(255), want: 255},
		// TS 29.244 V20.1.0 clause 8.2.176.
		{name: "SST", input: ie.NewSNSSAI(0x11, 0x223344), want: 0x11},
		// TS 29.244 V20.1.0 clause 8.2.2.
		{name: "SourceInterface", input: ie.NewSourceInterface(ie.SrcInterfaceAccess), want: ie.SrcInterfaceAccess},
		// TS 29.244 V20.1.0 clause 8.2.52.
		{name: "StartTime", input: ie.NewStartTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)), want: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)},
		// TS 29.244 V20.1.0 clause 8.2.124.
		{name: "SteeringFunctionality", input: ie.NewSteeringFunctionality(ie.SteeringFunctionalityATSSSLL), want: ie.SteeringFunctionalityATSSSLL},
		// TS 29.244 V20.1.0 clause 8.2.125.
		{name: "SteeringMode", input: ie.NewSteeringMode(ie.SteeringModeActiveStandby), want: ie.SteeringModeActiveStandby},
		// TS 29.244 V20.1.0 clause 8.2.87.
		{name: "SubsequentTimeQuota", input: ie.NewSubsequentTimeQuota(10 * time.Second), want: 10 * time.Second},
		// TS 29.244 V20.1.0 clause 8.2.17.
		{name: "SubsequentTimeThreshold", input: ie.NewSubsequentTimeThreshold(10 * time.Second), want: 10 * time.Second},
		// TS 29.244 V20.1.0 clause 8.2.100.
		{name: "SuggestedBufferingPacketsCount", input: ie.NewSuggestedBufferingPacketsCount(0x01), want: 0x01},
		// TS 29.244 V20.1.0 clause 8.2.118.
		{name: "TGPPInterfaceType", input: ie.NewTGPPInterfaceType(ie.TGPPInterfaceTypeS1U), want: ie.TGPPInterfaceTypeS1U},
		// TS 29.244 V20.1.0 clause 8.2.146.
		{name: "TimeDomainNumber", input: ie.NewTimeDomainNumber(255), want: 255},
		// TS 29.244 V20.1.0 clause 8.2.46.
		{name: "TimeOfFirstPacket", input: ie.NewTimeOfFirstPacket(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)), want: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)},
		// TS 29.244 V20.1.0 clause 8.2.47.
		{name: "TimeOfLastPacket", input: ie.NewTimeOfLastPacket(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)), want: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)},
		// TS 29.244 V20.1.0 clause 8.2.149.
		{name: "TimeOffsetMeasurement", input: ie.NewTimeOffsetMeasurement(10 * time.Second), want: 10 * time.Second},
		// TS 29.244 V20.1.0 clause 8.2.147.
		{name: "TimeOffsetThreshold", input: ie.NewTimeOffsetThreshold(10 * time.Second), want: 10 * time.Second},
		// TS 29.244 V20.1.0 clause 8.2.51.
		{name: "TimeQuota", input: ie.NewTimeQuota(10 * time.Second), want: 10 * time.Second},
		// TS 29.244 V20.1.0 clause 8.2.14.
		{name: "TimeThreshold", input: ie.NewTimeThreshold(10 * time.Second), want: 10 * time.Second},
		// TS 29.244 V20.1.0 clause 8.2.35.
		{name: "Timer", input: ie.NewTimer(20 * time.Hour), want: 20 * time.Hour},
		// TS 29.244 V20.1.0 clause 8.2.92.
		{name: "TrafficEndpointID", input: ie.NewTrafficEndpointID(0x01), want: 0x01},
		// TS 29.244 V20.1.0 clause 8.2.12.
		{name: "TransportLevelMarking", input: ie.NewTransportLevelMarking(0x1111), want: 0x1111},
		// TS 29.244 V20.1.0 clause 8.2.128.
		{name: "UEIPAddressPoolIdentityString", input: ie.NewUEIPAddressPoolIdentity("go-pfcp"), want: "go-pfcp"},
		// TS 29.244 V20.1.0 clause 8.2.59.
		{name: "UsageInformation", input: ie.NewUsageInformation(1, 1, 1, 1), want: 0x0f},
		// TS 29.244 V20.1.0 clause 8.2.83.
		{name: "UserPlaneInactivityTimer", input: ie.NewUserPlaneInactivityTimer(10 * time.Second), want: 10 * time.Second},
		// TS 29.244 V20.1.0 clause 8.2.182.
		{name: "UserPlaneNodeManagementInformationContainer", input: ie.NewUserPlaneNodeManagementInformationContainer([]byte("go-pfcp")), want: []byte("go-pfcp")},
		// TS 29.244 V20.1.0 clause 8.2.184.
		{name: "ValidityTimer", input: ie.NewValidityTimer(10 * time.Second), want: 10 * time.Second},
		// TS 29.244 V20.1.0 clause 8.2.84.
		{name: "ValueDigits", input: ie.NewMultiplier(0xffffffffffffffff, 0x11223344), want: uint64(0xffffffffffffffff)},
		// TS 29.244 V20.1.0 clause 8.2.126.
		{name: "Weight", input: ie.NewWeight(0x01), want: 0x01},
	}
	// Prefix boundaries follow the selected layouts: numeric/timestamp fixtures
	// contain one complete value; text strings have no malformed nonempty prefix.
	// Name and explicit-length fields must retain their selected complete value.
	for n := range values {
		tc := &values[n]
		switch tc.want.(type) {
		case string, []byte:
			tc.minimum = 1
		default:
			tc.minimum = len(tc.input.Payload)
		}
		switch tc.name {
		case "SMFSetID", "NodeID", "FlowDescription", "ForwardingPolicyIdentifier", "UEIPAddressPoolIdentityString":
			tc.minimum = len(tc.input.Payload)
		case "APNDNN":
			tc.minimum = 5
		case "NetworkInstanceFQDN":
			tc.minimum = 8
		case "DLBufferingSuggestedPacketCount":
			tc.minimum = 1
		case "FlowDirection", "NodeIDType", "OuterHeaderRemovalDescription", "RuleIDType":
			tc.minimum = 1
		case "PPI", "GTPUExtensionHeaderDeletion":
			tc.minimum = 2
		}
	}
	return append(values, helperExtendedSemanticValues()...)
}
