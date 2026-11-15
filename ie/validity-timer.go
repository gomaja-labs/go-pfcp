// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
	"time"
)

// NewValidityTimer creates a Validity Timer IE. The duration must be a whole
// number of seconds from 0 to 65535; another duration gives an IE whose
// encoding reports ErrInvalidValue. TS 29.244 V20.1.0 clause 8.2.184.
func NewValidityTimer(period time.Duration) *IE {
	if !durationFits(period, time.Second, 65535) {
		return newInvalidIE(ValidityTimer, ErrInvalidValue)
	}
	return newUint16ValIE(ValidityTimer, uint16(period/time.Second))
}

// ValidityTimer returns ValidityTimer in time.Duration if the type of IE matches.
func (i *IE) ValidityTimer() (time.Duration, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case ValidityTimer:
		if len(i.Payload) < 2 {
			return 0, io.ErrUnexpectedEOF
		}

		return time.Duration(binary.BigEndian.Uint16(i.Payload[0:2])) * time.Second, nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
