// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/gomaja-labs/go-pfcp/internal/utils"
)

//go:generate stringer -type=IEType
type IEType uint16

// IE Type definitions.
const (
	_                                                                IEType = 0
	CreatePDR                                                        IEType = 1
	PDI                                                              IEType = 2
	CreateFAR                                                        IEType = 3
	ForwardingParameters                                             IEType = 4
	DuplicatingParameters                                            IEType = 5
	CreateURR                                                        IEType = 6
	CreateQER                                                        IEType = 7
	CreatedPDR                                                       IEType = 8
	UpdatePDR                                                        IEType = 9
	UpdateFAR                                                        IEType = 10
	UpdateForwardingParameters                                       IEType = 11
	UpdateBARWithinSessionReportResponse                             IEType = 12
	UpdateURR                                                        IEType = 13
	UpdateQER                                                        IEType = 14
	RemovePDR                                                        IEType = 15
	RemoveFAR                                                        IEType = 16
	RemoveURR                                                        IEType = 17
	RemoveQER                                                        IEType = 18
	Cause                                                            IEType = 19
	SourceInterface                                                  IEType = 20
	FTEID                                                            IEType = 21
	NetworkInstance                                                  IEType = 22
	SDFFilter                                                        IEType = 23
	ApplicationID                                                    IEType = 24
	GateStatus                                                       IEType = 25
	MBR                                                              IEType = 26
	GBR                                                              IEType = 27
	QERCorrelationID                                                 IEType = 28
	Precedence                                                       IEType = 29
	TransportLevelMarking                                            IEType = 30
	VolumeThreshold                                                  IEType = 31
	TimeThreshold                                                    IEType = 32
	MonitoringTime                                                   IEType = 33
	SubsequentVolumeThreshold                                        IEType = 34
	SubsequentTimeThreshold                                          IEType = 35
	InactivityDetectionTime                                          IEType = 36
	ReportingTriggers                                                IEType = 37
	RedirectInformation                                              IEType = 38
	ReportType                                                       IEType = 39
	OffendingIE                                                      IEType = 40
	ForwardingPolicy                                                 IEType = 41
	DestinationInterface                                             IEType = 42
	UPFunctionFeatures                                               IEType = 43
	ApplyAction                                                      IEType = 44
	DownlinkDataServiceInformation                                   IEType = 45
	DownlinkDataNotificationDelay                                    IEType = 46
	DLBufferingDuration                                              IEType = 47
	DLBufferingSuggestedPacketCount                                  IEType = 48
	PFCPSMReqFlags                                                   IEType = 49
	PFCPSRRspFlags                                                   IEType = 50
	LoadControlInformation                                           IEType = 51
	SequenceNumber                                                   IEType = 52
	Metric                                                           IEType = 53
	OverloadControlInformation                                       IEType = 54
	Timer                                                            IEType = 55
	PDRID                                                            IEType = 56
	FSEID                                                            IEType = 57
	ApplicationIDsPFDs                                               IEType = 58
	PFDContext                                                       IEType = 59
	NodeID                                                           IEType = 60
	PFDContents                                                      IEType = 61
	MeasurementMethod                                                IEType = 62
	UsageReportTrigger                                               IEType = 63
	MeasurementPeriod                                                IEType = 64
	FQCSID                                                           IEType = 65
	VolumeMeasurement                                                IEType = 66
	DurationMeasurement                                              IEType = 67
	ApplicationDetectionInformation                                  IEType = 68
	TimeOfFirstPacket                                                IEType = 69
	TimeOfLastPacket                                                 IEType = 70
	QuotaHoldingTime                                                 IEType = 71
	DroppedDLTrafficThreshold                                        IEType = 72
	VolumeQuota                                                      IEType = 73
	TimeQuota                                                        IEType = 74
	StartTime                                                        IEType = 75
	EndTime                                                          IEType = 76
	QueryURR                                                         IEType = 77
	UsageReportWithinSessionModificationResponse                     IEType = 78
	UsageReportWithinSessionDeletionResponse                         IEType = 79
	UsageReportWithinSessionReportRequest                            IEType = 80
	URRID                                                            IEType = 81
	LinkedURRID                                                      IEType = 82
	DownlinkDataReport                                               IEType = 83
	OuterHeaderCreation                                              IEType = 84
	CreateBAR                                                        IEType = 85
	UpdateBARWithinSessionModificationRequest                        IEType = 86
	RemoveBAR                                                        IEType = 87
	BARID                                                            IEType = 88
	CPFunctionFeatures                                               IEType = 89
	UsageInformation                                                 IEType = 90
	ApplicationInstanceID                                            IEType = 91
	FlowInformation                                                  IEType = 92
	UEIPAddress                                                      IEType = 93
	PacketRate                                                       IEType = 94
	OuterHeaderRemoval                                               IEType = 95
	RecoveryTimeStamp                                                IEType = 96
	DLFlowLevelMarking                                               IEType = 97
	HeaderEnrichment                                                 IEType = 98
	ErrorIndicationReport                                            IEType = 99
	MeasurementInformation                                           IEType = 100
	NodeReportType                                                   IEType = 101
	UserPlanePathFailureReport                                       IEType = 102
	RemoteGTPUPeer                                                   IEType = 103
	URSEQN                                                           IEType = 104
	UpdateDuplicatingParameters                                      IEType = 105
	ActivatePredefinedRules                                          IEType = 106
	DeactivatePredefinedRules                                        IEType = 107
	FARID                                                            IEType = 108
	QERID                                                            IEType = 109
	OCIFlags                                                         IEType = 110
	PFCPAssociationReleaseRequest                                    IEType = 111
	GracefulReleasePeriod                                            IEType = 112
	PDNType                                                          IEType = 113
	FailedRuleID                                                     IEType = 114
	TimeQuotaMechanism                                               IEType = 115
	UserPlaneIPResourceInformation                                   IEType = 116
	UserPlaneInactivityTimer                                         IEType = 117
	AggregatedURRs                                                   IEType = 118
	Multiplier                                                       IEType = 119
	AggregatedURRID                                                  IEType = 120
	SubsequentVolumeQuota                                            IEType = 121
	SubsequentTimeQuota                                              IEType = 122
	RQI                                                              IEType = 123
	QFI                                                              IEType = 124
	QueryURRReference                                                IEType = 125
	AdditionalUsageReportsInformation                                IEType = 126
	CreateTrafficEndpoint                                            IEType = 127
	CreatedTrafficEndpoint                                           IEType = 128
	UpdateTrafficEndpoint                                            IEType = 129
	RemoveTrafficEndpoint                                            IEType = 130
	TrafficEndpointID                                                IEType = 131
	EthernetPacketFilter                                             IEType = 132
	MACAddress                                                       IEType = 133
	CTAG                                                             IEType = 134
	STAG                                                             IEType = 135
	Ethertype                                                        IEType = 136
	Proxying                                                         IEType = 137
	EthernetFilterID                                                 IEType = 138
	EthernetFilterProperties                                         IEType = 139
	SuggestedBufferingPacketsCount                                   IEType = 140
	UserID                                                           IEType = 141
	EthernetPDUSessionInformation                                    IEType = 142
	EthernetTrafficInformation                                       IEType = 143
	MACAddressesDetected                                             IEType = 144
	MACAddressesRemoved                                              IEType = 145
	EthernetInactivityTimer                                          IEType = 146
	AdditionalMonitoringTime                                         IEType = 147
	EventQuota                                                       IEType = 148
	EventThreshold                                                   IEType = 149
	SubsequentEventQuota                                             IEType = 150
	SubsequentEventThreshold                                         IEType = 151
	TraceInformation                                                 IEType = 152
	FramedRoute                                                      IEType = 153
	FramedRouting                                                    IEType = 154
	FramedIPv6Route                                                  IEType = 155
	EventTimeStamp                                                   IEType = 156
	AveragingWindow                                                  IEType = 157
	PagingPolicyIndicator                                            IEType = 158
	APNDNN                                                           IEType = 159
	TGPPInterfaceType                                                IEType = 160
	PFCPSRReqFlags                                                   IEType = 161
	PFCPAUReqFlags                                                   IEType = 162
	ActivationTime                                                   IEType = 163
	DeactivationTime                                                 IEType = 164
	CreateMAR                                                        IEType = 165
	TGPPAccessForwardingActionInformation                            IEType = 166
	NonTGPPAccessForwardingActionInformation                         IEType = 167
	RemoveMAR                                                        IEType = 168
	UpdateMAR                                                        IEType = 169
	MARID                                                            IEType = 170
	SteeringFunctionality                                            IEType = 171
	SteeringMode                                                     IEType = 172
	Weight                                                           IEType = 173
	Priority                                                         IEType = 174
	UpdateTGPPAccessForwardingActionInformation                      IEType = 175
	UpdateNonTGPPAccessForwardingActionInformation                   IEType = 176
	UEIPAddressPoolIdentity                                          IEType = 177
	AlternativeSMFIPAddress                                          IEType = 178
	PacketReplicationAndDetectionCarryOnInformation                  IEType = 179
	SMFSetID                                                         IEType = 180
	QuotaValidityTime                                                IEType = 181
	NumberOfReports                                                  IEType = 182
	PFCPSessionRetentionInformation                                  IEType = 183
	PFCPASRspFlags                                                   IEType = 184
	CPPFCPEntityIPAddress                                            IEType = 185
	PFCPSEReqFlags                                                   IEType = 186
	UserPlanePathRecoveryReport                                      IEType = 187
	IPMulticastAddressingInfo                                        IEType = 188
	JoinIPMulticastInformationWithinUsageReport                      IEType = 189
	LeaveIPMulticastInformationWithinUsageReport                     IEType = 190
	IPMulticastAddress                                               IEType = 191
	SourceIPAddress                                                  IEType = 192
	PacketRateStatus                                                 IEType = 193
	CreateBridgeInfoForTSC                                           IEType = 194
	CreatedBridgeInfoForTSC                                          IEType = 195
	DSTTPortNumber                                                   IEType = 196
	NWTTPortNumber                                                   IEType = 197
	TSNBridgeID                                                      IEType = 198
	TSCManagementInformationWithinSessionModificationRequest         IEType = 199
	TSCManagementInformationWithinSessionModificationResponse        IEType = 200
	TSCManagementInformationWithinSessionReportRequest               IEType = 201
	PortManagementInformationForTSCWithinSessionModificationRequest  IEType = 199 // Deprecated
	PortManagementInformationForTSCWithinSessionModificationResponse IEType = 200 // Deprecated
	PortManagementInformationForTSCWithinSessionReportRequest        IEType = 201 // Deprecated
	PortManagementInformationContainer                               IEType = 202
	ClockDriftControlInformation                                     IEType = 203
	RequestedClockDriftInformation                                   IEType = 204
	ClockDriftReport                                                 IEType = 205
	TSNTimeDomainNumber                                              IEType = 206
	TimeOffsetThreshold                                              IEType = 207
	CumulativeRateRatioThreshold                                     IEType = 208
	TimeOffsetMeasurement                                            IEType = 209
	CumulativeRateRatioMeasurement                                   IEType = 210
	RemoveSRR                                                        IEType = 211
	CreateSRR                                                        IEType = 212
	UpdateSRR                                                        IEType = 213
	SessionReport                                                    IEType = 214
	SRRID                                                            IEType = 215
	AccessAvailabilityControlInformation                             IEType = 216
	RequestedAccessAvailabilityInformation                           IEType = 217
	AccessAvailabilityReport                                         IEType = 218
	AccessAvailabilityInformation                                    IEType = 219
	ProvideATSSSControlInformation                                   IEType = 220
	ATSSSControlParameters                                           IEType = 221
	MPTCPControlInformation                                          IEType = 222
	ATSSSLLControlInformation                                        IEType = 223
	PMFControlInformation                                            IEType = 224
	MPTCPParameters                                                  IEType = 225
	ATSSSLLParameters                                                IEType = 226
	PMFParameters                                                    IEType = 227
	MPTCPAddressInformation                                          IEType = 228
	UELinkSpecificIPAddress                                          IEType = 229
	PMFAddressInformation                                            IEType = 230
	ATSSSLLInformation                                               IEType = 231
	DataNetworkAccessIdentifier                                      IEType = 232
	UEIPAddressPoolInformation                                       IEType = 233
	AveragePacketDelay                                               IEType = 234
	MinimumPacketDelay                                               IEType = 235
	MaximumPacketDelay                                               IEType = 236
	QoSReportTrigger                                                 IEType = 237
	GTPUPathQoSControlInformation                                    IEType = 238
	GTPUPathQoSReport                                                IEType = 239
	QoSInformationInGTPUPathQoSReport                                IEType = 240
	GTPUPathInterfaceType                                            IEType = 241
	QoSMonitoringPerQoSFlowControlInformation                        IEType = 242
	RequestedQoSMonitoring                                           IEType = 243
	ReportingFrequency                                               IEType = 244
	PacketDelayThresholds                                            IEType = 245
	MinimumWaitTime                                                  IEType = 246
	QoSMonitoringReport                                              IEType = 247
	QoSMonitoringMeasurement                                         IEType = 248
	MTEDTControlInformation                                          IEType = 249
	DLDataPacketsSize                                                IEType = 250
	QERControlIndications                                            IEType = 251
	PacketRateStatusReport                                           IEType = 252
	NFInstanceID                                                     IEType = 253
	EthernetContextInformation                                       IEType = 254
	RedundantTransmissionParameters                                  IEType = 255
	UpdatedPDR                                                       IEType = 256
	SNSSAI                                                           IEType = 257
	IPVersion                                                        IEType = 258
	PFCPASReqFlags                                                   IEType = 259
	DataStatus                                                       IEType = 260
	ProvideRDSConfigurationInformation                               IEType = 261
	RDSConfigurationInformation                                      IEType = 262
	QueryPacketRateStatusWithinSessionModificationRequest            IEType = 263
	PacketRateStatusReportWithinSessionModificationResponse          IEType = 264
	MPTCPApplicableIndication                                        IEType = 265
	BridgeManagementInformationContainer                             IEType = 266
	UEIPAddressUsageInformation                                      IEType = 267
	NumberOfUEIPAddresses                                            IEType = 268
	ValidityTimer                                                    IEType = 269
	RedundantTransmissionForwardingParameters                        IEType = 270
	TransportDelayReporting                                          IEType = 271
)

// IE represents an Information Element of PFCP messages.
type IE struct {
	Type         IEType
	EnterpriseID uint16
	Payload      []byte
	ChildIEs     []*IE
}

// New creates a new IE.
func New(itype IEType, data []byte) *IE {
	i := &IE{
		Type:    itype,
		Payload: data,
	}

	return i
}

// NewVendorSpecificIE creates a new vendor-specific IE.
func NewVendorSpecificIE(itype IEType, eid uint16, data []byte) *IE {
	i := &IE{
		Type:         itype,
		EnterpriseID: eid,
		Payload:      data,
	}

	return i
}

// NewGroupedIE creates a new grouped IE.
func NewGroupedIE(itype IEType, ies ...*IE) *IE {
	return newGroupedIE(itype, 0, ies...)
}

// NewVendorSpecificGroupedIE creates a new grouped IE.
func NewVendorSpecificGroupedIE(itype IEType, eid uint16, ies ...*IE) *IE {
	return newGroupedIE(itype, eid, ies...)
}

// NewUint8ValIE creates a new IE with uint8 value.
func NewUint8IE(itype IEType, v uint8) *IE {
	return newUint8ValIE(itype, v)
}

// NewUint16ValIE creates a new IE with uint16 value.
func NewUint16IE(itype IEType, v uint16) *IE {
	return newUint16ValIE(itype, v)
}

// NewUint32ValIE creates a new IE with uint32 value.
func NewUint32IE(itype IEType, v uint32) *IE {
	return newUint32ValIE(itype, v)
}

// NewUint64ValIE creates a new IE with uint64 value.
func NewUint64IE(itype IEType, v uint64) *IE {
	return newUint64ValIE(itype, v)
}

// NewStringIE creates a new IE with string value.
func NewStringIE(itype IEType, v string) *IE {
	return newStringIE(itype, v)
}

// NewFQDNIE creates a new IE with FQDN value.
func NewFQDNIE(itype IEType, v string) *IE {
	return newFQDNIE(itype, v)
}

// ValueAsUint8 returns the value of IE as uint8.
func (i *IE) ValueAsUint8() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	if i.IsGrouped() {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 1 {
		return 0, io.ErrUnexpectedEOF
	}

	return i.Payload[0], nil
}

// ValueAsUint16 returns the value of IE as uint16.
func (i *IE) ValueAsUint16() (uint16, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	if i.IsGrouped() {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	return binary.BigEndian.Uint16(i.Payload[0:2]), nil
}

// ValueAsUint32 returns the value of IE as uint32.
func (i *IE) ValueAsUint32() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	if i.IsGrouped() {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 4 {
		return 0, io.ErrUnexpectedEOF
	}

	return binary.BigEndian.Uint32(i.Payload[0:4]), nil
}

// ValueAsUint64 returns the value of IE as uint64.
func (i *IE) ValueAsUint64() (uint64, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	if i.IsGrouped() {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 8 {
		return 0, io.ErrUnexpectedEOF
	}

	return binary.BigEndian.Uint64(i.Payload[0:8]), nil
}

// ValueAsString returns the value of IE as string.
func (i *IE) ValueAsString() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}

	if i.IsGrouped() {
		return "", &InvalidTypeError{Type: i.Type}
	}

	return string(i.Payload), nil
}

// ValueAsFQDN returns the value of IE as string, decoded as FQDN.
func (i *IE) ValueAsFQDN() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}

	if i.IsGrouped() {
		return "", &InvalidTypeError{Type: i.Type}
	}

	return decodeFQDNValue(i.Payload)
}

// ValueAsGrouped returns the value of IE as grouped IE.
//
// This method returns the ChildIEs field if it is already parsed.
// Otherwise, it parses the Payload field and returns the result.
//
// Any IE, including a vendor-specific IE, can be explicitly interpreted as a group.
// The type-specific grouped helpers use this same source of truth.
func (i *IE) ValueAsGrouped() ([]*IE, error) {
	if i == nil {
		return nil, ErrIENotFound
	}
	if i.ChildIEs != nil {
		for n, child := range i.ChildIEs {
			if child != nil {
				continue
			}
			// Direct struct construction can introduce nil children. Encoding
			// skips them; helpers must use the same list without mutating it.
			children := make([]*IE, 0, len(i.ChildIEs)-1)
			children = append(children, i.ChildIEs[:n]...)
			for _, next := range i.ChildIEs[n+1:] {
				if next != nil {
					children = append(children, next)
				}
			}
			return children, nil
		}
		return i.ChildIEs, nil
	}
	return ParseMultiIEs(i.Payload)
}

// valueAs3GPPTimestamp uses era 1 for wire seconds below 2208988800, and
// era 0 otherwise. This selects the Unix window [1970-01-01, 2106-02-07T06:28:15Z].
// The wire carries no era, so earlier/later dates cannot be distinguished.
// RFC 5905 section 6; TS 29.244 V20.1.0 clause 8.2.65.
// It returns the value of IE as time.Time, in the format of
// timestamp IEs defined in 3GPP TS 29.244 as follows:
//
// "a UTC time. Octets 5 to 8 shall be encoded in the same format as the first four octets
// of the 64-bit timestamp format as defined in clause 6 of IETF RFC 5905 [12].
// NOTE: The encoding is defined as the time in seconds relative to 00:00:00 on 1 January 1900."
func (i *IE) valueAs3GPPTimestamp() (time.Time, error) {
	if i == nil {
		return time.Time{}, ErrIENotFound
	}

	if i.IsGrouped() {
		return time.Time{}, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 4 {
		return time.Time{}, io.ErrUnexpectedEOF
	}

	return time.Unix(int64(binary.BigEndian.Uint32(i.Payload[0:4])-2208988800), 0), nil
}

// Parse parses b into IE.
//
// Note that this function uses the given bytes directly, so not safe to use
// the buffer after calling this function. When you use the buffer somewhere
// else, copy it before calling this function.
func Parse(b []byte) (*IE, error) {
	i := &IE{}
	if err := i.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return i, nil
}

// ParseMultiIEs decodes multiple IEs at a time.
// This is easy and useful but slower than decoding one by one.
// When you don't know the number of IEs, this is the only way to decode them.
// See benchmarks in diameter_test.go for the detail.
//
// Note that this function uses the given bytes directly, so not safe to use
// the buffer after calling this function. When you use the buffer somewhere
// else, copy it before calling this function.
func ParseMultiIEs(b []byte) ([]*IE, error) {
	var ies []*IE
	for {
		if len(b) == 0 {
			break
		}

		i, err := Parse(b)
		if err != nil {
			return nil, err
		}
		ies = append(ies, i)
		b = b[i.MarshalLen():]
	}
	return ies, nil
}

// UnmarshalBinary parses b into IE.
// TS 29.244 V20.1.0 clause 8.1.1: Length excludes the four-octet header,
// but includes the Enterprise ID of a vendor-specific IE.
func (i *IE) UnmarshalBinary(b []byte) error {
	if i == nil {
		return ErrIENotFound
	}
	if len(b) < 4 {
		return io.ErrUnexpectedEOF
	}
	typ := IEType(binary.BigEndian.Uint16(b[:2]))
	length := int(binary.BigEndian.Uint16(b[2:4]))
	vendor := typ&0x8000 != 0
	if vendor && length < 2 {
		return ErrInvalidLength
	}
	end := 4 + length
	if len(b) < end {
		return io.ErrUnexpectedEOF
	}
	v := IE{Type: typ}
	offset := 4
	if vendor {
		v.EnterpriseID = binary.BigEndian.Uint16(b[4:6])
		offset = 6
	}
	payload := b[offset:end:end]
	if v.IsGrouped() {
		children, err := ParseMultiIEs(payload)
		if err != nil {
			return fmt.Errorf("grouped IE %d: %w", typ, err)
		}
		v.ChildIEs = children
		if v.ChildIEs == nil {
			v.ChildIEs = make([]*IE, 0)
		}
	} else {
		v.Payload = payload
	}
	*i = v
	return nil
}

// Marshal returns the byte sequence generated from an IE instance.
func (i *IE) Marshal() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}
	length := i.MarshalLen()
	if length-4 > 65535 {
		return nil, ErrInvalidLength
	}
	b := make([]byte, length)
	if err := i.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo writes an IE to b. It returns ErrInvalidLength for a value exceeding
// the 16-bit Length field or for a buffer shorter than MarshalLen().
func (i *IE) MarshalTo(b []byte) error {
	if i == nil {
		return ErrIENotFound
	}
	length := i.MarshalLen()
	if length-4 > 65535 || len(b) < length {
		return ErrInvalidLength
	}
	binary.BigEndian.PutUint16(b[:2], uint16(i.Type))
	binary.BigEndian.PutUint16(b[2:4], uint16(length-4))
	offset := 4
	if i.IsVendorSpecific() {
		binary.BigEndian.PutUint16(b[4:6], i.EnterpriseID)
		offset = 6
	}
	if i.ChildIEs != nil {
		for _, child := range i.ChildIEs {
			if child == nil {
				continue
			}
			if err := child.MarshalTo(b[offset:]); err != nil {
				return fmt.Errorf("child of IE %d: %w", i.Type, err)
			}
			offset += child.MarshalLen()
		}
	} else {
		copy(b[offset:length], i.Payload)
	}
	return nil
}

// MarshalLen returns the encoded length, including the IE header, or zero for nil.
func (i *IE) MarshalLen() int {
	if i == nil {
		return 0
	}
	length := 4
	if i.IsVendorSpecific() {
		length += 2
	}
	if i.ChildIEs != nil {
		for _, child := range i.ChildIEs {
			length += child.MarshalLen()
		}
	} else {
		length += len(i.Payload)
	}
	return length
}

// IsVendorSpecific reports whether an IE is vendor-specific or defined by 3gpp.
func (i *IE) IsVendorSpecific() bool {
	if i == nil {
		return false
	}

	// Spef: TS 29.244 8.1.1 Information Element Format
	// If the Bit 8 of Octet 1 is not set, this indicates that the IE is defined by 3GPP
	// and the EnterpriseID is absent. If Bit 8 of Octet 1 is set, this indicates that the
	// IE is defined by a vendor and the Enterprise ID is present identified by the Enterprise ID.
	return i.Type&0x8000 != 0
}

func newUint8ValIE(t IEType, v uint8) *IE {
	return New(t, []byte{v})
}

func newUint16ValIE(t IEType, v uint16) *IE {
	i := New(t, make([]byte, 2))
	binary.BigEndian.PutUint16(i.Payload, v)
	return i
}

func newUint32ValIE(t IEType, v uint32) *IE {
	i := New(t, make([]byte, 4))
	binary.BigEndian.PutUint32(i.Payload, v)
	return i
}

func newUint64ValIE(t IEType, v uint64) *IE {
	i := New(t, make([]byte, 8))
	binary.BigEndian.PutUint64(i.Payload, v)
	return i
}

func newStringIE(t IEType, v string) *IE {
	return New(t, []byte(v))
}

func newFQDNIE(t IEType, v string) *IE {
	encoded := utils.EncodeFQDN(v)
	if encoded == nil {
		return nil
	}
	return New(t, encoded)
}

func newGroupedIE(itype IEType, eid uint16, ies ...*IE) *IE {
	i := &IE{Type: itype, EnterpriseID: eid, ChildIEs: make([]*IE, 0, len(ies))}
	for _, child := range ies {
		if child != nil {
			i.ChildIEs = append(i.ChildIEs, child)
		}
	}
	return i
}

// new3GPPTimestampIE applies the era policy used by valueAs3GPPTimestamp.
// TS 29.244 V20.1.0 clause 8.2.65 and RFC 5905 section 6 specify seconds
// since 1900 modulo 2^32. With no era carried on the wire, this package chooses
// the Unix window [1970-01-01, 2106-02-07T06:28:15Z]. Times outside that window
// return nil. Fractional seconds are discarded.
func new3GPPTimestampIE(typ IEType, ts time.Time) *IE {
	seconds := ts.Unix()
	if seconds < 0 || seconds > 1<<32-1 {
		return nil
	}
	return newUint32ValIE(typ, uint32(seconds+2208988800))
}

func decodeFQDNValue(payload []byte) (string, error) {
	name := utils.DecodeFQDN(payload)
	if name == "" && !(len(payload) == 1 && payload[0] == 0) {
		return "", ErrMalformed
	}
	return name, nil
}
