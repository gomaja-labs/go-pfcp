// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
	"time"
)

// NewMeasurementPeriod creates a new MeasurementPeriod IE.
//
// the period should be within the range of uint32, otherwise it overflows.
func NewMeasurementPeriod(period time.Duration) *IE {
	return newUint32ValIE(MeasurementPeriod, uint32(period.Seconds()))
}

// MeasurementPeriod returns MeasurementPeriod in time.Duration if the type of IE matches.
func (i *IE) MeasurementPeriod() (time.Duration, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case MeasurementPeriod:
		if len(i.Payload) < 4 {
			return 0, io.ErrUnexpectedEOF
		}

		return time.Duration(binary.BigEndian.Uint32(i.Payload[0:4])) * time.Second, nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
