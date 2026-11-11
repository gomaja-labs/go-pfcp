// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package message encodes and decodes PFCP messages as specified by
// TS 29.244 V20.1.0.
//
// Message schemas, constructors, dispatch and type names are generated from the
// versioned catalogue under internal/spec. Constructors take (seq, ies...);
// session constructors take (seid, seq, ies...). They derive S from the message
// type and default to version 1 with no priority. Nil IEs are ignored. Use SetMP
// with a value from 0 to 15 to enable priority, and ClearMP to remove it. Encoding
// rejects sequence numbers above 0xffffff. Header.MessagePriority stores the
// unshifted value. The wire length is derived; Header has no stored Length field
// and messages have no SetLength method.
//
// Lists contain every occurrence in received order. For a singular field only
// the first occurrence is selected; further occurrences and unexpected types
// remain in OtherIEs (clauses 7.2.3.4, 7.6.9 and 7.6.10).
//
// New messages encode fields in table order followed by OtherIEs. Decoded
// messages use the received top-level IE type order as a template; within a
// type, current field/list order is used. Excess occurrences append in table
// order. This preserves unedited frames without promising identity-based
// positioning after edits among same-type IEs. Header.Payload supplies only the
// ordering template; values come from the current fields. Decoding borrows the
// input buffer, which must remain unchanged while the message is in use.
// Concurrent Marshal and MarshalTo calls are safe when callers keep the
// message, IEs and borrowed input unchanged and use separate output buffers.
//
// Parse accepts exactly one standalone message. ParseDatagram accepts an entire
// UDP payload and checks every frame and follow-on flag. MarshalDatagram sets
// follow-on flags in the output without changing the inputs. Callers must
// negotiate bundling and select a single peer before using it. Typed ParseX
// functions and UnmarshalBinary decode one frame and permit FO; they do not
// consume a surrounding datagram.
//
// Message decoding returns UnsupportedVersionError before interpreting IEs of
// an unsupported version. Its Header lets a caller send Version Not Supported.
// ParseHeader exposes the common header even for unsupported versions. Known
// message types reject a contradictory S flag; unknown types retain their
// received framing and decode as Generic without logging. Presence requirements,
// conditional applicability and peer capability negotiation remain application
// responsibilities.
package message
