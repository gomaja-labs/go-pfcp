// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
)

// NewUEIPAddressPoolIdentity creates a new UEIPAddressPoolIdentity IE.
func NewUEIPAddressPoolIdentity(id []byte) *IE {
	// TS 29.244 V20.1.0 clause 8.2.128: an opaque length-delimited pool identity.
	l := len(id)
	if l > 65533 {
		return newInvalidIE(UEIPAddressPoolIdentity, ErrInvalidValue)
	}
	i := New(UEIPAddressPoolIdentity, make([]byte, 2+l))

	binary.BigEndian.PutUint16(i.Payload[0:2], uint16(l))
	copy(i.Payload[2:], id)

	return i
}

// UEIPAddressPoolIdentity returns UEIPAddressPoolIdentity in []byte if the type of IE matches.
func (i *IE) UEIPAddressPoolIdentity() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case UEIPAddressPoolIdentity:
		// TS 29.244 V20.1.0 clause 8.2.128: two-octet length and identifier.
		if len(i.Payload) < 2 {
			return nil, io.ErrUnexpectedEOF
		}
		return i.payloadAtLeast(2 + int(binary.BigEndian.Uint16(i.Payload[:2])))
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// UEIPAddressPoolIdentityString returns an explicit string view of the opaque pool identity.
func (i *IE) UEIPAddressPoolIdentityString() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}

	v, err := i.UEIPAddressPoolIdentity()
	if err != nil {
		return "", err
	}

	if len(v) < 2 {
		return "", io.ErrUnexpectedEOF
	}

	idlen := int(binary.BigEndian.Uint16(v[0:2]))
	if len(v) < idlen+2 {
		return "", io.ErrUnexpectedEOF
	}

	return string(v[2 : idlen+2]), nil
}
