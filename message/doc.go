// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package message encodes and decodes PFCP messages as specified by
// 3GPP TS 29.244 V20.1.0.
//
// # Building and encoding
//
// Node constructors take (seq, ies...); session constructors take
// (seid, seq, ies...). They derive S from the message type and default to
// version 1 without priority. Nil IEs are ignored. For example,
// NewHeartbeatRequest(seq, ie.NewRecoveryTimeStamp(started)) constructs a
// heartbeat; NewSessionEstablishmentRequest(seid, seq, ies...) builds a session
// request. The application supplies the appropriate peer SEID, including zero
// where required for initial establishment (clause 7.2.2.4.1).
//
// Concrete messages offer Marshal, which allocates, and MarshalTo, which writes
// to a buffer of at least MarshalLen bytes. The Message interface exposes
// MarshalTo and MarshalLen. Lengths are derived from current content; Marshal
// never mutates the message. Standalone encoding always clears FO; only
// MarshalDatagram sets it in bundled output. Encoding propagates ie.EncodeError
// for an invalid IE, including one nested in a group.
//
// Sequence numbers are 24-bit values; encoding rejects values above 0xffffff.
// SetSequenceNumber does not allocate transaction identifiers or handle wrap.
// Responses use their request's sequence (clause 6.4). For session messages,
// SetMP with a value from 0 to 15 enables priority; ClearMP removes it.
// Header.MessagePriority holds the unshifted value (clause 7.2.2.4.1).
//
// # Decoding and ordering
//
// Parse accepts exactly one standalone message and returns a Message for a
// type switch. Typed ParseXxx functions, such as ParseHeartbeatRequest, and
// UnmarshalBinary decode one frame of the expected type and permit FO. They
// do not consume a surrounding datagram. Unknown message types decode as
// Generic. The application checks mandatory IEs, conditional applicability,
// peer capabilities and procedural state.
//
// Lists contain every occurrence in received order. Singular fields select the
// first occurrence; later occurrences and unexpected types remain in OtherIEs
// (clauses 7.2.3.4, 7.6.9 and 7.6.10). Message decoding frames all top-level IEs
// but can defer errors inside grouped values until a typed helper is called.
//
// New messages encode in table order followed by OtherIEs. Decoded messages
// use received top-level IE type order as a template; within a type, current
// field/list order is used. Excess occurrences append in table order. This
// preserves unedited IE ordering without promising identity-based positioning
// after edits among same-type IEs. Header.Payload supplies the ordering
// template; values come from current fields. Encoding clears header spare bits.
//
// # Bundled datagrams
//
// ParseDatagram reads an entire UDP payload and checks every frame and FO
// (follow-on) flag. MarshalDatagram sets FO on all but the last output frame
// without changing its inputs. Bundling requires session messages, prior peer
// support negotiation and a single destination peer (clauses 6.5 and 7.2.1A).
// Parse rejects FO with ErrBundledMessage; use ParseDatagram for UDP receives
// that may contain bundles. ParseDatagram returns successful frames even when
// others fail, including a decodable frame whose FO chain check fails.
// The joined error contains at most 16 *DatagramError details (zero-based Index,
// octet Offset, received Header and underlying Err), followed by a count of
// further failed frames. Later failures still allow valid frames to decode;
// omitted details are unavailable through errors.Is/As. Body failures allow later
// frames to decode; an incomplete header, inconsistent length or FO chain,
// invalid bundle framing, or unsupported version stops decoding. Always handle
// both the messages and the error. Transport size limits remain the caller's concern.
//
// # Version and error handling
//
// Only PFCP version 1 is implemented. UnsupportedVersionError is returned
// as soon as octet 1 exists, before any length or framing check. errors.As
// exposes its Header and
// received Version so an application can send NewVersionNotSupportedResponse.
// That response carries the highest supported version, 1, and only its header
// (clauses 7.6.2 and 7.4.4.7). The error Header contains only available fields;
// sequence and SEID are best-effort at version-1 offsets. ParseHeader follows
// the same version-first rule. Known message types reject contradictory S flags.
//
// An incomplete header returns io.ErrUnexpectedEOF; inconsistent declared
// framing returns ErrInvalidLength. Truncated IE framing wraps both
// ie.ErrInvalidLength and io.ErrUnexpectedEOF and identifies its available type
// through ie.DecodeError. Use errors.Is and errors.As, and check value-helper
// errors even after Parse succeeds.
//
// # Ownership and concurrency
//
// Decoding borrows the input buffer, which must remain unchanged while the
// message is in use. Copy a reusable receive buffer before parsing it.
// Decoding allocates in proportion to the number of frames and IEs; small
// frames can allocate much larger message structures. Bound datagram size and
// receive rate at the transport layer.
// Concurrent Marshal and MarshalTo calls are safe when callers keep the
// message, IEs and borrowed input unchanged and use separate output buffers.
// Synchronize mutations with all readers. The package does not log or provide
// sockets, retransmission, association state or session management.
package message
