// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"fmt"
	"io"
)

// NewForwardingPolicy creates a new ForwardingPolicy IE.
// An identifier exceeding 255 octets produces an IE that fails to marshal.
func NewForwardingPolicy(id []byte) *IE {
	// TS 29.244 V20.1.0 clause 8.2.23.
	l := len(id)
	if l > 255 {
		return newInvalidIE(ForwardingPolicy, fmt.Errorf("%w: policy identifier exceeds 255 octets", ErrInvalidValue))
	}
	i := New(ForwardingPolicy, make([]byte, 1+l))

	i.Payload[0] = uint8(l)
	copy(i.Payload[1:], id)

	return i
}

// ForwardingPolicy returns ForwardingPolicy in []byte if the type of IE matches.
func (i *IE) ForwardingPolicy() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case ForwardingPolicy:
		// TS 29.244 V20.1.0 clause 8.2.23: length and policy identifier.
		if len(i.Payload) < 1 {
			return nil, io.ErrUnexpectedEOF
		}
		return i.payloadAtLeast(1 + int(i.Payload[0]))
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// ForwardingPolicyIdentifier returns the opaque ForwardingPolicyIdentifier if the type of IE matches.
func (i *IE) ForwardingPolicyIdentifier() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	v, err := i.ForwardingPolicy()
	if err != nil {
		return nil, err
	}

	l := len(v)
	if l < 1 {
		return nil, io.ErrUnexpectedEOF
	}

	idlen := int(v[0])
	if l < idlen+1 {
		return nil, io.ErrUnexpectedEOF
	}

	return v[1 : idlen+1], nil
}
