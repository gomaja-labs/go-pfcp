// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

//go:generate go run ../internal/cmd/pfcpgen -target message

// Message is an interface that defines PFCP messages.
type Message interface {
	MarshalTo([]byte) error
	UnmarshalBinary(b []byte) error
	MarshalLen() int
	MessageType() uint8
	MessageTypeName() string
	Version() int
	SEID() uint64
	Sequence() uint32
	SetSequenceNumber(seq uint32)
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
// ErrBundledMessage; trailing bytes with FO=0 are an invalid length. Use
// ParseDatagram for UDP payloads that may contain bundles.
func Parse(b []byte) (Message, error) {
	if _, err := frameSize(b); err != nil {
		return nil, err
	}
	if b[0]&4 != 0 {
		return nil, ErrBundledMessage
	}
	return parseMessage(b)
}
