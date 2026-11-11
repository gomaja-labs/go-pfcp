// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// SessionModificationRequest is a SessionModificationRequest formed PFCP Header and its IEs above.
type SessionModificationRequest struct {
	*Header
	CPFSEID                        *ie.IE
	RemovePDR                      []*ie.IE
	RemoveFAR                      []*ie.IE
	RemoveURR                      []*ie.IE
	RemoveQER                      []*ie.IE
	RemoveBAR                      *ie.IE
	RemoveTrafficEndpoint          []*ie.IE
	CreatePDR                      []*ie.IE
	CreateFAR                      []*ie.IE
	CreateURR                      []*ie.IE
	CreateQER                      []*ie.IE
	CreateBAR                      *ie.IE
	CreateTrafficEndpoint          []*ie.IE
	UpdatePDR                      []*ie.IE
	UpdateFAR                      []*ie.IE
	UpdateURR                      []*ie.IE
	UpdateQER                      []*ie.IE
	UpdateBAR                      *ie.IE
	UpdateTrafficEndpoint          []*ie.IE
	PFCPSMReqFlags                 *ie.IE
	QueryURR                       []*ie.IE
	FQCSID                         *ie.IE
	UserPlaneInactivityTimer       *ie.IE
	QueryURRReference              *ie.IE
	TraceInformation               *ie.IE
	RemoveMAR                      []*ie.IE
	UpdateMAR                      []*ie.IE
	CreateMAR                      []*ie.IE
	NodeID                         *ie.IE
	TSCManagementInformation       *ie.IE
	RemoveSRR                      []*ie.IE
	CreateSRR                      []*ie.IE
	UpdateSRR                      []*ie.IE
	ProvideATSSSControlInformation *ie.IE
	EthernetContextInformation     *ie.IE
	AccessAvailabilityInformation  []*ie.IE
	QueryPacketRateStatus          []*ie.IE
	SNSSAI                         *ie.IE
	IEs                            []*ie.IE
}

// NewSessionModificationRequest creates a new SessionModificationRequest.
func NewSessionModificationRequest(seid uint64, seq uint32, ies ...*ie.IE) *SessionModificationRequest {
	m := &SessionModificationRequest{
		Header: NewHeader(MsgTypeSessionModificationRequest, seid, seq, nil),
	}

	m.addIEs(ies)

	m.SetLength()
	return m
}

// Marshal returns the byte sequence generated from a SessionModificationRequest.
func (m *SessionModificationRequest) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *SessionModificationRequest) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseSessionModificationRequest decodes a given byte sequence as a SessionModificationRequest.
func ParseSessionModificationRequest(b []byte) (*SessionModificationRequest, error) {
	m := &SessionModificationRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a SessionModificationRequest.
func (m *SessionModificationRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeSessionModificationRequest)
	if err != nil {
		return err
	}
	*m = SessionModificationRequest{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *SessionModificationRequest) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *SessionModificationRequest) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *SessionModificationRequest) MessageTypeName() string {
	return "Session Modification Request"
}

// SEID returns the SEID in uint64.
func (m *SessionModificationRequest) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *SessionModificationRequest) IsRequest() bool {
	return true
}

func (m *SessionModificationRequest) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 14+len(m.RemovePDR)+len(m.RemoveFAR)+len(m.RemoveURR)+len(m.RemoveQER)+len(m.RemoveTrafficEndpoint)+len(m.CreatePDR)+len(m.CreateFAR)+len(m.CreateURR)+len(m.CreateQER)+len(m.CreateTrafficEndpoint)+len(m.UpdatePDR)+len(m.UpdateFAR)+len(m.UpdateURR)+len(m.UpdateQER)+len(m.UpdateTrafficEndpoint)+len(m.QueryURR)+len(m.RemoveMAR)+len(m.UpdateMAR)+len(m.CreateMAR)+len(m.RemoveSRR)+len(m.CreateSRR)+len(m.UpdateSRR)+len(m.AccessAvailabilityInformation)+len(m.QueryPacketRateStatus)+len(m.IEs))
	ies = append(ies, m.CPFSEID)
	ies = append(ies, m.RemovePDR...)
	ies = append(ies, m.RemoveFAR...)
	ies = append(ies, m.RemoveURR...)
	ies = append(ies, m.RemoveQER...)
	ies = append(ies, m.RemoveBAR)
	ies = append(ies, m.RemoveTrafficEndpoint...)
	ies = append(ies, m.CreatePDR...)
	ies = append(ies, m.CreateFAR...)
	ies = append(ies, m.CreateURR...)
	ies = append(ies, m.CreateQER...)
	ies = append(ies, m.CreateBAR)
	ies = append(ies, m.CreateTrafficEndpoint...)
	ies = append(ies, m.UpdatePDR...)
	ies = append(ies, m.UpdateFAR...)
	ies = append(ies, m.UpdateURR...)
	ies = append(ies, m.UpdateQER...)
	ies = append(ies, m.UpdateBAR)
	ies = append(ies, m.UpdateTrafficEndpoint...)
	ies = append(ies, m.PFCPSMReqFlags)
	ies = append(ies, m.QueryURR...)
	ies = append(ies, m.FQCSID)
	ies = append(ies, m.UserPlaneInactivityTimer)
	ies = append(ies, m.QueryURRReference)
	ies = append(ies, m.TraceInformation)
	ies = append(ies, m.RemoveMAR...)
	ies = append(ies, m.UpdateMAR...)
	ies = append(ies, m.CreateMAR...)
	ies = append(ies, m.NodeID)
	ies = append(ies, m.TSCManagementInformation)
	ies = append(ies, m.RemoveSRR...)
	ies = append(ies, m.CreateSRR...)
	ies = append(ies, m.UpdateSRR...)
	ies = append(ies, m.ProvideATSSSControlInformation)
	ies = append(ies, m.EthernetContextInformation)
	ies = append(ies, m.AccessAvailabilityInformation...)
	ies = append(ies, m.QueryPacketRateStatus...)
	ies = append(ies, m.SNSSAI)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *SessionModificationRequest) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.FSEID:
			assignIE(&m.CPFSEID, i, &m.IEs)
		case ie.RemovePDR:
			m.RemovePDR = append(m.RemovePDR, i)
		case ie.RemoveFAR:
			m.RemoveFAR = append(m.RemoveFAR, i)
		case ie.RemoveURR:
			m.RemoveURR = append(m.RemoveURR, i)
		case ie.RemoveQER:
			m.RemoveQER = append(m.RemoveQER, i)
		case ie.RemoveBAR:
			assignIE(&m.RemoveBAR, i, &m.IEs)
		case ie.RemoveTrafficEndpoint:
			m.RemoveTrafficEndpoint = append(m.RemoveTrafficEndpoint, i)
		case ie.CreatePDR:
			m.CreatePDR = append(m.CreatePDR, i)
		case ie.CreateFAR:
			m.CreateFAR = append(m.CreateFAR, i)
		case ie.CreateURR:
			m.CreateURR = append(m.CreateURR, i)
		case ie.CreateQER:
			m.CreateQER = append(m.CreateQER, i)
		case ie.CreateBAR:
			assignIE(&m.CreateBAR, i, &m.IEs)
		case ie.CreateTrafficEndpoint:
			m.CreateTrafficEndpoint = append(m.CreateTrafficEndpoint, i)
		case ie.UpdatePDR:
			m.UpdatePDR = append(m.UpdatePDR, i)
		case ie.UpdateFAR:
			m.UpdateFAR = append(m.UpdateFAR, i)
		case ie.UpdateURR:
			m.UpdateURR = append(m.UpdateURR, i)
		case ie.UpdateQER:
			m.UpdateQER = append(m.UpdateQER, i)
		case ie.UpdateBARWithinSessionModificationRequest:
			assignIE(&m.UpdateBAR, i, &m.IEs)
		case ie.UpdateTrafficEndpoint:
			m.UpdateTrafficEndpoint = append(m.UpdateTrafficEndpoint, i)
		case ie.PFCPSMReqFlags:
			assignIE(&m.PFCPSMReqFlags, i, &m.IEs)
		case ie.QueryURR:
			m.QueryURR = append(m.QueryURR, i)
		case ie.FQCSID:
			assignIE(&m.FQCSID, i, &m.IEs)
		case ie.UserPlaneInactivityTimer:
			assignIE(&m.UserPlaneInactivityTimer, i, &m.IEs)
		case ie.QueryURRReference:
			assignIE(&m.QueryURRReference, i, &m.IEs)
		case ie.TraceInformation:
			assignIE(&m.TraceInformation, i, &m.IEs)
		case ie.RemoveMAR:
			m.RemoveMAR = append(m.RemoveMAR, i)
		case ie.UpdateMAR:
			m.UpdateMAR = append(m.UpdateMAR, i)
		case ie.CreateMAR:
			m.CreateMAR = append(m.CreateMAR, i)
		case ie.NodeID:
			assignIE(&m.NodeID, i, &m.IEs)
		case ie.TSCManagementInformationWithinSessionModificationRequest:
			assignIE(&m.TSCManagementInformation, i, &m.IEs)
		case ie.RemoveSRR:
			m.RemoveSRR = append(m.RemoveSRR, i)
		case ie.CreateSRR:
			m.CreateSRR = append(m.CreateSRR, i)
		case ie.UpdateSRR:
			m.UpdateSRR = append(m.UpdateSRR, i)
		case ie.ProvideATSSSControlInformation:
			assignIE(&m.ProvideATSSSControlInformation, i, &m.IEs)
		case ie.EthernetContextInformation:
			assignIE(&m.EthernetContextInformation, i, &m.IEs)
		case ie.AccessAvailabilityInformation:
			m.AccessAvailabilityInformation = append(m.AccessAvailabilityInformation, i)
		case ie.QueryPacketRateStatusWithinSessionModificationRequest:
			m.QueryPacketRateStatus = append(m.QueryPacketRateStatus, i)
		case ie.SNSSAI:
			assignIE(&m.SNSSAI, i, &m.IEs)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
