// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// SessionReportResponse is a SessionReportResponse formed PFCP Header and its IEs above.
type SessionReportResponse struct {
	*Header
	Cause                   *ie.IE
	OffendingIE             *ie.IE
	UpdateBAR               *ie.IE
	PFCPSRRspFlags          *ie.IE
	CPFSEID                 *ie.IE
	N4UFTEID                *ie.IE
	AlternativeSMFIPAddress *ie.IE
	IEs                     []*ie.IE
}

// NewSessionReportResponse creates a new SessionReportResponse.
func NewSessionReportResponse(seid uint64, seq uint32, ies ...*ie.IE) *SessionReportResponse {
	m := &SessionReportResponse{
		Header: NewHeader(MsgTypeSessionReportResponse, seid, seq, nil),
	}

	m.addIEs(ies)

	m.SetLength()
	return m
}

// Marshal returns the byte sequence generated from a SessionReportResponse.
func (m *SessionReportResponse) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *SessionReportResponse) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseSessionReportResponse decodes a given byte sequence as a SessionReportResponse.
func ParseSessionReportResponse(b []byte) (*SessionReportResponse, error) {
	m := &SessionReportResponse{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a SessionReportResponse.
func (m *SessionReportResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeSessionReportResponse)
	if err != nil {
		return err
	}
	*m = SessionReportResponse{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *SessionReportResponse) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *SessionReportResponse) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *SessionReportResponse) MessageTypeName() string {
	return "Session Report Response"
}

// SEID returns the SEID in uint64.
func (m *SessionReportResponse) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *SessionReportResponse) IsRequest() bool {
	return false
}

func (m *SessionReportResponse) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 7+len(m.IEs))
	ies = append(ies, m.Cause)
	ies = append(ies, m.OffendingIE)
	ies = append(ies, m.UpdateBAR)
	ies = append(ies, m.PFCPSRRspFlags)
	ies = append(ies, m.CPFSEID)
	ies = append(ies, m.N4UFTEID)
	ies = append(ies, m.AlternativeSMFIPAddress)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *SessionReportResponse) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.Cause:
			assignIE(&m.Cause, i, &m.IEs)
		case ie.OffendingIE:
			assignIE(&m.OffendingIE, i, &m.IEs)
		case ie.UpdateBARWithinSessionReportResponse:
			assignIE(&m.UpdateBAR, i, &m.IEs)
		case ie.PFCPSRRspFlags:
			assignIE(&m.PFCPSRRspFlags, i, &m.IEs)
		case ie.FSEID:
			assignIE(&m.CPFSEID, i, &m.IEs)
		case ie.FTEID:
			assignIE(&m.N4UFTEID, i, &m.IEs)
		case ie.AlternativeSMFIPAddress:
			assignIE(&m.AlternativeSMFIPAddress, i, &m.IEs)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
