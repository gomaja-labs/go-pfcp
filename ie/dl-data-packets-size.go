// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewDLDataPacketsSize creates a new DLDataPacketsSize IE.
func NewDLDataPacketsSize(size uint16) *IE {
	return newUint16ValIE(DLDataPacketsSize, size)
}

// DLDataPacketsSize returns DLDataPacketsSize in uint16 if the type of IE matches.
func (i *IE) DLDataPacketsSize() (uint16, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case DownlinkDataServiceInformation:
		f, e := i.DownlinkDataServiceInformation()
		if e != nil {
			return 0, e
		}
		if f.Flags&4 == 0 {
			return 0, ErrIENotFound
		}
		return f.DLDataPacketsSize, nil
	case DLDataPacketsSize:
		return i.ValueAsUint16()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
