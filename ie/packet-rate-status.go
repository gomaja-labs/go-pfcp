// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

// NewPacketRateStatus creates a PacketRateStatus IE. The timestamp uses the
// Unix window [1970-01-01, 2106-02-07T06:28:15.999999999Z], with nanosecond
// precision rounded to the nearest NTP fraction. Invalid arguments return an
// IE whose Marshal reports ErrInvalidValue. A no-limit record requires zero
// counters and time.Time{}. TS 29.244 V20.1.0 clause 8.2.139; RFC 5905 section 6.
func NewPacketRateStatus(flags uint8, ul, aul, dl, adl uint16, vtime time.Time) *IE {
	f := NewPacketRateStatusFields(flags, ul, aul, dl, adl, vtime)

	b, err := f.Marshal()
	if err != nil {
		return newInvalidIE(PacketRateStatus, err)
	}

	return New(PacketRateStatus, b)
}

// PacketRateStatus returns PacketRateStatus in structured format if the type of IE matches.
func (i *IE) PacketRateStatus() (*PacketRateStatusFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case PacketRateStatus:
		f, err := ParsePacketRateStatusFields(i.Payload)
		if err != nil {
			return nil, err
		}

		return f, nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// PacketRateStatusFields represents a f contained in PacketRateStatus IE.
type PacketRateStatusFields struct {
	// Flags contains the presence and option bits.
	Flags uint8
	// NumberOfRemainingUplinkPacketsAllowed is the number of remaining uplink packets allowed.
	NumberOfRemainingUplinkPacketsAllowed uint16
	// NumberOfRemainingAdditionalUplinkPacketsAllowed is the number of remaining additional uplink packets allowed.
	NumberOfRemainingAdditionalUplinkPacketsAllowed uint16
	// NumberOfRemainingDownlinkPacketsAllowed is the number of remaining downlink packets allowed.
	NumberOfRemainingDownlinkPacketsAllowed uint16
	// NumberOfRemainingAdditionalDownlinkPacketsAllowed is the number of remaining additional downlink packets allowed.
	NumberOfRemainingAdditionalDownlinkPacketsAllowed uint16
	// RateControlStatusValidityTime is the time until which the reported rate-control state is valid.
	RateControlStatusValidityTime time.Time
	// Retain fractions finer than time.Time can express when the decoded time
	// is unchanged. Editing the public time recomputes the NTP timestamp.
	wireTimestamp    uint64
	wireTime         time.Time
	hasWireTimestamp bool
}

// NewPacketRateStatusFields stores the supplied arguments without dropping absent fields.
// Marshal validates their flags and uses the era and fraction policy documented
// by NewPacketRateStatus (TS 29.244 V20.1.0 clause 8.2.139; RFC 5905 section 6).
func NewPacketRateStatusFields(flags uint8, ul, aul, dl, adl uint16, vtime time.Time) *PacketRateStatusFields {
	f := &PacketRateStatusFields{
		Flags:                                 flags,
		NumberOfRemainingUplinkPacketsAllowed: ul,
		NumberOfRemainingAdditionalUplinkPacketsAllowed:   aul,
		NumberOfRemainingDownlinkPacketsAllowed:           dl,
		NumberOfRemainingAdditionalDownlinkPacketsAllowed: adl,
		RateControlStatusValidityTime:                     vtime,
	}

	return f
}

// HasAPR reports whether APR flag is set.
func (f *PacketRateStatusFields) HasAPR() bool {
	if f == nil {
		return false
	}

	return has3rdBit(f.Flags)
}

// SetAPRFlag sets APR flag in PacketRateStatus.
func (f *PacketRateStatusFields) SetAPRFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x04
}

// HasDL reports whether DL flag is set.
func (f *PacketRateStatusFields) HasDL() bool {
	if f == nil {
		return false
	}

	return has2ndBit(f.Flags)
}

// SetDLFlag sets DL flag in PacketRateStatus.
func (f *PacketRateStatusFields) SetDLFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x02
}

// HasUL reports whether UL flag is set.
func (f *PacketRateStatusFields) HasUL() bool {
	if f == nil {
		return false
	}

	return has1stBit(f.Flags)
}

// SetULFlag sets UL flag in PacketRateStatus.
func (f *PacketRateStatusFields) SetULFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x01
}

// ParsePacketRateStatusFields parses b into PacketRateStatusFields.
func ParsePacketRateStatusFields(b []byte) (*PacketRateStatusFields, error) {
	f := &PacketRateStatusFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary parses b into IE.
func (f *PacketRateStatusFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = PacketRateStatusFields{}

	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	f.Flags = b[0] & 0x07
	offset := 1

	if f.HasUL() {
		if l < offset+2 {
			return io.ErrUnexpectedEOF
		}
		f.NumberOfRemainingUplinkPacketsAllowed = binary.BigEndian.Uint16(b[offset : offset+2])
		offset += 2

		if f.HasAPR() {
			if l < offset+2 {
				return io.ErrUnexpectedEOF
			}
			f.NumberOfRemainingAdditionalUplinkPacketsAllowed = binary.BigEndian.Uint16(b[offset : offset+2])
			offset += 2
		}
	}

	if f.HasDL() {
		if l < offset+2 {
			return io.ErrUnexpectedEOF
		}
		f.NumberOfRemainingDownlinkPacketsAllowed = binary.BigEndian.Uint16(b[offset : offset+2])
		offset += 2

		if f.HasAPR() {
			if l < offset+2 {
				return io.ErrUnexpectedEOF
			}
			f.NumberOfRemainingAdditionalDownlinkPacketsAllowed = binary.BigEndian.Uint16(b[offset : offset+2])
			offset += 2
		}
	}

	// TS 29.244 V20.1.0 clause 8.2.139: validity time exists only with UL or DL.
	if f.HasUL() || f.HasDL() {
		if l < offset+8 {
			return io.ErrUnexpectedEOF
		}
		value := binary.BigEndian.Uint64(b[offset : offset+8])
		// RFC 5905 section 6: high 32 bits are seconds, low 32 bits are fractions.
		seconds := int64(uint32(value>>32) - 2208988800)
		nanos := (uint64(uint32(value))*1000000000 + (1 << 31)) >> 32
		f.RateControlStatusValidityTime = time.Unix(seconds, int64(nanos)).UTC()
		f.wireTimestamp = value
		f.wireTime = f.RateControlStatusValidityTime
		f.hasWireTimestamp = true
	}

	return nil
}

// Marshal returns the serialized bytes of PacketRateStatusFields.
func (f *PacketRateStatusFields) Marshal() ([]byte, error) {
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
func (f *PacketRateStatusFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	// TS 29.244 V20.1.0 clause 8.2.139: do not discard unflagged arguments.
	if f.Flags&^uint8(7) != 0 {
		return fmt.Errorf("%w: packet-rate flags include spare bits", ErrInvalidValue)
	}
	if !f.HasUL() && (f.NumberOfRemainingUplinkPacketsAllowed != 0 || f.NumberOfRemainingAdditionalUplinkPacketsAllowed != 0) {
		return fmt.Errorf("%w: uplink counters require UL", ErrInvalidValue)
	}
	if !f.HasDL() && (f.NumberOfRemainingDownlinkPacketsAllowed != 0 || f.NumberOfRemainingAdditionalDownlinkPacketsAllowed != 0) {
		return fmt.Errorf("%w: downlink counters require DL", ErrInvalidValue)
	}
	if !f.HasAPR() && (f.NumberOfRemainingAdditionalUplinkPacketsAllowed != 0 || f.NumberOfRemainingAdditionalDownlinkPacketsAllowed != 0) {
		return fmt.Errorf("%w: additional counters require APR", ErrInvalidValue)
	}
	if !f.HasUL() && !f.HasDL() && !f.RateControlStatusValidityTime.IsZero() {
		return fmt.Errorf("%w: validity time requires UL or DL", ErrInvalidValue)
	}

	var timestamp uint64
	if f.HasUL() || f.HasDL() {
		if f.hasWireTimestamp && f.RateControlStatusValidityTime.Equal(f.wireTime) {
			timestamp = f.wireTimestamp
		} else {
			seconds := f.RateControlStatusValidityTime.Unix()
			if seconds < 0 || seconds > 1<<32-1 {
				return fmt.Errorf("%w: packet-rate timestamp is outside the Unix era window", ErrInvalidValue)
			}
			fraction := (uint64(f.RateControlStatusValidityTime.Nanosecond())<<32 + 500000000) / 1000000000
			timestamp = uint64(uint32(seconds+2208988800))<<32 | fraction
		}
	}

	if f.MarshalLen() > 65535 {
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
	offset := 1

	if f.HasUL() {
		if l < offset+2 {
			return io.ErrUnexpectedEOF
		}
		binary.BigEndian.PutUint16(b[offset:offset+2], f.NumberOfRemainingUplinkPacketsAllowed)
		offset += 2

		if f.HasAPR() {
			if l < offset+2 {
				return io.ErrUnexpectedEOF
			}
			binary.BigEndian.PutUint16(b[offset:offset+2], f.NumberOfRemainingAdditionalUplinkPacketsAllowed)
			offset += 2
		}
	}

	if f.HasDL() {
		if l < offset+2 {
			return io.ErrUnexpectedEOF
		}
		binary.BigEndian.PutUint16(b[offset:offset+2], f.NumberOfRemainingDownlinkPacketsAllowed)
		offset += 2

		if f.HasAPR() {
			if l < offset+2 {
				return io.ErrUnexpectedEOF
			}
			binary.BigEndian.PutUint16(b[offset:offset+2], f.NumberOfRemainingAdditionalDownlinkPacketsAllowed)
			offset += 2
		}
	}

	if f.HasUL() || f.HasDL() {
		if l < offset+8 {
			return io.ErrUnexpectedEOF
		}
		binary.BigEndian.PutUint64(b[offset:offset+8], timestamp)
	}

	return nil
}

// MarshalLen returns field length in integer.
func (f *PacketRateStatusFields) MarshalLen() int {
	if f == nil {
		return 0
	}

	l := 1
	if f.HasUL() || f.HasDL() {
		l += 8
	}
	if f.HasUL() {
		l += 2
		if f.HasAPR() {
			l += 2
		}
	}
	if f.HasDL() {
		l += 2
		if f.HasAPR() {
			l += 2
		}
	}

	return l
}

// Equal reports whether two PacketRateStatusFields values have the same encoding.
// This compares retained NTP fractions at their full wire precision.
func (f *PacketRateStatusFields) Equal(other *PacketRateStatusFields) bool {
	if f == nil || other == nil {
		return f == other
	}
	a, ea := f.Marshal()
	b, eb := other.Marshal()
	return ea == nil && eb == nil && bytes.Equal(a, b)
}
