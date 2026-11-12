package ie

import (
	"time"
)

// TrafficParameterThresholdFlags identifies the downlink jitter field in clause 8.2.220.
type TrafficParameterThresholdFlags uint8

// TrafficParameterThresholdDL selects DL.
const TrafficParameterThresholdDL TrafficParameterThresholdFlags = 1

// TrafficParameterThresholdFields is the structured value of Traffic Parameter Threshold.
// TS 29.244 V20.1.0 clause 8.2.220.
type TrafficParameterThresholdFields struct {
	// Flags contains the presence and option bits.
	Flags TrafficParameterThresholdFlags
	// DownlinkJitter is the downlink jitter.
	DownlinkJitter time.Duration
}

// NewTrafficParameterThreshold creates a Traffic Parameter Threshold IE from its fields.
// TS 29.244 V20.1.0 clause 8.2.220.
func NewTrafficParameterThreshold(fields TrafficParameterThresholdFields) *IE {
	b, err := fields.Marshal()
	return addedFieldsIE(TrafficParameterThreshold, b, err)
}

// TrafficParameterThreshold returns the structured Traffic Parameter Threshold value.
// TS 29.244 V20.1.0 clauses 8.2.220 and 8.1.2.
func (i *IE) TrafficParameterThreshold() (*TrafficParameterThresholdFields, error) {
	b, err := addedValuePayload(i, TrafficParameterThreshold, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseTrafficParameterThresholdFields(b)
	if err != nil {
		return nil, &DecodeError{Type: TrafficParameterThreshold, Err: err}
	}
	return f, nil
}

// ParseTrafficParameterThresholdFields decodes the known fields; unknown extension octets are ignored.
// TS 29.244 V20.1.0 clauses 8.2.220 and 8.1.2.
func ParseTrafficParameterThresholdFields(b []byte) (*TrafficParameterThresholdFields, error) {
	f := &TrafficParameterThresholdFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary decodes Traffic Parameter Threshold and resets previously decoded fields.
func (f *TrafficParameterThresholdFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = TrafficParameterThresholdFields{}
	r := &addedReader{b: b}
	f.Flags = TrafficParameterThresholdFlags(r.u8() & 1)
	if f.Flags&TrafficParameterThresholdDL != 0 {
		f.DownlinkJitter = time.Duration(r.u32()) * time.Millisecond
	}
	if r.err != nil {
		return r.err
	}

	if !f.valid() {
		return ErrInvalidValue
	}
	return nil
}
func (f *TrafficParameterThresholdFields) valid() bool {
	return f.Flags&^1 == 0 && addedDurationFits(f.DownlinkJitter, time.Millisecond, 0xffffffff) && (f.Flags&TrafficParameterThresholdDL != 0 || f.DownlinkJitter == 0)
}

// Marshal encodes the known fields with spare bits set to zero.
func (f *TrafficParameterThresholdFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	if !f.valid() {
		return nil, ErrInvalidValue
	}
	b := []byte{byte(f.Flags)}
	if f.Flags&TrafficParameterThresholdDL != 0 {
		b = addedPut32(b, uint32(f.DownlinkJitter/time.Millisecond))
	}
	if len(b) > 65535 {
		return nil, ErrInvalidValue
	}
	return b, nil
}

// MarshalLen returns the encoded value length, or zero for nil or invalid fields.
func (f *TrafficParameterThresholdFields) MarshalLen() int { b, _ := f.Marshal(); return len(b) }

// MarshalTo writes the encoded fields to b, rejecting an undersized buffer.
func (f *TrafficParameterThresholdFields) MarshalTo(b []byte) error {
	wire, err := f.Marshal()
	return addedCopyTo(b, wire, err)
}
