// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewQERControlIndications creates a new QERControlIndications IE.
func NewQERControlIndications(rcsr bool) *IE {
	// TS 29.244 V20.1.0 clause 8.2.174: only RCSR remains allocated.
	var flags uint8
	if rcsr {
		flags = 1
	}
	return newUint8ValIE(QERControlIndications, flags)
}

// QERControlIndications returns QERControlIndications in uint8 if the type of IE matches.
func (i *IE) QERControlIndications() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case QERControlIndications:
		// TS 29.244 V20.1.0 clause 8.2.174: ignore received spare/reserved bits.
		v, err := i.ValueAsUint8()
		return v & 1, err
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}

}

// HasRCSR reports whether an IE has RCSR bit.
func (i *IE) HasRCSR() bool {
	if i == nil {
		return false
	}

	v, err := i.QERControlIndications()
	if err != nil {
		return false
	}

	return has1stBit(v)
}
