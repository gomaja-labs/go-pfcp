// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
)

// FiveGSUserPlaneNodeIDFields contains the optional Unsigned64 node identifier.
// TS 29.244 V20.1.0 clause 8.2.143.
type FiveGSUserPlaneNodeIDFields struct {
	BID   bool
	Value uint64
}

// NewFiveGSUserPlaneNodeID creates a 5GS User Plane Node ID with explicit presence.
func NewFiveGSUserPlaneNodeID(bid bool, id uint64) *IE {
	if !bid && id != 0 {
		return newInvalidIE(FiveGSUserPlaneNodeID, ErrInvalidValue)
	}
	b := []byte{0}
	if bid {
		b = make([]byte, 9)
		b[0] = 1
		binary.BigEndian.PutUint64(b[1:], id)
	}
	return New(FiveGSUserPlaneNodeID, b)
}

// HasBID reports whether the Unsigned64 User Plane Node ID is present.
func (i *IE) HasBID() bool {
	return i != nil && i.Type == FiveGSUserPlaneNodeID && len(i.Payload) > 0 && i.Payload[0]&1 != 0
}

// FiveGSUserPlaneNodeID decodes the kind-independent eight-octet identifier.
func (i *IE) FiveGSUserPlaneNodeID() (*FiveGSUserPlaneNodeIDFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}
	if i.Type != FiveGSUserPlaneNodeID {
		return nil, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 1 {
		return nil, io.ErrUnexpectedEOF
	}
	f := &FiveGSUserPlaneNodeIDFields{BID: i.HasBID()}
	if f.BID {
		if len(i.Payload) < 9 {
			return nil, io.ErrUnexpectedEOF
		}
		f.Value = binary.BigEndian.Uint64(i.Payload[1:9])
	}
	return f, nil
}
