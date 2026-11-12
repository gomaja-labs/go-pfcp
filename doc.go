// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package pfcp documents a codec for the Packet Forwarding Control Protocol
// defined by 3GPP TS 29.244 V20.1.0. Consumers track main; API changes do not
// retain compatibility aliases.
//
// # Packages and coverage
//
// Use github.com/gomaja-labs/go-pfcp/ie to construct and read Information
// Elements, and github.com/gomaja-labs/go-pfcp/message to encode and decode
// messages and bundled UDP payloads. This root package contains documentation.
//
// All message types assigned by Table 7.3-1 have typed constructors and parsers.
// All assigned IE types in Table 8.1.2-1 have names and constructors: grouped
// IEs expose typed child fields, and leaf IEs expose value helpers. Reserved,
// spare and unassigned values have no named API; unknown values remain
// available through raw IEs and generic messages. Support means codec coverage,
// not validation of every procedure, mandatory IE, condition or peer capability.
// The versioned catalogue and its coverage checks are under internal/spec.
//
// # Application responsibilities
//
// The library does not open sockets, log messages, establish associations,
// retransmit requests or manage sessions. Applications provide transport,
// transaction state, capability negotiation and semantic validation.
//
// Decode failures distinguish incomplete headers, inconsistent framing and
// offending IE types. Use errors.Is for sentinel errors and errors.As for
// ie.DecodeError, ie.EncodeError, message.UnsupportedVersionError,
// message.InvalidLengthError and message.DatagramError. Reading
// an IE value can fail after message framing succeeds; check helper errors too.
// Invalid arguments retained by constructors surface when encoding.
//
// Decoding borrows the input bytes. Keep them unchanged, or copy them before
// decoding if a receive buffer will be reused. Messages, IEs and their slices
// are mutable and have no internal locking. Concurrent encoding is safe when
// the entire object graph and borrowed input are unchanged and each call uses
// its own output buffer; synchronize all mutations with readers.
package pfcp
