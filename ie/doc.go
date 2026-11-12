// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package ie builds, encodes and reads PFCP Information Elements from
// 3GPP TS 29.244 V20.1.0.
//
// # Building and reading values
//
// Prefer named constructors such as NewPDRID and NewNodeIDAddress. New accepts
// a type and raw value octets, without a TLV header. NewUint8IE, NewUint16IE,
// NewUint32IE and NewUint64IE encode unsigned integers in network byte order;
// NewStringIE stores string octets and NewFQDNIE encodes DNS labels. A named
// helper such as (*IE).PDRID checks the IE type; ValueAsUint16, ValueAsString
// and ValueAsFQDN explicitly interpret a raw value. Check every returned error.
//
// Marshal allocates an encoding; MarshalTo writes into a caller's buffer of
// at least MarshalLen bytes. Lengths are calculated, not stored. Parse and
// UnmarshalBinary require exactly one IE; ParseMultiIEs reads an entire
// sequence. Parsed leaf values borrow the input, which must remain unchanged.
// New also retains the supplied value slice without copying it.
//
// # Grouped IEs
//
// NewCreatePDR(children...) builds a group. Alternatively, populate
// CreatePDRFields and call its IE method. Decode typed children with
// (*IE).CreatePDR, then read fields such as fields.PDRID.PDRID(). This pattern
// applies to each generated XxxFields type. Fields remain *IE values so callers
// can distinguish absence, null length and an invalid value. Presence markers
// in field comments describe the specification; the codec does not enforce
// mandatory or conditional presence.
//
// Decoded known groups use ChildIEs with nil Payload. IsGrouped uses the
// specification's static type catalogue. OtherIEs in a typed group preserves
// unknown children and extra occurrences of singular children; the first
// occurrence selects a singular field (clauses 7.6.9 and 7.6.10). Typed fields'
// IE methods build table order followed by OtherIEs. Keep the original IE to
// retain its received child ordering. Nesting beyond 16 groups returns
// ErrNestingTooDeep; this is a library resource limit.
//
// # Null-length and vendor-specific IEs
//
// Clause 5.6.3 removes provisioned values using null-length IEs. For example,
// New(URRID, nil) inside NewUpdatePDR requests removal of the PDR's URRs.
// After successful decoding, recognize a standard null IE by len(i.Payload) == 0
// and len(i.ChildIEs) == 0; a missing IE is a nil pointer instead. Parsing accepts null length under
// clause 8.1.2, but a helper requiring a value reports an error.
//
// NewVendorSpecificIE accepts a type with the high bit set, an Enterprise ID
// and raw value octets. The Enterprise ID occupies two octets of its length
// (clause 8.1.1), so an empty vendor value is not a zero wire length.
// NewVendorSpecificGroupedIE builds vendor children explicitly. Vendor types
// are not in the static grouped catalogue: ValueAsGrouped can interpret any
// IE's value as concatenated children without registering its type globally.
//
// # Errors and equality
//
// DecodeError identifies the offending Type, enclosing Path and underlying
// error. Use errors.As to obtain it and errors.Is to inspect its cause, such
// as io.ErrUnexpectedEOF, ErrInvalidLength or ErrNestingTooDeep. A header
// shorter than four octets returns io.ErrUnexpectedEOF without an offending
// type. Message decoding uses ParseMessageIEs to retain malformed grouped
// values as raw Payload; their typed helpers report errors on access.
//
// Constructors can retain an unencodable argument in an IE. Encoding it, or a
// group or message containing it, returns EncodeError wrapping ErrInvalidValue.
// InvalidTypeError means a helper received an inappropriate type. ErrIENotFound
// reports a nil IE receiver or an absent requested child. Equal compares
// encoded octets; it also recognizes matching retained constructor errors.
//
// # Timestamps and forward compatibility
//
// Four-octet timestamps carry seconds without an era. This package selects the
// Unix-seconds window [0, 2^32-1], from 1970-01-01 through
// 2106-02-07T06:28:15Z, and discards fractional seconds when constructing them.
// Out-of-window times produce an encoding error. This is the library's era
// policy for TS 29.244 clause 8.2.65; the wire cannot distinguish other eras.
//
// Clause 8.1.2 requires receivers to ignore unknown trailing octets of
// extendable IEs and to ignore spare bits. Value helpers read known fields;
// raw Payload retains received octets. Callers supplying raw octets or flag
// masks must send spare bits as zero.
// Fixed-length helpers reject excess value octets; variable-length helpers
// follow their IE's length-delimited format. Rebuilding from interpreted values
// need not preserve unknown extensions or received spare bits.
//
// IEs and their backing slices are mutable. Concurrent readers and encoders
// require an unchanged object graph and input; MarshalTo calls need separate
// output buffers. The package does not log or manage transport state.
package ie
