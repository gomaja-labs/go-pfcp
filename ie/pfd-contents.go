// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
)

// NewPFDContents creates a new PFDContents IE.
func NewPFDContents(fd, url, dn, cp, dnp string, afd, aurl, adnp []string) *IE {
	f := NewPFDContentsFields(fd, url, dn, cp, dnp, afd, aurl, adnp)
	b, err := f.Marshal()
	if err != nil {
		return nil
	}

	return New(PFDContents, b)
}

// PFDContents returns PFDContents in structured format if the type of IE matches.
//
// This IE has a complex payload that costs much when parsing.
func (i *IE) PFDContents() (*PFDContentsFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case PFDContents:
		s, err := ParsePFDContentsFields(i.Payload)
		if err != nil {
			return nil, err
		}
		return s, nil
	case PFDContext:
		ies, err := i.PFDContext()
		if err != nil {
			return nil, err
		}
		for _, x := range ies {
			if x.Type == PFDContents {
				return x.PFDContents()
			}
		}
		return nil, ErrIENotFound
	case ApplicationIDsPFDs:
		ies, err := i.ApplicationIDsPFDs()
		if err != nil {
			return nil, err
		}
		for _, x := range ies {
			if x.Type == PFDContext {
				return x.PFDContents()
			}
		}
		return nil, ErrIENotFound
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// PFDContentsFields represents a fields contained in PFDContents IE.
type PFDContentsFields struct {
	Flags                           uint8
	FDLength                        uint16
	FlowDescription                 string
	URLLength                       uint16
	URL                             string
	DNLength                        uint16
	DomainName                      string
	CPLength                        uint16
	CustomPFDContent                string
	DNPLength                       uint16
	DomainNameProtocol              string
	AFDLength                       uint16
	AdditionalFlowDescription       []string
	AURLLength                      uint16
	AdditionalURL                   []string
	ADNPLength                      uint16
	AdditionalDomainNameAndProtocol []string
}

// NewPFDContentsFields creates a new NewPFDContentsFields.
func NewPFDContentsFields(fd, url, dn, cp, dnp string, afd, aurl, adnp []string) *PFDContentsFields {
	f := &PFDContentsFields{}
	if fd != "" {
		f.FlowDescription = fd
		f.FDLength = uint16(len([]byte(fd)))
		f.SetFDFlag()
	}

	if url != "" {
		f.URL = url
		f.URLLength = uint16(len([]byte(url)))
		f.SetURLFlag()
	}

	if dn != "" {
		f.DomainName = dn
		f.DNLength = uint16(len([]byte(dn)))
		f.SetDNFlag()
	}

	if cp != "" {
		f.CustomPFDContent = cp
		f.CPLength = uint16(len([]byte(cp)))
		f.SetCPFlag()
	}

	if dnp != "" {
		f.DomainNameProtocol = dnp
		f.DNPLength = uint16(len([]byte(dnp)))
		f.SetDNPFlag()
	}

	if len(afd) > 0 {
		l := 0
		for _, a := range afd {
			l += 2 + len([]byte(a))
		}

		f.AdditionalFlowDescription = afd
		f.AFDLength = uint16(l)
		f.SetAFDFlag()
	}

	if len(aurl) > 0 {
		l := 0
		for _, a := range aurl {
			l += 2 + len([]byte(a))
		}

		f.AdditionalURL = aurl
		f.AURLLength = uint16(l)
		f.SetAURLFlag()
	}

	if len(adnp) > 0 {
		l := 0
		for _, a := range adnp {
			l += 2 + len([]byte(a))
		}

		f.AdditionalDomainNameAndProtocol = adnp
		f.ADNPLength = uint16(l)
		f.SetADNPFlag()
	}

	return f
}

// HasADNP reports whether ADNP flag is set.
func (f *PFDContentsFields) HasADNP() bool {
	if f == nil {
		return false
	}

	return has8thBit(f.Flags)
}

// SetADNPFlag sets ADNP flag in PFDContents.
func (f *PFDContentsFields) SetADNPFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x80
}

// HasAURL reports whether AURL flag is set.
func (f *PFDContentsFields) HasAURL() bool {
	if f == nil {
		return false
	}

	return has7thBit(f.Flags)
}

// SetAURLFlag sets AURL flag in PFDContents.
func (f *PFDContentsFields) SetAURLFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x40
}

// HasAFD reports whether AFD flag is set.
func (f *PFDContentsFields) HasAFD() bool {
	if f == nil {
		return false
	}

	return has6thBit(f.Flags)
}

// SetAFDFlag sets AFD flag in PFDContents.
func (f *PFDContentsFields) SetAFDFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x20
}

// HasDNP reports whether DNP flag is set.
func (f *PFDContentsFields) HasDNP() bool {
	if f == nil {
		return false
	}

	return has5thBit(f.Flags)
}

// SetDNPFlag sets DNP flag in PFDContents.
func (f *PFDContentsFields) SetDNPFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x10
}

// HasCP reports whether CP flag is set.
func (f *PFDContentsFields) HasCP() bool {
	if f == nil {
		return false
	}

	return has4thBit(f.Flags)
}

// SetCPFlag sets CP flag in PFDContents.
func (f *PFDContentsFields) SetCPFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x08
}

// HasDN reports whether DN flag is set.
func (f *PFDContentsFields) HasDN() bool {
	if f == nil {
		return false
	}

	return has3rdBit(f.Flags)
}

// SetDNFlag sets DN flag in PFDContents.
func (f *PFDContentsFields) SetDNFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x04
}

// HasURL reports whether URL flag is set.
func (f *PFDContentsFields) HasURL() bool {
	if f == nil {
		return false
	}

	return has2ndBit(f.Flags)
}

// SetURLFlag sets URL flag in PFDContents.
func (f *PFDContentsFields) SetURLFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x02
}

// HasFD reports whether FD flag is set.
func (f *PFDContentsFields) HasFD() bool {
	if f == nil {
		return false
	}

	return has1stBit(f.Flags)
}

// SetFDFlag sets FD flag in PFDContents.
func (f *PFDContentsFields) SetFDFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x01
}

// ParsePFDContentsFields parses b into PFDContentsFields.
func ParsePFDContentsFields(b []byte) (*PFDContentsFields, error) {
	f := &PFDContentsFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary parses b into IE.
// TS 29.244 V20.1.0 clause 8.2.39: each present value, including each
// instance in an additional-values field, is preceded by a two-octet length.
func (f *PFDContentsFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = PFDContentsFields{}

	if len(b) < 2 {
		return io.ErrUnexpectedEOF
	}
	*f = PFDContentsFields{Flags: b[0]}
	offset := 2
	fields := []struct {
		bit    uint8
		length *uint16
		value  *string
		list   *[]string
	}{
		{1, &f.FDLength, &f.FlowDescription, nil},
		{2, &f.URLLength, &f.URL, nil},
		{4, &f.DNLength, &f.DomainName, nil},
		{8, &f.CPLength, &f.CustomPFDContent, nil},
		{16, &f.DNPLength, &f.DomainNameProtocol, nil},
		{32, &f.AFDLength, nil, &f.AdditionalFlowDescription},
		{64, &f.AURLLength, nil, &f.AdditionalURL},
		{128, &f.ADNPLength, nil, &f.AdditionalDomainNameAndProtocol},
	}
	for _, field := range fields {
		if f.Flags&field.bit == 0 {
			continue
		}
		value, err := takeLength16(b, &offset)
		if err != nil {
			return err
		}
		*field.length = uint16(len(value))
		if field.value != nil {
			*field.value = string(value)
			continue
		}
		for at := 0; at < len(value); {
			instance, err := takeLength16(value, &at)
			if err != nil {
				return err
			}
			*field.list = append(*field.list, string(instance))
		}
	}
	return nil
}

// Marshal returns the serialized bytes of PFDContentsFields.
func (f *PFDContentsFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}

	b := make([]byte, f.MarshalLen())
	if err := f.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (f *PFDContentsFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if f.MarshalLen() > 65535 {
		return ErrInvalidLength
	}
	if !f.validLengths() {
		return ErrInvalidLength
	}

	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}

	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	b[0] = f.Flags
	offset := 2 // 2nd octet is spare

	if f.HasFD() {
		binary.BigEndian.PutUint16(b[offset:offset+2], f.FDLength)
		copy(b[offset+2:offset+2+int(f.FDLength)], f.FlowDescription)
		offset += 2 + int(f.FDLength)
	}

	if f.HasURL() {
		binary.BigEndian.PutUint16(b[offset:offset+2], f.URLLength)
		copy(b[offset+2:offset+2+int(f.URLLength)], f.URL)
		offset += 2 + int(f.URLLength)
	}

	if f.HasDN() {
		binary.BigEndian.PutUint16(b[offset:offset+2], f.DNLength)
		copy(b[offset+2:offset+2+int(f.DNLength)], f.DomainName)
		offset += 2 + int(f.DNLength)
	}

	if f.HasCP() {
		binary.BigEndian.PutUint16(b[offset:offset+2], f.CPLength)
		copy(b[offset+2:offset+2+int(f.CPLength)], f.CustomPFDContent)
		offset += 2 + int(f.CPLength)
	}

	if f.HasDNP() {
		binary.BigEndian.PutUint16(b[offset:offset+2], f.DNPLength)
		copy(b[offset+2:offset+2+int(f.DNPLength)], f.DomainNameProtocol)
		offset += 2 + int(f.DNPLength)
	}

	if f.HasAFD() {
		binary.BigEndian.PutUint16(b[offset:offset+2], f.AFDLength)
		offset += 2
		for _, a := range f.AdditionalFlowDescription {
			l := len([]byte(a))
			binary.BigEndian.PutUint16(b[offset:offset+2], uint16(l))
			copy(b[offset+2:offset+2+l], a)
			offset += 2 + l
		}
	}

	if f.HasAURL() {
		binary.BigEndian.PutUint16(b[offset:offset+2], f.AURLLength)
		offset += 2
		for _, a := range f.AdditionalURL {
			l := len([]byte(a))
			binary.BigEndian.PutUint16(b[offset:offset+2], uint16(l))
			copy(b[offset+2:offset+2+l], a)
			offset += 2 + l
		}
	}

	if f.HasADNP() {
		binary.BigEndian.PutUint16(b[offset:offset+2], f.ADNPLength)
		offset += 2
		for _, a := range f.AdditionalDomainNameAndProtocol {
			l := len([]byte(a))
			binary.BigEndian.PutUint16(b[offset:offset+2], uint16(l))
			copy(b[offset+2:offset+2+l], a)
			offset += 2 + l
		}
	}

	return nil
}

// MarshalLen returns field length in integer.
func (f *PFDContentsFields) MarshalLen() int {
	if f == nil {
		return 0
	}

	l := 2 // Flag + Spare

	if f.HasFD() {
		l += 2 + int(f.FDLength)
	}

	if f.HasURL() {
		l += 2 + int(f.URLLength)
	}

	if f.HasDN() {
		l += 2 + int(f.DNLength)
	}

	if f.HasCP() {
		l += 2 + int(f.CPLength)
	}

	if f.HasDNP() {
		l += 2 + int(f.DNPLength)
	}

	if f.HasAFD() {
		l += 2 + int(f.AFDLength)
	}

	if f.HasAURL() {
		l += 2 + int(f.AURLLength)
	}

	if f.HasADNP() {
		l += 2 + int(f.ADNPLength)
	}

	return l
}

// Length metadata must agree with the values before any write. This also
// detects constructor input lengths that cannot fit in sixteen bits.
func (f *PFDContentsFields) validLengths() bool {
	for _, field := range []struct {
		present bool
		length  uint16
		value   string
	}{
		{f.HasFD(), f.FDLength, f.FlowDescription},
		{f.HasURL(), f.URLLength, f.URL},
		{f.HasDN(), f.DNLength, f.DomainName},
		{f.HasCP(), f.CPLength, f.CustomPFDContent},
		{f.HasDNP(), f.DNPLength, f.DomainNameProtocol},
	} {
		if field.present && len(field.value) != int(field.length) {
			return false
		}
	}
	for _, field := range []struct {
		present bool
		length  uint16
		values  []string
	}{
		{f.HasAFD(), f.AFDLength, f.AdditionalFlowDescription},
		{f.HasAURL(), f.AURLLength, f.AdditionalURL},
		{f.HasADNP(), f.ADNPLength, f.AdditionalDomainNameAndProtocol},
	} {
		if !field.present {
			continue
		}
		length := 0
		for _, value := range field.values {
			if len(value) > 65535 {
				return false
			}
			length += 2 + len(value)
		}
		if length != int(field.length) {
			return false
		}
	}
	return true
}
