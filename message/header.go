// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Errors returned for invalid PFCP header values or framing.
var (
	ErrInvalidLength    = errors.New("invalid PFCP message length")
	ErrInvalidHeader    = errors.New("invalid PFCP header")
	ErrNilMessage       = errors.New("nil PFCP message or header")
	ErrSequenceOverflow = errors.New("PFCP sequence number exceeds 24 bits")
	ErrInvalidPriority  = errors.New("PFCP priority requires S=1, MP=1 and a value from 0 to 15")
	ErrBundledMessage   = errors.New("PFCP follow-on message: use ParseDatagram")
)

// UnsupportedVersionError carries the header needed to answer with a Version
// Not Supported Response. No version-specific IEs have been decoded.
// TS 29.244 V20.1.0 clause 7.6.2.
type UnsupportedVersionError struct{ Header *Header }

func (e *UnsupportedVersionError) Error() string {
	if e == nil || e.Header == nil {
		return "unsupported PFCP version"
	}
	return fmt.Sprintf("unsupported PFCP version %d", e.Header.Version())
}

// Header represents a PFCP header. The wire length is always derived from the
// encoded content. MessagePriority is a value from 0 to 15,
// meaningful only when S and MP are set. Payload is the framed message body.
type Header struct {
	Flags           uint8 // Version, FO, MP, S; spare bits are ignored on receive.
	Type            uint8
	SEID            uint64
	SequenceNumber  uint32 // Values above 0xffffff are rejected when encoding.
	MessagePriority uint8  // Priority value, not a shifted wire octet.
	Payload         []byte
}

// NewHeader creates a version 1 header. S is derived from the message type;
// node messages ignore seid. Unknown types default to session framing.
// TS 29.244 V20.1.0 clause 7.2.2.4.1 and Table 7.3-1.
func NewHeader(mtype uint8, seid uint64, seq uint32, payload []byte) *Header {
	h := &Header{Flags: 0x21, Type: mtype, SEID: seid, SequenceNumber: seq, Payload: payload}
	if isNodeType(mtype) {
		h.Flags = 0x20
		h.SEID = 0
	}
	return h
}

// Node/session framing follows TS 29.244 V20.1.0 Table 7.3-1.
func isNodeType(t uint8) bool       { return t >= 1 && t <= 17 }
func knownSessionType(t uint8) bool { return t >= 50 && t <= 57 }

// HasFO reports whether another message follows in the datagram.
func (h *Header) HasFO() bool { return h != nil && h.Flags&4 != 0 }

// HasMP reports whether the priority is present.
func (h *Header) HasMP() bool { return h != nil && h.Flags&2 != 0 }

// HasSEID reports whether the SEID field is present.
func (h *Header) HasSEID() bool { return h != nil && h.Flags&1 != 0 }

// SetFO changes only the follow-on bit. Prefer MarshalDatagram when bundling.
func (h *Header) SetFO(on bool) {
	if h == nil {
		return
	}
	if on {
		h.Flags |= 4
	} else {
		h.Flags &^= 4
	}
}

func (h *Header) validate() error {
	if h == nil {
		return ErrNilMessage
	}
	if h.Version() != 1 {
		return &UnsupportedVersionError{Header: h}
	}
	if h.SequenceNumber > 0xffffff {
		return ErrSequenceOverflow
	}
	if isNodeType(h.Type) && h.HasSEID() || knownSessionType(h.Type) && !h.HasSEID() {
		return ErrInvalidHeader
	}
	if !h.HasSEID() && (h.HasFO() || h.HasMP()) {
		return ErrInvalidHeader
	}
	if h.HasMP() && h.MessagePriority > 15 {
		return ErrInvalidPriority
	}
	if h.MarshalLen()-4 > 0xffff {
		return ErrInvalidLength
	}
	return nil
}

// Marshal returns the complete header and payload encoding.
func (h *Header) Marshal() ([]byte, error) {
	if err := h.validate(); err != nil {
		return nil, err
	}
	b := make([]byte, h.MarshalLen())
	if err := h.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo encodes the header and payload, computing Message Length and zeroing
// spare bits. It returns io.ErrShortBuffer without writing if b is too short.
// TS 29.244 V20.1.0 clauses 7.2.2.1 through 7.2.2.4.1.
func (h *Header) MarshalTo(b []byte) error {
	if err := h.validate(); err != nil {
		return err
	}
	n := h.MarshalLen()
	if len(b) < n {
		return io.ErrShortBuffer
	}
	b[0], b[1] = h.Flags&^0x18, h.Type
	binary.BigEndian.PutUint16(b[2:4], uint16(n-4))
	offset := 4
	if h.HasSEID() {
		binary.BigEndian.PutUint64(b[offset:offset+8], h.SEID)
		offset += 8
	}
	copy(b[offset:offset+3], uint32To24(h.SequenceNumber))
	b[offset+3] = 0
	if h.HasSEID() && h.HasMP() {
		b[offset+3] = h.MessagePriority << 4
	}
	copy(b[offset+4:n], h.Payload)
	return nil
}

// ParseHeader decodes exactly one framed message, including its payload. FO may
// be set; this low-level API does not validate the surrounding datagram or version.
// Use Parse or ParseDatagram for version-aware message decoding.
func ParseHeader(b []byte) (*Header, error) {
	h := &Header{}
	if err := h.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return h, nil
}

// frameSize validates the advertised extent without reading beyond it.
// TS 29.244 V20.1.0 clause 7.2.2.4.1: Length excludes the first four octets.
func frameSize(b []byte) (int, error) {
	if len(b) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	minimum := 8
	if b[0]&1 != 0 {
		minimum = 16
	}
	// TS 29.244 V20.1.0 clause 7.6.3 distinguishes a short respective
	// header from an inconsistent Length field after a complete header.
	if len(b) < minimum {
		return 0, io.ErrUnexpectedEOF
	}
	length := int(binary.BigEndian.Uint16(b[2:4]))
	if length < minimum-4 {
		return 0, ErrInvalidLength
	}
	n := length + 4
	if len(b) < n {
		return 0, ErrInvalidLength
	}
	return n, nil
}

// UnmarshalBinary decodes exactly one frame. Unsupported versions remain
// readable here so an application can inspect their common header. Payload
// borrows b; keep b unchanged while using decoded data. Spare bits
// are ignored. A failed decode leaves the receiver unchanged.
func (h *Header) UnmarshalBinary(b []byte) error {
	if h == nil {
		return ErrNilMessage
	}
	n, err := frameSize(b)
	if err != nil {
		return err
	}
	if n != len(b) {
		return ErrInvalidLength
	}
	v := Header{Flags: b[0] &^ 0x18, Type: b[1]}
	offset := 4
	if v.HasSEID() {
		v.SEID = binary.BigEndian.Uint64(b[offset : offset+8])
		offset += 8
	}
	v.SequenceNumber = uint24To32(b[offset : offset+3])
	if v.HasSEID() && v.HasMP() {
		v.MessagePriority = b[offset+3] >> 4
	}
	v.Payload = b[offset+4 : n : n]
	*h = v
	return nil
}

// MarshalLen returns the complete encoded size, or zero for a nil header.
func (h *Header) MarshalLen() int {
	if h == nil {
		return 0
	}
	n := 8 + len(h.Payload)
	if h.HasSEID() {
		n += 8
	}
	return n
}

// Version returns the three-bit version, or zero for nil.
func (h *Header) Version() int {
	if h == nil {
		return 0
	}
	return int(h.Flags >> 5)
}

// MessageType returns the wire message type, or zero for nil.
func (h *Header) MessageType() uint8 {
	if h == nil {
		return 0
	}
	return h.Type
}
func (h *Header) seid() uint64 {
	if !h.HasSEID() {
		return 0
	}
	return h.SEID
}

// SetSEID sets the session endpoint identifier. Node message types reject it.
func (h *Header) SetSEID(seid uint64) error {
	if h == nil {
		return ErrNilMessage
	}
	if isNodeType(h.Type) {
		return ErrInvalidHeader
	}
	h.Flags |= 1
	h.SEID = seid
	return nil
}

// Sequence returns the sequence number, or zero for nil.
func (h *Header) Sequence() uint32 {
	if h == nil {
		return 0
	}
	return h.SequenceNumber
}

// SetSequenceNumber assigns the sequence number. Encoding rejects values above
// 0xffffff rather than silently truncating them.
func (h *Header) SetSequenceNumber(seq uint32) {
	if h != nil {
		h.SequenceNumber = seq
	}
}

// SetMP enables priority with an unshifted value from 0 to 15. It changes no
// neighbouring flag. Node headers reject priority. Invalid input leaves h unchanged.
// TS 29.244 V20.1.0 clause 7.2.2.4.1.
func (h *Header) SetMP(priority uint8) error {
	if h == nil {
		return ErrNilMessage
	}
	if !h.HasSEID() || priority > 15 {
		return ErrInvalidPriority
	}
	h.Flags |= 2
	h.MessagePriority = priority
	return nil
}

// ClearMP removes priority and clears its stored value.
func (h *Header) ClearMP() {
	if h != nil {
		h.Flags &^= 2
		h.MessagePriority = 0
	}
}

// MP returns the priority value; it returns zero when priority is absent.
func (h *Header) MP() uint8 {
	if !h.HasSEID() || !h.HasMP() {
		return 0
	}
	return h.MessagePriority
}

// String returns a readable summary of header fields.
func (h *Header) String() string {
	if h == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{Version: %d, FO: %t, MP: %t, S: %t, Type: %d, Length: %d, SEID: %d, Sequence: %d, Priority: %d}", h.Version(), h.HasFO(), h.HasMP(), h.HasSEID(), h.Type, h.MarshalLen()-4, h.SEID, h.SequenceNumber, h.MP())
}
