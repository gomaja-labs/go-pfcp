// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
)

// Time Unit definitions
const (
	TimeUnitMinute   uint8 = 0
	TimeUnit6Minutes uint8 = 1
	TimeUnitHour     uint8 = 2
	TimeUnitDay      uint8 = 3
	TimeUnitWeek     uint8 = 4
)

// NewPacketRate creates a new PacketRate IE.
func NewPacketRate(flags uint8, ulunit uint8, ulpackets uint16, dlunit uint8, dlpackets uint16) *IE {
	fields := NewPacketRateFields(flags, ulunit, ulpackets, dlunit, dlpackets)
	b, err := fields.Marshal()
	if err != nil {
		return newInvalidIE(PacketRate, ErrInvalidValue)
	}

	return New(PacketRate, b)
}

// PacketRate returns PacketRate in *PacketRateFields if the type of IE matches.
func (i *IE) PacketRate() (*PacketRateFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case PacketRate:
		f, err := ParsePacketRateFields(i.Payload)
		if err != nil {
			return nil, err
		}
		return f, nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}

}

// HasDLPR reports whether an IE has DLPR bit.
func (i *IE) HasDLPR() bool {
	if i == nil {
		return false
	}

	if i.Type != PacketRate {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has2ndBit(i.Payload[0])
}

// HasULPR reports whether an IE has ULPR bit.
func (i *IE) HasULPR() bool {
	if i == nil {
		return false
	}

	if i.Type != PacketRate {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has1stBit(i.Payload[0])
}

// PacketRateFields represents a fields contained in PacketRate IE.
type PacketRateFields struct {
	Flags                        uint8
	UplinkTimeUnit               uint8
	DownlinkTimeUnit             uint8
	UplinkPacketRate             uint16
	DownlinkPacketRate           uint16
	AdditionalUplinkTimeUnit     uint8
	AdditionalDownlinkTimeUnit   uint8
	AdditionalUplinkPacketRate   uint16
	AdditionalDownlinkPacketRate uint16
}

// NewPacketRateFields creates a new PacketRateFields.
func NewPacketRateFields(flags uint8, ulunit uint8, ulpackets uint16, dlunit uint8, dlpackets uint16) *PacketRateFields {
	f := &PacketRateFields{Flags: flags, UplinkTimeUnit: ulunit, UplinkPacketRate: ulpackets, DownlinkTimeUnit: dlunit, DownlinkPacketRate: dlpackets}

	return f
}

// ParsePacketRateFields parses b into PacketRateFields.
func ParsePacketRateFields(b []byte) (*PacketRateFields, error) {
	f := &PacketRateFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary parses b into IE.
func (f *PacketRateFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = PacketRateFields{}

	l := len(b)
	// TS 29.244 V20.1.0 clause 8.2.63: one fixed flags octet.
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	f.Flags = b[0] & 7
	offset := 1

	if has1stBit(f.Flags) {
		if l < offset+3 {
			return io.ErrUnexpectedEOF
		}
		f.UplinkTimeUnit = packetRateTimeUnit(b[offset])
		f.UplinkPacketRate = binary.BigEndian.Uint16(b[offset+1 : offset+3])
		offset += 3
	}

	if has2ndBit(f.Flags) {
		if l < offset+3 {
			return io.ErrUnexpectedEOF
		}
		f.DownlinkTimeUnit = packetRateTimeUnit(b[offset])
		f.DownlinkPacketRate = binary.BigEndian.Uint16(b[offset+1 : offset+3])
		offset += 3
	}

	// TS 29.244 V20.1.0 clause 8.2.63: APRC uses the same direction flags.
	for _, v := range []struct {
		bit  uint8
		unit *uint8
		rate *uint16
	}{{1, &f.AdditionalUplinkTimeUnit, &f.AdditionalUplinkPacketRate}, {2, &f.AdditionalDownlinkTimeUnit, &f.AdditionalDownlinkPacketRate}} {
		if f.Flags&4 != 0 && f.Flags&v.bit != 0 {
			if l < offset+3 {
				return io.ErrUnexpectedEOF
			}
			*v.unit = packetRateTimeUnit(b[offset])
			*v.rate = binary.BigEndian.Uint16(b[offset+1 : offset+3])
			offset += 3
		}
	}
	return nil
}

// Marshal returns the serialized bytes of PacketRateFields.
func (f *PacketRateFields) Marshal() ([]byte, error) {
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
func (f *PacketRateFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if f.Flags == 0 || f.Flags&^uint8(7) != 0 || f.Flags&1 == 0 && (f.UplinkTimeUnit != 0 || f.UplinkPacketRate != 0) || f.Flags&2 == 0 && (f.DownlinkTimeUnit != 0 || f.DownlinkPacketRate != 0) || f.Flags&5 != 5 && (f.AdditionalUplinkTimeUnit != 0 || f.AdditionalUplinkPacketRate != 0) || f.Flags&6 != 6 && (f.AdditionalDownlinkTimeUnit != 0 || f.AdditionalDownlinkPacketRate != 0) {
		return ErrInvalidValue
	}
	if f.UplinkTimeUnit > 4 || f.DownlinkTimeUnit > 4 || f.AdditionalUplinkTimeUnit > 4 || f.AdditionalDownlinkTimeUnit > 4 {
		return ErrInvalidValue
	}
	if f.MarshalLen() > 65535 {
		return ErrInvalidLength
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}

	l := len(b)
	// TS 29.244 V20.1.0 clause 8.2.63: one fixed flags octet.
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	b[0] = f.Flags & 7
	offset := 1

	if has1stBit(f.Flags) {
		b[offset] = f.UplinkTimeUnit & 7
		binary.BigEndian.PutUint16(b[offset+1:offset+3], f.UplinkPacketRate)
		offset += 3
	}

	if has2ndBit(f.Flags) {
		b[offset] = f.DownlinkTimeUnit & 7
		binary.BigEndian.PutUint16(b[offset+1:offset+3], f.DownlinkPacketRate)
		offset += 3
	}

	for _, v := range []struct {
		bit  uint8
		unit uint8
		rate uint16
	}{{1, f.AdditionalUplinkTimeUnit, f.AdditionalUplinkPacketRate}, {2, f.AdditionalDownlinkTimeUnit, f.AdditionalDownlinkPacketRate}} {
		if f.Flags&4 != 0 && f.Flags&v.bit != 0 {
			b[offset] = v.unit & 7
			binary.BigEndian.PutUint16(b[offset+1:offset+3], v.rate)
			offset += 3
		}
	}
	return nil
}

// MarshalLen returns field length in integer.
func (f *PacketRateFields) MarshalLen() int {
	if f == nil {
		return 0
	}

	l := 1
	if has1stBit(f.Flags) {
		l += 3
	}
	if has2ndBit(f.Flags) {
		l += 3
	}

	if f.Flags&4 != 0 {
		if f.Flags&1 != 0 {
			l += 3
		}
		if f.Flags&2 != 0 {
			l += 3
		}
	}
	return l
}

// NewPacketRateFromFields encodes all current fields of PacketRate.
func NewPacketRateFromFields(fields *PacketRateFields) *IE {
	b, err := fields.Marshal()
	if err != nil {
		return newInvalidIE(PacketRate, ErrInvalidValue)
	}
	return New(PacketRate, b)
}

// TS 29.244 V20.1.0 Table 8.2.63.1 interprets reserved time units as minutes.
func packetRateTimeUnit(v uint8) uint8 {
	v &= 7
	if v > 4 {
		return 0
	}
	return v
}
