// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
)

// NewQoSMonitoringMeasurement creates a new QoSMonitoringMeasurement IE.
func NewQoSMonitoringMeasurement(flags uint8, dl, ul, rp uint32) *IE {
	fields := NewQoSMonitoringMeasurementFields(flags, dl, ul, rp)

	b, err := fields.Marshal()
	if err != nil {
		return newInvalidIE(QoSMonitoringMeasurement, ErrInvalidValue)
	}

	return New(QoSMonitoringMeasurement, b)
}

// QoSMonitoringMeasurement returns QoSMonitoringMeasurement in structured format if the type of IE matches.
func (i *IE) QoSMonitoringMeasurement() (*QoSMonitoringMeasurementFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case QoSMonitoringMeasurement:
		fields, err := ParseQoSMonitoringMeasurementFields(i.Payload)
		if err != nil {
			return nil, err
		}

		return fields, nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// QoSMonitoringMeasurementFields represents a fields contained in QoSMonitoringMeasurement IE.
type QoSMonitoringMeasurementFields struct {
	Flags                         uint8
	DownlinkPacketDelay           uint32
	UplinkPacketDelay             uint32
	RoundTripPacketDelay          uint32
	DownlinkCongestionInformation uint16 // hundredths of one percent
	UplinkCongestionInformation   uint16 // hundredths of one percent
	AverageDownlinkDataRate       uint64 // kbps, five octets
	AverageUplinkDataRate         uint64 // kbps, five octets
}

// NewQoSMonitoringMeasurementFields creates a new NewQoSMonitoringMeasurementFields.
func NewQoSMonitoringMeasurementFields(flags uint8, dl, ul, rp uint32) *QoSMonitoringMeasurementFields {
	return &QoSMonitoringMeasurementFields{
		Flags:                flags,
		DownlinkPacketDelay:  dl,
		UplinkPacketDelay:    ul,
		RoundTripPacketDelay: rp,
	}
}

// HasPLMF reports whether PLMF flag is set.
func (f *QoSMonitoringMeasurementFields) HasPLMF() bool {
	if f == nil {
		return false
	}

	return has4thBit(f.Flags)
}

// SetPLMFFlag sets PLMF flag in QoSMonitoringMeasurement.
func (f *QoSMonitoringMeasurementFields) SetPLMFFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x08
}

// HasRP reports whether RP flag is set.
func (f *QoSMonitoringMeasurementFields) HasRP() bool {
	if f == nil {
		return false
	}

	return has3rdBit(f.Flags)
}

// SetRPFlag sets RP flag in QoSMonitoringMeasurement.
func (f *QoSMonitoringMeasurementFields) SetRPFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x04
}

// HasUL reports whether UL flag is set.
func (f *QoSMonitoringMeasurementFields) HasUL() bool {
	if f == nil {
		return false
	}

	return has2ndBit(f.Flags)
}

// SetULFlag sets UL flag in QoSMonitoringMeasurement.
func (f *QoSMonitoringMeasurementFields) SetULFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x02
}

// HasDL reports whether DL flag is set.
func (f *QoSMonitoringMeasurementFields) HasDL() bool {
	if f == nil {
		return false
	}

	return has1stBit(f.Flags)
}

// SetDLFlag sets DL flag in QoSMonitoringMeasurement.
func (f *QoSMonitoringMeasurementFields) SetDLFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x01
}

// ParseQoSMonitoringMeasurementFields parses b into QoSMonitoringMeasurementFields.
func ParseQoSMonitoringMeasurementFields(b []byte) (*QoSMonitoringMeasurementFields, error) {
	f := &QoSMonitoringMeasurementFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary parses b into IE.
func (f *QoSMonitoringMeasurementFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = QoSMonitoringMeasurementFields{}

	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	f.Flags = b[0]
	offset := 1

	if f.HasDL() {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		f.DownlinkPacketDelay = binary.BigEndian.Uint32(b[offset : offset+4])
		offset += 4
	}

	if f.HasUL() {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		f.UplinkPacketDelay = binary.BigEndian.Uint32(b[offset : offset+4])
		offset += 4
	}

	if f.HasRP() {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		f.RoundTripPacketDelay = binary.BigEndian.Uint32(b[offset : offset+4])
		offset += 4
	}

	for _, item := range []struct {
		bit uint8
		dst *uint16
	}{{0x10, &f.DownlinkCongestionInformation}, {0x20, &f.UplinkCongestionInformation}} {
		if f.Flags&item.bit != 0 {
			if l < offset+2 {
				return io.ErrUnexpectedEOF
			}
			*item.dst = binary.BigEndian.Uint16(b[offset : offset+2])
			offset += 2
		}
	}
	for _, item := range []struct {
		bit uint8
		dst *uint64
	}{{0x40, &f.AverageDownlinkDataRate}, {0x80, &f.AverageUplinkDataRate}} {
		if f.Flags&item.bit != 0 {
			if l < offset+5 {
				return io.ErrUnexpectedEOF
			}
			for _, v := range b[offset : offset+5] {
				*item.dst = *item.dst<<8 | uint64(v)
			}
			offset += 5
		}
	}
	return f.validate()
}

// Marshal returns the serialized bytes of QoSMonitoringMeasurementFields.
func (f *QoSMonitoringMeasurementFields) Marshal() ([]byte, error) {
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
func (f *QoSMonitoringMeasurementFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if err := f.validate(); err != nil {
		return err
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

	if f.HasDL() {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		binary.BigEndian.PutUint32(b[offset:offset+4], f.DownlinkPacketDelay)
		offset += 4
	}

	if f.HasUL() {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		binary.BigEndian.PutUint32(b[offset:offset+4], f.UplinkPacketDelay)
		offset += 4
	}

	if f.HasRP() {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		binary.BigEndian.PutUint32(b[offset:offset+4], f.RoundTripPacketDelay)
		offset += 4
	}

	for _, item := range []struct {
		bit uint8
		v   uint16
	}{{0x10, f.DownlinkCongestionInformation}, {0x20, f.UplinkCongestionInformation}} {
		if f.Flags&item.bit != 0 {
			binary.BigEndian.PutUint16(b[offset:offset+2], item.v)
			offset += 2
		}
	}
	for _, item := range []struct {
		bit uint8
		v   uint64
	}{{0x40, f.AverageDownlinkDataRate}, {0x80, f.AverageUplinkDataRate}} {
		if f.Flags&item.bit != 0 {
			for k := 4; k >= 0; k-- {
				b[offset+k] = uint8(item.v)
				item.v >>= 8
			}
			offset += 5
		}
	}
	return nil
}

// MarshalLen returns field length in integer.
func (f *QoSMonitoringMeasurementFields) MarshalLen() int {
	if f == nil {
		return 0
	}

	l := 1
	if f.HasDL() {
		l += 4
	}
	if f.HasUL() {
		l += 4
	}
	if f.HasRP() {
		l += 4
	}
	if f.Flags&0x10 != 0 {
		l += 2
	}
	if f.Flags&0x20 != 0 {
		l += 2
	}
	if f.Flags&0x40 != 0 {
		l += 5
	}
	if f.Flags&0x80 != 0 {
		l += 5
	}
	return l
}

// NewQoSMonitoringMeasurementFromFields encodes all current fields of QoSMonitoringMeasurement.
func NewQoSMonitoringMeasurementFromFields(fields *QoSMonitoringMeasurementFields) *IE {
	b, err := fields.Marshal()
	if err != nil {
		return newInvalidIE(QoSMonitoringMeasurement, ErrInvalidValue)
	}
	return New(QoSMonitoringMeasurement, b)
}

// TS 29.244 V20.1.0 clause 8.2.171: congestion is 0..10000 and
// data rates occupy 40 bits. Failure excludes all packet delay fields.
func (f *QoSMonitoringMeasurementFields) validate() error {
	if f.Flags&1 == 0 && f.DownlinkPacketDelay != 0 || f.Flags&2 == 0 && f.UplinkPacketDelay != 0 || f.Flags&4 == 0 && f.RoundTripPacketDelay != 0 || f.Flags&16 == 0 && f.DownlinkCongestionInformation != 0 || f.Flags&32 == 0 && f.UplinkCongestionInformation != 0 || f.Flags&64 == 0 && f.AverageDownlinkDataRate != 0 || f.Flags&128 == 0 && f.AverageUplinkDataRate != 0 {
		return ErrInvalidValue
	}

	if f.Flags == 0 || f.Flags&8 != 0 && f.Flags&7 != 0 {
		return ErrInvalidValue
	}
	if f.DownlinkCongestionInformation > 10000 || f.UplinkCongestionInformation > 10000 || f.AverageDownlinkDataRate>>40 != 0 || f.AverageUplinkDataRate>>40 != 0 {
		return ErrInvalidValue
	}
	return nil
}
