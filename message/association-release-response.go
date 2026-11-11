// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// AssociationReleaseResponse is a AssociationReleaseResponse formed PFCP Header and its IEs above.
type AssociationReleaseResponse struct {
	*Header
	NodeID *ie.IE
	Cause  *ie.IE
	IEs    []*ie.IE
}

// NewAssociationReleaseResponse creates a new AssociationReleaseResponse.
func NewAssociationReleaseResponse(seq uint32, id, cause *ie.IE, ies ...*ie.IE) *AssociationReleaseResponse {
	m := &AssociationReleaseResponse{
		Header: NewHeader(MsgTypeAssociationReleaseResponse, 0, seq, nil),
		NodeID: id,
		Cause:  cause,
		IEs:    ies,
	}
	m.SetLength()

	return m
}

// Marshal returns the byte sequence generated from a AssociationReleaseResponse.
func (m *AssociationReleaseResponse) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *AssociationReleaseResponse) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseAssociationReleaseResponse decodes a given byte sequence as a AssociationReleaseResponse.
func ParseAssociationReleaseResponse(b []byte) (*AssociationReleaseResponse, error) {
	m := &AssociationReleaseResponse{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a AssociationReleaseResponse.
func (m *AssociationReleaseResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeAssociationReleaseResponse)
	if err != nil {
		return err
	}
	*m = AssociationReleaseResponse{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *AssociationReleaseResponse) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *AssociationReleaseResponse) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *AssociationReleaseResponse) MessageTypeName() string {
	return "Association Release Response"
}

// SEID returns the SEID in uint64.
func (m *AssociationReleaseResponse) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *AssociationReleaseResponse) IsRequest() bool {
	return false
}

func (m *AssociationReleaseResponse) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 2+len(m.IEs))
	ies = append(ies, m.NodeID)
	ies = append(ies, m.Cause)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *AssociationReleaseResponse) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.NodeID:
			assignIE(&m.NodeID, i, &m.IEs)
		case ie.Cause:
			assignIE(&m.Cause, i, &m.IEs)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
