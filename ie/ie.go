// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/gomaja-labs/go-pfcp/internal/utils"
)

//go:generate go run ../internal/cmd/pfcpgen -target ie
//go:generate go run ../internal/testutil/registrygen
type IEType uint16

// IE represents an Information Element of PFCP messages.
type IE struct {
	Type         IEType
	EnterpriseID uint16
	Payload      []byte
	ChildIEs     []*IE

	// err is set by a constructor whose arguments cannot be encoded. Marshal
	// and MarshalTo report it as an *EncodeError, so an invalid argument is
	// neither sent nor silently dropped from a message or grouped IE.
	err error
}

// newInvalidIE returns an IE of type t that cannot be encoded because of err.
func newInvalidIE(t IEType, err error) *IE {
	return &IE{Type: t, err: err}
}

// Equal reports whether i and o encode to the same octets. Two IEs that
// cannot be encoded are equal when their errors read the same. A nil IE is
// equal only to another nil IE. Equal makes IEs comparable with go-cmp
// despite the unexported constructor error.
func (i *IE) Equal(o *IE) bool {
	if i == nil || o == nil {
		return i == o
	}
	if i.err != nil || o.err != nil {
		return i.err != nil && o.err != nil && i.Type == o.Type && i.err.Error() == o.err.Error()
	}
	a, errA := i.Marshal()
	b, errB := o.Marshal()
	if errA != nil || errB != nil {
		return false
	}
	return bytes.Equal(a, b)
}

// New creates a new IE.
func New(itype IEType, data []byte) *IE {
	i := &IE{
		Type:    itype,
		Payload: data,
	}

	return i
}

// NewVendorSpecificIE creates a new vendor-specific IE.
func NewVendorSpecificIE(itype IEType, eid uint16, data []byte) *IE {
	i := &IE{
		Type:         itype,
		EnterpriseID: eid,
		Payload:      data,
	}

	return i
}

// NewGroupedIE creates a new grouped IE.
func NewGroupedIE(itype IEType, ies ...*IE) *IE {
	return newGroupedIE(itype, 0, ies...)
}

// NewVendorSpecificGroupedIE creates a new grouped IE.
func NewVendorSpecificGroupedIE(itype IEType, eid uint16, ies ...*IE) *IE {
	return newGroupedIE(itype, eid, ies...)
}

// NewUint8ValIE creates a new IE with uint8 value.
func NewUint8IE(itype IEType, v uint8) *IE {
	return newUint8ValIE(itype, v)
}

// NewUint16ValIE creates a new IE with uint16 value.
func NewUint16IE(itype IEType, v uint16) *IE {
	return newUint16ValIE(itype, v)
}

// NewUint32ValIE creates a new IE with uint32 value.
func NewUint32IE(itype IEType, v uint32) *IE {
	return newUint32ValIE(itype, v)
}

// NewUint64ValIE creates a new IE with uint64 value.
func NewUint64IE(itype IEType, v uint64) *IE {
	return newUint64ValIE(itype, v)
}

// NewStringIE creates a new IE with string value.
func NewStringIE(itype IEType, v string) *IE {
	return newStringIE(itype, v)
}

// NewFQDNIE creates a new IE with FQDN value.
func NewFQDNIE(itype IEType, v string) *IE {
	return newFQDNIE(itype, v)
}

// ValueAsUint8 returns the value of IE as uint8.
func (i *IE) ValueAsUint8() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	if i.IsGrouped() {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 1 {
		return 0, io.ErrUnexpectedEOF
	}

	return i.Payload[0], nil
}

// ValueAsUint16 returns the value of IE as uint16.
func (i *IE) ValueAsUint16() (uint16, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	if i.IsGrouped() {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	return binary.BigEndian.Uint16(i.Payload[0:2]), nil
}

// ValueAsUint32 returns the value of IE as uint32.
func (i *IE) ValueAsUint32() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	if i.IsGrouped() {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 4 {
		return 0, io.ErrUnexpectedEOF
	}

	return binary.BigEndian.Uint32(i.Payload[0:4]), nil
}

// ValueAsUint64 returns the value of IE as uint64.
func (i *IE) ValueAsUint64() (uint64, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	if i.IsGrouped() {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 8 {
		return 0, io.ErrUnexpectedEOF
	}

	return binary.BigEndian.Uint64(i.Payload[0:8]), nil
}

// ValueAsString returns the value of IE as string.
func (i *IE) ValueAsString() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}

	if i.IsGrouped() {
		return "", &InvalidTypeError{Type: i.Type}
	}

	// TS 29.244 V20.1.0 clause 8.1.2 permits a null IE as a removal
	// instruction; it does not supply a string value.
	if len(i.Payload) == 0 {
		return "", io.ErrUnexpectedEOF
	}

	return string(i.Payload), nil
}

// ValueAsFQDN returns the value of IE as string, decoded as FQDN.
func (i *IE) ValueAsFQDN() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}

	if i.IsGrouped() {
		return "", &InvalidTypeError{Type: i.Type}
	}

	return decodeFQDNValue(i.Payload)
}

// ValueAsGrouped returns the value of IE as grouped IE.
//
// This method returns the ChildIEs field if it is already parsed.
// Otherwise, it parses the Payload field and returns the result.
//
// Any IE, including a vendor-specific IE, can be explicitly interpreted as a group.
// The type-specific grouped helpers use this same source of truth.
func (i *IE) ValueAsGrouped() ([]*IE, error) {
	if i == nil {
		return nil, ErrIENotFound
	}
	if i.ChildIEs != nil {
		for n, child := range i.ChildIEs {
			if child != nil {
				continue
			}
			// Direct struct construction can introduce nil children. Encoding
			// skips them; helpers must use the same list without mutating it.
			children := make([]*IE, 0, len(i.ChildIEs)-1)
			children = append(children, i.ChildIEs[:n]...)
			for _, next := range i.ChildIEs[n+1:] {
				if next != nil {
					children = append(children, next)
				}
			}
			return children, nil
		}
		return i.ChildIEs, nil
	}
	return ParseMultiIEs(i.Payload)
}

// valueAs3GPPTimestamp uses era 1 for wire seconds below 2208988800, and
// era 0 otherwise. This selects the Unix window [1970-01-01, 2106-02-07T06:28:15Z].
// The wire carries no era, so earlier/later dates cannot be distinguished.
// RFC 5905 section 6; TS 29.244 V20.1.0 clause 8.2.65.
// It returns the value of IE as time.Time, in the format of
// timestamp IEs defined in 3GPP TS 29.244 as follows:
//
// "a UTC time. Octets 5 to 8 shall be encoded in the same format as the first four octets
// of the 64-bit timestamp format as defined in clause 6 of IETF RFC 5905 [12].
// NOTE: The encoding is defined as the time in seconds relative to 00:00:00 on 1 January 1900."
func (i *IE) valueAs3GPPTimestamp() (time.Time, error) {
	if i == nil {
		return time.Time{}, ErrIENotFound
	}

	if i.IsGrouped() {
		return time.Time{}, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 4 {
		return time.Time{}, io.ErrUnexpectedEOF
	}

	return time.Unix(int64(binary.BigEndian.Uint32(i.Payload[0:4])-2208988800), 0), nil
}

// Parse parses b into IE.
//
// Note that this function uses the given bytes directly, so not safe to use
// the buffer after calling this function. When you use the buffer somewhere
// else, copy it before calling this function.
func Parse(b []byte) (*IE, error) {
	i := &IE{}
	if err := i.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return i, nil
}

// maxGroupedDepth bounds how many grouped IEs may enclose an IE the decoder
// follows. TS 29.244 V20.1.0 nests grouped IEs at most four levels below a
// message (clause 7 tables, e.g. Create PDR > PDI > Ethernet Packet Filter),
// so deeper input is not a PFCP structure; without a bound one 64 KiB
// datagram could force thousands of decoding levels.
const maxGroupedDepth = 16

// ParseMultiIEs decodes a sequence of IEs that fills b exactly, such as the
// IEs of a message or the value of a grouped IE.
//
// The returned IEs reference b: modifying b afterwards changes them. Copy b
// first if it is reused.
func ParseMultiIEs(b []byte) ([]*IE, error) {
	return decodeIEs(b, 0)
}

// decodeIEs decodes the IEs filling b. depth counts the grouped IEs that
// enclose them.
func decodeIEs(b []byte, depth int) ([]*IE, error) {
	var ies []*IE
	for len(b) > 0 {
		i, n, err := decodeIE(b, depth)
		if err != nil {
			return nil, err
		}
		ies = append(ies, i)
		b = b[n:]
	}
	return ies, nil
}

// decodeIE decodes the IE at the start of b and returns it with the number of
// octets it occupies, so that a caller advances by the declared length rather
// than by re-measuring the decoded IE.
// TS 29.244 V20.1.0 clause 8.1.1: Length excludes the four-octet header,
// but includes the Enterprise ID of a vendor-specific IE.
func decodeIE(b []byte, depth int) (*IE, int, error) {
	if len(b) < 4 {
		return nil, 0, io.ErrUnexpectedEOF
	}
	typ := IEType(binary.BigEndian.Uint16(b[:2]))
	length := int(binary.BigEndian.Uint16(b[2:4]))
	vendor := typ&0x8000 != 0
	if vendor && length < 2 {
		return nil, 0, &DecodeError{Type: typ, Err: ErrInvalidLength}
	}
	end := 4 + length
	if len(b) < end {
		return nil, 0, &DecodeError{Type: typ, Err: io.ErrUnexpectedEOF}
	}
	v := &IE{Type: typ}
	offset := 4
	if vendor {
		v.EnterpriseID = binary.BigEndian.Uint16(b[4:6])
		offset = 6
	}
	payload := b[offset:end:end]
	if v.IsGrouped() {
		if depth >= maxGroupedDepth {
			return nil, 0, &DecodeError{Type: typ, Err: ErrNestingTooDeep}
		}
		children, err := decodeIEs(payload, depth+1)
		if err != nil {
			return nil, 0, &DecodeError{Type: typ, Err: err}
		}
		if children == nil {
			children = make([]*IE, 0)
		}
		v.ChildIEs = children
	} else {
		v.Payload = payload
	}
	return v, end, nil
}

// UnmarshalBinary parses b, which must hold exactly one IE, into i.
// Octets after the IE's declared length are an error (ErrInvalidLength):
// use ParseMultiIEs for a sequence of IEs.
func (i *IE) UnmarshalBinary(b []byte) error {
	if i == nil {
		return ErrIENotFound
	}
	v, n, err := decodeIE(b, 0)
	if err != nil {
		return err
	}
	if n != len(b) {
		return &DecodeError{Type: v.Type, Err: ErrInvalidLength}
	}
	*i = *v
	return nil
}

// Marshal returns the byte sequence generated from an IE instance.
func (i *IE) Marshal() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}
	length := i.MarshalLen()
	if length-4 > 65535 {
		return nil, ErrInvalidLength
	}
	b := make([]byte, length)
	if err := i.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo writes an IE to b. It returns ErrInvalidLength for a value exceeding
// the 16-bit Length field or for a buffer shorter than MarshalLen().
func (i *IE) MarshalTo(b []byte) error {
	if i == nil {
		return ErrIENotFound
	}
	if i.err != nil {
		return &EncodeError{Type: i.Type, Err: i.err}
	}
	length := i.MarshalLen()
	if length-4 > 65535 || len(b) < length {
		return ErrInvalidLength
	}
	binary.BigEndian.PutUint16(b[:2], uint16(i.Type))
	binary.BigEndian.PutUint16(b[2:4], uint16(length-4))
	offset := 4
	if i.IsVendorSpecific() {
		binary.BigEndian.PutUint16(b[4:6], i.EnterpriseID)
		offset = 6
	}
	if i.ChildIEs != nil {
		for _, child := range i.ChildIEs {
			if child == nil {
				continue
			}
			if err := child.MarshalTo(b[offset:]); err != nil {
				return fmt.Errorf("child of IE %d: %w", i.Type, err)
			}
			offset += child.MarshalLen()
		}
	} else {
		copy(b[offset:length], i.Payload)
	}
	return nil
}

// MarshalLen returns the encoded length, including the IE header, or zero for nil.
func (i *IE) MarshalLen() int {
	if i == nil {
		return 0
	}
	length := 4
	if i.IsVendorSpecific() {
		length += 2
	}
	if i.ChildIEs != nil {
		for _, child := range i.ChildIEs {
			length += child.MarshalLen()
		}
	} else {
		length += len(i.Payload)
	}
	return length
}

// IsVendorSpecific reports whether an IE is vendor-specific or defined by 3gpp.
func (i *IE) IsVendorSpecific() bool {
	if i == nil {
		return false
	}

	// Spef: TS 29.244 8.1.1 Information Element Format
	// If the Bit 8 of Octet 1 is not set, this indicates that the IE is defined by 3GPP
	// and the EnterpriseID is absent. If Bit 8 of Octet 1 is set, this indicates that the
	// IE is defined by a vendor and the Enterprise ID is present identified by the Enterprise ID.
	return i.Type&0x8000 != 0
}

func newUint8ValIE(t IEType, v uint8) *IE {
	return New(t, []byte{v})
}

func newUint16ValIE(t IEType, v uint16) *IE {
	i := New(t, make([]byte, 2))
	binary.BigEndian.PutUint16(i.Payload, v)
	return i
}

func newUint32ValIE(t IEType, v uint32) *IE {
	i := New(t, make([]byte, 4))
	binary.BigEndian.PutUint32(i.Payload, v)
	return i
}

func newUint64ValIE(t IEType, v uint64) *IE {
	i := New(t, make([]byte, 8))
	binary.BigEndian.PutUint64(i.Payload, v)
	return i
}

func newStringIE(t IEType, v string) *IE {
	return New(t, []byte(v))
}

func newFQDNIE(t IEType, v string) *IE {
	encoded := utils.EncodeFQDN(v)
	if encoded == nil {
		return newInvalidIE(t, invalidFQDN(v))
	}
	return New(t, encoded)
}

// invalidFQDN describes a name EncodeFQDN rejects (RFC 1035 clause 3.1:
// labels of 1 to 63 octets, 255 octets in total).
func invalidFQDN(v string) error {
	return fmt.Errorf("%w: %q is not a domain name (RFC 1035 clause 3.1)", ErrInvalidValue, v)
}

func newGroupedIE(itype IEType, eid uint16, ies ...*IE) *IE {
	i := &IE{Type: itype, EnterpriseID: eid, ChildIEs: make([]*IE, 0, len(ies))}
	for _, child := range ies {
		if child != nil {
			i.ChildIEs = append(i.ChildIEs, child)
		}
	}
	return i
}

// new3GPPTimestampIE applies the era policy used by valueAs3GPPTimestamp.
// TS 29.244 V20.1.0 clause 8.2.65 and RFC 5905 section 6 specify seconds
// since 1900 modulo 2^32. With no era carried on the wire, this package chooses
// the Unix window [1970-01-01, 2106-02-07T06:28:15Z]. A time outside that
// window gives an IE whose Marshal reports ErrInvalidValue. Fractional seconds
// are discarded.
func new3GPPTimestampIE(typ IEType, ts time.Time) *IE {
	seconds := ts.Unix()
	if seconds < 0 || seconds > 1<<32-1 {
		return newInvalidIE(typ, fmt.Errorf("%w: %s is outside 1970-01-01T00:00:00Z to 2106-02-07T06:28:15Z",
			ErrInvalidValue, ts.UTC().Format(time.RFC3339)))
	}
	return newUint32ValIE(typ, uint32(seconds+2208988800))
}

func decodeFQDNValue(payload []byte) (string, error) {
	name := utils.DecodeFQDN(payload)
	if name == "" && (len(payload) != 1 || payload[0] != 0) {
		return "", ErrMalformed
	}
	return name, nil
}
