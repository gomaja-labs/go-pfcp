// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// NodeReportRequest is a NodeReportRequest formed PFCP Header and its IEs above.
type NodeReportRequest struct {
	*Header
	NodeID                      *ie.IE
	NodeReportType              *ie.IE
	UserPlanePathFailureReport  *ie.IE
	UserPlanePathRecoveryReport *ie.IE
	ClockDriftReport            []*ie.IE
	GTPUPathQoSReport           []*ie.IE
	IEs                         []*ie.IE
}

// NewNodeReportRequest creates a new NodeReportRequest.
func NewNodeReportRequest(seq uint32, ies ...*ie.IE) *NodeReportRequest {
	m := &NodeReportRequest{
		Header: NewHeader(MsgTypeNodeReportRequest, 0, seq, nil),
	}

	m.addIEs(ies)

	m.SetLength()
	return m
}

// Marshal returns the byte sequence generated from a NodeReportRequest.
func (m *NodeReportRequest) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *NodeReportRequest) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseNodeReportRequest decodes a given byte sequence as a NodeReportRequest.
func ParseNodeReportRequest(b []byte) (*NodeReportRequest, error) {
	m := &NodeReportRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a NodeReportRequest.
func (m *NodeReportRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeNodeReportRequest)
	if err != nil {
		return err
	}
	*m = NodeReportRequest{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *NodeReportRequest) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *NodeReportRequest) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *NodeReportRequest) MessageTypeName() string {
	return "Node Report Request"
}

// SEID returns the SEID in uint64.
func (m *NodeReportRequest) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *NodeReportRequest) IsRequest() bool {
	return true
}

func (m *NodeReportRequest) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 4+len(m.ClockDriftReport)+len(m.GTPUPathQoSReport)+len(m.IEs))
	ies = append(ies, m.NodeID)
	ies = append(ies, m.NodeReportType)
	ies = append(ies, m.UserPlanePathFailureReport)
	ies = append(ies, m.UserPlanePathRecoveryReport)
	ies = append(ies, m.ClockDriftReport...)
	ies = append(ies, m.GTPUPathQoSReport...)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *NodeReportRequest) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.NodeID:
			assignIE(&m.NodeID, i, &m.IEs)
		case ie.NodeReportType:
			assignIE(&m.NodeReportType, i, &m.IEs)
		case ie.UserPlanePathFailureReport:
			assignIE(&m.UserPlanePathFailureReport, i, &m.IEs)
		case ie.UserPlanePathRecoveryReport:
			assignIE(&m.UserPlanePathRecoveryReport, i, &m.IEs)
		case ie.ClockDriftReport:
			m.ClockDriftReport = append(m.ClockDriftReport, i)
		case ie.GTPUPathQoSReport:
			m.GTPUPathQoSReport = append(m.GTPUPathQoSReport, i)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
