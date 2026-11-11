// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewTransportLevelMarking creates a new TransportLevelMarking IE.
func NewTransportLevelMarking(tos uint16) *IE {
	return newUint16ValIE(TransportLevelMarking, tos)
}

// TransportLevelMarking returns TransportLevelMarking in uint16 if the type of IE matches.
func (i *IE) TransportLevelMarking() (uint16, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case TransportLevelMarking:
		return i.ValueAsUint16()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}

}
