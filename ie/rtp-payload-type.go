package ie

import (
	"fmt"
	"io"
)

// NewRTPPayloadType creates a RTP Payload Type IE.
// TS 29.244 V20.1.0 clause 8.2.236.
func NewRTPPayloadType(v uint8) *IE {
	if v > 127 {
		return invalidRecentIE(RTPPayloadType, "RTP Payload Type")
	}
	return New(RTPPayloadType, []byte{v})
}

// RTPPayloadType returns the RTP Payload Type value.
// TS 29.244 V20.1.0 clause 8.2.236; clause 8.1.2 governs spare bits and extensions.
func (i *IE) RTPPayloadType() (uint8, error) {
	b, err := recentIEPayload(i, RTPPayloadType, 1, true)
	if err != nil {
		return 0, err
	}
	v := b[0] & 255

	if v > 127 {
		return 0, &DecodeError{Type: RTPPayloadType, Err: ErrInvalidValue}
	}
	return v, nil
}

// recentIEPayload enforces the length policy of TS 29.244 V20.1.0 clause
// 8.1.2: a null IE is decoded by Parse, but cannot supply a value; fixed
// values have exactly their specified length, extendable values may grow.
func recentIEPayload(i *IE, typ IEType, minimum int, fixed bool) ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}
	if i.Type != typ {
		return nil, &InvalidTypeError{Type: i.Type}
	}
	if i.err != nil {
		return nil, i.err
	}
	if len(i.Payload) < minimum {
		return nil, &DecodeError{Type: typ, Err: io.ErrUnexpectedEOF}
	}
	if fixed && len(i.Payload) != minimum {
		return nil, &DecodeError{Type: typ, Err: ErrInvalidLength}
	}
	return i.Payload, nil
}

func invalidRecentIE(typ IEType, field string) *IE {
	return newInvalidIE(typ, fmt.Errorf("%w: %s", ErrInvalidValue, field))
}

func recentFieldsIE(typ IEType, f interface{ Marshal() ([]byte, error) }) *IE {
	b, err := f.Marshal()
	if err != nil {
		return newInvalidIE(typ, fmt.Errorf("%w: %w", ErrInvalidValue, err))
	}
	return New(typ, b)
}

// recentUint and putRecentUint handle the 24- and 40-bit unsigned fields in
// addition to conventional widths, in network octet order.
func recentUint(b []byte) uint64 {
	var v uint64
	for _, x := range b {
		v = v<<8 | uint64(x)
	}
	return v
}
func putRecentUint(b []byte, v uint64) {
	for n := len(b) - 1; n >= 0; n-- {
		b[n] = byte(v)
		v >>= 8
	}
}
