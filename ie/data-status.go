// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewDataStatus creates a new DataStatus IE.
func NewDataStatus(flag uint8) *IE {
	return newUint8ValIE(DataStatus, flag)
}

// DataStatus returns DataStatus in uint8 if the type of IE matches.
func (i *IE) DataStatus() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case DataStatus:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
