// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
)

// NewUEIPAddressPoolIdentity creates a new UEIPAddressPoolIdentity IE.
// The octet string is carried unchanged; no character set is imposed or checked.
func NewUEIPAddressPoolIdentity(id string) *IE {
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

// UEIPAddressPoolIdentity returns the pool identifier octet string unchanged.
// No character set is imposed or checked. TS 29.244 V20.1.0 clause 8.2.128.
func (i *IE) UEIPAddressPoolIdentity() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}
	if i.Type != UEIPAddressPoolIdentity {
		return "", &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 2 {
		return "", io.ErrUnexpectedEOF
	}
	n := int(binary.BigEndian.Uint16(i.Payload[:2]))
	b, err := i.payloadAtLeast(2 + n)
	if err != nil {
		return "", err
	}
	return string(b[2 : 2+n]), nil
}

// UEIPAddressPoolIdentityString returns the explicit string view of the pool identifier.
// The octet string is unchanged; no character set is imposed or checked.
// TS 29.244 V20.1.0 clause 8.2.128.
func (i *IE) UEIPAddressPoolIdentityString() (string, error) {
	return i.UEIPAddressPoolIdentity()
}
