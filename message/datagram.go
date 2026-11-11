// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import "io"

// ParseDatagram decodes every message in one UDP payload. Each advertised frame
// must fit exactly, FO must link all but the last frame, and only session
// messages may be bundled. Errors return no partially decoded bundle.
// TS 29.244 V20.1.0 clauses 6.5, 7.2.1A and 7.2.2.4.1.
func ParseDatagram(b []byte) ([]Message, error) {
	if len(b) == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	var messages []Message
	for len(b) > 0 {
		n, err := frameSize(b)
		if err != nil {
			return nil, err
		}
		frame := b[:n]
		h, err := parseMessageHeader(frame, 0)
		if err != nil {
			return nil, err
		}
		if (h.HasFO() || len(messages) > 0) && !h.HasSEID() {
			return nil, ErrInvalidHeader
		}
		if h.HasFO() && n == len(b) {
			return nil, io.ErrUnexpectedEOF
		}
		if !h.HasFO() && n != len(b) {
			return nil, ErrInvalidLength
		}
		m, err := parseMessage(frame)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
		b = b[n:]
	}
	return messages, nil
}

// MarshalDatagram encodes one or more messages for one peer into a UDP payload.
// It sets FO in the output, without changing the inputs' FO flags. Bundling
// requires session messages and prior peer support negotiation; callers must
// ensure all messages target the same peer. Transport size and IP fragmentation
// policy are the caller's responsibility.
// TS 29.244 V20.1.0 clauses 6.5 and 7.2.1A.
func MarshalDatagram(messages ...Message) ([]byte, error) {
	if len(messages) == 0 {
		return nil, ErrInvalidLength
	}
	var b []byte
	for n, m := range messages {
		if m == nil {
			return nil, ErrNilMessage
		}
		size := m.MarshalLen()
		if size < 8 {
			return nil, ErrNilMessage
		}
		frame := make([]byte, size)
		if err := m.MarshalTo(frame); err != nil {
			return nil, err
		}
		h, err := parseMessageHeader(frame, 0)
		if err != nil {
			return nil, err
		}
		if len(messages) > 1 && !h.HasSEID() {
			return nil, ErrInvalidHeader
		}
		frame[0] &^= 4
		if n < len(messages)-1 {
			frame[0] |= 4
		}
		b = append(b, frame...)
	}
	return b, nil
}
