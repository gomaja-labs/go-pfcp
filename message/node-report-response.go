// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// NodeReportResponse is a NodeReportResponse formed PFCP Header and its IEs above.
type NodeReportResponse struct {
	*Header
	NodeID      *ie.IE
	Cause       *ie.IE
	OffendingIE *ie.IE
	IEs         []*ie.IE
}

// NewNodeReportResponse creates a new NodeReportResponse.
func NewNodeReportResponse(seq uint32, id, cause, offending *ie.IE, ies ...*ie.IE) *NodeReportResponse {
	m := &NodeReportResponse{
		Header:      NewHeader(MsgTypeNodeReportResponse, 0, seq, nil),
		NodeID:      id,
		Cause:       cause,
		OffendingIE: offending,
		IEs:         ies,
	}
	m.SetLength()

	return m
}

// Marshal returns the byte sequence generated from a NodeReportResponse.
func (m *NodeReportResponse) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *NodeReportResponse) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseNodeReportResponse decodes a given byte sequence as a NodeReportResponse.
func ParseNodeReportResponse(b []byte) (*NodeReportResponse, error) {
	m := &NodeReportResponse{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a NodeReportResponse.
func (m *NodeReportResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeNodeReportResponse)
	if err != nil {
		return err
	}
	*m = NodeReportResponse{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *NodeReportResponse) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *NodeReportResponse) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *NodeReportResponse) MessageTypeName() string {
	return "Node Report Response"
}

// SEID returns the SEID in uint64.
func (m *NodeReportResponse) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *NodeReportResponse) IsRequest() bool {
	return false
}

func (m *NodeReportResponse) messageIEs() []*ie.IE {
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

func (m *NodeReportResponse) addIEs(ies []*ie.IE) {
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
