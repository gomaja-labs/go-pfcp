// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
	"time"
)

// NewUserPlaneInactivityTimer creates a User Plane Inactivity Timer IE. The
// duration must be a whole number of seconds from 0 to 4294967295, where zero
// stops user plane inactivity detection and reporting; another duration gives
// an IE whose encoding reports ErrInvalidValue. TS 29.244 V20.1.0 clause
// 8.2.83.
func NewUserPlaneInactivityTimer(period time.Duration) *IE {
	if !durationFits(period, time.Second, 4294967295) {
		return newInvalidIE(UserPlaneInactivityTimer, ErrInvalidValue)
	}
	return newUint32ValIE(UserPlaneInactivityTimer, uint32(period/time.Second))
}

// UserPlaneInactivityTimer returns UserPlaneInactivityTimer in time.Duration if the type of IE matches.
func (i *IE) UserPlaneInactivityTimer() (time.Duration, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	if i.Type != UserPlaneInactivityTimer {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 4 {
		return 0, io.ErrUnexpectedEOF
	}

	return time.Duration(binary.BigEndian.Uint32(i.Payload[0:4])) * time.Second, nil
}
