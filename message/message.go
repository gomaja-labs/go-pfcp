// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

//go:generate go run ../internal/cmd/pfcpgen -target message

// Message is an interface that defines PFCP messages.
type Message interface {
	// MarshalTo encodes the message into a buffer of at least MarshalLen bytes.
	MarshalTo([]byte) error
	// UnmarshalBinary replaces the receiver with one decoded frame. Decoding
	// borrows b; keep b unchanged while the message is in use.
	UnmarshalBinary(b []byte) error
	// MarshalLen returns the encoded size in octets.
	MarshalLen() int
	// MessageType returns the PFCP message type.
	MessageType() uint8
	// MessageTypeName returns the specification name of the message type.
	MessageTypeName() string
	// Version returns the received or constructed PFCP version.
	Version() int
	// SEID returns the session endpoint identifier, or zero for node messages.
	SEID() uint64
	// Sequence returns the transaction sequence number.
	Sequence() uint32
	// SetSequenceNumber sets the sequence; encoding checks its 24-bit range.
	SetSequenceNumber(seq uint32)
	// IsRequest reports whether the message type is a request.
	IsRequest() bool
}

// parseMessage decodes a framed message whose datagram boundaries were checked.
func parseMessage(b []byte) (Message, error) {
	if _, err := parseMessageHeader(b, 0); err != nil {
		return nil, err
	}

	m := newMessage(b[1])

	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}

// Parse decodes exactly one standalone PFCP message. FO=1 is rejected with
// ErrBundledMessage; trailing bytes with FO=0 return InvalidLengthError. Use
// ParseDatagram for UDP payloads that may contain bundles.
// A complete header with inconsistent Length returns InvalidLengthError;
// an incomplete header returns io.ErrUnexpectedEOF (TS 29.244 V20.1.0 clause 7.6.3).
// Decoding may borrow b; keep b unchanged while the result is in use.
func Parse(b []byte) (Message, error) {
	if _, err := frameSize(b); err != nil {
		return nil, err
	}
	if b[0]&4 != 0 {
		return nil, ErrBundledMessage
	}
	return parseMessage(b)
}
