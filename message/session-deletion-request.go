// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// SessionDeletionRequest is a SessionDeletionRequest formed PFCP Header and its IEs above.
type SessionDeletionRequest struct {
	*Header
	IEs []*ie.IE
}

// NewSessionDeletionRequest creates a new SessionDeletionRequest.
func NewSessionDeletionRequest(seid uint64, seq uint32, ies ...*ie.IE) *SessionDeletionRequest {
	m := &SessionDeletionRequest{
		Header: NewHeader(MsgTypeSessionDeletionRequest, seid, seq, nil),
		IEs:    ies,
	}
	m.SetLength()

	return m
}

// Marshal returns the byte sequence generated from a SessionDeletionRequest.
func (m *SessionDeletionRequest) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *SessionDeletionRequest) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseSessionDeletionRequest decodes a given byte sequence as a SessionDeletionRequest.
func ParseSessionDeletionRequest(b []byte) (*SessionDeletionRequest, error) {
	m := &SessionDeletionRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a SessionDeletionRequest.
func (m *SessionDeletionRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeSessionDeletionRequest)
	if err != nil {
		return err
	}
	*m = SessionDeletionRequest{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.IEs = append(m.IEs, ies...)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *SessionDeletionRequest) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *SessionDeletionRequest) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *SessionDeletionRequest) MessageTypeName() string {
	return "Session Deletion Request"
}

// SEID returns the SEID in uint64.
func (m *SessionDeletionRequest) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *SessionDeletionRequest) IsRequest() bool {
	return true
}

func (m *SessionDeletionRequest) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, len(m.IEs))
	ies = append(ies, m.IEs...)
	return ies
}
