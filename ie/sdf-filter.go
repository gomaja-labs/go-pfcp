// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
)

// NewSDFFilter creates a new SDFFilter IE.
func NewSDFFilter(fd string, ttc, spi, fl []byte, fid uint32) *IE {
	f := NewSDFFilterFields(fd, ttc, spi, fl, fid)
	b, err := f.Marshal()
	if err != nil {
		return newInvalidIE(SDFFilter, ErrInvalidValue)
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
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// SDFFilterFields represents a fields contained in SDFFilter IE.
type SDFFilterFields struct {
	// Flags contains the presence and option bits.
	Flags uint8
	// FDLength is the encoded flow description length in octets.
	FDLength uint16
	// FlowDescription is the flow description.
	FlowDescription string
	// ToSTrafficClass contains the ToS/traffic class and mask octets.
	ToSTrafficClass []byte
	// SecurityParameterIndex contains the four security parameter index octets.
	SecurityParameterIndex []byte
	// FlowLabel contains the three flow label octets.
	FlowLabel []byte
	// SDFFilterID is the sdf filter ID.
	SDFFilterID uint32
	// MultiplexedMedia is the multiplexed media.
	MultiplexedMedia []SDFMultiplexedMedia
}

// NewSDFFilterFields creates a new NewSDFFilterFields.
func NewSDFFilterFields(fd string, ttc, spi, fl []byte, fid uint32) *SDFFilterFields {
	f := &SDFFilterFields{}
	if fd != "" {
		f.FlowDescription = fd
		f.FDLength = uint16(len([]byte(fd)))
		f.SetFDFlag()
	}

	if len(ttc) != 0 {
		f.ToSTrafficClass = ttc
		f.SetTTCFlag()
	}

	if len(spi) != 0 {
		f.SecurityParameterIndex = spi
		f.SetSPIFlag()
	}

	if len(fl) != 0 {
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
// Decoding may borrow b; keep b unchanged while the result is in use.
func ParseSDFFilterFields(b []byte) (*SDFFilterFields, error) {
	f := &SDFFilterFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary parses b into IE.
// Decoding may borrow b; keep b unchanged while the result is in use.
func (f *SDFFilterFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = SDFFilterFields{}

	l := len(b)
	if l < 2 {
		return io.ErrUnexpectedEOF
	}

	f.Flags = b[0] & 0x3f
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
		f.ToSTrafficClass = append([]byte(nil), b[offset:offset+2]...)
		offset += 2
	}

	if f.HasSPI() {
		if len(b[offset:]) < 4 {
			return io.ErrUnexpectedEOF
		}
		f.SecurityParameterIndex = append([]byte(nil), b[offset:offset+4]...)
		offset += 4
	}

	if f.HasFL() {
		if len(b[offset:]) < 3 {
			return io.ErrUnexpectedEOF
		}
		f.FlowLabel = append([]byte(nil), b[offset:offset+3]...)
		f.FlowLabel[0] &= 0x0f // TS 29.244 V20.1.0 clauses 8.2.5 and 8.1.2.
		offset += 3
	}

	if f.HasBID() {
		if len(b[offset:]) < 4 {
			return io.ErrUnexpectedEOF
		}
		f.SDFFilterID = binary.BigEndian.Uint32(b[offset : offset+4])
		offset += 4
	}

	if f.Flags&0x20 != 0 {
		if len(b) < offset+1 {
			return io.ErrUnexpectedEOF
		}
		n := int(b[offset])
		offset++
		if n == 0 {
			return ErrInvalidLength
		}
		f.MultiplexedMedia = make([]SDFMultiplexedMedia, 0, n)
		for range n {
			v, e := takeLength16(b, &offset)
			if e != nil {
				return e
			}
			entry := SDFMultiplexedMedia{}
			if e := entry.UnmarshalBinary(v); e != nil {
				return e
			}
			f.MultiplexedMedia = append(f.MultiplexedMedia, entry)
		}
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
	if f.Flags&^uint8(0x3f) != 0 || f.Flags&1 == 0 && (f.FDLength != 0 || f.FlowDescription != "") || f.Flags&2 == 0 && len(f.ToSTrafficClass) != 0 || f.Flags&4 == 0 && len(f.SecurityParameterIndex) != 0 || f.Flags&8 == 0 && len(f.FlowLabel) != 0 || f.Flags&16 == 0 && f.SDFFilterID != 0 {
		return ErrInvalidValue
	}
	if f.Flags&0x20 != 0 && (len(f.MultiplexedMedia) == 0 || len(f.MultiplexedMedia) > 255) {
		return ErrInvalidLength
	}
	if f.Flags&0x20 == 0 && len(f.MultiplexedMedia) != 0 {
		return ErrInvalidValue
	}
	if f.MarshalLen() > 65535 {
		return ErrInvalidLength
	}
	if f.HasFD() && len(f.FlowDescription) != int(f.FDLength) || f.HasTTC() && len(f.ToSTrafficClass) != 2 || f.HasSPI() && len(f.SecurityParameterIndex) != 4 || f.HasFL() && len(f.FlowLabel) != 3 {
		return ErrInvalidLength
	}
	if f.HasFL() && f.FlowLabel[0]&0xf0 != 0 {
		return ErrInvalidValue
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}

	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	b[0] = f.Flags & 0x3f
	b[1] = 0
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
		offset += 4
	}

	if f.Flags&0x20 != 0 {
		b[offset] = uint8(len(f.MultiplexedMedia))
		offset++
		for k := range f.MultiplexedMedia {
			entry := &f.MultiplexedMedia[k]
			n := entry.MarshalLen()
			binary.BigEndian.PutUint16(b[offset:offset+2], uint16(n))
			offset += 2
			if e := entry.MarshalTo(b[offset : offset+n]); e != nil {
				return e
			}
			offset += n
		}
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

	if f.Flags&0x20 != 0 {
		l++
		for k := range f.MultiplexedMedia {
			l += 2 + f.MultiplexedMedia[k].MarshalLen()
		}
	}
	return l
}

// NewSDFFilterFromFields encodes all current fields of SDFFilter.
func NewSDFFilterFromFields(fields *SDFFilterFields) *IE {
	b, err := fields.Marshal()
	if err != nil {
		return newInvalidIE(SDFFilter, ErrInvalidValue)
	}
	return New(SDFFilter, b)
}

// HasSMMII reports whether multiplexed media identification is present.
func (f *SDFFilterFields) HasSMMII() bool { return f != nil && f.Flags&0x20 != 0 }

// SetSMMIIFlag enables the multiplexed media list.
func (f *SDFFilterFields) SetSMMIIFlag() {
	if f != nil {
		f.Flags |= 0x20
	}
}
