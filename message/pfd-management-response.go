// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// PFDManagementResponse is a PFDManagementResponse formed PFCP Header and its IEs above.
//
// TODO: add Node ID IE.
type PFDManagementResponse struct {
	*Header
	Cause       *ie.IE
	OffendingIE *ie.IE
	IEs         []*ie.IE
}

// NewPFDManagementResponse creates a new PFDManagementResponse.
func NewPFDManagementResponse(seq uint32, cause, offending *ie.IE, ies ...*ie.IE) *PFDManagementResponse {
	m := &PFDManagementResponse{
		Header:      NewHeader(MsgTypePFDManagementResponse, 0, seq, nil),
		Cause:       cause,
		OffendingIE: offending,
		IEs:         ies,
	}
	m.SetLength()

	return m
}

// Marshal returns the byte sequence generated from a PFDManagementResponse.
func (m *PFDManagementResponse) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *PFDManagementResponse) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParsePFDManagementResponse decodes a given byte sequence as a PFDManagementResponse.
func ParsePFDManagementResponse(b []byte) (*PFDManagementResponse, error) {
	m := &PFDManagementResponse{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a PFDManagementResponse.
func (m *PFDManagementResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypePFDManagementResponse)
	if err != nil {
		return err
	}
	*m = PFDManagementResponse{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *PFDManagementResponse) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *PFDManagementResponse) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *PFDManagementResponse) MessageTypeName() string {
	return "PFD Management Response"
}

// SEID returns the SEID in uint64.
func (m *PFDManagementResponse) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *PFDManagementResponse) IsRequest() bool {
	return false
}

func (m *PFDManagementResponse) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 2+len(m.IEs))
	ies = append(ies, m.Cause)
	ies = append(ies, m.OffendingIE)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *PFDManagementResponse) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.Cause:
			assignIE(&m.Cause, i, &m.IEs)
		case ie.OffendingIE:
			assignIE(&m.OffendingIE, i, &m.IEs)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
