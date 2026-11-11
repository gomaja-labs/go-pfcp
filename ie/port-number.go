// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// TS 29.244 V20.1.0 clause 8.2.141.

// NewPortNumber creates a new PortNumber IE.
func NewPortNumber(port uint32) *IE {
	return newUint32ValIE(PortNumber, port)
}

// PortNumber returns PortNumber in uint32 if the type of IE matches.
func (i *IE) PortNumber() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case PortNumber:
		return i.ValueAsUint32()
	case CreatedBridgeRouterInfo:
		ies, err := i.CreatedBridgeRouterInfo()
		if err != nil {
			return 0, err
		}
		for _, x := range ies {
			if x.Type == PortNumber {
				return x.PortNumber()
			}
		}
		return 0, ErrIENotFound
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}

}
