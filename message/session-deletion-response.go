// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// SessionDeletionResponse is a SessionDeletionResponse formed PFCP Header and its IEs above.
//
// TODO: add Packet Rate Status Report and Session Report IE.
type SessionDeletionResponse struct {
	*Header
	Cause                             *ie.IE
	OffendingIE                       *ie.IE
	LoadControlInformation            *ie.IE
	OverloadControlInformation        *ie.IE
	UsageReport                       []*ie.IE
	AdditionalUsageReportsInformation *ie.IE
	IEs                               []*ie.IE
}

// NewSessionDeletionResponse creates a new SessionDeletionResponse.
func NewSessionDeletionResponse(seid uint64, seq uint32, ies ...*ie.IE) *SessionDeletionResponse {
	m := &SessionDeletionResponse{
		Header: NewHeader(MsgTypeSessionDeletionResponse, seid, seq, nil),
	}

	m.addIEs(ies)

	m.SetLength()
	return m
}

// Marshal returns the byte sequence generated from a SessionDeletionResponse.
func (m *SessionDeletionResponse) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *SessionDeletionResponse) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseSessionDeletionResponse decodes a given byte sequence as a SessionDeletionResponse.
func ParseSessionDeletionResponse(b []byte) (*SessionDeletionResponse, error) {
	m := &SessionDeletionResponse{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a SessionDeletionResponse.
func (m *SessionDeletionResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeSessionDeletionResponse)
	if err != nil {
		return err
	}
	*m = SessionDeletionResponse{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *SessionDeletionResponse) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *SessionDeletionResponse) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *SessionDeletionResponse) MessageTypeName() string {
	return "Session Deletion Response"
}

// SEID returns the SEID in uint64.
func (m *SessionDeletionResponse) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *SessionDeletionResponse) IsRequest() bool {
	return false
}

func (m *SessionDeletionResponse) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 5+len(m.UsageReport)+len(m.IEs))
	ies = append(ies, m.Cause)
	ies = append(ies, m.OffendingIE)
	ies = append(ies, m.LoadControlInformation)
	ies = append(ies, m.OverloadControlInformation)
	ies = append(ies, m.UsageReport...)
	ies = append(ies, m.AdditionalUsageReportsInformation)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *SessionDeletionResponse) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.Cause:
			assignIE(&m.Cause, i, &m.IEs)
		case ie.OffendingIE:
			assignIE(&m.OffendingIE, i, &m.IEs)
		case ie.LoadControlInformation:
			assignIE(&m.LoadControlInformation, i, &m.IEs)
		case ie.OverloadControlInformation:
			assignIE(&m.OverloadControlInformation, i, &m.IEs)
		case ie.UsageReportWithinSessionDeletionResponse:
			m.UsageReport = append(m.UsageReport, i)
		case ie.AdditionalUsageReportsInformation:
			assignIE(&m.AdditionalUsageReportsInformation, i, &m.IEs)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
