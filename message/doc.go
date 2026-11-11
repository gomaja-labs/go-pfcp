// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package message encodes and decodes PFCP messages as specified by
// TS 29.244 V20.1.0.
//
// Parse accepts exactly one standalone message. ParseDatagram accepts an entire
// UDP payload and checks every frame and follow-on flag. MarshalDatagram sets
// follow-on flags in the output; callers must negotiate bundling and select a
// single peer before using it. Typed ParseX functions and UnmarshalBinary decode
// one framed message and permit FO; they do not consume a surrounding datagram.
//
// Constructors derive S from the message type and default to version 1 with no
// priority. Session constructors take (seid, seq, ies...). Use SetMP with a value
// from 0 to 15 to enable priority, and ClearMP to remove it. Encoding rejects
// sequence numbers above 0xffffff. Header.MessagePriority stores the unshifted
// value, and Header.Length is informational: encoding always computes length.
//
// Message decoding returns UnsupportedVersionError before interpreting IEs of
// an unsupported version. Its Header lets a caller send Version Not Supported.
// ParseHeader exposes the common header even for unsupported versions. Known
// message types reject a contradictory S flag; unknown types retain their
// received framing for forward compatibility and decode as Generic without
// logging. The application remains responsible for protocol error procedures,
// mandatory IE checks and negotiated features.
//
// Decoded payloads borrow the input buffer. Keep the buffer unchanged while
// using the decoded message. Appending to a header payload cannot overwrite
// another frame. Encoding reads current IE fields, using the received IE type
// order as a template. Occurrences of each type use their current field/list
// order; excess occurrences append in constructor order. Unedited messages
// retain their wire order. Additional occurrences of a named singleton are
// retained in IEs.
package message
