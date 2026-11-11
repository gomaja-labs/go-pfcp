// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// HeartbeatRequest is a HeartbeatRequest formed PFCP Header and its IEs above.
type HeartbeatRequest struct {
	*Header
	RecoveryTimeStamp *ie.IE
	SourceIPAddress   *ie.IE
	IEs               []*ie.IE
}

// NewHeartbeatRequest creates a new HeartbeatRequest.
func NewHeartbeatRequest(seq uint32, ts, ip *ie.IE, ies ...*ie.IE) *HeartbeatRequest {
	m := &HeartbeatRequest{
		Header:            NewHeader(MsgTypeHeartbeatRequest, 0, seq, nil),
		RecoveryTimeStamp: ts,
		SourceIPAddress:   ip,
		IEs:               ies,
	}
	m.SetLength()

	return m
}

// Marshal returns the byte sequence generated from a HeartbeatRequest.
func (m *HeartbeatRequest) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *HeartbeatRequest) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseHeartbeatRequest decodes a given byte sequence as a HeartbeatRequest.
func ParseHeartbeatRequest(b []byte) (*HeartbeatRequest, error) {
	m := &HeartbeatRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a HeartbeatRequest.
func (m *HeartbeatRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeHeartbeatRequest)
	if err != nil {
		return err
	}
	*m = HeartbeatRequest{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *HeartbeatRequest) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *HeartbeatRequest) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *HeartbeatRequest) MessageTypeName() string {
	return "Heartbeat Request"
}

// SEID returns the SEID in uint64.
func (m *HeartbeatRequest) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *HeartbeatRequest) IsRequest() bool {
	return true
}

func (m *HeartbeatRequest) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 2+len(m.IEs))
	ies = append(ies, m.RecoveryTimeStamp)
	ies = append(ies, m.SourceIPAddress)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *HeartbeatRequest) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.RecoveryTimeStamp:
			assignIE(&m.RecoveryTimeStamp, i, &m.IEs)
		case ie.SourceIPAddress:
			assignIE(&m.SourceIPAddress, i, &m.IEs)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
