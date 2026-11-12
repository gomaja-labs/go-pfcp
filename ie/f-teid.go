// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// NewFTEID creates a new FTEID IE.
func NewFTEID(flags uint8, teid uint32, v4, v6 net.IP, chid uint8) *IE {
	fields := NewFTEIDFields(flags, teid, v4, v6, chid)

	b, err := fields.Marshal()
	if err != nil {
		return newInvalidIE(FTEID, fmt.Errorf("%w: %w", ErrInvalidValue, err))
	}

	return New(FTEID, b)
}

// FTEID returns FTEID in structured format if the type of IE matches.
func (i *IE) FTEID() (*FTEIDFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case FTEID:
		fields, err := ParseFTEIDFields(i.Payload)
		if err != nil {
			return nil, err
		}

		return fields, nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// FTEIDFields represents a fields contained in FTEID IE.
type FTEIDFields struct {
	// Flags contains the presence and option bits.
	Flags uint8
	// TEID is the tunnel endpoint identifier.
	TEID uint32
	// IPv4Address is the IPv4 address, when present.
	IPv4Address net.IP
	// IPv6Address is the IPv6 address, when present.
	IPv6Address net.IP
	// ChooseID is the choose identifier, when its presence bit is set.
	ChooseID uint8
}

// NewFTEIDFields creates a new NewFTEIDFields.
func NewFTEIDFields(flags uint8, teid uint32, v4, v6 net.IP, chid uint8) *FTEIDFields {
	return &FTEIDFields{
		Flags:       flags,
		TEID:        teid,
		IPv4Address: v4,
		IPv6Address: v6,
		ChooseID:    chid,
	}
}

// HasChID reports whether CHID flag is set.
func (f *FTEIDFields) HasChID() bool {
	if f == nil {
		return false
	}

	return has4thBit(f.Flags)
}

// SetChIDFlag sets CHID flag in FTEID.
func (f *FTEIDFields) SetChIDFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x08
}

// HasCh reports whether CH flag is set.
func (f *FTEIDFields) HasCh() bool {
	if f == nil {
		return false
	}

	return has3rdBit(f.Flags)
}

// SetChFlag sets CH flag in FTEID.
func (f *FTEIDFields) SetChFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x04
}

// HasIPv6 reports whether IPv6 flag is set.
func (f *FTEIDFields) HasIPv6() bool {
	if f == nil {
		return false
	}

	return has2ndBit(f.Flags)
}

// SetIPv6Flag sets IPv6 flag in FTEID.
func (f *FTEIDFields) SetIPv6Flag() {
	if f == nil {
		return
	}

	f.Flags |= 0x02
}

// HasIPv4 reports whether IPv4 flag is set.
func (f *FTEIDFields) HasIPv4() bool {
	if f == nil {
		return false
	}

	return has1stBit(f.Flags)
}

// SetIPv4Flag sets IPv4 flag in FTEID.
func (f *FTEIDFields) SetIPv4Flag() {
	if f == nil {
		return
	}

	f.Flags |= 0x01
}

// ParseFTEIDFields parses b into FTEIDFields.
func ParseFTEIDFields(b []byte) (*FTEIDFields, error) {
	f := &FTEIDFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary parses b into IE.
func (f *FTEIDFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = FTEIDFields{}

	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	// TS 29.244 V20.1.0 clauses 8.1.2 and 8.2.3: ignore received spare bits.
	f.Flags = b[0] & 15
	offset := 1

	if f.HasCh() {
		if f.HasChID() {
			if l < offset+1 {
				return io.ErrUnexpectedEOF
			}
			f.ChooseID = b[offset]
		}
		return nil
	}

	if l < offset+4 {
		return io.ErrUnexpectedEOF
	}
	f.TEID = binary.BigEndian.Uint32(b[offset : offset+4])
	offset += 4

	if f.HasIPv4() {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		f.IPv4Address = net.IP(b[offset : offset+4])
		offset += 4
	}
	if f.HasIPv6() {
		if l < offset+16 {
			return io.ErrUnexpectedEOF
		}
		f.IPv6Address = net.IP(b[offset : offset+16])
	}

	return nil
}

// Marshal returns the serialized bytes of FTEIDFields.
func (f *FTEIDFields) Marshal() ([]byte, error) {
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
func (f *FTEIDFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	// TS 29.244 V20.1.0 clause 8.2.3: require an address-family flag,
	// reserve CHID for CH, and send the spare bits as zero.
	if f.Flags&^uint8(15) != 0 || f.Flags&3 == 0 || f.HasChID() && !f.HasCh() {
		return ErrInvalidValue
	}
	// With CH clear each flagged address is mandatory.
	if !f.HasCh() && (f.HasIPv4() && f.IPv4Address.To4() == nil || f.HasIPv6() && (f.IPv6Address.To16() == nil || f.IPv6Address.To4() != nil)) {
		return ErrInvalidValue
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

	if f.HasCh() {
		if f.HasChID() {
			if l < offset+1 {
				return io.ErrUnexpectedEOF
			}
			b[offset] = f.ChooseID
		}
		return nil
	}

	if l < offset+4 {
		return io.ErrUnexpectedEOF
	}
	binary.BigEndian.PutUint32(b[offset:offset+4], f.TEID)
	offset += 4

	if f.HasIPv4() && f.IPv4Address != nil {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		copy(b[offset:offset+4], f.IPv4Address.To4())
		offset += 4
	}
	if f.HasIPv6() && f.IPv6Address != nil {
		if l < offset+16 {
			return io.ErrUnexpectedEOF
		}
		copy(b[offset:offset+16], f.IPv6Address.To16())
	}

	return nil
}

// MarshalLen returns field length in integer.
func (f *FTEIDFields) MarshalLen() int {
	if f == nil {
		return 0
	}

	l := 1
	if f.HasCh() {
		if f.HasChID() {
			l++
		}
		return l
	}

	l += 4
	if f.HasIPv4() {
		l += 4
	}
	if f.HasIPv6() {
		l += 16
	}
	return l
}
