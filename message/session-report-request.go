// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// SessionReportRequest is a SessionReportRequest formed PFCP Header and its IEs above.
//
// TODO: rename PortManagementInformationForTSC => TSCManagementInformation
type SessionReportRequest struct {
	*Header
	ReportType                        *ie.IE
	DownlinkDataReport                *ie.IE
	UsageReport                       []*ie.IE
	ErrorIndicationReport             *ie.IE
	LoadControlInformation            *ie.IE
	OverloadControlInformation        *ie.IE
	AdditionalUsageReportsInformation *ie.IE
	PFCPSRReqFlags                    *ie.IE
	OldCPFSEID                        *ie.IE
	PacketRateStatusReport            *ie.IE
	PortManagementInformationForTSC   *ie.IE
	SessionReport                     []*ie.IE
	IEs                               []*ie.IE
}

// NewSessionReportRequest creates a new SessionReportRequest.
func NewSessionReportRequest(seid uint64, seq uint32, ies ...*ie.IE) *SessionReportRequest {
	m := &SessionReportRequest{
		Header: NewHeader(MsgTypeSessionReportRequest, seid, seq, nil),
	}

	m.addIEs(ies)

	m.SetLength()
	return m
}

// Marshal returns the byte sequence generated from a SessionReportRequest.
func (m *SessionReportRequest) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *SessionReportRequest) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseSessionReportRequest decodes a given byte sequence as a SessionReportRequest.
func ParseSessionReportRequest(b []byte) (*SessionReportRequest, error) {
	m := &SessionReportRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a SessionReportRequest.
func (m *SessionReportRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeSessionReportRequest)
	if err != nil {
		return err
	}
	*m = SessionReportRequest{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *SessionReportRequest) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *SessionReportRequest) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *SessionReportRequest) MessageTypeName() string {
	return "Session Report Request"
}

// SEID returns the SEID in uint64.
func (m *SessionReportRequest) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *SessionReportRequest) IsRequest() bool {
	return true
}

func (m *SessionReportRequest) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 10+len(m.UsageReport)+len(m.SessionReport)+len(m.IEs))
	ies = append(ies, m.ReportType)
	ies = append(ies, m.DownlinkDataReport)
	ies = append(ies, m.UsageReport...)
	ies = append(ies, m.ErrorIndicationReport)
	ies = append(ies, m.LoadControlInformation)
	ies = append(ies, m.OverloadControlInformation)
	ies = append(ies, m.AdditionalUsageReportsInformation)
	ies = append(ies, m.PFCPSRReqFlags)
	ies = append(ies, m.OldCPFSEID)
	ies = append(ies, m.PacketRateStatusReport)
	ies = append(ies, m.PortManagementInformationForTSC)
	ies = append(ies, m.SessionReport...)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *SessionReportRequest) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.ReportType:
			assignIE(&m.ReportType, i, &m.IEs)
		case ie.DownlinkDataReport:
			assignIE(&m.DownlinkDataReport, i, &m.IEs)
		case ie.UsageReportWithinSessionReportRequest:
			m.UsageReport = append(m.UsageReport, i)
		case ie.ErrorIndicationReport:
			assignIE(&m.ErrorIndicationReport, i, &m.IEs)
		case ie.LoadControlInformation:
			assignIE(&m.LoadControlInformation, i, &m.IEs)
		case ie.OverloadControlInformation:
			assignIE(&m.OverloadControlInformation, i, &m.IEs)
		case ie.AdditionalUsageReportsInformation:
			assignIE(&m.AdditionalUsageReportsInformation, i, &m.IEs)
		case ie.PFCPSRReqFlags:
			assignIE(&m.PFCPSRReqFlags, i, &m.IEs)
		case ie.FSEID:
			assignIE(&m.OldCPFSEID, i, &m.IEs)
		case ie.PacketRateStatusReport:
			assignIE(&m.PacketRateStatusReport, i, &m.IEs)
		case ie.PortManagementInformationForTSCWithinSessionReportRequest:
			assignIE(&m.PortManagementInformationForTSC, i, &m.IEs)
		case ie.SessionReport:
			m.SessionReport = append(m.SessionReport, i)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
