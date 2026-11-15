// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
	"time"
)

// NewMinimumWaitTime creates a Minimum Wait Time IE. The duration must be a
// whole number of seconds from 0 to 4294967295; another duration gives an IE
// whose encoding reports ErrInvalidValue. TS 29.244 V20.1.0 clause 8.2.170.
func NewMinimumWaitTime(period time.Duration) *IE {
	if !durationFits(period, time.Second, 4294967295) {
		return newInvalidIE(MinimumWaitTime, ErrInvalidValue)
	}
	return newUint32ValIE(MinimumWaitTime, uint32(period/time.Second))
}

// MinimumWaitTime returns MinimumWaitTime in time.Duration if the type of IE matches.
func (i *IE) MinimumWaitTime() (time.Duration, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case MinimumWaitTime:
		if len(i.Payload) < 4 {
			return 0, io.ErrUnexpectedEOF
		}

		return time.Duration(binary.BigEndian.Uint32(i.Payload[0:4])) * time.Second, nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
