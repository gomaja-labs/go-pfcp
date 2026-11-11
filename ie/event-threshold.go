// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewEventThreshold creates a new EventThreshold IE.
func NewEventThreshold(quota uint32) *IE {
	return newUint32ValIE(EventThreshold, quota)
}

// EventThreshold returns EventThreshold in uint32 if the type of IE matches.
func (i *IE) EventThreshold() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case EventThreshold:
		return i.ValueAsUint32()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
