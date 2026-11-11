// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"time"
)

// NewEndTime creates a new EndTime IE.
// It returns nil outside the supported Unix-seconds window [0, 2^32-1];
// fractional seconds are discarded. See RFC 5905 section 6 and the package era policy.
func NewEndTime(ts time.Time) *IE {
	return new3GPPTimestampIE(EndTime, ts)
}

// EndTime returns EndTime in time.Time if the type of IE matches.
func (i *IE) EndTime() (time.Time, error) {
	if i == nil {
		return time.Time{}, ErrIENotFound
	}

	switch i.Type {
	case EndTime:
		return i.valueAs3GPPTimestamp()
	case UsageReportWithinSessionModificationResponse,
		UsageReportWithinSessionDeletionResponse,
		UsageReportWithinSessionReportRequest:
		ies, err := i.UsageReport()
		if err != nil {
			return time.Time{}, err
		}
		for _, x := range ies {
			if x.Type == EndTime {
				return x.EndTime()
			}
		}
		return time.Time{}, ErrIENotFound
	default:
		return time.Time{}, &InvalidTypeError{Type: i.Type}
	}

}
