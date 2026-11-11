// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import "io"

// NewQFI creates a new QFI IE.
func NewQFI(qfi uint8) *IE {
	return newUint8ValIE(QFI, qfi)
}

// QFI returns QFI in uint8 if the type of IE matches.
func (i *IE) QFI() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case QFI:
		return i.ValueAsUint8()
	case DownlinkDataServiceInformation:
		// TS 29.244 V20.1.0 clause 8.2.27: PPI and QFI are independently optional.
		if len(i.Payload) < 1 {
			return 0, io.ErrUnexpectedEOF
		}
		if !i.HasQFI() {
			return 0, ErrIENotFound
		}
		offset := 1
		if i.HasPPI() {
			offset++
		}
		if len(i.Payload) <= offset {
			return 0, io.ErrUnexpectedEOF
		}
		return i.Payload[offset], nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
