// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// SessionSetDeletionResponse is a SessionSetDeletionResponse formed PFCP Header and its IEs above.
type SessionSetDeletionResponse struct {
	*Header
	NodeID      *ie.IE
	Cause       *ie.IE
	OffendingIE *ie.IE
	IEs         []*ie.IE
}

// NewSessionSetDeletionResponse creates a new SessionSetDeletionResponse.
func NewSessionSetDeletionResponse(seq uint32, id, cause, offending *ie.IE, ies ...*ie.IE) *SessionSetDeletionResponse {
	m := &SessionSetDeletionResponse{
		Header:      NewHeader(MsgTypeSessionSetDeletionResponse, 0, seq, nil),
		NodeID:      id,
		Cause:       cause,
		OffendingIE: offending,
		IEs:         ies,
	}
	m.SetLength()

	return m
}

// Marshal returns the byte sequence generated from a SessionSetDeletionResponse.
func (m *SessionSetDeletionResponse) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *SessionSetDeletionResponse) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseSessionSetDeletionResponse decodes a given byte sequence as a SessionSetDeletionResponse.
func ParseSessionSetDeletionResponse(b []byte) (*SessionSetDeletionResponse, error) {
	m := &SessionSetDeletionResponse{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a SessionSetDeletionResponse.
func (m *SessionSetDeletionResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeSessionSetDeletionResponse)
	if err != nil {
		return err
	}
	*m = SessionSetDeletionResponse{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *SessionSetDeletionResponse) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *SessionSetDeletionResponse) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *SessionSetDeletionResponse) MessageTypeName() string {
	return "Node Report Response"
}

// SEID returns the SEID in uint64.
func (m *SessionSetDeletionResponse) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *SessionSetDeletionResponse) IsRequest() bool {
	return false
}

func (m *SessionSetDeletionResponse) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 3+len(m.IEs))
	ies = append(ies, m.NodeID)
	ies = append(ies, m.Cause)
	ies = append(ies, m.OffendingIE)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *SessionSetDeletionResponse) addIEs(ies []*ie.IE) {
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
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
