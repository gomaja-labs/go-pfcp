// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewSubsequentEventThreshold creates a new SubsequentEventThreshold IE.
func NewSubsequentEventThreshold(quota uint32) *IE {
	return newUint32ValIE(SubsequentEventThreshold, quota)
}

// SubsequentEventThreshold returns SubsequentEventThreshold in uint32 if the type of IE matches.
func (i *IE) SubsequentEventThreshold() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case SubsequentEventThreshold:
		return i.ValueAsUint32()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
