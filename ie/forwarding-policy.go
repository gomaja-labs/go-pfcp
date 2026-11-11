// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import "io"

// NewForwardingPolicy creates a new ForwardingPolicy IE.
// It returns nil if id exceeds the 255-octet limit.
func NewForwardingPolicy(id string) *IE {
	// TS 29.244 V20.1.0 clause 8.2.23.
	l := len(id)
	if l > 255 {
		return nil
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
		return i.Payload, nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// ForwardingPolicyIdentifier returns ForwardingPolicyIdentifier in string if the type of IE matches.
func (i *IE) ForwardingPolicyIdentifier() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}

	v, err := i.ForwardingPolicy()
	if err != nil {
		return "", err
	}

	l := len(v)
	if l < 1 {
		return "", io.ErrUnexpectedEOF
	}

	idlen := int(v[0])
	if l < idlen+1 {
		return "", io.ErrUnexpectedEOF
	}

	return string(v[1 : idlen+1]), nil
}
