// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// TS 29.244 V20.1.0 clause 8.2.143.

import (
	"io"
	"net"
)

// NewFiveGSUserPlaneNodeID creates a new FiveGSUserPlaneNodeID IE.
// TS 29.244 V20.1.0 clause 8.2.143 instead defines type 198 as the
// 5GS User Plane Node ID with an eight-octet value. This constructor's
// MAC-based API does not represent that current definition.
func NewFiveGSUserPlaneNodeID(mac net.HardwareAddr) *IE {
	if mac == nil {
		return New(FiveGSUserPlaneNodeID, []byte{0x00})
	}
	l := len(mac) + 1
	b := make([]byte, l)
	b[0] = 0x01
	copy(b[1:l], mac)
	return New(FiveGSUserPlaneNodeID, b)
}

// HasMAC reports whether an IE has MAC bit.
func (i *IE) HasMAC() bool {
	if i == nil {
		return false
	}

	switch i.Type {
	case FiveGSUserPlaneNodeID:
		if len(i.Payload) == 0 {
			return false
		}
		return has1stBit(i.Payload[0])
	default:
		return false
	}
}

// FiveGSUserPlaneNodeID returns FiveGSUserPlaneNodeID in net.HardwareAddr if the type of IE matches.
func (i *IE) FiveGSUserPlaneNodeID() (net.HardwareAddr, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case FiveGSUserPlaneNodeID:
		if len(i.Payload) < 1 {
			return nil, io.ErrUnexpectedEOF
		}

		if !has1stBit(i.Payload[0]) {
			return nil, nil
		}
		if (len(i.Payload) != 7) && (len(i.Payload) < 9) {
			return nil, io.ErrUnexpectedEOF
		}
		if len(i.Payload) == 7 {
			return net.HardwareAddr(i.Payload[1:7]), nil
		}
		return net.HardwareAddr(i.Payload[1:9]), nil
	case CreatedBridgeRouterInfo:
		ies, err := i.CreatedBridgeRouterInfo()
		if err != nil {
			return nil, err
		}
		for _, x := range ies {
			if x.Type == FiveGSUserPlaneNodeID {
				return x.FiveGSUserPlaneNodeID()
			}
		}
		return nil, ErrIENotFound
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}
