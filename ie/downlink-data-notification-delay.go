// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"
	"time"
)

// NewDownlinkDataNotificationDelay creates a Downlink Data Notification Delay
// IE. The delay must be a whole multiple of 50 milliseconds from 0 to 12.75
// seconds, where zero clears a previously set delay condition; another duration
// gives an IE whose encoding reports ErrInvalidValue. TS 29.244 V20.1.0 clause
// 8.2.28.
func NewDownlinkDataNotificationDelay(delay time.Duration) *IE {
	if !durationFits(delay, 50*time.Millisecond, 255) {
		return newInvalidIE(DownlinkDataNotificationDelay, ErrInvalidValue)
	}
	return newUint8ValIE(DownlinkDataNotificationDelay, uint8(delay/(50*time.Millisecond)))
}

// DownlinkDataNotificationDelay returns DownlinkDataNotificationDelay in time.Duration if the type of IE matches.
func (i *IE) DownlinkDataNotificationDelay() (time.Duration, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case DownlinkDataNotificationDelay:
		if len(i.Payload) < 1 {
			return 0, io.ErrUnexpectedEOF
		}

		return time.Duration(int64(i.Payload[0]) * 50000000), nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
