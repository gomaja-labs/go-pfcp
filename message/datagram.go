// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"errors"
	"fmt"
	"io"
)

// DatagramError identifies a failed frame in a UDP payload.
// Index is zero-based; Offset counts octets from the start of the datagram.
// Header is available when the complete version-1 header was decoded, or from
// UnsupportedVersionError when only a future-version prefix could be read.
type DatagramError struct {
	// Index is the zero-based frame index.
	Index int
	// Offset is the frame start in octets from the datagram start.
	Offset int
	// Header is the received header when it can be decoded.
	Header *Header
	// Err is the underlying decoding or framing failure.
	Err error
}

// Error describes the frame position and its failure.
func (e *DatagramError) Error() string {
	return fmt.Sprintf("PFCP frame %d at octet %d: %v", e.Index, e.Offset, e.Err)
}

// Unwrap exposes the underlying framing, version or IE decoding failure.
func (e *DatagramError) Unwrap() error { return e.Err }

// ParseDatagram decodes messages in one UDP payload, borrowing b. Keep b unchanged
// while using the results. It returns every successfully decoded message and a
// joined error containing one *DatagramError per failed frame, in wire order.
// A malformed body does not prevent decoding later frames. A framing failure
// (incomplete header, inconsistent length or FO chain, invalid bundle framing),
// or an unsupported version whose layout is unknown, stops further decoding.
// Applications must process both the returned messages and the error.
// TS 29.244 V20.1.0 clauses 6.5, 7.2.1A, 7.2.2.4.1 and 7.6.2–7.6.7.
func ParseDatagram(b []byte) ([]Message, error) {
	var messages []Message
	var failures []error
	index, offset := 0, 0
	fail := func(h *Header, err error) {
		failures = append(failures, &DatagramError{Index: index, Offset: offset, Header: h, Err: err})
	}
	if len(b) == 0 {
		fail(nil, io.ErrUnexpectedEOF)
	}
	for len(b) > 0 {
		n, err := frameSize(b)
		if err != nil {
			var v *UnsupportedVersionError
			var h *Header
			if errors.As(err, &v) {
				h = v.Header
			}
			fail(h, err)
			break
		}
		frame := b[:n:n]
		h, err := ParseHeader(frame)
		if err != nil {
			fail(nil, err)
			break
		}
		if (h.HasFO() || index > 0) && !h.HasSEID() {
			fail(h, ErrInvalidHeader)
			break
		}
		if h.HasFO() && n == len(b) {
			fail(h, io.ErrUnexpectedEOF)
			break
		}
		if !h.HasFO() && n != len(b) {
			fail(h, ErrInvalidLength)
			break
		}
		m, err := parseMessage(frame)
		if err != nil {
			fail(h, err)
		} else {
			messages = append(messages, m)
		}
		b = b[n:]
		offset += n
		index++
	}
	return messages, errors.Join(failures...)
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
