// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"errors"
	"fmt"
)

// Error definitions.
var (
	ErrTooShortToParse = errors.New("too short to decode as GTP")
	ErrInvalidLength   = errors.New("length value is invalid")

	ErrInvalidType = errors.New("invalid type")
	ErrIENotFound  = errors.New("could not find the specified IE in a grouped IE")

	ErrMalformed = errors.New("malformed IE")

	ErrElementNotFound = errors.New("element not found")
)

// DecodeError identifies the IE whose declared value could not be decoded.
// It unwraps to the underlying error, including io.ErrUnexpectedEOF and
// ErrInvalidLength. An incomplete four-octet header has no DecodeError because
// its type and length cannot yet be read together.
// TS 29.244 V20.1.0 clauses 7.6.3 and 7.6.7 require the offending IE type
// when an invalid IE length is reported.
type DecodeError struct {
	Type IEType
	Err  error
}

// Error describes the IE and the underlying decoding failure.
func (e *DecodeError) Error() string {
	return fmt.Sprintf("decode IE %d: %v", e.Type, e.Err)
}

// Unwrap returns the underlying decoding failure.
func (e *DecodeError) Unwrap() error { return e.Err }

// IEType returns the offending IE type.
func (e *DecodeError) IEType() IEType { return e.Type }

// InvalidTypeError indicates the type of IE is invalid.
type InvalidTypeError struct {
	Type IEType
}

// Error returns message with the invalid type given.
func (e *InvalidTypeError) Error() string {
	return fmt.Sprintf("got invalid type: %d", e.Type)
}

// InvalidNodeIDError indicates the NodeID value is invalid.
type InvalidNodeIDError struct {
	ID uint8
}

// Error returns message with the invalid NodeID given.
func (e *InvalidNodeIDError) Error() string {
	return fmt.Sprintf("got invalid NodeID: %d", e.ID)
}
