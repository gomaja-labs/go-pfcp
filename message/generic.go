// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"fmt"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// Generic is a Generic formed PFCP Header and its IEs above.
type Generic struct {
	*Header
	IEs []*ie.IE
}

// NewGeneric creates a new Generic.
func NewGeneric(msgType uint8, seid uint64, seq uint32, ies ...*ie.IE) *Generic {
	m := &Generic{
		Header: NewHeader(msgType, seid, seq, nil),
		IEs:    ies,
	}
	m.SetLength()

	return m
}

// Marshal returns the byte sequence generated from a Generic.
func (m *Generic) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *Generic) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseGeneric decodes a given byte sequence as a Generic.
func ParseGeneric(b []byte) (*Generic, error) {
	m := &Generic{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a Generic.
func (m *Generic) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, 0)
	if err != nil {
		return err
	}
	*m = Generic{Header: h}

	m.IEs, err = ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}
	return nil
}

// MarshalLen returns the serial length of Data.
func (m *Generic) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *Generic) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *Generic) MessageTypeName() string {
	if m == nil || m.Header == nil {
		return "Unknown (0)"
	}
	return fmt.Sprintf("Unknown (%d)", m.Header.Type)
}

// SEID returns the SEID in uint64.
func (m *Generic) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// AddIE add IEs to Generic type of PFCP message and update Length field.
func (m *Generic) AddIE(ies ...*ie.IE) {
	if m == nil {
		return
	}
	m.IEs = append(m.IEs, ies...)
	m.SetLength()
}

// IsRequest reports whether the message is a request.
func (m *Generic) IsRequest() bool {
	return false
}

func (m *Generic) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, len(m.IEs))
	ies = append(ies, m.IEs...)
	return ies
}
