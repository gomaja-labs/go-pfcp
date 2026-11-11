// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// AssociationReleaseRequest is a AssociationReleaseRequest formed PFCP Header and its IEs above.
type AssociationReleaseRequest struct {
	*Header
	NodeID *ie.IE
	IEs    []*ie.IE
}

// NewAssociationReleaseRequest creates a new AssociationReleaseRequest.
func NewAssociationReleaseRequest(seq uint32, id *ie.IE, ies ...*ie.IE) *AssociationReleaseRequest {
	m := &AssociationReleaseRequest{
		Header: NewHeader(MsgTypeAssociationReleaseRequest, 0, seq, nil),
		NodeID: id,
		IEs:    ies,
	}
	m.SetLength()

	return m
}

// Marshal returns the byte sequence generated from a AssociationReleaseRequest.
func (m *AssociationReleaseRequest) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *AssociationReleaseRequest) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseAssociationReleaseRequest decodes a given byte sequence as a AssociationReleaseRequest.
func ParseAssociationReleaseRequest(b []byte) (*AssociationReleaseRequest, error) {
	m := &AssociationReleaseRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a AssociationReleaseRequest.
func (m *AssociationReleaseRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeAssociationReleaseRequest)
	if err != nil {
		return err
	}
	*m = AssociationReleaseRequest{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *AssociationReleaseRequest) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *AssociationReleaseRequest) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *AssociationReleaseRequest) MessageTypeName() string {
	return "Association Release Request"
}

// SEID returns the SEID in uint64.
func (m *AssociationReleaseRequest) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *AssociationReleaseRequest) IsRequest() bool {
	return true
}

func (m *AssociationReleaseRequest) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 1+len(m.IEs))
	ies = append(ies, m.NodeID)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *AssociationReleaseRequest) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.NodeID:
			assignIE(&m.NodeID, i, &m.IEs)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
