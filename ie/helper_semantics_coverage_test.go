package ie_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// The owned source files are the nongrouped IE files in the type 1–273
// review, together with the generic value helpers in ie.go. Reflection
// enumerates the current exported method set; this list defines ownership,
// and never supplies the helper numerator or hides an unregistered method.
var helperOwnedSources = map[string]bool{
	"access-availability-information.go":      true,
	"activate-predefined-rules.go":            true,
	"activation-time.go":                      true,
	"additional-usage-reports-information.go": true,
	"aggregated-urr-id.go":                    true,
	"alternative-smf-ip-address.go":           true,
	"apn-dnn.go":                              true,
	"application-id.go":                       true,
	"application-instance-id.go":              true,
	"apply-action.go":                         true,
	"atsss-ll-control-information.go":         true,
	"atsss-ll-information.go":                 true,
	"average-packet-delay.go":                 true,
	"averaging-window.go":                     true,
	"bar-id.go":                               true,
	"c-tag.go":                                true,
	"cause.go":                                true,
	"cp-function-features.go":                 true,
	"cp-pfcp-entity-ip-address.go":            true,
	"create-bridge-router-info.go":            true,
	"cumulative-rate-ratio-measurement.go":    true,
	"cumulative-rate-ratio-threshold.go":      true,
	"data-network-access-identifier.go":       true,
	"data-status.go":                          true,
	"deactivate-predefined-rules.go":          true,
	"deactivation-time.go":                    true,
	"destination-interface.go":                true,
	"dl-buffering-duration.go":                true,
	"dl-buffering-suggested-packet-count.go":  true,
	"dl-data-packets-size.go":                 true,
	"dl-flow-level-marking.go":                true,
	"downlink-data-notification-delay.go":     true,
	"downlink-data-service-information.go":    true,
	"dropped-dl-traffic-threshold.go":         true,
	"duration-measurement.go":                 true,
	"end-time.go":                             true,
	"ethernet-filter-id.go":                   true,
	"ethernet-filter-properties.go":           true,
	"ethernet-inactivity-timer.go":            true,
	"ethernet-pdu-session-information.go":     true,
	"ethertype.go":                            true,
	"event-quota.go":                          true,
	"event-threshold.go":                      true,
	"event-time-stamp.go":                     true,
	"f-seid.go":                               true,
	"f-teid.go":                               true,
	"failed-rule-id.go":                       true,
	"far-id.go":                               true,
	"five-gs-user-plane-node-id.go":           true,
	"flow-information.go":                     true,
	"forwarding-policy.go":                    true,
	"fq-csid.go":                              true,
	"framed-ipv6-route.go":                    true,
	"framed-route.go":                         true,
	"framed-routing.go":                       true,
	"gate-status.go":                          true,
	"gbr.go":                                  true,
	"graceful-release-period.go":              true,
	"gtp-u-path-interface-type.go":            true,
	"header-enrichment.go":                    true,
	"ie.go":                                   true,
	"inactivity-detection-time.go":            true,
	"ip-multicast-address.go":                 true,
	"ip-version.go":                           true,
	"link-specific-multipath-ip-address.go":   true,
	"linked-urr-id.go":                        true,
	"mac-address-detected.go":                 true,
	"mac-address-removed.go":                  true,
	"mac-address.go":                          true,
	"mar-id.go":                               true,
	"maximum-packet-delay.go":                 true,
	"mbr.go":                                  true,
	"measurement-information.go":              true,
	"measurement-method.go":                   true,
	"measurement-period.go":                   true,
	"metric.go":                               true,
	"minimum-packet-delay.go":                 true,
	"minimum-wait-time.go":                    true,
	"monitoring-time.go":                      true,
	"mptcp-address-information.go":            true,
	"mptcp-control-information.go":            true,
	"mt-edt-control-information.go":           true,
	"multipath-applicable-indication.go":      true,
	"multiplier.go":                           true,
	"network-instance.go":                     true,
	"nf-instance-id.go":                       true,
	"node-id.go":                              true,
	"node-report-type.go":                     true,
	"number-of-reports.go":                    true,
	"number-of-ue-ip-addresses.go":            true,
	"nw-tt-port-number.go":                    true,
	"oci-flags.go":                            true,
	"offending-ie.go":                         true,
	"outer-header-creation.go":                true,
	"outer-header-removal.go":                 true,
	"packet-delay-thresholds.go":              true,
	"packet-rate-status.go":                   true,
	"packet-rate.go":                          true,
	"packet-replication-and-detection-carry-on-information.go": true,
	"paging-policy-indicator.go":                               true,
	"pdn-type.go":                                              true,
	"pdr-id.go":                                                true,
	"pfcp-association-release-request.go":                      true,
	"pfcpasreq-flags.go":                                       true,
	"pfcpasrsp-flags.go":                                       true,
	"pfcpaureq-flags.go":                                       true,
	"pfcpsereq-flags.go":                                       true,
	"pfcpsmreq-flags.go":                                       true,
	"pfcpsrreq-flags.go":                                       true,
	"pfcpsrrsp-flags.go":                                       true,
	"pfd-contents.go":                                          true,
	"pmf-address-information.go":                               true,
	"pmf-control-information.go":                               true,
	"port-management-information-container.go":                 true,
	"port-number.go":                                           true,
	"precedence.go":                                            true,
	"priority.go":                                              true,
	"proxying.go":                                              true,
	"qer-control-indications.go":                               true,
	"qer-correlation-id.go":                                    true,
	"qer-id.go":                                                true,
	"qfi.go":                                                   true,
	"qos-monitoring-measurement.go":                            true,
	"qos-report-trigger.go":                                    true,
	"query-urr-reference.go":                                   true,
	"quota-holding-time.go":                                    true,
	"quota-validity-time.go":                                   true,
	"rds-configuration-information.go":                         true,
	"recovery-time-stamp.go":                                   true,
	"redirect-information.go":                                  true,
	"remote-gtp-u-peer.go":                                     true,
	"report-type.go":                                           true,
	"reporting-frequency.go":                                   true,
	"reporting-triggers.go":                                    true,
	"requested-access-availability-information.go":             true,
	"requested-clock-drift-information.go":                     true,
	"requested-qos-monitoring.go":                              true,
	"rqi.go":                                                   true,
	"s-nssai.go":                                               true,
	"s-tag.go":                                                 true,
	"sdf-filter.go":                                            true,
	"sequence-number.go":                                       true,
	"smf-set-id.go":                                            true,
	"source-interface.go":                                      true,
	"source-ip-address.go":                                     true,
	"srr-id.go":                                                true,
	"start-time.go":                                            true,
	"steering-functionality.go":                                true,
	"steering-mode.go":                                         true,
	"subsequent-event-quota.go":                                true,
	"subsequent-event-threshold.go":                            true,
	"subsequent-time-quota.go":                                 true,
	"subsequent-time-threshold.go":                             true,
	"subsequent-volume-quota.go":                               true,
	"subsequent-volume-threshold.go":                           true,
	"suggested-buffering-packets-count.go":                     true,
	"tgpp-interface-type.go":                                   true,
	"time-domain-number.go":                                    true,
	"time-of-first-packet.go":                                  true,
	"time-of-last-packet.go":                                   true,
	"time-offset-measurement.go":                               true,
	"time-offset-threshold.go":                                 true,
	"time-quota-mechanism.go":                                  true,
	"time-quota.go":                                            true,
	"time-threshold.go":                                        true,
	"timer.go":                                                 true,
	"trace-information.go":                                     true,
	"traffic-endpoint-id.go":                                   true,
	"transport-level-marking.go":                               true,
	"ue-ip-address-pool-identity.go":                           true,
	"ue-ip-address.go":                                         true,
	"up-function-features.go":                                  true,
	"ur-seqn.go":                                               true,
	"urr-id.go":                                                true,
	"usage-information.go":                                     true,
	"usage-report-trigger.go":                                  true,
	"user-id.go":                                               true,
	"user-plane-inactivity-timer.go":                           true,
	"user-plane-node-management-information-container.go":      true,
	"validity-timer.go":                                        true,
	"volume-measurement.go":                                    true,
	"volume-quota.go":                                          true,
	"volume-threshold.go":                                      true,
	"weight.go":                                                true,
}
var helperSpecClauses = map[ie.IEType]string{
	19:  "8.2.1",
	20:  "8.2.2",
	21:  "8.2.3",
	22:  "8.2.4",
	23:  "8.2.5",
	24:  "8.2.6",
	25:  "8.2.7",
	26:  "8.2.8",
	27:  "8.2.9",
	28:  "8.2.10",
	29:  "8.2.11",
	30:  "8.2.12",
	31:  "8.2.13",
	32:  "8.2.14",
	33:  "8.2.15",
	34:  "8.2.16",
	35:  "8.2.17",
	36:  "8.2.18",
	37:  "8.2.19",
	39:  "8.2.21",
	40:  "8.2.22",
	41:  "8.2.23",
	42:  "8.2.24",
	43:  "8.2.25",
	44:  "8.2.26",
	45:  "8.2.27",
	46:  "8.2.28",
	47:  "8.2.29",
	48:  "8.2.30",
	49:  "8.2.31",
	50:  "8.2.32",
	52:  "8.2.33",
	53:  "8.2.34",
	55:  "8.2.35",
	56:  "8.2.36",
	57:  "8.2.37",
	60:  "8.2.38",
	61:  "8.2.39",
	62:  "8.2.40",
	63:  "8.2.41",
	64:  "8.2.42",
	65:  "8.2.43",
	66:  "8.2.44",
	67:  "8.2.45",
	69:  "8.2.46",
	70:  "8.2.47",
	71:  "8.2.48",
	72:  "8.2.49",
	73:  "8.2.50",
	74:  "8.2.51",
	75:  "8.2.52",
	76:  "8.2.53",
	81:  "8.2.54",
	82:  "8.2.55",
	84:  "8.2.56",
	88:  "8.2.57",
	89:  "8.2.58",
	90:  "8.2.59",
	91:  "8.2.60",
	92:  "8.2.61",
	93:  "8.2.62",
	94:  "8.2.63",
	95:  "8.2.64",
	96:  "8.2.65",
	97:  "8.2.66",
	98:  "8.2.67",
	100: "8.2.68",
	101: "8.2.69",
	103: "8.2.70",
	104: "8.2.71",
	106: "8.2.72",
	107: "8.2.73",
	108: "8.2.74",
	109: "8.2.75",
	110: "8.2.76",
	111: "8.2.77",
	112: "8.2.78",
	113: "8.2.79",
	114: "8.2.80",
	115: "8.2.81",
	117: "8.2.83",
	119: "8.2.84",
	120: "8.2.85",
	121: "8.2.86",
	122: "8.2.87",
	123: "8.2.88",
	124: "8.2.89",
	125: "8.2.90",
	126: "8.2.91",
	131: "8.2.92",
	133: "8.2.93",
	134: "8.2.94",
	135: "8.2.95",
	136: "8.2.96",
	137: "8.2.97",
	138: "8.2.98",
	139: "8.2.99",
	140: "8.2.100",
	141: "8.2.101",
	142: "8.2.102",
	144: "8.2.103",
	145: "8.2.104",
	146: "8.2.105",
	148: "8.2.112",
	149: "8.2.113",
	150: "8.2.106",
	151: "8.2.107",
	152: "8.2.108",
	153: "8.2.109",
	154: "8.2.110",
	155: "8.2.111",
	156: "8.2.114",
	157: "8.2.115",
	158: "8.2.116",
	159: "8.2.117",
	160: "8.2.118",
	161: "8.2.119",
	162: "8.2.120",
	163: "8.2.121",
	164: "8.2.122",
	170: "8.2.123",
	171: "8.2.124",
	172: "8.2.125",
	173: "8.2.126",
	174: "8.2.127",
	177: "8.2.128",
	178: "8.2.129",
	179: "8.2.130",
	180: "8.2.131",
	181: "8.2.132",
	182: "8.2.133",
	184: "8.2.134",
	185: "8.2.135",
	186: "8.2.136",
	191: "8.2.137",
	192: "8.2.138",
	193: "8.2.139",
	194: "8.2.140",
	196: "8.2.141",
	197: "8.2.142",
	198: "8.2.143",
	202: "8.2.144",
	204: "8.2.145",
	206: "8.2.146",
	207: "8.2.147",
	208: "8.2.148",
	209: "8.2.149",
	210: "8.2.150",
	215: "8.2.151",
	217: "8.2.152",
	219: "8.2.153",
	222: "8.2.154",
	223: "8.2.155",
	224: "8.2.156",
	228: "8.2.157",
	229: "8.2.158",
	230: "8.2.159",
	231: "8.2.160",
	232: "8.2.161",
	234: "8.2.162",
	235: "8.2.163",
	236: "8.2.164",
	237: "8.2.165",
	241: "8.2.166",
	243: "8.2.167",
	244: "8.2.168",
	245: "8.2.169",
	246: "8.2.170",
	248: "8.2.171",
	249: "8.2.172",
	250: "8.2.173",
	251: "8.2.174",
	253: "8.2.175",
	257: "8.2.176",
	258: "8.2.177",
	259: "8.2.178",
	260: "8.2.179",
	262: "8.2.180",
	265: "8.2.181",
	266: "8.2.182",
	268: "8.2.183",
	269: "8.2.184",
}

type helperCoverageRow struct {
	Method             string      `json:"method"`
	Source             string      `json:"source"`
	Signature          string      `json:"signature"`
	Owned              bool        `json:"owned"`
	DirectValue        bool        `json:"direct_value"`
	WrongType          bool        `json:"wrong_type"`
	NilReceiver        bool        `json:"nil_receiver"`
	Null               bool        `json:"null"`
	NonemptyTruncation bool        `json:"nonempty_truncation"`
	TruncationNote     string      `json:"truncation_note,omitempty"`
	Test               string      `json:"test,omitempty"`
	Contexts           []ie.IEType `json:"ie_types,omitempty"`
	Clauses            []string    `json:"clauses,omitempty"`
	Exclusion          string      `json:"exclusion,omitempty"`
}

func helperCoverageExclusions() map[string]string {
	return map[string]string{
		"HasMODE":  "TS 29.244 V20.1.0 clause 8.2.174 makes the old MODE bit spare; obsolete API awaiting removal.",
		"HasNORD":  "TS 29.244 V20.1.0 clause 8.2.174 makes the old NORD bit spare; obsolete API awaiting removal.",
		"HasSESRL": "TS 29.244 V20.1.0 clause 8.2.168 reserves the former SESRL bit; obsolete API awaiting removal.",
	}
}

// PFCP_HELPER_COVERAGE_JSON writes an optional machine-readable inventory.
// The report records registered semantic assertions, not their run status;
// the TestHelperSemantic* result is the pass/fail evidence for the matrix.
func TestHelperSemanticCoverageReport(t *testing.T) {
	rows := map[string]*helperCoverageRow{}
	typ := reflect.TypeFor[*ie.IE]()
	exclusions := helperCoverageExclusions()
	for n := 0; n < typ.NumMethod(); n++ {
		method := typ.Method(n)
		source, _ := runtime.FuncForPC(method.Func.Pointer()).FileLine(0)
		file := filepath.Base(source)
		row := &helperCoverageRow{Method: method.Name, Source: file, Signature: method.Type.String(), Owned: helperOwnedSources[file]}
		if method.Type.NumIn() != 1 {
			row.Exclusion = "Requires arguments; outside the no-argument value helper matrix."
		} else if reason := exclusions[method.Name]; reason != "" {
			row.Exclusion = reason
		} else if !row.Owned {
			row.Exclusion = "Outside the owned nongrouped type 1–273 source files (grouped API, new IE type, or general operation)."
		}
		rows[method.Name] = row
	}
	for _, tc := range helperSemanticValues() {
		row := rows[tc.name]
		if row == nil {
			t.Errorf("registered value helper missing: %s", tc.name)
			continue
		}
		row.DirectValue = true
		row.WrongType = true
		row.NilReceiver = true
		row.Null = true
		row.Test = "TestHelperSemanticValues + TestHelperSemanticContracts"
		row.NonemptyTruncation = tc.minimum > 1
		if tc.minimum <= 1 {
			row.TruncationNote = "No shorter nonempty prefix is malformed for this one-octet value or unframed text/octet string; null covers its empty prefix."
		}
		row.Contexts = append(row.Contexts, tc.input.Type)
		if clause := helperSpecClauses[tc.input.Type]; clause != "" {
			row.Clauses = append(row.Clauses, clause)
		}
	}
	for _, tc := range helperSemanticFlags() {
		row := rows[tc.name]
		if row == nil {
			t.Errorf("registered flag helper missing: %s", tc.name)
			continue
		}
		row.DirectValue = true
		row.WrongType = true
		row.NilReceiver = true
		row.Null = true
		row.Test = "TestHelperSemanticFlags"
		row.NonemptyTruncation = row.NonemptyTruncation || tc.octet > 0
		if !row.NonemptyTruncation {
			row.TruncationNote = "Flag view reads the first payload octet; no nonempty prefix precedes that octet. Full selected field layouts are tested by the corresponding value helper."
		} else {
			row.TruncationNote = "Payload ending before the selected flag octet returns false."
		}
		row.Contexts = append(row.Contexts, tc.typ)
		row.Clauses = append(row.Clauses, tc.clause)
	}
	for name, test := range map[string]string{"IsAllocatedByCPFunction": "TestHelperSemanticAllocation", "IsAllocatedByUPFunction": "TestHelperSemanticAllocation", "MustAPNDNN": "TestHelperSemanticMustAPNDNN", "ValidateApplyAction": "TestHelperSemanticValidateApplyAction"} {
		row := rows[name]
		if row == nil {
			t.Errorf("registered special helper missing: %s", name)
			continue
		}
		row.DirectValue = true
		row.WrongType = true
		row.NilReceiver = true
		row.Null = true
		row.Test = test
		row.NonemptyTruncation = name != "ValidateApplyAction"
		if !row.NonemptyTruncation {
			row.TruncationNote = "One required action octet: the only shorter prefix is null."
		}
	}
	for name, test := range map[string]string{"Marshal": "TestHelperSemanticHeaderOperations", "MarshalLen": "TestHelperSemanticHeaderOperations", "IsVendorSpecific": "TestHelperSemanticHeaderOperations", "ValueAsGrouped": "TestHelperSemanticValueAsGrouped"} {
		row := rows[name]
		row.DirectValue = true
		row.WrongType = true
		row.NilReceiver = true
		row.Null = true
		row.Test = test
		row.NonemptyTruncation = name == "ValueAsGrouped"
		row.TruncationNote = "General operation: every IE type is accepted; null values are valid. Header operations preserve opaque payload prefixes; ValueAsGrouped rejects incomplete child encodings."
	}

	var report struct {
		Scope            string              `json:"scope"`
		EvidenceBoundary string              `json:"evidence_boundary"`
		ReflectedMethods int                 `json:"reflected_exported_methods"`
		OwnedNoArg       int                 `json:"owned_noarg_helpers"`
		DirectHelpers    int                 `json:"direct_contract_helpers"`
		Excluded         int                 `json:"excluded_owned_noarg_helpers"`
		Missing          int                 `json:"missing_owned_noarg_helpers"`
		Rows             []helperCoverageRow `json:"helpers"`
	}
	report.Scope = "Nongrouped IE source files owned by the type 1–273 review and generic value methods from ie.go. Fields methods and top-level functions are outside this *IE method inventory."
	report.EvidenceBoundary = "Direct expected-value assertions plus wrong type, nil, null and applicable malformed nonempty prefixes. Counts identify registered assertions, not pass status. Selected IE type contexts only; the report does not claim every alternate supported type, every value, wire oracle equivalence, or complete protocol conformance."
	for n := 0; n < typ.NumMethod(); n++ {
		method := typ.Method(n)
		row := rows[method.Name]
		report.ReflectedMethods++
		if row.Owned && method.Type.NumIn() == 1 {
			report.OwnedNoArg++
			if row.DirectValue {
				report.DirectHelpers++
			} else if row.Exclusion != "" {
				report.Excluded++
			} else {
				report.Missing++
				t.Errorf("owned helper has no direct semantic matrix: %s (%s)", row.Method, row.Source)
			}
		}
		report.Rows = append(report.Rows, *row)
	}
	t.Logf("reflection: %d exported methods; owned %d no-argument helpers; %d direct; %d excluded; %d missing", report.ReflectedMethods, report.OwnedNoArg, report.DirectHelpers, report.Excluded, report.Missing)
	if path := os.Getenv("PFCP_HELPER_COVERAGE_JSON"); path != "" {
		data, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(path, append(data, '\n'), 0600); err != nil {
			t.Fatal(err)
		}
	}
}
