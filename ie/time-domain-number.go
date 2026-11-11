// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// TS 29.244 V20.1.0 clause 8.2.146.

// NewTimeDomainNumber creates a new TimeDomainNumber IE.
func NewTimeDomainNumber(num uint8) *IE {
	return newUint8ValIE(TimeDomainNumber, num)
}

// TimeDomainNumber returns TimeDomainNumber in uint8 if the type of IE matches.
func (i *IE) TimeDomainNumber() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case TimeDomainNumber:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
