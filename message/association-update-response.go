// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// AssociationUpdateResponse is a AssociationUpdateResponse formed PFCP Header and its IEs above.
type AssociationUpdateResponse struct {
	*Header
	NodeID             *ie.IE
	Cause              *ie.IE
	UPFunctionFeatures *ie.IE
	CPFunctionFeatures *ie.IE
	IEs                []*ie.IE
}

// NewAssociationUpdateResponse creates a new AssociationUpdateResponse.
func NewAssociationUpdateResponse(seq uint32, ies ...*ie.IE) *AssociationUpdateResponse {
	m := &AssociationUpdateResponse{
		Header: NewHeader(MsgTypeAssociationUpdateResponse, 0, seq, nil),
	}

	m.addIEs(ies)

	m.SetLength()
	return m
}

// Marshal returns the byte sequence generated from a AssociationUpdateResponse.
func (m *AssociationUpdateResponse) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *AssociationUpdateResponse) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseAssociationUpdateResponse decodes a given byte sequence as a AssociationUpdateResponse.
func ParseAssociationUpdateResponse(b []byte) (*AssociationUpdateResponse, error) {
	m := &AssociationUpdateResponse{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a AssociationUpdateResponse.
func (m *AssociationUpdateResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeAssociationUpdateResponse)
	if err != nil {
		return err
	}
	*m = AssociationUpdateResponse{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *AssociationUpdateResponse) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *AssociationUpdateResponse) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *AssociationUpdateResponse) MessageTypeName() string {
	return "Association Update Request"
}

// SEID returns the SEID in uint64.
func (m *AssociationUpdateResponse) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *AssociationUpdateResponse) IsRequest() bool {
	return false
}

func (m *AssociationUpdateResponse) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 4+len(m.IEs))
	ies = append(ies, m.NodeID)
	ies = append(ies, m.Cause)
	ies = append(ies, m.UPFunctionFeatures)
	ies = append(ies, m.CPFunctionFeatures)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *AssociationUpdateResponse) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.NodeID:
			assignIE(&m.NodeID, i, &m.IEs)
		case ie.Cause:
			assignIE(&m.Cause, i, &m.IEs)
		case ie.UPFunctionFeatures:
			assignIE(&m.UPFunctionFeatures, i, &m.IEs)
		case ie.CPFunctionFeatures:
			assignIE(&m.CPFunctionFeatures, i, &m.IEs)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
