// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
	"time"
)

// NewTimeOffsetMeasurement creates a TimeOffsetMeasurement IE from any signed nanosecond duration.
// TS 29.244 V20.1.0 clause 8.2.149: signed64 nanoseconds, including negative offsets.
func NewTimeOffsetMeasurement(measurement time.Duration) *IE {
	return newUint64ValIE(TimeOffsetMeasurement, uint64(measurement.Nanoseconds()))
}

// TimeOffsetMeasurement returns TimeOffsetMeasurement in time.Duration if the type of IE matches.
func (i *IE) TimeOffsetMeasurement() (time.Duration, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	if len(i.Payload) < 8 {
		return 0, io.ErrUnexpectedEOF
	}

	switch i.Type {
	case TimeOffsetMeasurement:
		return time.Duration(binary.BigEndian.Uint64(i.Payload[0:8])), nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
