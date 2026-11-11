// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// PFDManagementRequest is a PFDManagementRequest formed PFCP Header and its IEs above.
type PFDManagementRequest struct {
	*Header
	ApplicationIDsPFDs []*ie.IE
	IEs                []*ie.IE
}

// NewPFDManagementRequest creates a new PFDManagementRequest.
func NewPFDManagementRequest(seq uint32, ies ...*ie.IE) *PFDManagementRequest {
	m := &PFDManagementRequest{
		Header: NewHeader(MsgTypePFDManagementRequest, 0, seq, nil),
	}

	m.addIEs(ies)

	m.SetLength()
	return m
}

// Marshal returns the byte sequence generated from a PFDManagementRequest.
func (m *PFDManagementRequest) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *PFDManagementRequest) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParsePFDManagementRequest decodes a given byte sequence as a PFDManagementRequest.
func ParsePFDManagementRequest(b []byte) (*PFDManagementRequest, error) {
	m := &PFDManagementRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a PFDManagementRequest.
func (m *PFDManagementRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypePFDManagementRequest)
	if err != nil {
		return err
	}
	*m = PFDManagementRequest{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *PFDManagementRequest) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *PFDManagementRequest) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *PFDManagementRequest) MessageTypeName() string {
	return "PFD Management Request"
}

// SEID returns the SEID in uint64.
func (m *PFDManagementRequest) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *PFDManagementRequest) IsRequest() bool {
	return true
}

func (m *PFDManagementRequest) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, len(m.ApplicationIDsPFDs)+len(m.IEs))
	ies = append(ies, m.ApplicationIDsPFDs...)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *PFDManagementRequest) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.ApplicationIDsPFDs:
			m.ApplicationIDsPFDs = append(m.ApplicationIDsPFDs, i)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
