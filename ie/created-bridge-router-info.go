// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// TS 29.244 V20.1.0 Table 7.5.3.6-1.

// NewCreatedBridgeRouterInfo creates a new CreatedBridgeRouterInfo IE.
func NewCreatedBridgeRouterInfo(ies ...*IE) *IE {
	return newGroupedIE(CreatedBridgeRouterInfo, 0, ies...)
}

// CreatedBridgeRouterInfo returns the IEs above CreatedBridgeRouterInfo if the type of IE matches.
func (i *IE) CreatedBridgeRouterInfo() ([]*IE, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	if i.Type != CreatedBridgeRouterInfo {
		return nil, &InvalidTypeError{Type: i.Type}
	}

	return i.ValueAsGrouped()
}
