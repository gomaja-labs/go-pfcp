package ie

import (
	"time"
)

// N6JitterMeasurementFlags identifies the downlink fields in clause 8.2.222.
type N6JitterMeasurementFlags uint8

const N6JitterDL N6JitterMeasurementFlags = 1

// N6JitterMeasurementFields is the structured value of N6 Jitter Measurement.
// TS 29.244 V20.1.0 clause 8.2.222.
type N6JitterMeasurementFields struct {
	Flags                                        N6JitterMeasurementFlags
	DLPeriodicity, LowerDLJitter, HigherDLJitter time.Duration
}

// NewN6JitterMeasurement creates a N6 Jitter Measurement IE from its fields.
// TS 29.244 V20.1.0 clause 8.2.222.
func NewN6JitterMeasurement(fields N6JitterMeasurementFields) *IE {
	b, err := fields.Marshal()
	return addedFieldsIE(N6JitterMeasurement, b, err)
}

// N6JitterMeasurement returns the structured N6 Jitter Measurement value.
// TS 29.244 V20.1.0 clauses 8.2.222 and 8.1.2.
func (i *IE) N6JitterMeasurement() (*N6JitterMeasurementFields, error) {
	b, err := addedValuePayload(i, N6JitterMeasurement, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseN6JitterMeasurementFields(b)
	if err != nil {
		return nil, &DecodeError{Type: N6JitterMeasurement, Err: err}
	}
	return f, nil
}

// ParseN6JitterMeasurementFields decodes the known fields; unknown extension octets are ignored.
// TS 29.244 V20.1.0 clauses 8.2.222 and 8.1.2.
func ParseN6JitterMeasurementFields(b []byte) (*N6JitterMeasurementFields, error) {
	f := &N6JitterMeasurementFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary decodes N6 Jitter Measurement and resets previously decoded fields.
func (f *N6JitterMeasurementFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = N6JitterMeasurementFields{}
	r := &addedReader{b: b}
	f.Flags = N6JitterMeasurementFlags(r.u8() & 1)
	if f.Flags&N6JitterDL != 0 {
		f.DLPeriodicity = time.Duration(r.u32()) * time.Millisecond
		f.LowerDLJitter = time.Duration(int32(r.u32())) * time.Millisecond
		f.HigherDLJitter = time.Duration(int32(r.u32())) * time.Millisecond
	}
	if r.err != nil {
		return r.err
	}

	if !f.valid() {
		return ErrInvalidValue
	}
	return nil
}
func (f *N6JitterMeasurementFields) valid() bool {
	return f.Flags&^1 == 0 && addedDurationFits(f.DLPeriodicity, time.Millisecond, 0xffffffff) && addedSignedMilliseconds(f.LowerDLJitter) && addedSignedMilliseconds(f.HigherDLJitter) && (f.Flags&N6JitterDL != 0 || (f.DLPeriodicity == 0 && f.LowerDLJitter == 0 && f.HigherDLJitter == 0))
}

// Marshal encodes the known fields with spare bits set to zero.
func (f *N6JitterMeasurementFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	if !f.valid() {
		return nil, ErrInvalidValue
	}
	b := []byte{byte(f.Flags)}
	if f.Flags&N6JitterDL != 0 {
		b = addedPut32(b, uint32(f.DLPeriodicity/time.Millisecond))
		b = addedPut32(b, uint32(int32(f.LowerDLJitter/time.Millisecond)))
		b = addedPut32(b, uint32(int32(f.HigherDLJitter/time.Millisecond)))
	}
	if len(b) > 65535 {
		return nil, ErrInvalidValue
	}
	return b, nil
}

// MarshalLen returns the encoded value length, or zero for nil or invalid fields.
func (f *N6JitterMeasurementFields) MarshalLen() int { b, _ := f.Marshal(); return len(b) }

// MarshalTo writes the encoded fields to b, rejecting an undersized buffer.
func (f *N6JitterMeasurementFields) MarshalTo(b []byte) error {
	wire, err := f.Marshal()
	return addedCopyTo(b, wire, err)
}
func addedSignedMilliseconds(v time.Duration) bool {
	return v%time.Millisecond == 0 && v/time.Millisecond >= -0x80000000 && v/time.Millisecond <= 0x7fffffff
}
