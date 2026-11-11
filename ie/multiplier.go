// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
)

// NewMultiplier creates a new Multiplier IE.
func NewMultiplier(val uint64, exp uint32) *IE {
	i := New(Multiplier, make([]byte, 12))
	binary.BigEndian.PutUint64(i.Payload[0:8], val)
	binary.BigEndian.PutUint32(i.Payload[8:12], exp)

	return i
}

// Multiplier returns Multiplier in []byte if the type of IE matches.
func (i *IE) Multiplier() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case Multiplier:
		if err := i.validateFixedValue(); err != nil {
			return nil, err
		}
		return i.Payload, nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// ValueDigits returns ValueDigits in uint64 if the type of IE matches.
func (i *IE) ValueDigits() (uint64, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case Multiplier:
		return i.ValueAsUint64()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// Exponent returns Exponent in uint32 if the type of IE matches.
func (i *IE) Exponent() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case Multiplier:
		if err := i.validateFixedValue(); err != nil {
			return 0, err
		}
		if len(i.Payload) < 12 {
			return 0, io.ErrUnexpectedEOF
		}

		return binary.BigEndian.Uint32(i.Payload[8:12]), nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
