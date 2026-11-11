// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
	"time"
)

// NewValidityTimer creates a new ValidityTimer IE.
//
// the period should be within the range of uint16, otherwise it overflows.
func NewValidityTimer(period time.Duration) *IE {
	return newUint16ValIE(ValidityTimer, uint16(period.Seconds()))
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
