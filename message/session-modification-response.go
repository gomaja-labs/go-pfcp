// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// SessionModificationResponse is a SessionModificationResponse formed PFCP Header and its IEs above.
//
// TODO: add Packet Rate Status Report IE.
//
// TODO: rename CreatedBridgeInfoForTSC => TSCManagementInformation
type SessionModificationResponse struct {
	*Header
	Cause                             *ie.IE
	OffendingIE                       *ie.IE
	CreatedPDR                        []*ie.IE
	LoadControlInformation            *ie.IE
	OverloadControlInformation        *ie.IE
	UsageReport                       []*ie.IE
	FailedRuleID                      *ie.IE
	AdditionalUsageReportsInformation *ie.IE
	CreatedUpdatedTrafficEndpoint     []*ie.IE
	CreatedBridgeInfoForTSC           *ie.IE
	ATSSSControlParameters            *ie.IE
	UpdatedPDR                        []*ie.IE
	PacketRateStatusReport            []*ie.IE
	IEs                               []*ie.IE
}

// NewSessionModificationResponse creates a new SessionModificationResponse.
func NewSessionModificationResponse(seid uint64, seq uint32, ies ...*ie.IE) *SessionModificationResponse {
	m := &SessionModificationResponse{
		Header: NewHeader(MsgTypeSessionModificationResponse, seid, seq, nil),
	}

	m.addIEs(ies)

	m.SetLength()
	return m
}

// Marshal returns the byte sequence generated from a SessionModificationResponse.
func (m *SessionModificationResponse) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *SessionModificationResponse) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseSessionModificationResponse decodes a given byte sequence as a SessionModificationResponse.
func ParseSessionModificationResponse(b []byte) (*SessionModificationResponse, error) {
	m := &SessionModificationResponse{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a SessionModificationResponse.
func (m *SessionModificationResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeSessionModificationResponse)
	if err != nil {
		return err
	}
	*m = SessionModificationResponse{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *SessionModificationResponse) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *SessionModificationResponse) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *SessionModificationResponse) MessageTypeName() string {
	return "Session Modification Response"
}

// SEID returns the SEID in uint64.
func (m *SessionModificationResponse) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *SessionModificationResponse) IsRequest() bool {
	return false
}

func (m *SessionModificationResponse) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 8+len(m.CreatedPDR)+len(m.UsageReport)+len(m.CreatedUpdatedTrafficEndpoint)+len(m.UpdatedPDR)+len(m.PacketRateStatusReport)+len(m.IEs))
	ies = append(ies, m.Cause)
	ies = append(ies, m.OffendingIE)
	ies = append(ies, m.CreatedPDR...)
	ies = append(ies, m.LoadControlInformation)
	ies = append(ies, m.OverloadControlInformation)
	ies = append(ies, m.UsageReport...)
	ies = append(ies, m.FailedRuleID)
	ies = append(ies, m.AdditionalUsageReportsInformation)
	ies = append(ies, m.CreatedUpdatedTrafficEndpoint...)
	ies = append(ies, m.CreatedBridgeInfoForTSC)
	ies = append(ies, m.ATSSSControlParameters)
	ies = append(ies, m.UpdatedPDR...)
	ies = append(ies, m.PacketRateStatusReport...)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *SessionModificationResponse) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.Cause:
			assignIE(&m.Cause, i, &m.IEs)
		case ie.OffendingIE:
			assignIE(&m.OffendingIE, i, &m.IEs)
		case ie.CreatedPDR:
			m.CreatedPDR = append(m.CreatedPDR, i)
		case ie.LoadControlInformation:
			assignIE(&m.LoadControlInformation, i, &m.IEs)
		case ie.OverloadControlInformation:
			assignIE(&m.OverloadControlInformation, i, &m.IEs)
		case ie.UsageReportWithinSessionModificationResponse:
			m.UsageReport = append(m.UsageReport, i)
		case ie.FailedRuleID:
			assignIE(&m.FailedRuleID, i, &m.IEs)
		case ie.AdditionalUsageReportsInformation:
			assignIE(&m.AdditionalUsageReportsInformation, i, &m.IEs)
		case ie.CreatedTrafficEndpoint:
			m.CreatedUpdatedTrafficEndpoint = append(m.CreatedUpdatedTrafficEndpoint, i)
		case ie.CreatedBridgeInfoForTSC:
			assignIE(&m.CreatedBridgeInfoForTSC, i, &m.IEs)
		case ie.ATSSSControlParameters:
			assignIE(&m.ATSSSControlParameters, i, &m.IEs)
		case ie.UpdatedPDR:
			m.UpdatedPDR = append(m.UpdatedPDR, i)
		case ie.PacketRateStatusReport:
			m.PacketRateStatusReport = append(m.PacketRateStatusReport, i)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
