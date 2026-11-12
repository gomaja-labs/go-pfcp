package ie_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestValueHelperRejectsGroupedParent(t *testing.T) {
	pdr := ie.NewCreatePDR(ie.NewPDRID(42))
	if _, err := pdr.PDRID(); err == nil {
		t.Fatal("PDRID accepted grouped parent")
	} else {
		var invalid *ie.InvalidTypeError
		if !errors.As(err, &invalid) || invalid.Type != ie.CreatePDR {
			t.Fatalf("wrong error: %v", err)
		}
	}
}

// Every value method whose historical grouped traversal was removed rejects
// all grouped parent types, including when their children contain a value.
func TestValueHelpersRejectAllGroupedParents(t *testing.T) {
	methods := []string{"ATSSSLLControlInformation", "ATSSSLLInformation", "AccessAvailabilityInformation", "ActivatePredefinedRules", "ActivationTime", "AggregatedURRID", "ApplicationID", "ApplicationInstanceID", "ApplyAction", "AveragePacketDelay", "AveragingWindow", "BARID", "CPPFCPEntityIPAddress", "CTAG", "CumulativeRateRatioThreshold", "DLBufferingDuration", "DLBufferingSuggestedPacketCount", "DLDataPacketsSize", "DLFlowLevelMarking", "DataNetworkAccessIdentifier", "DataStatus", "DeactivatePredefinedRules", "DeactivationTime", "DestinationInterface", "DownlinkDataNotificationDelay", "DownlinkDataServiceInformation", "DroppedDLTrafficThreshold", "DurationMeasurement", "EndTime", "EthernetFilterID", "EthernetFilterProperties", "EthernetInactivityTimer", "EthernetPDUSessionInformation", "Ethertype", "EventQuota", "EventThreshold", "EventTimeStamp", "Exponent", "FARID", "FTEID", "FiveGSUserPlaneNodeID", "FlowDescription", "FlowDirection", "FlowInformation", "ForwardingPolicy", "FramedIPv6Route", "FramedRoute", "FramedRouting", "GBR", "GTPUPathInterfaceType", "GateStatus", "HasBUFF", "HasDROP", "HasDROTH", "HasEMRRE", "HasENVCL", "HasEVEQU", "HasEVETH", "HasIMMER", "HasIPMJL", "HasLIUSA", "HasMACAR", "HasMONIT", "HasPERIO", "HasQUHTI", "HasQUVTI", "HasREEMR", "HasSTART", "HasSTOPT", "HasTEBUR", "HasTERMR", "HasTIMQU", "HasTIMTH", "HasUPINT", "HasVOLQU", "HasVOLTH", "HeaderEnrichment", "IPMulticastAddress", "IPVersion", "InactivityDetectionTime", "LinkSpecificMultipathIPAddress", "LinkedURRID", "MACAddress", "MACAddressesDetected", "MACAddressesRemoved", "MARID", "MBR", "MPTCPAddressInformation", "MPTCPControlInformation", "MTEDTControlInformation", "MaximumPacketDelay", "MeasurementInformation", "MeasurementMethod", "MeasurementPeriod", "Metric", "MinimumPacketDelay", "MinimumWaitTime", "MonitoringTime", "MultipathApplicableIndication", "Multiplier", "NFInstanceID", "NWTTPortNumber", "NetworkInstance", "NumberOfReports", "NumberOfUEIPAddresses", "OCIFlags", "OuterHeaderCreation", "OuterHeaderRemoval", "PDRID", "PFCPSMReqFlags", "PFDContents", "PMFAddressInformation", "PMFControlInformation", "PacketDelayThresholds", "PacketRate", "PacketRateStatus", "PacketReplicationAndDetectionCarryOnInformation", "PagingPolicyIndicator", "PortManagementInformationContainer", "PortNumber", "Precedence", "Priority", "Proxying", "QERControlIndications", "QERCorrelationID", "QERID", "QFI", "QoSMonitoringMeasurement", "QoSReportTrigger", "QueryURRReference", "QuotaHoldingTime", "QuotaValidityTime", "RDSConfigurationInformation", "RQI", "RedirectInformation", "RemoteGTPUPeer", "ReportingFrequency", "ReportingTriggers", "RequestedAccessAvailabilityInformation", "RequestedClockDriftInformation", "RequestedQoSMonitoring", "SDFFilter", "SNSSAI", "SRRID", "STAG", "SequenceNumber", "SourceIPAddress", "SourceInterface", "StartTime", "SteeringFunctionality", "SteeringMode", "SubsequentEventQuota", "SubsequentEventThreshold", "SubsequentTimeQuota", "SubsequentTimeThreshold", "SubsequentVolumeQuota", "SubsequentVolumeThreshold", "SuggestedBufferingPacketsCount", "TGPPInterfaceType", "TimeDomainNumber", "TimeOfFirstPacket", "TimeOfLastPacket", "TimeOffsetThreshold", "TimeQuota", "TimeQuotaMechanism", "TimeThreshold", "Timer", "TrafficEndpointID", "TransportLevelMarking", "UEIPAddress", "UEIPAddressPoolIdentity", "URRID", "URSEQN", "UsageInformation", "UsageReportTrigger", "UserPlaneNodeManagementInformationContainer", "ValidityTimer", "ValueDigits", "VolumeMeasurement", "VolumeQuota", "VolumeThreshold", "Weight"}
	for _, tc := range catalogueGrouped() {
		parent := tc.newIE(ie.NewPDRID(42), ie.NewNetworkInstance("example"))
		for _, name := range methods {
			method := reflect.ValueOf(parent).MethodByName(name)
			if !method.IsValid() {
				t.Fatal("missing method", name)
			}
			args := make([]reflect.Value, method.Type().NumIn())
			for j := range args {
				args[j] = reflect.Zero(method.Type().In(j))
			}
			results := method.Call(args)
			last := results[len(results)-1]
			if last.Type() == reflect.TypeFor[error]() {
				var invalid *ie.InvalidTypeError
				if last.IsNil() || !errors.As(last.Interface().(error), &invalid) || invalid.Type != tc.typ {
					t.Fatalf("%s accepted %s", name, tc.typ)
				}
			} else if last.Kind() == reflect.Bool && last.Bool() {
				t.Fatalf("%s accepted %s", name, tc.typ)
			}
		}
	}
}
