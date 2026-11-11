// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// VersionNotSupportedResponse is a VersionNotSupportedResponse formed PFCP Header and its IEs above.
type VersionNotSupportedResponse struct {
	*Header
	IEs []*ie.IE
}

// NewVersionNotSupportedResponse creates a new VersionNotSupportedResponse.
func NewVersionNotSupportedResponse(seq uint32, ies ...*ie.IE) *VersionNotSupportedResponse {
	m := &VersionNotSupportedResponse{
		Header: NewHeader(MsgTypeVersionNotSupportedResponse, 0, seq, nil),
		IEs:    ies,
	}
	m.SetLength()

	return m
}

// Marshal returns the byte sequence generated from a VersionNotSupportedResponse.
func (m *VersionNotSupportedResponse) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *VersionNotSupportedResponse) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseVersionNotSupportedResponse decodes a given byte sequence as a VersionNotSupportedResponse.
func ParseVersionNotSupportedResponse(b []byte) (*VersionNotSupportedResponse, error) {
	m := &VersionNotSupportedResponse{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a VersionNotSupportedResponse.
func (m *VersionNotSupportedResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeVersionNotSupportedResponse)
	if err != nil {
		return err
	}
	*m = VersionNotSupportedResponse{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.IEs = append(m.IEs, ies...)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *VersionNotSupportedResponse) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *VersionNotSupportedResponse) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *VersionNotSupportedResponse) MessageTypeName() string {
	return "Version Not Supported Response"
}

// SEID returns the SEID in uint64.
func (m *VersionNotSupportedResponse) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *VersionNotSupportedResponse) IsRequest() bool {
	return false
}

func (m *VersionNotSupportedResponse) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, len(m.IEs))
	ies = append(ies, m.IEs...)
	return ies
}
