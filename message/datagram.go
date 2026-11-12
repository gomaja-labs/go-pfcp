// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"errors"
	"fmt"
)

// DatagramError identifies a failed frame in a UDP payload.
// Index is zero-based; Offset counts octets from the start of the datagram.
// Header is nil only when the input ended before the frame's complete header
// (TS 29.244 V20.1.0 clause 7.6.3: silently discard it). For an unsupported
// version it holds the readable prefix; see UnsupportedVersionError.
type DatagramError struct {
	// Index is the zero-based frame index.
	Index int
	// Offset is the frame start in octets from the datagram start.
	Offset int
	// Header is the received header or unsupported-version prefix; nil means
	// the input ended before the complete frame header.
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
// joined error containing at most 16 *DatagramError details, in wire order,
// followed by a count of further failed frames when details were omitted.
// A malformed body does not prevent decoding later frames. A framing failure
// (incomplete header, inconsistent length, invalid bundle framing) or an
// unsupported version whose layout is unknown stops further decoding.
// A frame with FO=0 that does not end the datagram has an inconsistent length:
// it is reported as InvalidLengthError and neither decoded nor returned.
// A frame with FO=1 that ends the datagram is decoded, and the frame it
// announces is reported missing at the datagram end with a nil Header.
// Applications must process both the returned messages and the error.
// TS 29.244 V20.1.0 clauses 6.5, 7.2.1A, 7.2.2.4.1 and 7.6.2–7.6.7.
func ParseDatagram(b []byte) ([]Message, error) {
	var messages []Message
	var failures []error
	index, offset, omitted := 0, 0, 0
	fail := func(h *Header, err error) {
		if len(failures) < 16 {
			failures = append(failures, &DatagramError{Index: index, Offset: offset, Header: h, Err: err})
		} else {
			omitted++
		}
	}
	for {
		n, err := frameSize(b)
		if err != nil {
			var v *UnsupportedVersionError
			var l *InvalidLengthError
			var h *Header
			if errors.As(err, &v) {
				h = v.Header
			} else if errors.As(err, &l) {
				h = l.Header
			}
			fail(h, err)
			break
		}
		frame := b[:n:n]
		h := decodeHeaderFields(frame)
		h.Payload = frame[h.MarshalLen():n:n]
		if (h.HasFO() || index > 0) && !h.HasSEID() {
			fail(&h, ErrInvalidHeader)
			break
		}
		if !h.HasFO() && n != len(b) {
			err := invalidLength(b)
			fail(err.Header, err)
			break
		}
		// Clause 6.5: handle each framed message as if sent on its own.
		m, err := parseMessage(frame)
		if err != nil {
			fail(&h, err)
		} else {
			messages = append(messages, m)
		}
		if !h.HasFO() {
			break
		}
		b = b[n:]
		offset += n
		index++
	}
	if omitted != 0 {
		failures = append(failures, fmt.Errorf("%d further frames failed", omitted))
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
