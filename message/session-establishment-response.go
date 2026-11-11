// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// SessionEstablishmentResponse is a SessionEstablishmentResponse formed PFCP Header and its IEs above.
//
// TODO: add RDS configuration information IE.
type SessionEstablishmentResponse struct {
	*Header
	NodeID                     *ie.IE
	Cause                      *ie.IE
	OffendingIE                *ie.IE
	UPFSEID                    *ie.IE
	CreatedPDR                 []*ie.IE
	LoadControlInformation     *ie.IE
	OverloadControlInformation *ie.IE
	FQCSID                     *ie.IE
	FailedRuleID               *ie.IE
	CreatedTrafficEndpoint     []*ie.IE
	CreatedBridgeInfoForTSC    *ie.IE
	ATSSSControlParameters     *ie.IE
	IEs                        []*ie.IE
}

// NewSessionEstablishmentResponse creates a new SessionEstablishmentResponse.
func NewSessionEstablishmentResponse(seid uint64, seq uint32, ies ...*ie.IE) *SessionEstablishmentResponse {
	m := &SessionEstablishmentResponse{
		Header: NewHeader(MsgTypeSessionEstablishmentResponse, seid, seq, nil),
	}

	m.addIEs(ies)

	m.SetLength()
	return m
}

// Marshal returns the byte sequence generated from a SessionEstablishmentResponse.
func (m *SessionEstablishmentResponse) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *SessionEstablishmentResponse) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseSessionEstablishmentResponse decodes a given byte sequence as a SessionEstablishmentResponse.
func ParseSessionEstablishmentResponse(b []byte) (*SessionEstablishmentResponse, error) {
	m := &SessionEstablishmentResponse{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a SessionEstablishmentResponse.
func (m *SessionEstablishmentResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeSessionEstablishmentResponse)
	if err != nil {
		return err
	}
	*m = SessionEstablishmentResponse{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *SessionEstablishmentResponse) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *SessionEstablishmentResponse) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *SessionEstablishmentResponse) MessageTypeName() string {
	return "Session Establishment Response"
}

// SEID returns the SEID in uint64.
func (m *SessionEstablishmentResponse) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *SessionEstablishmentResponse) IsRequest() bool {
	return false
}

func (m *SessionEstablishmentResponse) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 10+len(m.CreatedPDR)+len(m.CreatedTrafficEndpoint)+len(m.IEs))
	ies = append(ies, m.NodeID)
	ies = append(ies, m.Cause)
	ies = append(ies, m.OffendingIE)
	ies = append(ies, m.UPFSEID)
	ies = append(ies, m.CreatedPDR...)
	ies = append(ies, m.LoadControlInformation)
	ies = append(ies, m.OverloadControlInformation)
	ies = append(ies, m.FQCSID)
	ies = append(ies, m.FailedRuleID)
	ies = append(ies, m.CreatedTrafficEndpoint...)
	ies = append(ies, m.CreatedBridgeInfoForTSC)
	ies = append(ies, m.ATSSSControlParameters)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *SessionEstablishmentResponse) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.NodeID:
			assignIE(&m.NodeID, i, &m.IEs)
		case ie.Cause:
			assignIE(&m.Cause, i, &m.IEs)
		case ie.OffendingIE:
			assignIE(&m.OffendingIE, i, &m.IEs)
		case ie.FSEID:
			assignIE(&m.UPFSEID, i, &m.IEs)
		case ie.CreatedPDR:
			m.CreatedPDR = append(m.CreatedPDR, i)
		case ie.LoadControlInformation:
			assignIE(&m.LoadControlInformation, i, &m.IEs)
		case ie.OverloadControlInformation:
			assignIE(&m.OverloadControlInformation, i, &m.IEs)
		case ie.FQCSID:
			assignIE(&m.FQCSID, i, &m.IEs)
		case ie.FailedRuleID:
			assignIE(&m.FailedRuleID, i, &m.IEs)
		case ie.CreatedTrafficEndpoint:
			m.CreatedTrafficEndpoint = append(m.CreatedTrafficEndpoint, i)
		case ie.CreatedBridgeInfoForTSC:
			assignIE(&m.CreatedBridgeInfoForTSC, i, &m.IEs)
		case ie.ATSSSControlParameters:
			assignIE(&m.ATSSSControlParameters, i, &m.IEs)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
