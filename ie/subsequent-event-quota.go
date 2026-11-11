// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewSubsequentEventQuota creates a new SubsequentEventQuota IE.
func NewSubsequentEventQuota(quota uint32) *IE {
	return newUint32ValIE(SubsequentEventQuota, quota)
}

// SubsequentEventQuota returns SubsequentEventQuota in uint32 if the type of IE matches.
func (i *IE) SubsequentEventQuota() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case SubsequentEventQuota:
		return i.ValueAsUint32()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
