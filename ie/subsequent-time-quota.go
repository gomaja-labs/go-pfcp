// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
	"time"
)

// NewSubsequentTimeQuota creates a Subsequent Time Quota IE. The duration must
// be a whole number of seconds from 0 to 4294967295; another duration gives an
// IE whose encoding reports ErrInvalidValue. TS 29.244 V20.1.0 clause 8.2.87.
func NewSubsequentTimeQuota(period time.Duration) *IE {
	if !durationFits(period, time.Second, 4294967295) {
		return newInvalidIE(SubsequentTimeQuota, ErrInvalidValue)
	}
	return newUint32ValIE(SubsequentTimeQuota, uint32(period/time.Second))
}

// SubsequentTimeQuota returns SubsequentTimeQuota in time.Duration if the type of IE matches.
func (i *IE) SubsequentTimeQuota() (time.Duration, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case SubsequentTimeQuota:
		// TS 29.244 V20.1.0 clause 8.2.87: four-octet seconds.
		if len(i.Payload) < 4 {
			return 0, io.ErrUnexpectedEOF
		}
		return time.Duration(binary.BigEndian.Uint32(i.Payload[0:4])) * time.Second, nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
