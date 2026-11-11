// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"
	"math"
)

// NewDLBufferingSuggestedPacketCount creates a new DLBufferingSuggestedPacketCount IE.
func NewDLBufferingSuggestedPacketCount(count uint16) *IE {
	if count <= math.MaxUint8 {
		return newUint8ValIE(DLBufferingSuggestedPacketCount, uint8(count))
	}
	return newUint16ValIE(DLBufferingSuggestedPacketCount, count)
}

// DLBufferingSuggestedPacketCount returns DLBufferingSuggestedPacketCount in uint16 if the type of IE matches.
func (i *IE) DLBufferingSuggestedPacketCount() (uint16, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case DLBufferingSuggestedPacketCount:
		if len(i.Payload) == 1 {
			return uint16(i.Payload[0]), nil
		}

		if len(i.Payload) >= 2 {
			return i.ValueAsUint16()
		}

		return 0, io.ErrUnexpectedEOF
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
