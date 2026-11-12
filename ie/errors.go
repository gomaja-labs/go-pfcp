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
	// ErrInvalidLength reports an IE length inconsistent with its framing or value format.
	ErrInvalidLength = errors.New("length value is invalid")

	// ErrInvalidType reports an inappropriate IE type.
	ErrInvalidType = errors.New("invalid type")
	// ErrIENotFound reports an absent IE, including a nil receiver.
	ErrIENotFound = errors.New("IE not found")

	// ErrMalformed reports a malformed IE value.
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
// ErrInvalidLength. Incomplete headers and declared values beyond the available
// octets wrap both ErrInvalidLength and io.ErrUnexpectedEOF (clause 7.6.7).
// A value helper requiring more value octets returns io.ErrUnexpectedEOF without
// ErrInvalidLength (clause 7.6.8).
// TS 29.244 V20.1.0 clauses 7.6.3 and 7.6.7 require the offending IE type
// when an invalid IE length is reported.
type DecodeError struct {
	// Type is the innermost IE whose header or value could not be decoded.
	// An incomplete header carries its type when two octets are available;
	// otherwise Type is zero, with enclosing groups retained in Path.
	Type IEType
	// Path lists enclosing grouped types, outermost first; it excludes Type.
	Path []IEType
	// Err is the underlying failure, never another DecodeError.
	Err error
}

// Error describes the IE and the underlying decoding failure.
func (e *DecodeError) Error() string {
	if len(e.Path) != 0 {
		return fmt.Sprintf("decode IE %d in groups %v: %v", e.Type, e.Path, e.Err)
	}
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
	// Type identifies the IE that could not be encoded.
	Type IEType
	// Err is the underlying failure.
	Err error
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
	// Type is the unexpected IE type.
	Type IEType
}

// Error returns message with the invalid type given.
func (e *InvalidTypeError) Error() string {
	return fmt.Sprintf("got invalid type: %d", e.Type)
}

// InvalidNodeIDError indicates the NodeID value is invalid.
type InvalidNodeIDError struct {
	// ID is the unsupported Node ID kind.
	ID uint8
}

// Error returns message with the invalid NodeID given.
func (e *InvalidNodeIDError) Error() string {
	return fmt.Sprintf("got invalid NodeID: %d", e.ID)
}

// withDecodePath preserves one offending IE, adding its enclosing group.
// TS 29.244 V20.1.0 clause 7.6.7 requires the offending type, not a parent.
func withDecodePath(parent IEType, err error) error {
	var leaf *DecodeError
	if errors.As(err, &leaf) {
		path := make([]IEType, 1, len(leaf.Path)+1)
		path[0] = parent
		path = append(path, leaf.Path...)
		return &DecodeError{Type: leaf.Type, Path: path, Err: leaf.Err}
	}
	return &DecodeError{Type: parent, Err: err}
}
