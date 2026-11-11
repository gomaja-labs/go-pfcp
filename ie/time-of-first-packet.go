// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"time"
)

// NewTimeOfFirstPacket creates a new TimeOfFirstPacket IE.
// It returns nil outside the supported Unix-seconds window [0, 2^32-1];
// fractional seconds are discarded. See RFC 5905 section 6 and the package era policy.
func NewTimeOfFirstPacket(ts time.Time) *IE {
	return new3GPPTimestampIE(TimeOfFirstPacket, ts)
}

// TimeOfFirstPacket returns TimeOfFirstPacket in time.Time if the type of IE matches.
func (i *IE) TimeOfFirstPacket() (time.Time, error) {
	if i == nil {
		return time.Time{}, ErrIENotFound
	}

	switch i.Type {
	case TimeOfFirstPacket:
		return i.valueAs3GPPTimestamp()
	case UsageReportWithinSessionModificationResponse,
		UsageReportWithinSessionDeletionResponse,
		UsageReportWithinSessionReportRequest:
		ies, err := i.UsageReport()
		if err != nil {
			return time.Time{}, err
		}
		for _, x := range ies {
			if x.Type == TimeOfFirstPacket {
				return x.TimeOfFirstPacket()
			}
		}
		return time.Time{}, ErrIENotFound
	default:
		return time.Time{}, &InvalidTypeError{Type: i.Type}
	}
}
