// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"time"
)

// NewMonitoringTime creates a new MonitoringTime IE.
// It returns nil outside the supported Unix-seconds window [0, 2^32-1];
// fractional seconds are discarded. See RFC 5905 section 6 and the package era policy.
func NewMonitoringTime(ts time.Time) *IE {
	return new3GPPTimestampIE(MonitoringTime, ts)
}

// MonitoringTime returns MonitoringTime in time.Time if the type of IE matches.
func (i *IE) MonitoringTime() (time.Time, error) {
	if i == nil {
		return time.Time{}, ErrIENotFound
	}

	switch i.Type {
	case MonitoringTime:
		return i.valueAs3GPPTimestamp()
	default:
		return time.Time{}, &InvalidTypeError{Type: i.Type}
	}
}
