// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// TS 29.244 V20.1.0 clause 8.2.140.

// NewCreateBridgeRouterInfo creates a new CreateBridgeRouterInfo IE.
func NewCreateBridgeRouterInfo(bii uint8) *IE {
	return newUint8ValIE(CreateBridgeRouterInfo, bii&0x03)
}

// CreateBridgeRouterInfo returns CreateBridgeRouterInfo in uint8 if the type of IE matches.
func (i *IE) CreateBridgeRouterInfo() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case CreateBridgeRouterInfo:
		return i.valueAsUint8Masked(0x03)
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// HasBII reports whether an IE has BII bit.
func (i *IE) HasBII() bool {
	if i == nil {
		return false
	}

	v, err := i.CreateBridgeRouterInfo()
	if err != nil {
		return false
	}

	return has1stBit(v)
}

// HasRII reports whether the RII flag is set.
// TS 29.244 V20.1.0 clause 8.2.140, octet 5 bit 2.
func (i *IE) HasRII() bool { return payloadFlag(i, CreateBridgeRouterInfo, 0, 0x02) }
