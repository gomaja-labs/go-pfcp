// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// Interface definitions.
const (
	// DstInterfaceAccess selects Access.
	DstInterfaceAccess uint8 = 0
	// DstInterfaceCore selects Core.
	DstInterfaceCore uint8 = 1
	// DstInterfaceSGiLANN6LAN selects the SGi-LAN/N6-LAN interface.
	DstInterfaceSGiLANN6LAN uint8 = 2
	// DstInterfaceCPFunction selects CP Function.
	DstInterfaceCPFunction uint8 = 3
	// DstInterfaceLIFunction selects LI Function.
	DstInterfaceLIFunction uint8 = 4
	// DstInterface5GVNInternal selects the internal 5G VN interface.
	DstInterface5GVNInternal uint8 = 5
)

// NewDestinationInterface creates a new DestinationInterface IE.
func NewDestinationInterface(intf uint8) *IE {
	return newUint8ValIE(DestinationInterface, intf)
}

// DestinationInterface returns DestinationInterface in uint8 if the type of IE matches.
func (i *IE) DestinationInterface() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case DestinationInterface:
		// TS 29.244 V20.1.0 clauses 8.1.2 and 8.2.24: ignore spare bits.
		return i.valueAsUint8Masked(0x0f)
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
