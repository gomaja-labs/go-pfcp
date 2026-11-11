// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"time"
)

// NewDeactivationTime creates a new DeactivationTime IE.
// It returns nil outside the supported Unix-seconds window [0, 2^32-1];
// fractional seconds are discarded. See RFC 5905 section 6 and the package era policy.
func NewDeactivationTime(ts time.Time) *IE {
	return new3GPPTimestampIE(DeactivationTime, ts)
}

// DeactivationTime returns DeactivationTime in time.Time if the type of IE matches.
func (i *IE) DeactivationTime() (time.Time, error) {
	if i == nil {
		return time.Time{}, ErrIENotFound
	}

	switch i.Type {
	case DeactivationTime:
		return i.valueAs3GPPTimestamp()
	case CreatePDR:
		ies, err := i.CreatePDR()
		if err != nil {
			return time.Time{}, err
		}
		for _, x := range ies {
			if x.Type == DeactivationTime {
				return x.DeactivationTime()
			}
		}
		return time.Time{}, ErrIENotFound
	case UpdatePDR:
		ies, err := i.UpdatePDR()
		if err != nil {
			return time.Time{}, err
		}
		for _, x := range ies {
			if x.Type == DeactivationTime {
				return x.DeactivationTime()
			}
		}
		return time.Time{}, ErrIENotFound
	default:
		return time.Time{}, &InvalidTypeError{Type: i.Type}
	}
}
