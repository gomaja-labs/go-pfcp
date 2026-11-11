// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewPDI creates a new PDI IE.
func NewPDI(ies ...*IE) *IE {
	return newGroupedIE(PDI, 0, ies...)
}

// PDI returns the IEs above PDI if the type of IE matches.
func (i *IE) PDI() ([]*IE, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case PDI:
		return i.ValueAsGrouped()
	case CreatePDR:
		ies, err := i.CreatePDR()
		if err != nil {
			return nil, err
		}
		for _, x := range ies {
			if x.Type == PDI {
				return x.PDI()
			}
		}
		return nil, ErrIENotFound
	case UpdatePDR:
		ies, err := i.UpdatePDR()
		if err != nil {
			return nil, err
		}
		for _, x := range ies {
			if x.Type == PDI {
				return x.PDI()
			}
		}
		return nil, ErrIENotFound
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}
