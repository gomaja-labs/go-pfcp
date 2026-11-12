// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"encoding/binary"
	"io"

	"github.com/gomaja-labs/go-pfcp/ie"
)

type messageFields interface {
	MarshalTo([]byte) error
	MarshalLen() int
}

func marshalMessage(m messageFields) ([]byte, error) {
	n := m.MarshalLen()
	if n == 0 {
		return nil, ErrNilMessage
	}
	if n-4 > 0xffff {
		return nil, ErrInvalidLength
	}
	b := make([]byte, n)
	if err := m.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

func messageLen(h *Header, ies []*ie.IE) int {
	if h == nil {
		return 0
	}
	n := 8
	if h.HasSEID() {
		n += 8
	}
	for _, i := range ies {
		if i != nil {
			n += i.MarshalLen()
		}
	}
	return n
}

// orderedIEs uses the received top-level type order as a template. Within a
// type, occurrences follow their current field/list order. Excess occurrences
// append in constructor order. Payload supplies only the type-order template;
// all values are encoded from the current exported IE fields.
// TS 29.244 V20.1.0 clauses 7.2.1 and 7.2.3.4.
func orderedIEs(h *Header, ies []*ie.IE) []*ie.IE {
	if len(h.Payload) == 0 {
		return ies
	}
	queues := make(map[ie.IEType][]*ie.IE)
	for _, i := range ies {
		if i != nil {
			queues[i.Type] = append(queues[i.Type], i)
		}
	}
	result := make([]*ie.IE, 0, len(ies))
	used := make(map[*ie.IE]int)
	for b := h.Payload; len(b) >= 4; {
		typ := ie.IEType(binary.BigEndian.Uint16(b[:2]))
		n := 4 + int(binary.BigEndian.Uint16(b[2:4]))
		if n > len(b) {
			break
		}
		if q := queues[typ]; len(q) > 0 {
			result = append(result, q[0])
			used[q[0]]++
			queues[typ] = q[1:]
		}
		b = b[n:]
	}
	for _, i := range ies {
		if i == nil {
			continue
		}
		if used[i] > 0 {
			used[i]--
			continue
		}
		result = append(result, i)
	}
	return result
}

func marshalMessageTo(h *Header, ies []*ie.IE, b []byte) error {
	if h == nil {
		return ErrNilMessage
	}
	n := messageLen(h, ies)
	if n-4 > 0xffff {
		return ErrInvalidLength
	}
	if len(b) < n {
		return io.ErrShortBuffer
	}
	// Validate before encoding any IEs and do not trust the stale payload length.
	frame := *h
	frame.Payload = nil
	frame.Flags &^= 4
	if err := frame.validate(); err != nil {
		return err
	}
	payload := make([]byte, n-frame.MarshalLen())
	offset := 0
	for _, i := range orderedIEs(h, ies) {
		if i == nil {
			continue
		}
		size := i.MarshalLen()
		if err := i.MarshalTo(payload[offset : offset+size]); err != nil {
			return err
		}
		offset += size
	}
	frame.Payload = payload
	if err := frame.MarshalTo(b); err != nil {
		return err
	}
	return nil
}

func parseMessageHeader(b []byte, expected uint8) (*Header, error) {
	h, err := ParseHeader(b)
	if err != nil {
		return nil, err
	}
	if err := h.validate(); err != nil {
		return nil, err
	}
	if expected != 0 && h.Type != expected {
		return nil, ErrInvalidHeader
	}
	return h, nil
}

// assignIE retains additional occurrences without changing the public field
// layout. The first occurrence occupies the named field, the rest remain in OtherIEs.
// TS 29.244 V20.1.0 clause 7.6.10.
func assignIE(dst **ie.IE, value *ie.IE, extras *[]*ie.IE) {
	if *dst == nil {
		*dst = value
	} else {
		*extras = append(*extras, value)
	}
}
