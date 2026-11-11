// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// AssociationSetupResponse is a AssociationSetupResponse formed PFCP Header and its IEs above.
type AssociationSetupResponse struct {
	*Header
	NodeID                         *ie.IE
	Cause                          *ie.IE
	RecoveryTimeStamp              *ie.IE
	UPFunctionFeatures             *ie.IE
	CPFunctionFeatures             *ie.IE
	UserPlaneIPResourceInformation []*ie.IE
	AlternativeSMFIPAddress        []*ie.IE
	PFCPASRspFlags                 *ie.IE
	UEIPAddressPoolInformation     []*ie.IE
	GTPUPathQoSControlInformation  []*ie.IE
	ClockDriftControlInformation   []*ie.IE
	UPFInstanceID                  *ie.IE
	IEs                            []*ie.IE
}

// NewAssociationSetupResponse creates a new AssociationSetupResponse.
func NewAssociationSetupResponse(seq uint32, ies ...*ie.IE) *AssociationSetupResponse {
	m := &AssociationSetupResponse{
		Header: NewHeader(MsgTypeAssociationSetupResponse, 0, seq, nil),
	}

	m.addIEs(ies)

	m.SetLength()
	return m
}

// Marshal returns the byte sequence generated from a AssociationSetupResponse.
func (m *AssociationSetupResponse) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *AssociationSetupResponse) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseAssociationSetupResponse decodes a given byte sequence as a AssociationSetupResponse.
func ParseAssociationSetupResponse(b []byte) (*AssociationSetupResponse, error) {
	m := &AssociationSetupResponse{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a AssociationSetupResponse.
func (m *AssociationSetupResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeAssociationSetupResponse)
	if err != nil {
		return err
	}
	*m = AssociationSetupResponse{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *AssociationSetupResponse) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *AssociationSetupResponse) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *AssociationSetupResponse) MessageTypeName() string {
	return "Association Setup Response"
}

// SEID returns the SEID in uint64.
func (m *AssociationSetupResponse) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *AssociationSetupResponse) IsRequest() bool {
	return false
}

func (m *AssociationSetupResponse) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 7+len(m.UserPlaneIPResourceInformation)+len(m.AlternativeSMFIPAddress)+len(m.UEIPAddressPoolInformation)+len(m.GTPUPathQoSControlInformation)+len(m.ClockDriftControlInformation)+len(m.IEs))
	ies = append(ies, m.NodeID)
	ies = append(ies, m.Cause)
	ies = append(ies, m.RecoveryTimeStamp)
	ies = append(ies, m.UPFunctionFeatures)
	ies = append(ies, m.CPFunctionFeatures)
	ies = append(ies, m.UserPlaneIPResourceInformation...)
	ies = append(ies, m.AlternativeSMFIPAddress...)
	ies = append(ies, m.PFCPASRspFlags)
	ies = append(ies, m.UEIPAddressPoolInformation...)
	ies = append(ies, m.GTPUPathQoSControlInformation...)
	ies = append(ies, m.ClockDriftControlInformation...)
	ies = append(ies, m.UPFInstanceID)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *AssociationSetupResponse) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.NodeID:
			assignIE(&m.NodeID, i, &m.IEs)
		case ie.Cause:
			assignIE(&m.Cause, i, &m.IEs)
		case ie.RecoveryTimeStamp:
			assignIE(&m.RecoveryTimeStamp, i, &m.IEs)
		case ie.UPFunctionFeatures:
			assignIE(&m.UPFunctionFeatures, i, &m.IEs)
		case ie.CPFunctionFeatures:
			assignIE(&m.CPFunctionFeatures, i, &m.IEs)
		case ie.UserPlaneIPResourceInformation:
			m.UserPlaneIPResourceInformation = append(m.UserPlaneIPResourceInformation, i)
		case ie.AlternativeSMFIPAddress:
			m.AlternativeSMFIPAddress = append(m.AlternativeSMFIPAddress, i)
		case ie.PFCPASRspFlags:
			assignIE(&m.PFCPASRspFlags, i, &m.IEs)
		case ie.UEIPAddressPoolInformation:
			m.UEIPAddressPoolInformation = append(m.UEIPAddressPoolInformation, i)
		case ie.GTPUPathQoSControlInformation:
			m.GTPUPathQoSControlInformation = append(m.GTPUPathQoSControlInformation, i)
		case ie.ClockDriftControlInformation:
			m.ClockDriftControlInformation = append(m.ClockDriftControlInformation, i)
		case ie.NFInstanceID:
			assignIE(&m.UPFInstanceID, i, &m.IEs)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
