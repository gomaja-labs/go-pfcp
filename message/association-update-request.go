// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-pfcp/ie"
)

// AssociationUpdateRequest is a AssociationUpdateRequest formed PFCP Header and its IEs above.
type AssociationUpdateRequest struct {
	*Header
	NodeID                        *ie.IE
	UPFunctionFeatures            *ie.IE
	CPFunctionFeatures            *ie.IE
	PFCPAssociationReleaseRequest *ie.IE
	GracefulReleasePeriod         *ie.IE
	PFCPAUReqFlags                *ie.IE
	AlternativeSMFIPAddress       []*ie.IE
	ClockDriftControlInformation  []*ie.IE
	UEIPAddressPoolInformation    []*ie.IE
	GTPUPathQoSControlInformation []*ie.IE
	UEIPAddressUsageInformation   []*ie.IE
	IEs                           []*ie.IE
}

// NewAssociationUpdateRequest creates a new AssociationUpdateRequest.
func NewAssociationUpdateRequest(seq uint32, ies ...*ie.IE) *AssociationUpdateRequest {
	m := &AssociationUpdateRequest{
		Header: NewHeader(MsgTypeAssociationUpdateRequest, 0, seq, nil),
	}

	m.addIEs(ies)

	m.SetLength()
	return m
}

// Marshal returns the byte sequence generated from a AssociationUpdateRequest.
func (m *AssociationUpdateRequest) Marshal() ([]byte, error) {
	return marshalMessage(m)
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *AssociationUpdateRequest) MarshalTo(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	return marshalMessageTo(m.Header, m.messageIEs(), b)
}

// ParseAssociationUpdateRequest decodes a given byte sequence as a AssociationUpdateRequest.
func ParseAssociationUpdateRequest(b []byte) (*AssociationUpdateRequest, error) {
	m := &AssociationUpdateRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes a given byte sequence as a AssociationUpdateRequest.
func (m *AssociationUpdateRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	h, err := parseMessageHeader(b, MsgTypeAssociationUpdateRequest)
	if err != nil {
		return err
	}
	*m = AssociationUpdateRequest{Header: h}

	ies, err := ie.ParseMultiIEs(m.Header.Payload)
	if err != nil {
		return err
	}

	m.addIEs(ies)

	return nil
}

// MarshalLen returns the serial length of Data.
func (m *AssociationUpdateRequest) MarshalLen() int {
	if m == nil {
		return 0
	}
	return messageLen(m.Header, m.messageIEs())
}

// SetLength sets the length in Length field.
func (m *AssociationUpdateRequest) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.Length = uint16(m.MarshalLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *AssociationUpdateRequest) MessageTypeName() string {
	return "Association Update Request"
}

// SEID returns the SEID in uint64.
func (m *AssociationUpdateRequest) SEID() uint64 {
	if m == nil {
		return 0
	}
	return m.Header.seid()
}

// IsRequest reports whether the message is a request.
func (m *AssociationUpdateRequest) IsRequest() bool {
	return true
}

func (m *AssociationUpdateRequest) messageIEs() []*ie.IE {
	if m == nil {
		return nil
	}
	ies := make([]*ie.IE, 0, 6+len(m.AlternativeSMFIPAddress)+len(m.ClockDriftControlInformation)+len(m.UEIPAddressPoolInformation)+len(m.GTPUPathQoSControlInformation)+len(m.UEIPAddressUsageInformation)+len(m.IEs))
	ies = append(ies, m.NodeID)
	ies = append(ies, m.UPFunctionFeatures)
	ies = append(ies, m.CPFunctionFeatures)
	ies = append(ies, m.PFCPAssociationReleaseRequest)
	ies = append(ies, m.GracefulReleasePeriod)
	ies = append(ies, m.PFCPAUReqFlags)
	ies = append(ies, m.AlternativeSMFIPAddress...)
	ies = append(ies, m.ClockDriftControlInformation...)
	ies = append(ies, m.UEIPAddressPoolInformation...)
	ies = append(ies, m.GTPUPathQoSControlInformation...)
	ies = append(ies, m.UEIPAddressUsageInformation...)
	ies = append(ies, m.IEs...)
	return ies
}

func (m *AssociationUpdateRequest) addIEs(ies []*ie.IE) {
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.NodeID:
			assignIE(&m.NodeID, i, &m.IEs)
		case ie.UPFunctionFeatures:
			assignIE(&m.UPFunctionFeatures, i, &m.IEs)
		case ie.CPFunctionFeatures:
			assignIE(&m.CPFunctionFeatures, i, &m.IEs)
		case ie.PFCPAssociationReleaseRequest:
			assignIE(&m.PFCPAssociationReleaseRequest, i, &m.IEs)
		case ie.GracefulReleasePeriod:
			assignIE(&m.GracefulReleasePeriod, i, &m.IEs)
		case ie.PFCPAUReqFlags:
			assignIE(&m.PFCPAUReqFlags, i, &m.IEs)
		case ie.AlternativeSMFIPAddress:
			m.AlternativeSMFIPAddress = append(m.AlternativeSMFIPAddress, i)
		case ie.ClockDriftControlInformation:
			m.ClockDriftControlInformation = append(m.ClockDriftControlInformation, i)
		case ie.UEIPAddressPoolInformation:
			m.UEIPAddressPoolInformation = append(m.UEIPAddressPoolInformation, i)
		case ie.GTPUPathQoSControlInformation:
			m.GTPUPathQoSControlInformation = append(m.GTPUPathQoSControlInformation, i)
		case ie.UEIPAddressUsageInformation:
			m.UEIPAddressUsageInformation = append(m.UEIPAddressUsageInformation, i)
		default:
			m.IEs = append(m.IEs, i)
		}
	}
}
