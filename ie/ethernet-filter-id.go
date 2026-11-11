// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewEthernetFilterID creates a new EthernetFilterID IE.
func NewEthernetFilterID(id uint32) *IE {
	return newUint32ValIE(EthernetFilterID, id)
}

// EthernetFilterID returns EthernetFilterID in uint32 if the type of IE matches.
func (i *IE) EthernetFilterID() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case EthernetFilterID:
		return i.ValueAsUint32()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}

}
