// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewTrafficEndpointID creates a new TrafficEndpointID IE.
func NewTrafficEndpointID(id uint8) *IE {
	return newUint8ValIE(TrafficEndpointID, id)
}

// TrafficEndpointID returns TrafficEndpointID in uint8 if the type of IE matches.
func (i *IE) TrafficEndpointID() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case TrafficEndpointID:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
