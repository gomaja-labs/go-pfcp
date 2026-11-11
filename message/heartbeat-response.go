// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// HeartbeatResponse is a HeartbeatResponse formed PFCP Header and its IEs above.
type HeartbeatResponse struct {
	*Header
	RecoveryTimeStamp *ie.IE
	IEs               []*ie.IE
}

// NewHeartbeatResponse creates a new HeartbeatResponse.
func NewHeartbeatResponse(seq uint32, ts *ie.IE, ies ...*ie.IE) *HeartbeatResponse {
	m := &HeartbeatResponse{
		Header:            NewHeader(MsgTypeHeartbeatResponse, 0, seq, nil),
		RecoveryTimeStamp: ts,
		IEs:               ies,
	}
	m.SetLength()

	return m
}

// Marshal returns the byte sequence generated from a HeartbeatResponse.
func (m *HeartbeatResponse) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *HeartbeatResponse) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseHeartbeatResponse decodes a given byte sequence as a HeartbeatResponse.
func ParseHeartbeatResponse(b []byte) (*HeartbeatResponse, error) {
	m := &HeartbeatResponse{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a HeartbeatResponse.
func (m *HeartbeatResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeHeartbeatResponse)
	if err != nil {
		return err
	}
	*m = HeartbeatResponse{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *HeartbeatResponse) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *HeartbeatResponse) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *HeartbeatResponse) MessageTypeName() string {
	return "Heartbeat Response"
}

// SEID returns the SEID in uint64.
func (m *HeartbeatResponse) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *HeartbeatResponse) IsRequest() bool {
	return false
}

func (m *HeartbeatResponse) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 1+len(m.IEs))
	ies = append(ies, m.RecoveryTimeStamp)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *HeartbeatResponse) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.RecoveryTimeStamp:
			assignIE(&m.RecoveryTimeStamp, i, &m.IEs)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
