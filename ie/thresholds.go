package ie

import (
	"time"
)

// ThresholdFlags identifies the supplied thresholds in clause 8.2.196.
type ThresholdFlags uint8

const (
	ThresholdRTT ThresholdFlags = 1
	ThresholdPLR ThresholdFlags = 2
)

// ThresholdsFields is the structured value of Thresholds.
// TS 29.244 V20.1.0 clause 8.2.196.
type ThresholdsFields struct {
	Flags          ThresholdFlags
	RTT            time.Duration
	PacketLossRate uint8
}

// NewThresholds creates a Thresholds IE from its fields.
// TS 29.244 V20.1.0 clause 8.2.196.
func NewThresholds(fields ThresholdsFields) *IE {
	b, err := fields.Marshal()
	return addedFieldsIE(Thresholds, b, err)
}

// Thresholds returns the structured Thresholds value.
// TS 29.244 V20.1.0 clauses 8.2.196 and 8.1.2.
func (i *IE) Thresholds() (*ThresholdsFields, error) {
	b, err := addedValuePayload(i, Thresholds, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseThresholdsFields(b)
	if err != nil {
		return nil, &DecodeError{Type: Thresholds, Err: err}
	}
	return f, nil
}

// ParseThresholdsFields decodes the known fields; unknown extension octets are ignored.
// TS 29.244 V20.1.0 clauses 8.2.196 and 8.1.2.
func ParseThresholdsFields(b []byte) (*ThresholdsFields, error) {
	f := &ThresholdsFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary decodes Thresholds and resets previously decoded fields.
func (f *ThresholdsFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = ThresholdsFields{}
	r := &addedReader{b: b}
	f.Flags = ThresholdFlags(r.u8() & 3)
	if f.Flags&ThresholdRTT != 0 {
		f.RTT = time.Duration(r.u16()) * time.Millisecond
	}
	if f.Flags&ThresholdPLR != 0 {
		f.PacketLossRate = r.u8()
		if f.PacketLossRate > 100 {
			f.PacketLossRate = 0
		}
	}
	if r.err != nil {
		return r.err
	}

	if !f.valid() {
		return ErrInvalidValue
	}
	return nil
}
func (f *ThresholdsFields) valid() bool {
	return f.Flags >= 1 && f.Flags <= 3 && addedDurationFits(f.RTT, time.Millisecond, 65535) && f.PacketLossRate <= 100 && (f.Flags&ThresholdRTT != 0 || f.RTT == 0) && (f.Flags&ThresholdPLR != 0 || f.PacketLossRate == 0)
}

// Marshal encodes the known fields with spare bits set to zero.
func (f *ThresholdsFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	if !f.valid() {
		return nil, ErrInvalidValue
	}
	b := []byte{byte(f.Flags)}
	if f.Flags&ThresholdRTT != 0 {
		b = addedPut16(b, uint16(f.RTT/time.Millisecond))
	}
	if f.Flags&ThresholdPLR != 0 {
		b = append(b, f.PacketLossRate)
	}
	if len(b) > 65535 {
		return nil, ErrInvalidValue
	}
	return b, nil
}

// MarshalLen returns the encoded value length, or zero for nil or invalid fields.
func (f *ThresholdsFields) MarshalLen() int { b, _ := f.Marshal(); return len(b) }

// MarshalTo writes the encoded fields to b, rejecting an undersized buffer.
func (f *ThresholdsFields) MarshalTo(b []byte) error {
	wire, err := f.Marshal()
	return addedCopyTo(b, wire, err)
}
