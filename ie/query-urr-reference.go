// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewQueryURRReference creates a new QueryURRReference IE.
func NewQueryURRReference(id uint32) *IE {
	return newUint32ValIE(QueryURRReference, id)
}

// QueryURRReference returns QueryURRReference in uint32 if the type of IE matches.
func (i *IE) QueryURRReference() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case QueryURRReference:
		return i.ValueAsUint32()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
