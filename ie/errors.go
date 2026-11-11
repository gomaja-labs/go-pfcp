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
	ErrInvalidLength = errors.New("length value is invalid")

	ErrInvalidType = errors.New("invalid type")
	ErrIENotFound  = errors.New("could not find the specified IE in a grouped IE")

	ErrMalformed = errors.New("malformed IE")

	// ErrInvalidValue reports a constructor argument that cannot be encoded,
	// such as a malformed domain name or a time outside the 3GPP timestamp
	// window. The constructor still returns an IE; Marshal reports the error.
	ErrInvalidValue = errors.New("value cannot be encoded")

	// ErrNestingTooDeep reports grouped IEs nested deeper than any structure
	// TS 29.244 defines; the decoder stops there instead of following them.
	ErrNestingTooDeep = errors.New("grouped IEs are nested too deeply")
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

// EncodeError identifies the IE that could not be encoded. It unwraps to the
// underlying error, such as ErrInvalidValue for an argument a constructor could
// not encode.
type EncodeError struct {
	Type IEType
	Err  error
}

// Error describes the IE and the underlying encoding failure.
func (e *EncodeError) Error() string {
	return fmt.Sprintf("encode IE %d: %v", e.Type, e.Err)
}

// Unwrap returns the underlying encoding failure.
func (e *EncodeError) Unwrap() error { return e.Err }

// IEType returns the offending IE type.
func (e *EncodeError) IEType() IEType { return e.Type }

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
