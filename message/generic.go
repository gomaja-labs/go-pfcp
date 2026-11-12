// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"fmt"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// Generic retains the header and IEs of a message without a generated schema.
type Generic struct {
	// Header holds the PFCP framing and received IE ordering template.
	*Header
	// OtherIEs retains received IEs outside the named fields.
	OtherIEs []*ie.IE
}

// NewGeneric creates a Generic, copying the variadic slice and ignoring nil IEs.
func NewGeneric(msgType uint8, seid uint64, seq uint32, ies ...*ie.IE) *Generic {
	m := &Generic{
		Header: NewHeader(msgType, seid, seq, nil),
	}

	m.AddIE(ies...)
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
// Complete headers with inconsistent Length return InvalidLengthError
// (TS 29.244 V20.1.0 clause 7.6.3).
// Decoding may borrow b; keep b unchanged while the result is in use.
func ParseGeneric(b []byte) (*Generic, error) {
	m := &Generic{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a Generic.
// Complete headers with inconsistent Length return InvalidLengthError;
// a failed decode leaves the receiver unchanged (TS 29.244 V20.1.0 clause 7.6.3).
// Decoding may borrow b; keep b unchanged while the result is in use.
func (m *Generic) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, 0)
	if err != nil {
		return err
	}
	ies, err := ie.ParseMessageIEs(h.Payload)
	if err != nil {
		return err
	}
	*m = Generic{Header: h, OtherIEs: ies}
	return nil
}

// MarshalLen returns the serial length of Data.
func (m *Generic) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// MessageTypeName returns the name of protocol.
func (m *Generic) MessageTypeName() string {
	if m == nil || m.Header == nil {
		return ""
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

// AddIE appends non-nil IEs without writing into a slice shared with a copy.
func (m *Generic) AddIE(ies ...*ie.IE) {
	if m == nil {
		return
	}
	children := m.OtherIEs[:len(m.OtherIEs):len(m.OtherIEs)]
	for _, i := range ies {
		if i != nil {
			children = append(children, i)
		}
	}
	m.OtherIEs = children
}

// IsRequest classifies the received wire Type field through Header.IsRequest,
// including known request types decoded by ParseGeneric.
func (m *Generic) IsRequest() bool {
	return m != nil && m.Header.IsRequest()
}

func (m *Generic) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, len(m.OtherIEs))
	ies = append(ies, m.OtherIEs...)
	return ies
}

// Version returns the received version, or zero for nil.
func (m *Generic) Version() int {
	if m == nil {
		return 0
	}
	return m.Header.Version()
}

// MessageType returns the received type, or zero for nil.
func (m *Generic) MessageType() uint8 {
	if m == nil {
		return 0
	}
	return m.Header.MessageType()
}

// Sequence returns the sequence number, or zero for nil.
func (m *Generic) Sequence() uint32 {
	if m == nil {
		return 0
	}
	return m.Header.Sequence()
}

// SetSequenceNumber assigns the sequence number unless m is nil.
func (m *Generic) SetSequenceNumber(seq uint32) {
	if m != nil {
		m.Header.SetSequenceNumber(seq)
	}
}

// HasFO reports whether another frame follows, or false for nil.
func (m *Generic) HasFO() bool { return m != nil && m.Header.HasFO() }

// HasMP reports whether priority is present, or false for nil.
func (m *Generic) HasMP() bool { return m != nil && m.Header.HasMP() }

// HasSEID reports whether the SEID is present, or false for nil.
func (m *Generic) HasSEID() bool { return m != nil && m.Header.HasSEID() }

// SetSEID assigns the endpoint identifier, or returns ErrNilMessage for nil.
func (m *Generic) SetSEID(seid uint64) error {
	if m == nil {
		return ErrNilMessage
	}
	return m.Header.SetSEID(seid)
}

// SetMP sets priority, or returns ErrNilMessage for nil.
func (m *Generic) SetMP(priority uint8) error {
	if m == nil {
		return ErrNilMessage
	}
	return m.Header.SetMP(priority)
}

// ClearMP removes priority unless m is nil.
func (m *Generic) ClearMP() {
	if m != nil {
		m.Header.ClearMP()
	}
}

// MP returns priority, or zero for nil.
func (m *Generic) MP() uint8 {
	if m == nil {
		return 0
	}
	return m.Header.MP()
}

// String summarizes the header, or returns an empty string for nil.
func (m *Generic) String() string {
	if m == nil {
		return ""
	}
	return m.Header.String()
}
