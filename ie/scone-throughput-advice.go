package ie

import (
	"io"
)

// SCONEThroughputAdvice flags and values, TS 29.244 V20.1.0 clause 8.2.286.
const (
	// SCONEThroughputAdviceULSI selects ULSI.
	SCONEThroughputAdviceULSI uint8 = 1
	// SCONEThroughputAdviceDLSI selects DLSI.
	SCONEThroughputAdviceDLSI uint8 = 2
)

// NewSCONEThroughputAdvice creates a SCONE Throughput Advice IE.
// TS 29.244 V20.1.0 clause 8.2.286.
func NewSCONEThroughputAdvice(flags uint8, uplinkThroughputAdvice uint32, downlinkThroughputAdvice uint32) *IE {
	return recentFieldsIE(SCONEThroughputAdvice, NewSCONEThroughputAdviceFields(flags, uplinkThroughputAdvice, downlinkThroughputAdvice))
}

// SCONEThroughputAdvice returns the SCONE Throughput Advice fields.
// TS 29.244 V20.1.0 clause 8.2.286.
func (i *IE) SCONEThroughputAdvice() (*SCONEThroughputAdviceFields, error) {
	b, err := recentIEPayload(i, SCONEThroughputAdvice, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseSCONEThroughputAdviceFields(b)
	if err != nil {
		return nil, &DecodeError{Type: SCONEThroughputAdvice, Err: err}
	}
	return f, nil
}

// SCONEThroughputAdviceFields models TS 29.244 V20.1.0 clause 8.2.286.
// Flags select optional fields; nonzero values for absent fields are rejected.
// Rates are integer kilobits per second (1000 bits/s).
type SCONEThroughputAdviceFields struct {
	// Flags contains the presence and option bits.
	Flags uint8
	// UplinkThroughputAdvice is the uplink throughput advice.
	UplinkThroughputAdvice uint32
	// DownlinkThroughputAdvice is the downlink throughput advice.
	DownlinkThroughputAdvice uint32
}

// NewSCONEThroughputAdviceFields stores all arguments; Marshal validates their representation.
func NewSCONEThroughputAdviceFields(flags uint8, uplinkThroughputAdvice uint32, downlinkThroughputAdvice uint32) *SCONEThroughputAdviceFields {
	return &SCONEThroughputAdviceFields{flags, uplinkThroughputAdvice, downlinkThroughputAdvice}
}

// ParseSCONEThroughputAdviceFields decodes fields, ignoring spare bits and trailing extension octets.
func ParseSCONEThroughputAdviceFields(b []byte) (*SCONEThroughputAdviceFields, error) {
	f := &SCONEThroughputAdviceFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}
func (f *SCONEThroughputAdviceFields) validate() error {
	if f == nil {
		return ErrIENotFound
	}
	if f.Flags & ^uint8(3) != 0 {
		return ErrInvalidValue
	}
	if f.Flags&1 == 0 && f.UplinkThroughputAdvice != 0 {
		return ErrInvalidValue
	}
	if f.Flags&2 == 0 && f.DownlinkThroughputAdvice != 0 {
		return ErrInvalidValue
	}
	return nil
}

// MarshalLen returns the known field size in octets.
func (f *SCONEThroughputAdviceFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	n := 1
	if f.Flags&1 != 0 {
		n += 4
	}
	if f.Flags&2 != 0 {
		n += 4
	}
	return n
}

// Marshal serializes fields as defined in clause 8.2.286.
func (f *SCONEThroughputAdviceFields) Marshal() ([]byte, error) {
	if err := f.validate(); err != nil {
		return nil, err
	}
	b := make([]byte, f.MarshalLen())
	if err := f.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo writes fields to b, rejecting invalid values and short buffers.
func (f *SCONEThroughputAdviceFields) MarshalTo(b []byte) error {
	if err := f.validate(); err != nil {
		return err
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}
	b[0] = f.Flags
	offset := 1
	if f.Flags&1 != 0 {
		putRecentUint(b[offset:offset+4], uint64(f.UplinkThroughputAdvice))
		offset += 4
	}
	if f.Flags&2 != 0 {
		putRecentUint(b[offset:offset+4], uint64(f.DownlinkThroughputAdvice))
	}
	return nil
}

// UnmarshalBinary replaces the receiver with the decoded known fields.
func (f *SCONEThroughputAdviceFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = SCONEThroughputAdviceFields{}
	if len(b) < 1 {
		return io.ErrUnexpectedEOF
	}
	f.Flags = b[0] & 3
	offset := 1
	if f.Flags&1 != 0 {
		if len(b)-offset < 4 {
			return io.ErrUnexpectedEOF
		}
		f.UplinkThroughputAdvice = uint32(recentUint(b[offset : offset+4]))
		offset += 4
	}
	if f.Flags&2 != 0 {
		if len(b)-offset < 4 {
			return io.ErrUnexpectedEOF
		}
		f.DownlinkThroughputAdvice = uint32(recentUint(b[offset : offset+4]))
	}
	return f.validate()
}
