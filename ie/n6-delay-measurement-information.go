package ie

import (
	"io"
	"time"
)

// N6DelayMeasurementInformation flags and values, TS 29.244 V20.1.0 clause 8.2.247.
const (
	N6DelayMeasurementInformationMIND uint8 = 1
	N6DelayMeasurementInformationMAXD uint8 = 2
)

// NewN6DelayMeasurementInformation creates a N6 Delay Measurement Information IE.
// TS 29.244 V20.1.0 clause 8.2.247.
func NewN6DelayMeasurementInformation(flags uint8, averagePacketDelay time.Duration, minimumPacketDelay time.Duration, maximumPacketDelay time.Duration) *IE {
	return recentFieldsIE(N6DelayMeasurementInformation, NewN6DelayMeasurementInformationFields(flags, averagePacketDelay, minimumPacketDelay, maximumPacketDelay))
}

// N6DelayMeasurementInformation returns the N6 Delay Measurement Information fields.
// TS 29.244 V20.1.0 clause 8.2.247.
func (i *IE) N6DelayMeasurementInformation() (*N6DelayMeasurementInformationFields, error) {
	b, err := recentIEPayload(i, N6DelayMeasurementInformation, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseN6DelayMeasurementInformationFields(b)
	if err != nil {
		return nil, &DecodeError{Type: N6DelayMeasurementInformation, Err: err}
	}
	return f, nil
}

// N6DelayMeasurementInformationFields models TS 29.244 V20.1.0 clause 8.2.247.
// Flags select optional fields; nonzero values for absent fields are rejected.
// Packet delays must be nonnegative integral milliseconds fitting in 32 bits.
type N6DelayMeasurementInformationFields struct {
	Flags              uint8
	AveragePacketDelay time.Duration
	MinimumPacketDelay time.Duration
	MaximumPacketDelay time.Duration
}

// NewN6DelayMeasurementInformationFields stores all arguments; Marshal validates their representation.
func NewN6DelayMeasurementInformationFields(flags uint8, averagePacketDelay time.Duration, minimumPacketDelay time.Duration, maximumPacketDelay time.Duration) *N6DelayMeasurementInformationFields {
	return &N6DelayMeasurementInformationFields{flags, averagePacketDelay, minimumPacketDelay, maximumPacketDelay}
}

// ParseN6DelayMeasurementInformationFields decodes fields, ignoring spare bits and trailing extension octets.
func ParseN6DelayMeasurementInformationFields(b []byte) (*N6DelayMeasurementInformationFields, error) {
	f := &N6DelayMeasurementInformationFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}
func (f *N6DelayMeasurementInformationFields) validate() error {
	if f == nil {
		return ErrIENotFound
	}
	if f.Flags & ^uint8(3) != 0 {
		return ErrInvalidValue
	}
	if f.AveragePacketDelay < 0 || f.AveragePacketDelay%time.Millisecond != 0 || f.AveragePacketDelay/time.Millisecond > 4294967295 {
		return ErrInvalidValue
	}
	if f.Flags&1 == 0 && f.MinimumPacketDelay != 0 {
		return ErrInvalidValue
	}
	if f.MinimumPacketDelay < 0 || f.MinimumPacketDelay%time.Millisecond != 0 || f.MinimumPacketDelay/time.Millisecond > 4294967295 {
		return ErrInvalidValue
	}
	if f.Flags&2 == 0 && f.MaximumPacketDelay != 0 {
		return ErrInvalidValue
	}
	if f.MaximumPacketDelay < 0 || f.MaximumPacketDelay%time.Millisecond != 0 || f.MaximumPacketDelay/time.Millisecond > 4294967295 {
		return ErrInvalidValue
	}
	return nil
}

// MarshalLen returns the known field size in octets.
func (f *N6DelayMeasurementInformationFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	n := 1
	n += 4
	if f.Flags&1 != 0 {
		n += 4
	}
	if f.Flags&2 != 0 {
		n += 4
	}
	return n
}

// Marshal serializes fields as defined in clause 8.2.247.
func (f *N6DelayMeasurementInformationFields) Marshal() ([]byte, error) {
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
func (f *N6DelayMeasurementInformationFields) MarshalTo(b []byte) error {
	if err := f.validate(); err != nil {
		return err
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}
	b[0] = f.Flags
	offset := 1
	putRecentUint(b[offset:offset+4], uint64(f.AveragePacketDelay/time.Millisecond))
	offset += 4
	if f.Flags&1 != 0 {
		putRecentUint(b[offset:offset+4], uint64(f.MinimumPacketDelay/time.Millisecond))
		offset += 4
	}
	if f.Flags&2 != 0 {
		putRecentUint(b[offset:offset+4], uint64(f.MaximumPacketDelay/time.Millisecond))
	}
	return nil
}

// UnmarshalBinary replaces the receiver with the decoded known fields.
func (f *N6DelayMeasurementInformationFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = N6DelayMeasurementInformationFields{}
	if len(b) < 1 {
		return io.ErrUnexpectedEOF
	}
	f.Flags = b[0] & 3
	offset := 1
	if len(b)-offset < 4 {
		return io.ErrUnexpectedEOF
	}
	f.AveragePacketDelay = time.Duration(recentUint(b[offset:offset+4])) * time.Millisecond
	offset += 4
	if f.Flags&1 != 0 {
		if len(b)-offset < 4 {
			return io.ErrUnexpectedEOF
		}
		f.MinimumPacketDelay = time.Duration(recentUint(b[offset:offset+4])) * time.Millisecond
		offset += 4
	}
	if f.Flags&2 != 0 {
		if len(b)-offset < 4 {
			return io.ErrUnexpectedEOF
		}
		f.MaximumPacketDelay = time.Duration(recentUint(b[offset:offset+4])) * time.Millisecond
	}
	return f.validate()
}
