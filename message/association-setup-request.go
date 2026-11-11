// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// AssociationSetupRequest is a AssociationSetupRequest formed PFCP Header and its IEs above.
type AssociationSetupRequest struct {
	*Header
	NodeID                          *ie.IE
	RecoveryTimeStamp               *ie.IE
	UPFunctionFeatures              *ie.IE
	CPFunctionFeatures              *ie.IE
	UserPlaneIPResourceInformation  []*ie.IE
	AlternativeSMFIPAddress         []*ie.IE
	SMFSetID                        *ie.IE
	PFCPSessionRetentionInformation *ie.IE
	UEIPAddressPoolInformation      []*ie.IE
	GTPUPathQoSControlInformation   []*ie.IE
	ClockDriftControlInformation    []*ie.IE
	UPFInstanceID                   *ie.IE
	PFCPASReqFlags                  *ie.IE
	IEs                             []*ie.IE
}

// NewAssociationSetupRequest creates a new AssociationSetupRequest.
func NewAssociationSetupRequest(seq uint32, ies ...*ie.IE) *AssociationSetupRequest {
	m := &AssociationSetupRequest{
		Header: NewHeader(MsgTypeAssociationSetupRequest, 0, seq, nil),
	}

	m.addIEs(ies)

	m.SetLength()
	return m
}

// Marshal returns the byte sequence generated from a AssociationSetupRequest.
func (m *AssociationSetupRequest) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *AssociationSetupRequest) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseAssociationSetupRequest decodes a given byte sequence as a AssociationSetupRequest.
func ParseAssociationSetupRequest(b []byte) (*AssociationSetupRequest, error) {
	m := &AssociationSetupRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a AssociationSetupRequest.
func (m *AssociationSetupRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeAssociationSetupRequest)
	if err != nil {
		return err
	}
	*m = AssociationSetupRequest{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *AssociationSetupRequest) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *AssociationSetupRequest) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *AssociationSetupRequest) MessageTypeName() string {
	return "Association Setup Request"
}

// SEID returns the SEID in uint64.
func (m *AssociationSetupRequest) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *AssociationSetupRequest) IsRequest() bool {
	return true
}

func (m *AssociationSetupRequest) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 8+len(m.UserPlaneIPResourceInformation)+len(m.AlternativeSMFIPAddress)+len(m.UEIPAddressPoolInformation)+len(m.GTPUPathQoSControlInformation)+len(m.ClockDriftControlInformation)+len(m.IEs))
	ies = append(ies, m.NodeID)
	ies = append(ies, m.RecoveryTimeStamp)
	ies = append(ies, m.UPFunctionFeatures)
	ies = append(ies, m.CPFunctionFeatures)
	ies = append(ies, m.UserPlaneIPResourceInformation...)
	ies = append(ies, m.AlternativeSMFIPAddress...)
	ies = append(ies, m.SMFSetID)
	ies = append(ies, m.PFCPSessionRetentionInformation)
	ies = append(ies, m.UEIPAddressPoolInformation...)
	ies = append(ies, m.GTPUPathQoSControlInformation...)
	ies = append(ies, m.ClockDriftControlInformation...)
	ies = append(ies, m.UPFInstanceID)
	ies = append(ies, m.PFCPASReqFlags)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *AssociationSetupRequest) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.NodeID:
			assignIE(&m.NodeID, i, &m.IEs)
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
		case ie.SMFSetID:
			assignIE(&m.SMFSetID, i, &m.IEs)
		case ie.PFCPSessionRetentionInformation:
			assignIE(&m.PFCPSessionRetentionInformation, i, &m.IEs)
		case ie.UEIPAddressPoolInformation:
			m.UEIPAddressPoolInformation = append(m.UEIPAddressPoolInformation, i)
		case ie.GTPUPathQoSControlInformation:
			m.GTPUPathQoSControlInformation = append(m.GTPUPathQoSControlInformation, i)
		case ie.ClockDriftControlInformation:
			m.ClockDriftControlInformation = append(m.ClockDriftControlInformation, i)
		case ie.NFInstanceID:
			assignIE(&m.UPFInstanceID, i, &m.IEs)
		case ie.PFCPASReqFlags:
			assignIE(&m.PFCPASReqFlags, i, &m.IEs)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
