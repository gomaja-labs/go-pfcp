package ie_test

import (
	"bytes"
	"errors"
	"net/netip"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// FuzzAddedIEValues exercises typed helpers and Fields parsers with arbitrary
// input, then checks that their known-field encoding is stable.
func FuzzAddedIEValues(f *testing.F) {
	for _, b := range [][]byte{{0}, {1, 255}, {3, 0, 0, 0}, {255, 255, 255, 255, 255}} {
		f.Add(b)
	}
	// CP IP Address, TS 29.244 V20.1.0 clause 8.2.199: receive is
	// tolerant of a mapped IPv6 value that typed send validation rejects.
	f.Add(append([]byte{16, 1}, netip.MustParseAddr("::ffff:192.0.2.1").AsSlice()...))
	names := []string{"L2TPUserAuthentication", "LNSAddress", "IPAddressAndPortNumberReplacement", "DSCPToPPIMappingInformation", "OffendingIEInformation", "RATType", "TunnelPreference", "CallingNumber", "CalledNumber", "L2TPSessionIndications", "DNSServerAddress", "NBNSServerAddress", "MaximumReceiveUnit", "Thresholds", "SteeringModeIndicator", "GroupID", "CPIPAddress", "DNSQueryResponseFilter", "EventNotificationURI", "NotificationCorrelationID", "ReportingFlags", "PredefinedRulesName", "MBSSessionIdentifier", "MulticastTransportInformation", "MBSN4mbReqFlags", "LocalIngressTunnel", "MBSUnicastParametersID", "MBSN4RespFlags", "TunnelPassword", "AreaSessionID", "PFCPSDRspFlags", "QERIndications", "VendorSpecificNodeReportType", "ConfiguredTimeDomain", "Metadata", "TrafficParameterThreshold", "DLPeriodicity", "N6JitterMeasurement", "TrafficParameterMeasurementIndication", "ULPeriodicity", "MPQUICControlInformation", "MPQUICAddressInformation", "TransportMode", "ReportingSuggestionInfo", "TLContainer", "MeasurementIndication", "HPLMNSNSSAI", "MediaTransportProtocol", "RTPHeaderExtensionType", "RTPHeaderExtensionID"}
	types := []ie.IEType{ie.L2TPUserAuthentication, ie.LNSAddress, ie.IPAddressAndPortNumberReplacement, ie.DSCPToPPIMappingInformation, ie.OffendingIEInformation, ie.RATType, ie.TunnelPreference, ie.CallingNumber, ie.CalledNumber, ie.L2TPSessionIndications, ie.DNSServerAddress, ie.NBNSServerAddress, ie.MaximumReceiveUnit, ie.Thresholds, ie.SteeringModeIndicator, ie.GroupID, ie.CPIPAddress, ie.DNSQueryResponseFilter, ie.EventNotificationURI, ie.NotificationCorrelationID, ie.ReportingFlags, ie.PredefinedRulesName, ie.MBSSessionIdentifier, ie.MulticastTransportInformation, ie.MBSN4mbReqFlags, ie.LocalIngressTunnel, ie.MBSUnicastParametersID, ie.MBSN4RespFlags, ie.TunnelPassword, ie.AreaSessionID, ie.PFCPSDRspFlags, ie.QERIndications, ie.VendorSpecificNodeReportType, ie.ConfiguredTimeDomain, ie.Metadata, ie.TrafficParameterThreshold, ie.DLPeriodicity, ie.N6JitterMeasurement, ie.TrafficParameterMeasurementIndication, ie.ULPeriodicity, ie.MPQUICControlInformation, ie.MPQUICAddressInformation, ie.TransportMode, ie.ReportingSuggestionInfo, ie.TLContainer, ie.MeasurementIndication, ie.HPLMNSNSSAI, ie.MediaTransportProtocol, ie.RTPHeaderExtensionType, ie.RTPHeaderExtensionID}
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) == 0 {
			return
		}
		n := int(data[0]) % len(names)
		payload := data[1:]
		i := ie.New(types[n], payload)
		_, _ = addedGetter(i, names[n])
		_, _ = i.Marshal()
		_ = i.MarshalTo(nil)
		var parse func([]byte) (addedFields, error)
		switch types[n] {
		case ie.L2TPUserAuthentication:
			parse = func(b []byte) (addedFields, error) { return ie.ParseL2TPUserAuthenticationFields(b) }
		case ie.IPAddressAndPortNumberReplacement:
			parse = func(b []byte) (addedFields, error) { return ie.ParseIPAddressAndPortNumberReplacementFields(b) }
		case ie.DSCPToPPIMappingInformation:
			parse = func(b []byte) (addedFields, error) { return ie.ParseDSCPToPPIMappingInformationFields(b) }
		case ie.OffendingIEInformation:
			parse = func(b []byte) (addedFields, error) { return ie.ParseOffendingIEInformationFields(b) }
		case ie.Thresholds:
			parse = func(b []byte) (addedFields, error) { return ie.ParseThresholdsFields(b) }
		case ie.CPIPAddress:
			parse = func(b []byte) (addedFields, error) { return ie.ParseCPIPAddressFields(b) }
		case ie.DNSQueryResponseFilter:
			parse = func(b []byte) (addedFields, error) { return ie.ParseDNSQueryResponseFilterFields(b) }
		case ie.MBSSessionIdentifier:
			parse = func(b []byte) (addedFields, error) { return ie.ParseMBSSessionIdentifierFields(b) }
		case ie.MulticastTransportInformation:
			parse = func(b []byte) (addedFields, error) { return ie.ParseMulticastTransportInformationFields(b) }
		case ie.LocalIngressTunnel:
			parse = func(b []byte) (addedFields, error) { return ie.ParseLocalIngressTunnelFields(b) }
		case ie.VendorSpecificNodeReportType:
			parse = func(b []byte) (addedFields, error) { return ie.ParseVendorSpecificNodeReportTypeFields(b) }
		case ie.TrafficParameterThreshold:
			parse = func(b []byte) (addedFields, error) { return ie.ParseTrafficParameterThresholdFields(b) }
		case ie.N6JitterMeasurement:
			parse = func(b []byte) (addedFields, error) { return ie.ParseN6JitterMeasurementFields(b) }
		case ie.MPQUICAddressInformation:
			parse = func(b []byte) (addedFields, error) { return ie.ParseMPQUICAddressInformationFields(b) }
		case ie.ReportingSuggestionInfo:
			parse = func(b []byte) (addedFields, error) { return ie.ParseReportingSuggestionInfoFields(b) }
		case ie.HPLMNSNSSAI:
			parse = func(b []byte) (addedFields, error) { return ie.ParseHPLMNSNSSAIFields(b) }
		}
		if parse == nil {
			return
		}
		v, err := parse(payload)
		if err != nil {
			return
		}
		b, err := v.Marshal()
		// A mapped address is valid on receive, but binary typed encoders
		// reject it on send. Keep checking every other decoded value.
		fields := reflect.ValueOf(v).Elem()
		for n := 0; n < fields.NumField(); n++ {
			if addr, ok := fields.Field(n).Interface().(netip.Addr); ok && addr.Is4In6() {
				if !errors.Is(err, ie.ErrInvalidValue) {
					t.Fatalf("mapped binary send: %v", err)
				}
				return
			}
		}
		if err != nil {
			t.Fatalf("successful decoder produced invalid fields: %v", err)
		}
		again, err := parse(b)
		if err != nil || !reflect.DeepEqual(v, again) {
			t.Fatalf("known fields changed: %v", err)
		}
		c, err := again.Marshal()
		if err != nil || !bytes.Equal(b, c) {
			t.Fatalf("known encoding changed: %v", err)
		}
		_ = v.UnmarshalBinary(nil)
		_ = v.UnmarshalBinary(payload)
	})
}
