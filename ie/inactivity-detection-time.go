// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewInactivityDetectionTime creates a new InactivityDetectionTime IE.
func NewInactivityDetectionTime(threshold uint32) *IE {
	return newUint32ValIE(InactivityDetectionTime, threshold)
}

// InactivityDetectionTime returns InactivityDetectionTime in uint32 if the type of IE matches.
func (i *IE) InactivityDetectionTime() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case InactivityDetectionTime:
		return i.ValueAsUint32()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
