// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// Interface definitions.
const (
	SrcInterfaceAccess       uint8 = 0
	SrcInterfaceCore         uint8 = 1
	SrcInterfaceSGiLANN6LAN  uint8 = 2
	SrcInterfaceCPFunction   uint8 = 3
	SrcInterface5GVNInternal uint8 = 4
)

// NewSourceInterface creates a new SourceInterface IE.
func NewSourceInterface(intf uint8) *IE {
	return newUint8ValIE(SourceInterface, intf)
}

// SourceInterface returns SourceInterface in uint8 if the type of IE matches.
func (i *IE) SourceInterface() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case SourceInterface:
		// TS 29.244 V20.1.0 clauses 8.1.2 and 8.2.2: ignore spare bits.
		return i.valueAsUint8Masked(0x0f)
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
