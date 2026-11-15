// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
	"time"
)

// NewMeasurementPeriod creates a Measurement Period IE. The duration must be a
// whole number of seconds from 0 to 4294967295; another duration gives an IE
// whose encoding reports ErrInvalidValue. TS 29.244 V20.1.0 clause 8.2.42.
func NewMeasurementPeriod(period time.Duration) *IE {
	if !durationFits(period, time.Second, 4294967295) {
		return newInvalidIE(MeasurementPeriod, ErrInvalidValue)
	}
	return newUint32ValIE(MeasurementPeriod, uint32(period/time.Second))
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
