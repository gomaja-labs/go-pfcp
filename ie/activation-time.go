// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"time"
)

// NewActivationTime creates a new ActivationTime IE.
// It returns nil outside the supported Unix-seconds window [0, 2^32-1];
// fractional seconds are discarded. See RFC 5905 section 6 and the package era policy.
func NewActivationTime(ts time.Time) *IE {
	return new3GPPTimestampIE(ActivationTime, ts)
}

// ActivationTime returns ActivationTime in time.Time if the type of IE matches.
func (i *IE) ActivationTime() (time.Time, error) {
	if i == nil {
		return time.Time{}, ErrIENotFound
	}

	switch i.Type {
	case ActivationTime:
		return i.valueAs3GPPTimestamp()
	default:
		return time.Time{}, &InvalidTypeError{Type: i.Type}
	}
}
