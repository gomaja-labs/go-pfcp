// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewMeasurementMethod creates a new MeasurementMethod IE.
func NewMeasurementMethod(event, volum, durat int) *IE {
	return newUint8ValIE(MeasurementMethod, uint8((event<<2)|(volum<<1)|(durat)))
}

// MeasurementMethod returns MeasurementMethod in uint8 if the type of IE matches.
func (i *IE) MeasurementMethod() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case MeasurementMethod:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// HasEVENT reports whether an IE has EVENT bit.
func (i *IE) HasEVENT() bool {
	if i == nil {
		return false
	}

	v, err := i.MeasurementMethod()
	if err != nil {
		return false
	}

	return has3rdBit(v)
}

// HasVOLUM reports whether an IE has VOLUM bit.
func (i *IE) HasVOLUM() bool {
	if i == nil {
		return false
	}

	v, err := i.MeasurementMethod()
	if err != nil {
		return false
	}

	return has2ndBit(v)
}

// HasDURAT reports whether an IE has DURAT bit.
func (i *IE) HasDURAT() bool {
	if i == nil {
		return false
	}

	v, err := i.MeasurementMethod()
	if err != nil {
		return false
	}

	return has1stBit(v)
}
