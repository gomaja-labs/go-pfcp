// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"time"
)

// BaseTimeIntervalType definitions.
const (
	// BTITCTP selects CTP.
	BTITCTP uint8 = 0
	// BTITDTP selects DTP.
	BTITDTP uint8 = 1
)

// NewTimeQuotaMechanism creates a new TimeQuotaMechanism IE.
func NewTimeQuotaMechanism(btit uint8, bti time.Duration) *IE {
	b := make([]byte, 5)
	b[0] = btit & 0x03
	binary.BigEndian.PutUint32(b[1:5], uint32(bti.Seconds()))

	return New(TimeQuotaMechanism, b)
}

// TimeQuotaMechanism returns TimeQuotaMechanism in []byte if the type of IE matches.
func (i *IE) TimeQuotaMechanism() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case TimeQuotaMechanism:
		// TS 29.244 V20.1.0 clause 8.2.81: type and four-octet interval.
		return i.payloadAtLeast(5)
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}
