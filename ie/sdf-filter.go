// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
)

// NewSDFFilter creates a new SDFFilter IE.
func NewSDFFilter(fd, ttc, spi, fl string, fid uint32) *IE {
	f := NewSDFFilterFields(fd, ttc, spi, fl, fid)
	b, err := f.Marshal()
	if err != nil {
		return nil
	}

	return New(SDFFilter, b)
}

// SDFFilter returns SDFFilter in structured format if the type of IE matches.
//
// This IE has a complex payload that costs much when parsing.
func (i *IE) SDFFilter() (*SDFFilterFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case SDFFilter:
		fields, err := ParseSDFFilterFields(i.Payload)
		if err != nil {
			return nil, err
		}

		return fields, nil
	case CreatePDR:
		ies, err := i.CreatePDR()
		if err != nil {
			return nil, err
		}
		for _, x := range ies {
			if x.Type == PDI {
				return x.SDFFilter()
			}
		}
		return nil, ErrIENotFound
	case PDI:
		ies, err := i.PDI()
		if err != nil {
			return nil, err
		}
		for _, x := range ies {
			if x.Type == SDFFilter {
				return x.SDFFilter()
			}
		}
		return nil, ErrIENotFound
	case EthernetPacketFilter:
		ies, err := i.EthernetPacketFilter()
		if err != nil {
			return nil, err
		}
		for _, x := range ies {
			if x.Type == SDFFilter {
				return x.SDFFilter()
			}
		}
		return nil, ErrIENotFound
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// SDFFilterFields represents a fields contained in SDFFilter IE.
type SDFFilterFields struct {
	Flags                  uint8
	FDLength               uint16
	FlowDescription        string
	ToSTrafficClass        string // 2 octets
	SecurityParameterIndex string // 4 octets
	FlowLabel              string // 3 octets
	SDFFilterID            uint32
}

// NewSDFFilterFields creates a new NewSDFFilterFields.
func NewSDFFilterFields(fd, ttc, spi, fl string, fid uint32) *SDFFilterFields {
	f := &SDFFilterFields{}
	if fd != "" {
		f.FlowDescription = fd
		f.FDLength = uint16(len([]byte(fd)))
		f.SetFDFlag()
	}

	if ttc != "" {
		f.ToSTrafficClass = ttc
		f.SetTTCFlag()
	}

	if spi != "" {
		f.SecurityParameterIndex = spi
		f.SetSPIFlag()
	}

	if fl != "" {
		f.FlowLabel = fl
		f.SetFLFlag()
	}

	if fid != 0 {
		f.SDFFilterID = fid
		f.SetBIDFlag()
	}

	return f
}

// HasBID reports whether BID flag is set.
func (f *SDFFilterFields) HasBID() bool {
	if f == nil {
		return false
	}

	return has5thBit(f.Flags)
}

// SetBIDFlag sets BID flag in SDFFilter.
func (f *SDFFilterFields) SetBIDFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x10
}

// HasFL reports whether FL flag is set.
func (f *SDFFilterFields) HasFL() bool {
	if f == nil {
		return false
	}

	return has4thBit(f.Flags)
}

// SetFLFlag sets FL flag in SDFFilter.
func (f *SDFFilterFields) SetFLFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x08
}

// HasSPI reports whether SPI flag is set.
func (f *SDFFilterFields) HasSPI() bool {
	if f == nil {
		return false
	}

	return has3rdBit(f.Flags)
}

// SetSPIFlag sets SPI flag in SDFFilter.
func (f *SDFFilterFields) SetSPIFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x04
}

// HasTTC reports whether TTC flag is set.
func (f *SDFFilterFields) HasTTC() bool {
	if f == nil {
		return false
	}

	return has2ndBit(f.Flags)
}

// SetTTCFlag sets TTC flag in SDFFilter.
func (f *SDFFilterFields) SetTTCFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x02
}

// HasFD reports whether FD flag is set.
func (f *SDFFilterFields) HasFD() bool {
	if f == nil {
		return false
	}

	return has1stBit(f.Flags)
}

// SetFDFlag sets FD flag in SDFFilter.
func (f *SDFFilterFields) SetFDFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x01
}

// ParseSDFFilterFields parses b into SDFFilterFields.
func ParseSDFFilterFields(b []byte) (*SDFFilterFields, error) {
	f := &SDFFilterFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary parses b into IE.
func (f *SDFFilterFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = SDFFilterFields{}

	l := len(b)
	if l < 2 {
		return io.ErrUnexpectedEOF
	}

	f.Flags = b[0]
	offset := 2 // 2nd octet is spare

	// TS 29.244 V20.1.0 clause 8.2.5.
	if f.HasFD() {
		value, err := takeLength16(b, &offset)
		if err != nil {
			return err
		}
		f.FDLength = uint16(len(value))
		f.FlowDescription = string(value)
	}

	if f.HasTTC() {
		if len(b[offset:]) < 2 {
			return io.ErrUnexpectedEOF
		}
		f.ToSTrafficClass = string(b[offset : offset+2])
		offset += 2
	}

	if f.HasSPI() {
		if len(b[offset:]) < 4 {
			return io.ErrUnexpectedEOF
		}
		f.SecurityParameterIndex = string(b[offset : offset+4])
		offset += 4
	}

	if f.HasFL() {
		if len(b[offset:]) < 3 {
			return io.ErrUnexpectedEOF
		}
		f.FlowLabel = string(b[offset : offset+3])
		offset += 3
	}

	if f.HasBID() {
		if len(b[offset:]) < 4 {
			return io.ErrUnexpectedEOF
		}
		f.SDFFilterID = binary.BigEndian.Uint32(b[offset : offset+4])
	}

	return nil
}

// Marshal returns the serialized bytes of SDFFilterFields.
func (f *SDFFilterFields) Marshal() ([]byte, error) {
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
func (f *SDFFilterFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if f.MarshalLen() > 65535 {
		return ErrInvalidLength
	}
	if f.HasFD() && len(f.FlowDescription) != int(f.FDLength) || f.HasTTC() && len(f.ToSTrafficClass) != 2 || f.HasSPI() && len(f.SecurityParameterIndex) != 4 || f.HasFL() && len(f.FlowLabel) != 3 {
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

	if f.HasTTC() {
		copy(b[offset:offset+2], f.ToSTrafficClass)
		offset += 2
	}

	if f.HasSPI() {
		copy(b[offset:offset+4], f.SecurityParameterIndex)
		offset += 4
	}

	if f.HasFL() {
		copy(b[offset:offset+3], f.FlowLabel)
		offset += 3
	}

	if f.HasBID() {
		binary.BigEndian.PutUint32(b[offset:offset+4], f.SDFFilterID)
	}

	return nil
}

// MarshalLen returns field length in integer.
func (f *SDFFilterFields) MarshalLen() int {
	if f == nil {
		return 0
	}

	l := 2 // Flag + Spare

	if f.HasFD() {
		l += 2 + int(f.FDLength)
	}

	if f.HasTTC() {
		l += 2
	}

	if f.HasSPI() {
		l += 4
	}

	if f.HasFL() {
		l += 3
	}

	if f.HasBID() {
		l += 4
	}

	return l
}
