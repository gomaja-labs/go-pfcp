// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
	"time"
)

// NewQuotaHoldingTime creates a Quota Holding Time IE. The duration must be a
// whole number of seconds from 0 to 4294967295; another duration gives an IE
// whose encoding reports ErrInvalidValue. TS 29.244 V20.1.0 clause 8.2.48.
func NewQuotaHoldingTime(duration time.Duration) *IE {
	if !durationFits(duration, time.Second, 4294967295) {
		return newInvalidIE(QuotaHoldingTime, ErrInvalidValue)
	}
	return newUint32ValIE(QuotaHoldingTime, uint32(duration/time.Second))
}

// QuotaHoldingTime returns QuotaHoldingTime in time.Duration if the type of IE matches.
func (i *IE) QuotaHoldingTime() (time.Duration, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case QuotaHoldingTime:
		if len(i.Payload) < 4 {
			return 0, io.ErrUnexpectedEOF
		}

		return time.Duration(binary.BigEndian.Uint32(i.Payload[0:4])) * time.Second, nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
