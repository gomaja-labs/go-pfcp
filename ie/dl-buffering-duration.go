// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"
	"time"
)

// NewDLBufferingDuration creates a DL Buffering Duration IE. Zero stops the
// timer and InfiniteTimer makes it infinite. A finite duration must be exactly
// 1 to 31 units of 2 seconds, 1 minute, 10 minutes, 1 hour or 10 hours; another
// duration gives an IE whose encoding reports ErrInvalidValue. TS 29.244
// V20.1.0 clause 8.2.29.
func NewDLBufferingDuration(duration time.Duration) *IE {
	return newTimerIE(DLBufferingDuration, duration)
}

// DLBufferingDuration returns DLBufferingDuration in time.Duration if the type of IE matches.
func (i *IE) DLBufferingDuration() (time.Duration, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case DLBufferingDuration:
		if len(i.Payload) < 1 {
			return 0, io.ErrUnexpectedEOF
		}

		return decodeTimerOctet(i.Payload[0]), nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
