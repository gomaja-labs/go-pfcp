// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// SessionSetDeletionRequest is a SessionSetDeletionRequest formed PFCP Header and its IEs above.
type SessionSetDeletionRequest struct {
	*Header
	NodeID *ie.IE
	FQCSID *ie.IE
	IEs    []*ie.IE
}

// NewSessionSetDeletionRequest creates a new SessionSetDeletionRequest.
func NewSessionSetDeletionRequest(seq uint32, id, csid *ie.IE, ies ...*ie.IE) *SessionSetDeletionRequest {
	m := &SessionSetDeletionRequest{
		Header: NewHeader(MsgTypeSessionSetDeletionRequest, 0, seq, nil),
		NodeID: id,
		FQCSID: csid,
		IEs:    ies,
	}
	m.SetLength()

	return m
}

// Marshal returns the byte sequence generated from a SessionSetDeletionRequest.
func (m *SessionSetDeletionRequest) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *SessionSetDeletionRequest) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseSessionSetDeletionRequest decodes a given byte sequence as a SessionSetDeletionRequest.
func ParseSessionSetDeletionRequest(b []byte) (*SessionSetDeletionRequest, error) {
	m := &SessionSetDeletionRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a SessionSetDeletionRequest.
func (m *SessionSetDeletionRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeSessionSetDeletionRequest)
	if err != nil {
		return err
	}
	*m = SessionSetDeletionRequest{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *SessionSetDeletionRequest) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *SessionSetDeletionRequest) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *SessionSetDeletionRequest) MessageTypeName() string {
	return "Session Set Deletion Request"
}

// SEID returns the SEID in uint64.
func (m *SessionSetDeletionRequest) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *SessionSetDeletionRequest) IsRequest() bool {
	return true
}

func (m *SessionSetDeletionRequest) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 2+len(m.IEs))
	ies = append(ies, m.NodeID)
	ies = append(ies, m.FQCSID)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *SessionSetDeletionRequest) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.NodeID:
			assignIE(&m.NodeID, i, &m.IEs)
		case ie.FQCSID:
			assignIE(&m.FQCSID, i, &m.IEs)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
