// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
	"time"
)

// NewAveragePacketDelay creates an Average Packet Delay IE. The duration must
// be a whole number of milliseconds from 0 to 4294967295; another duration
// gives an IE whose encoding reports ErrInvalidValue. TS 29.244 V20.1.0 clause
// 8.2.162.
func NewAveragePacketDelay(delay time.Duration) *IE {
	if !durationFits(delay, time.Millisecond, 4294967295) {
		return newInvalidIE(AveragePacketDelay, ErrInvalidValue)
	}
	return newUint32ValIE(AveragePacketDelay, uint32(delay/time.Millisecond))
}

// AveragePacketDelay returns AveragePacketDelay in time.Duration if the type of IE matches.
func (i *IE) AveragePacketDelay() (time.Duration, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case AveragePacketDelay:
		if len(i.Payload) < 4 {
			return 0, io.ErrUnexpectedEOF
		}

		return time.Duration(binary.BigEndian.Uint32(i.Payload[0:4])) * time.Millisecond, nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
