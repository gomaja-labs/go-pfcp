// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"fmt"
	"io"
)

// NewNFInstanceID creates a new NFInstanceID IE; id must contain 16 octets.
func NewNFInstanceID(id []byte) *IE {
	// TS 29.244 V20.1.0 clause 8.2.175: exactly sixteen octets.
	if len(id) != 16 {
		return newInvalidIE(NFInstanceID, fmt.Errorf("%w: NF Instance ID requires 16 octets", ErrInvalidValue))
	}
	return New(NFInstanceID, id[:16:16])
}

// NFInstanceID returns NFInstanceID in []byte if the type of IE matches.
func (i *IE) NFInstanceID() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case NFInstanceID:
		if err := i.validateFixedValue(); err != nil {
			return nil, err
		}
		if len(i.Payload) < 16 {
			return nil, io.ErrUnexpectedEOF
		}

		return i.Payload[:16], nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}
