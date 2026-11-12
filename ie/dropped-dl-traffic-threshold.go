// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
)

// DroppedDLTrafficThresholdFields contains independently present packet and byte limits.
// TS 29.244 V20.1.0 clause 8.2.49.
type DroppedDLTrafficThresholdFields struct {
	// Flags contains the presence and option bits.
	Flags uint8
	// DownlinkPackets is the downlink packets.
	DownlinkPackets uint64
	// DownlinkDataBytes is the downlink data bytes.
	DownlinkDataBytes uint64
}

// NewDroppedDLTrafficThreshold creates independently present Unsigned64 thresholds.
func NewDroppedDLTrafficThreshold(dlpa, dlby bool, packets, bytes uint64) *IE {
	f := &DroppedDLTrafficThresholdFields{DownlinkPackets: packets, DownlinkDataBytes: bytes}
	if dlpa {
		f.Flags |= 1
	}
	if dlby {
		f.Flags |= 2
	}
	b, e := f.Marshal()
	if e != nil {
		return newInvalidIE(DroppedDLTrafficThreshold, ErrInvalidValue)
	}
	return New(DroppedDLTrafficThreshold, b)
}

// DroppedDLTrafficThreshold returns both limits and their presence flags.
func (i *IE) DroppedDLTrafficThreshold() (*DroppedDLTrafficThresholdFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}
	if i.Type != DroppedDLTrafficThreshold {
		return nil, &InvalidTypeError{Type: i.Type}
	}
	return ParseDroppedDLTrafficThresholdFields(i.Payload)
}

// ParseDroppedDLTrafficThresholdFields decodes the thresholds.
func ParseDroppedDLTrafficThresholdFields(b []byte) (*DroppedDLTrafficThresholdFields, error) {
	f := &DroppedDLTrafficThresholdFields{}
	if e := f.UnmarshalBinary(b); e != nil {
		return nil, e
	}
	return f, nil
}

// UnmarshalBinary decodes thresholds and ignores trailing extension octets.
func (f *DroppedDLTrafficThresholdFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = DroppedDLTrafficThresholdFields{}
	if len(b) < 1 {
		return io.ErrUnexpectedEOF
	}
	f.Flags = b[0] & 3
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}
	o := 1
	if f.Flags&1 != 0 {
		f.DownlinkPackets = binary.BigEndian.Uint64(b[o : o+8])
		o += 8
	}
	if f.Flags&2 != 0 {
		f.DownlinkDataBytes = binary.BigEndian.Uint64(b[o : o+8])
	}
	return nil
}

// Marshal serializes the thresholds.
func (f *DroppedDLTrafficThresholdFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	b := make([]byte, f.MarshalLen())
	if e := f.MarshalTo(b); e != nil {
		return nil, e
	}
	return b, nil
}

// MarshalLen returns the encoded size.
func (f *DroppedDLTrafficThresholdFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	n := 1
	if f.Flags&1 != 0 {
		n += 8
	}
	if f.Flags&2 != 0 {
		n += 8
	}
	return n
}

// MarshalTo writes both counters according to their presence flags.
func (f *DroppedDLTrafficThresholdFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if f.Flags&^uint8(3) != 0 || f.Flags&1 == 0 && f.DownlinkPackets != 0 || f.Flags&2 == 0 && f.DownlinkDataBytes != 0 {
		return ErrInvalidValue
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}
	b[0] = f.Flags
	o := 1
	if f.Flags&1 != 0 {
		binary.BigEndian.PutUint64(b[o:o+8], f.DownlinkPackets)
		o += 8
	}
	if f.Flags&2 != 0 {
		binary.BigEndian.PutUint64(b[o:o+8], f.DownlinkDataBytes)
	}
	return nil
}

// HasDLBY reports presence of the byte threshold.
func (i *IE) HasDLBY() bool { f, e := i.DroppedDLTrafficThreshold(); return e == nil && f.Flags&2 != 0 }

// HasDLPA reports presence of the packet threshold.
func (i *IE) HasDLPA() bool { f, e := i.DroppedDLTrafficThreshold(); return e == nil && f.Flags&1 != 0 }
