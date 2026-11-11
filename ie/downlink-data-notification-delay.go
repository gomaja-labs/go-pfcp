// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"
	"time"
)

// NewDownlinkDataNotificationDelay creates a new DownlinkDataNotificationDelay IE.
func NewDownlinkDataNotificationDelay(delay time.Duration) *IE {
	// TS29.244 8.2.28 Downlink Data Notification Delay
	// Delay Value in integer multiples of 50 millisecs, or zero
	return newUint8ValIE(DownlinkDataNotificationDelay, uint8(delay/50000000))
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
