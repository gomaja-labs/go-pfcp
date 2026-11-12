// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"time"
)

// NewActivationTime creates a new ActivationTime IE.
// The wire has no era: decoding selects the Unix-seconds window [0, 2^32-1].
// A date outside that window returns an IE whose Marshal reports ErrInvalidValue.
// The IE carries whole seconds; the fractional part is discarded. RFC 5905 section 6; TS 29.244 V20.1.0 clause 8.2.65.
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
