// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// Framed-Routing definitions.
//
// Ref: https://tools.ietf.org/html/rfc2865#section-5.10
const (
	// FramedRoutingNone selects None.
	FramedRoutingNone uint32 = 0
	// FramedRoutingSendRoutingPackets selects Send Routing Packets.
	FramedRoutingSendRoutingPackets uint32 = 1
	// FramedRoutingListenForRoutingPackets selects Listen For Routing Packets.
	FramedRoutingListenForRoutingPackets uint32 = 2
	// FramedRoutingSendAndListen selects Send And Listen.
	FramedRoutingSendAndListen uint32 = 3
)

// NewFramedRouting creates a new FramedRouting IE.
func NewFramedRouting(routing uint32) *IE {
	return newUint32ValIE(FramedRouting, routing)
}

// FramedRouting returns FramedRouting in uint32 if the type of IE matches.
func (i *IE) FramedRouting() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case FramedRouting:
		return i.ValueAsUint32()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
