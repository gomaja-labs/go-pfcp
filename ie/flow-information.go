// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
)

// FlowDirection definitions.
const (
	FlowDirectionUnspecified   uint8 = 0
	FlowDirectionDownlink      uint8 = 1
	FlowDirectionUplink        uint8 = 2
	FlowDirectionBidirectional uint8 = 3
)

// NewFlowInformation creates a new FlowInformation IE.
func NewFlowInformation(dir uint8, desc string) *IE {
	d := []byte(desc)
	l := len(d)

	i := New(FlowInformation, make([]byte, 3+l))
	i.Payload[0] = dir
	binary.BigEndian.PutUint16(i.Payload[1:3], uint16(l))
	copy(i.Payload[3:], d)

	return i
}

// FlowInformation returns FlowInformation in []byte if the type of IE matches.
func (i *IE) FlowInformation() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case FlowInformation:
		// TS 29.244 V20.1.0 clause 8.2.61: direction, length, description.
		if len(i.Payload) < 3 {
			return nil, io.ErrUnexpectedEOF
		}
		return i.payloadAtLeast(3 + int(binary.BigEndian.Uint16(i.Payload[1:3])))
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// FlowDirection returns FlowDirection in uint8 if the type of IE matches.
func (i *IE) FlowDirection() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case FlowInformation:
		if len(i.Payload) < 1 {
			return 0, io.ErrUnexpectedEOF
		}
		// TS 29.244 V20.1.0 Table 8.2.61-1: received spare values
		// 4..7 are interpreted as unspecified; spare high bits are ignored.
		direction := i.Payload[0] & 0x07
		if direction > FlowDirectionBidirectional {
			direction = FlowDirectionUnspecified
		}
		return direction, nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// FlowDescription returns FlowDescription in string if the type of IE matches.
func (i *IE) FlowDescription() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}

	switch i.Type {
	case FlowInformation:
		// TS 29.244 V20.1.0 clause 8.2.61: the length excludes the
		// direction octet and the two-octet length field.
		if len(i.Payload) < 3 {
			return "", io.ErrUnexpectedEOF
		}
		l := int(binary.BigEndian.Uint16(i.Payload[1:3]))
		if len(i.Payload) < 3+l {
			return "", io.ErrUnexpectedEOF
		}

		return string(i.Payload[3 : 3+l]), nil
	default:
		return "", &InvalidTypeError{Type: i.Type}
	}
}
