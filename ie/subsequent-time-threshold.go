// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
	"time"
)

// NewSubsequentTimeThreshold creates a Subsequent Time Threshold IE. The
// duration must be a whole number of seconds from 0 to 4294967295; another
// duration gives an IE whose encoding reports ErrInvalidValue. TS 29.244
// V20.1.0 clause 8.2.17.
func NewSubsequentTimeThreshold(t time.Duration) *IE {
	if !durationFits(t, time.Second, 4294967295) {
		return newInvalidIE(SubsequentTimeThreshold, ErrInvalidValue)
	}
	return newUint32ValIE(SubsequentTimeThreshold, uint32(t/time.Second))
}

// SubsequentTimeThreshold returns SubsequentTimeThreshold in time.Duration if the type of IE matches.
func (i *IE) SubsequentTimeThreshold() (time.Duration, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case SubsequentTimeThreshold:
		if len(i.Payload) < 4 {
			return 0, io.ErrUnexpectedEOF
		}

		t := binary.BigEndian.Uint32(i.Payload[0:4])
		return time.Duration(t) * time.Second, nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
