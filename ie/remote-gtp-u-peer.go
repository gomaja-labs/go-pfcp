// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"time"
)

// NewRemoteGTPUPeer creates a new RemoteGTPUPeer IE.
func NewRemoteGTPUPeer(flags uint8, v4, v6 string, di uint8, ni []byte, recoveryTimestamp time.Time) *IE {
	fields := NewRemoteGTPUPeerFields(flags, v4, v6, di, ni, recoveryTimestamp)
	b, err := fields.Marshal()
	if err != nil {
		return newInvalidIE(RemoteGTPUPeer, ErrInvalidValue)
	}

	return New(RemoteGTPUPeer, b)
}

// RemoteGTPUPeer returns RemoteGTPUPeer in *RemoteGTPUPeerFields if the type of IE matches.
func (i *IE) RemoteGTPUPeer() (*RemoteGTPUPeerFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case RemoteGTPUPeer:
		fields, err := ParseRemoteGTPUPeerFields(i.Payload)
		if err != nil {
			return nil, err
		}

		return fields, nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// HasNI reports whether an IE has NI bit.
func (i *IE) HasNI() bool {
	if i == nil {
		return false
	}

	if i.Type != RemoteGTPUPeer {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has4thBit(i.Payload[0])
}

// HasDI reports whether an IE has DI bit.
func (i *IE) HasDI() bool {
	if i == nil {
		return false
	}

	if i.Type != RemoteGTPUPeer {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has3rdBit(i.Payload[0])
}

// RemoteGTPUPeerFields represents a fields contained in RemoteGTPUPeer IE.
type RemoteGTPUPeerFields struct {
	// Flags contains the presence and option bits.
	Flags uint8
	// IPv4Address is the IPv4 address, when present.
	IPv4Address net.IP
	// IPv6Address is the IPv6 address, when present.
	IPv6Address net.IP
	// DILength is the encoded destination interface length in octets.
	DILength uint16
	// DestinationInterface is the destination interface.
	DestinationInterface uint8
	// NILength is the encoded network instance length in octets.
	NILength uint16
	// NetworkInstance is the encoded Network Instance value (clause 8.2.4).
	NetworkInstance []byte
	// RecoveryTimestamp is the recovery timestamp of the peer GTP-U entity.
	RecoveryTimestamp time.Time
	invalidArgument   bool
}

// NewRemoteGTPUPeerFields creates a new RemoteGTPUPeerFields.
func NewRemoteGTPUPeerFields(flags uint8, v4, v6 string, di uint8, ni []byte, recoveryTimestamp time.Time) *RemoteGTPUPeerFields {
	f := &RemoteGTPUPeerFields{Flags: flags, RecoveryTimestamp: recoveryTimestamp}

	f.IPv4Address = net.ParseIP(v4).To4()
	f.IPv6Address = net.ParseIP(v6).To16()
	f.invalidArgument = v4 != "" && f.IPv4Address == nil || v6 != "" && f.IPv6Address == nil

	f.DestinationInterface = di
	if flags&4 != 0 {
		f.DILength = 1
	}
	f.NetworkInstance = ni
	f.NILength = uint16(len(ni))

	return f
}

// ParseRemoteGTPUPeerFields parses b into RemoteGTPUPeerFields.
func ParseRemoteGTPUPeerFields(b []byte) (*RemoteGTPUPeerFields, error) {
	f := &RemoteGTPUPeerFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary parses b into IE.
func (f *RemoteGTPUPeerFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = RemoteGTPUPeerFields{}

	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	f.Flags = b[0] & 0x1f
	offset := 1

	if has2ndBit(f.Flags) {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		f.IPv4Address = net.IP(b[offset : offset+4]).To4()
		offset += 4
	}

	if has1stBit(f.Flags) {
		if l < offset+16 {
			return io.ErrUnexpectedEOF
		}
		f.IPv6Address = net.IP(b[offset : offset+16]).To16()
		offset += 16
	}

	if has3rdBit(f.Flags) {
		if l < offset+2 {
			return io.ErrUnexpectedEOF
		}
		f.DILength = binary.BigEndian.Uint16(b[offset : offset+2])
		offset += 2

		if l < offset+int(f.DILength) {
			return io.ErrUnexpectedEOF
		}
		// TS 29.244 V20.1.0 clauses 8.2.70 and 8.2.24: the first
		// Destination Interface octet is mandatory; extension octets may follow.
		if f.DILength == 0 {
			return ErrInvalidLength
		}
		f.DestinationInterface = b[offset] & 0x0f
		offset += int(f.DILength)
	}

	if has4thBit(f.Flags) {
		if l < offset+2 {
			return io.ErrUnexpectedEOF
		}
		f.NILength = binary.BigEndian.Uint16(b[offset : offset+2])
		offset += 2

		if l < offset+int(f.NILength) {
			return io.ErrUnexpectedEOF
		}
		f.NetworkInstance = b[offset : offset+int(f.NILength)]
		offset += int(f.NILength)
	}

	// TS 29.244 V20.1.0 clause 8.2.70 references the timestamp in 8.2.114.
	if f.Flags&0x10 != 0 {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		v, e := New(RecoveryTimeStamp, b[offset:offset+4]).valueAs3GPPTimestamp()
		if e != nil {
			return e
		}
		f.RecoveryTimestamp = v
	}
	return nil
}

// Marshal returns the serialized bytes of RemoteGTPUPeerFields.
func (f *RemoteGTPUPeerFields) Marshal() ([]byte, error) {
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
func (f *RemoteGTPUPeerFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if f.invalidArgument || f.Flags&3 == 0 || f.Flags&^uint8(31) != 0 || f.Flags&2 == 0 && len(f.IPv4Address) != 0 || f.Flags&1 == 0 && len(f.IPv6Address) != 0 || f.Flags&4 == 0 && (f.DILength != 0 || f.DestinationInterface != 0) || f.Flags&8 == 0 && (f.NILength != 0 || len(f.NetworkInstance) != 0) || f.Flags&16 == 0 && !f.RecoveryTimestamp.IsZero() || f.DestinationInterface > 15 {
		return ErrInvalidValue
	}
	if f.MarshalLen() > 65535 {
		return ErrInvalidLength
	}
	if has3rdBit(f.Flags) && f.DILength == 0 {
		return ErrInvalidLength
	}
	if has4thBit(f.Flags) && len(f.NetworkInstance) != int(f.NILength) {
		return ErrInvalidLength
	}

	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}

	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	b[0] = f.Flags & 0x1f
	offset := 1

	if has2ndBit(f.Flags) {
		if f.IPv4Address.To4() == nil {
			return ErrMalformed
		}
		copy(b[offset:offset+4], f.IPv4Address.To4())
		offset += 4
	}

	if has1stBit(f.Flags) {
		if f.IPv6Address.To16() == nil {
			return ErrMalformed
		}
		copy(b[offset:offset+16], f.IPv6Address.To16())
		offset += 16
	}

	if has3rdBit(f.Flags) {
		binary.BigEndian.PutUint16(b[offset:offset+2], f.DILength)
		offset += 2
		// TS 29.244 V20.1.0 clauses 8.1.2 and 8.2.24: zero spare extension octets.
		clear(b[offset : offset+int(f.DILength)])
		b[offset] = f.DestinationInterface
		offset += int(f.DILength)
	}

	if has4thBit(f.Flags) {
		binary.BigEndian.PutUint16(b[offset:offset+2], f.NILength)
		offset += 2
		copy(b[offset:offset+int(f.NILength)], f.NetworkInstance)
		offset += int(f.NILength)
	}

	if f.Flags&0x10 != 0 {
		v := NewRecoveryTimeStamp(f.RecoveryTimestamp)
		if _, err := v.Marshal(); err != nil {
			return err
		}
		copy(b[offset:offset+4], v.Payload)
	}
	return nil
}

// MarshalLen returns field length in integer.
func (f *RemoteGTPUPeerFields) MarshalLen() int {
	if f == nil {
		return 0
	}

	l := 1
	if has2ndBit(f.Flags) {
		l += 4
	}
	if has1stBit(f.Flags) {
		l += 16
	}
	if has3rdBit(f.Flags) {
		l += 2 + int(f.DILength)
	}
	if has4thBit(f.Flags) {
		l += 2 + int(f.NILength)
	}

	if f.Flags&0x10 != 0 {
		l += 4
	}
	return l
}

// NewRemoteGTPUPeerFromFields encodes all current fields of RemoteGTPUPeer.
func NewRemoteGTPUPeerFromFields(fields *RemoteGTPUPeerFields) *IE {
	b, err := fields.Marshal()
	if err != nil {
		return newInvalidIE(RemoteGTPUPeer, ErrInvalidValue)
	}
	return New(RemoteGTPUPeer, b)
}

// Equal compares the complete peer fields, including invalid constructor state.
func (f *RemoteGTPUPeerFields) Equal(other *RemoteGTPUPeerFields) bool {
	if f == nil || other == nil {
		return f == other
	}
	return f.Flags == other.Flags && f.IPv4Address.Equal(other.IPv4Address) && f.IPv6Address.Equal(other.IPv6Address) && f.DILength == other.DILength && f.DestinationInterface == other.DestinationInterface && f.NILength == other.NILength && bytes.Equal(f.NetworkInstance, other.NetworkInstance) && f.RecoveryTimestamp.Equal(other.RecoveryTimestamp) && f.invalidArgument == other.invalidArgument
}
