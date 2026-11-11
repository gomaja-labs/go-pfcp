// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// SessionEstablishmentRequest is a SessionEstablishmentRequest formed PFCP Header and its IEs above.
//
// TODO: add S-NSSAI and Provide RDS configuration information.
type SessionEstablishmentRequest struct {
	*Header
	NodeID                             *ie.IE
	CPFSEID                            *ie.IE
	CreatePDR                          []*ie.IE
	CreateFAR                          []*ie.IE
	CreateURR                          []*ie.IE
	CreateQER                          []*ie.IE
	CreateBAR                          *ie.IE
	CreateTrafficEndpoint              []*ie.IE
	PDNType                            *ie.IE
	FQCSID                             *ie.IE
	UserPlaneInactivityTimer           *ie.IE
	UserID                             *ie.IE
	TraceInformation                   *ie.IE
	APNDNN                             *ie.IE
	CreateMAR                          []*ie.IE
	PFCPSEReqFlags                     *ie.IE
	CreateBridgeInfoForTSC             *ie.IE
	CreateSRR                          []*ie.IE
	ProvideATSSSControlInformation     *ie.IE
	RecoveryTimeStamp                  *ie.IE
	SNSSAI                             *ie.IE
	ProvideRDSConfigurationInformation *ie.IE
	IEs                                []*ie.IE
}

// NewSessionEstablishmentRequest creates a new SessionEstablishmentRequest.
func NewSessionEstablishmentRequest(seid uint64, seq uint32, ies ...*ie.IE) *SessionEstablishmentRequest {
	m := &SessionEstablishmentRequest{
		Header: NewHeader(MsgTypeSessionEstablishmentRequest, seid, seq, nil),
	}

	m.addIEs(ies)

	m.SetLength()
	return m
}

// Marshal returns the byte sequence generated from a SessionEstablishmentRequest.
func (m *SessionEstablishmentRequest) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *SessionEstablishmentRequest) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseSessionEstablishmentRequest decodes a given byte sequence as a SessionEstablishmentRequest.
func ParseSessionEstablishmentRequest(b []byte) (*SessionEstablishmentRequest, error) {
	m := &SessionEstablishmentRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a SessionEstablishmentRequest.
func (m *SessionEstablishmentRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeSessionEstablishmentRequest)
	if err != nil {
		return err
	}
	*m = SessionEstablishmentRequest{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *SessionEstablishmentRequest) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *SessionEstablishmentRequest) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *SessionEstablishmentRequest) MessageTypeName() string {
	return "Session Establishment Request"
}

// SEID returns the SEID in uint64.
func (m *SessionEstablishmentRequest) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *SessionEstablishmentRequest) IsRequest() bool {
	return true
}

func (m *SessionEstablishmentRequest) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 15+len(m.CreatePDR)+len(m.CreateFAR)+len(m.CreateURR)+len(m.CreateQER)+len(m.CreateTrafficEndpoint)+len(m.CreateMAR)+len(m.CreateSRR)+len(m.IEs))
	ies = append(ies, m.NodeID)
	ies = append(ies, m.CPFSEID)
	ies = append(ies, m.CreatePDR...)
	ies = append(ies, m.CreateFAR...)
	ies = append(ies, m.CreateURR...)
	ies = append(ies, m.CreateQER...)
	ies = append(ies, m.CreateBAR)
	ies = append(ies, m.CreateTrafficEndpoint...)
	ies = append(ies, m.PDNType)
	ies = append(ies, m.FQCSID)
	ies = append(ies, m.UserPlaneInactivityTimer)
	ies = append(ies, m.UserID)
	ies = append(ies, m.TraceInformation)
	ies = append(ies, m.APNDNN)
	ies = append(ies, m.CreateMAR...)
	ies = append(ies, m.PFCPSEReqFlags)
	ies = append(ies, m.CreateBridgeInfoForTSC)
	ies = append(ies, m.CreateSRR...)
	ies = append(ies, m.ProvideATSSSControlInformation)
	ies = append(ies, m.RecoveryTimeStamp)
	ies = append(ies, m.SNSSAI)
	ies = append(ies, m.ProvideRDSConfigurationInformation)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *SessionEstablishmentRequest) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.NodeID:
			assignIE(&m.NodeID, i, &m.IEs)
		case ie.FSEID:
			assignIE(&m.CPFSEID, i, &m.IEs)
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
		case ie.PDNType:
			assignIE(&m.PDNType, i, &m.IEs)
		case ie.FQCSID:
			assignIE(&m.FQCSID, i, &m.IEs)
		case ie.UserPlaneInactivityTimer:
			assignIE(&m.UserPlaneInactivityTimer, i, &m.IEs)
		case ie.UserID:
			assignIE(&m.UserID, i, &m.IEs)
		case ie.TraceInformation:
			assignIE(&m.TraceInformation, i, &m.IEs)
		case ie.APNDNN:
			assignIE(&m.APNDNN, i, &m.IEs)
		case ie.CreateMAR:
			m.CreateMAR = append(m.CreateMAR, i)
		case ie.PFCPSEReqFlags:
			assignIE(&m.PFCPSEReqFlags, i, &m.IEs)
		case ie.CreateBridgeInfoForTSC:
			assignIE(&m.CreateBridgeInfoForTSC, i, &m.IEs)
		case ie.CreateSRR:
			m.CreateSRR = append(m.CreateSRR, i)
		case ie.ProvideATSSSControlInformation:
			assignIE(&m.ProvideATSSSControlInformation, i, &m.IEs)
		case ie.RecoveryTimeStamp:
			assignIE(&m.RecoveryTimeStamp, i, &m.IEs)
		case ie.SNSSAI:
			assignIE(&m.SNSSAI, i, &m.IEs)
		case ie.ProvideRDSConfigurationInformation:
			assignIE(&m.ProvideRDSConfigurationInformation, i, &m.IEs)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
