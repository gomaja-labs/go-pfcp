// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"time"
)

// NewEventTimeStamp creates a new EventTimeStamp IE.
// It returns nil outside the supported Unix-seconds window [0, 2^32-1];
// fractional seconds are discarded. See RFC 5905 section 6 and the package era policy.
func NewEventTimeStamp(ts time.Time) *IE {
	return new3GPPTimestampIE(EventTimeStamp, ts)
}

// EventTimeStamp returns EventTimeStamp in time.Time if the type of IE matches.
func (i *IE) EventTimeStamp() (time.Time, error) {
	if i == nil {
		return time.Time{}, ErrIENotFound
	}

	switch i.Type {
	case EventTimeStamp:
		return i.valueAs3GPPTimestamp()
	case UsageReportWithinSessionReportRequest:
		ies, err := i.UsageReport()
		if err != nil {
			return time.Time{}, err
		}
		for _, x := range ies {
			if x.Type == EventTimeStamp {
				return x.EventTimeStamp()
			}
		}
		return time.Time{}, ErrIENotFound
	case GTPUPathQoSReport:
		ies, err := i.GTPUPathQoSReport()
		if err != nil {
			return time.Time{}, err
		}
		for _, x := range ies {
			if x.Type == EventTimeStamp {
				return x.EventTimeStamp()
			}
		}
		return time.Time{}, ErrIENotFound
	case QoSMonitoringReport:
		ies, err := i.QoSMonitoringReport()
		if err != nil {
			return time.Time{}, err
		}
		for _, x := range ies {
			if x.Type == EventTimeStamp {
				return x.EventTimeStamp()
			}
		}
		return time.Time{}, ErrIENotFound
	case ClockDriftReport:
		ies, err := i.ClockDriftReport()
		if err != nil {
			return time.Time{}, err
		}

		for _, x := range ies {
			if x.Type == EventTimeStamp {
				return x.EventTimeStamp()
			}
		}
		return time.Time{}, ErrIENotFound
	default:
		return time.Time{}, &InvalidTypeError{Type: i.Type}
	}

}
