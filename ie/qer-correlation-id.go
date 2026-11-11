// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewQERCorrelationID creates a new QERCorrelationID IE.
func NewQERCorrelationID(id uint32) *IE {
	return newUint32ValIE(QERCorrelationID, id)
}

// QERCorrelationID returns QERCorrelationID in uint32 if the type of IE matches.
func (i *IE) QERCorrelationID() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case QERCorrelationID:
		return i.ValueAsUint32()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
